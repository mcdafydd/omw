//go:generate go run -tags generate gen.go

package backend

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/pkg/errors"
	hook "github.com/robotn/gohook"
	"github.com/zserge/lorca"
)

type formatType int

const (
	// FormatJSON indicates that user requested JSON report format output
	FormatJSON = iota
	// FormatText indicates that user requested text template report format output
	FormatText
)

func (d formatType) String() string {
	return [...]string{"JSON", "Text"}[d]
}

// TemplateString defines the template used to output a Report() with FormatText
var TemplateString = `{{define "Entry"}}
({{- .Duration}}) {{.Start.Hour}}:{{.Start.Minute}}-{{.Ts.Hour}}:{{.Ts.Minute}} -- {{.Task -}}
{{end}}

Report Start: {{.From}}
Report End: {{.To}}
Total Task Hours: {{.TaskHrs}}
Total Break Hours: {{.BrkHrs}}
Total Ignore Hours: {{.IgnoreHrs}}
{{$day := "" }}
{{range .Entries}}
{{- if ne $day .Start.Weekday.String}}
{{$day = .Start.Weekday.String}}

----------------------- {{$day}}, {{.Start.Year}}-{{.Start.Month}}-{{.Start.Day}} -----------------------
{{end -}}
{{- template "Entry" .}}
{{- end -}}
`

// Backend represents the context and configuration of every instance of the omw command
// Immediate commands (like omw add, omw report), immediately affect the timesheet
// Long-running commands (like omw server), maintain a context
type Backend struct {
	ctx    context.Context
	config *config
	fp     *os.File
	worker *worker
}

// Entry describes a single entry in the timesheet
// Used by omw report
type Entry struct {
	Start    time.Time     `json:"startTime"`
	Ts       time.Time     `json:"timestamp"`
	Duration time.Duration `json:"duration"`
	Task     string        `json:"task"`
	Ignore   bool          `json:"ignore"`
	Brk      bool          `json:"break"`
}

// Report describes a report
// previous is only used during report calculation to
// populate Entry.Duration
type Report struct {
	From      time.Time     `json:"reportFrom"`
	To        time.Time     `json:"reportTo"`
	IgnoreHrs time.Duration `json:"ignoreTotalHours"`
	BrkHrs    time.Duration `json:"breakTotalHours"`
	TaskHrs   time.Duration `json:"taskTotalHours"`
	Entries   []Entry       `json:"entries"`
	previous  *time.Time
}

type config struct {
	omwDir  string
	omwFile string
}

// Go types that are bound to the Lorca UI must be thread-safe, because each
// binding is executed in its own goroutine. In this simple case we may use
// atomic operations, but for more complex cases one should use proper
// synchronization.
type worker struct {
	sync.Mutex
	cmd            string
	bounds         *lorca.Bounds
	ui             lorca.UI
	leftShiftDown  bool
	rightShiftDown bool
}

// Add appends the current time and task to your timesheet
func (b *Backend) Add(args []string) {
	if b.worker != nil {
		b.worker.Lock()
		defer b.worker.Unlock()
	}
	task := strings.Join(args, " ")
	b.addEntry(task)
	return
}

// Close cleans up before exiting
func (b *Backend) Close() (err error) {
	if b.worker != nil {
		b.worker.Lock()
		defer b.worker.Unlock()
	}
	if b.fp != nil {
		err = b.fp.Close()
	}
	return err
}

// Create an instance of the structures that operate on OMW data
func Create(fp *os.File, omwDir, omwFile string) *Backend {
	return &Backend{
		ctx: context.Background(),
		config: &config{
			omwDir:  omwDir,
			omwFile: omwFile,
		},
		fp:     fp,
		worker: nil,
	}
}

// Edit opens your current timesheet in your default editor or
// in the editor specified by the EDITOR environment variable
func (b *Backend) Edit() error {
	if b.worker != nil {
		b.worker.Lock()
		defer b.worker.Unlock()
	}
	editor := DefaultEditor
	if preferred := os.Getenv("EDITOR"); preferred != "" {
		editor = preferred
	}
	argv := []string{b.config.omwFile}
	cmd := exec.CommandContext(b.ctx, editor, argv...)
	// should work if run from terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return runCommand(cmd)
}

// Hello appends a newline and then another line to end of timesheet with current time
// and the word "Hello".  Meant to be run at the beginning of a new work day
func (b *Backend) Hello() {
	if b.worker != nil {
		b.worker.Lock()
		defer b.worker.Unlock()
	}
	b.fp.WriteString("\n")
	b.addEntry("hello")
	return
}

// Report outputs various report formats to specified type (for now - just text)
// We add 24 hours to the parsed end time so that when a user specifies
// --from 2019-01-01 --to 2019-01-02
// that translates to "report on tasks that occurred between 2019-01-01 00:00
// and "2019-01-03 00:00"
func (b *Backend) Report(start, end string, format string) (output string, report Report, err error) {
	layout := "2006-1-2" // should support optional leading zeros
	layoutEvent := "2006-1-2 15:04"
	if b.worker != nil {
		b.worker.Lock()
		defer b.worker.Unlock()
	}
	report.From, err = time.Parse(layout, start)
	if err != nil {
		return "", report, err
	}
	report.To, err = time.Parse(layout, end)
	if err != nil {
		return "", report, err
	}
	report.To = report.To.Add(24 * time.Hour)
	r, err := os.Open(b.config.omwFile)
	if err != nil {
		return "", report, err
	}
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		// Indicates line is missing required information
		if len(line) <= 2 {
			continue
		}
		ts, err := time.Parse(layoutEvent, strings.Join(line[:2], " "))
		if err != nil {
			continue
		}
		// Indicates task timestamp is outside the requested time period
		if ts.Before(report.From) || ts.After(report.To) {
			continue
		}
		entry, err := b.parseEntry(strings.Join(line[2:], " "))
		if err != nil {
			continue
		}
		entry.Ts = ts
		if err != nil {
			continue
		}
		// Should indicate first task in requested report time period
		if report.previous == nil {
			report.previous = &entry.Ts
			entry.Start = entry.Ts
			report.Entries = append(report.Entries, *entry)
			continue
		}
		// For now, we explicitly assume that a new day restarts the duration calculation
		// We may change the marker from new day to first entry of "hello" on a given day
		// to better allow tracking tasks that extend from a previous day into a new day
		if entry.Ts.Day() != (*report.previous).Day() {
			report.previous = &entry.Ts
			entry.Start = entry.Ts
		}
		entry.Start = *report.previous
		entry.Duration = entry.Ts.Sub(*report.previous)

		*report.previous = entry.Ts
		// Use else if to make it clear we only process the event's
		// duration one time
		if entry.Ignore == false && entry.Brk == false {
			report.TaskHrs += entry.Duration
		} else if entry.Ignore == true && entry.Brk == false {
			report.IgnoreHrs += entry.Duration
		} else if entry.Ignore == false && entry.Brk == true {
			report.BrkHrs += entry.Duration
		} else if entry.Ignore == true && entry.Brk == true {
			return "", report, errors.New("Entry has both break and ignore set to true.  Something's wrong")
		}
		report.Entries = append(report.Entries, *entry)
	}
	f := FormatText
	if format == "json" {
		f = FormatJSON
	}
	output, err = b.formatReport(report, formatType(f))
	return output, report, err
}

// Run does the following:
// 1. Creates the Lorca object
// 2. Loads the Chrome interface and HTML/JS content
// 3. Starts the hotkey listener
func (b *Backend) Run(args []string) error {
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}
	ui, err := lorca.New("", "", 480, 200, args...)
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	// A simple way to know when UI is ready (uses body.onload event in JS)
	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	// Create and bind Go object to the UI
	b.worker = &worker{ui: ui, cmd: ""}
	ui.Bind("OmwAdd", b.Add)
	ui.Bind("OmwEdit", b.Edit)
	ui.Bind("OmwHello", b.Hello)
	ui.Bind("OmwReport", b.Report)
	ui.Bind("OmwStretch", b.Stretch)
	ui.Bind("minimize", b.worker.Minimize)
	ui.Bind("restore", b.worker.Restore)

	// Load HTML.
	// You may also use `data:text/html,<base64>` approach to load initial HTML,
	// e.g: ui.Load("data:text/html," + url.PathEscape(html))

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(FS))
	ui.Load(fmt.Sprintf("http://%s", ln.Addr()))
	// You may use console.log to debug your JS code, it will be printed via
	// log.Println(). Also exceptions are printed in a similar manner.
	/*ui.Eval(`
		console.log('Multiple values:', [1, false, {"x":5}]);
	`)*/

	// Wait until the interrupt signal arrives or browser window is closed
	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)

	// start hook
	hotkey := hook.Start()
	// end hook
	defer hook.End()

	eventLoop(b.worker, &sigc, ui, &hotkey)

	return nil
}

// Stretch append current timestamp to end of timesheet and copy previous task
// fp is opened in append mode, so seek to beginning of file first
func (b *Backend) Stretch() error {
	if b.worker != nil {
		b.worker.Lock()
		defer b.worker.Unlock()
	}
	buf, err := ioutil.ReadFile(b.config.omwFile)
	if err != nil {
		errors.Wrapf(err, "Error reading %s", b.config.omwFile)
		return err
	}
	_, lastLine, err := bufio.ScanLines(buf, false)
	if err != nil {
		return err
	}
	if len(lastLine) <= 2 {
		return errors.New("Missing task description")
	}
	lastEntry := strings.Fields(string(lastLine))[2:]
	return b.addEntry(strings.Join(lastEntry, " "))
}

// addEntry seeks to end of file and appends a formatted string
// will create a new empty file if file is missing
func (b *Backend) addEntry(s string) error {
	fp, err := os.OpenFile(b.config.omwFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		errors.Wrapf(err, "Can't open or create %s", b.config.omwFile)
		return err
	}
	ts := time.Now()
	tsFmt := fmt.Sprintf("%d-%d-%d %d:%d",
		ts.Year(),
		ts.Month(),
		ts.Day(),
		ts.Hour(),
		ts.Minute(),
	)
	entry := fmt.Sprintf("%s\t%s\n", tsFmt, s)
	fp.WriteString(entry)
	fp.Close()
	return nil
}

func (b *Backend) parseEntry(s string) (*Entry, error) {
	re := regexp.MustCompile(`(?P<task>[a-zA-Z0-9,._+:@%/-]*) ?(?P<mod>\*\*\*?)*`)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return nil, errors.New("Invalid string")
	}
	entry := &Entry{
		Task: matches[1],
	}
	if matches[2] == "**" {
		entry.Brk = true
	}
	if matches[2] == "***" {
		entry.Ignore = true
	}
	return entry, nil
}

func (b *Backend) formatReport(report Report, format formatType) (string, error) {
	if format == FormatJSON {
		output, err := json.Marshal(report)
		return string(output), err
	}

	reportTmpl, err := template.New("report").Parse(TemplateString)
	if err != nil {
		return "", err
	}
	err = reportTmpl.Execute(os.Stdout, report)
	if err != nil {
		panic(err)
	}
	return "", nil
}

// Minimize Hides the application window
// Saves the current window lorca.Bounds
func (c *worker) Minimize() {
	c.Lock()
	defer c.Unlock()
	bounds, err := c.ui.Bounds()
	if err != nil {
		log.Println("[ERROR] Minimize.Bounds(): ", err)
		return
	}
	c.bounds = &bounds

	c.bounds.WindowState = lorca.WindowStateMinimized
	err = c.ui.SetBounds(*c.bounds)
	if err != nil {
		log.Println("[ERROR] Minimize.SetBounds(): ", err)
		return
	}
}

// Restore Restores previous visible window state after Minimize()
func (c *worker) Restore() {
	c.Lock()
	defer c.Unlock()
	bounds, err := c.ui.Bounds()
	if err != nil {
		log.Println("[ERROR] Minimize.Bounds(): ", err)
		return
	}
	c.bounds = &bounds

	c.bounds.WindowState = lorca.WindowStateNormal
	err = c.ui.SetBounds(*c.bounds)
	if err != nil {
		log.Println("[ERROR] Restore.SetBounds() WindowStateNormal: ", err)
		return
	}
}

// runCommand Executes cmd and handles any output
func runCommand(cmd *exec.Cmd) error {
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// eventLoop is the main loop that handles global hotkey events
func eventLoop(c *worker, sigc *chan os.Signal, ui lorca.UI, hotkey *chan hook.Event) {
	// main event loop
	keepLooping := true
	for keepLooping {
		select {
		case <-*sigc:
			keepLooping = false
			break
		case <-ui.Done():
			keepLooping = false
			break
		case ev := <-*hotkey:
			if ev.Rawcode == 65505 && ev.Kind == hook.KeyDown {
				fmt.Printf("Got left shift down = %#v\n", ev)
				c.leftShiftDown = true
			}
			if ev.Rawcode == 65506 && ev.Kind == hook.KeyDown {
				c.rightShiftDown = true
			}
			if ev.Rawcode == 65505 && ev.Kind == hook.KeyUp {
				c.leftShiftDown = false
			}
			if ev.Rawcode == 65506 && ev.Kind == hook.KeyUp {
				c.rightShiftDown = false
			}
			if c.leftShiftDown && c.rightShiftDown {
				log.Println("Got hotkey - restoring command window")
				c.Restore()
			}
		}
	}

	select {
	case <-*sigc:
	case <-ui.Done():
	}
}

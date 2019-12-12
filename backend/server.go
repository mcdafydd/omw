package backend

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"
)

type formatType int

const (
	// FormatFC indicates that user requested FC report format output
	FormatFC = iota
	// FormatJSON indicates that user requested JSON report format output
	FormatJSON = iota
	// FormatText indicates that user requested text template report format output
	FormatText
)

func (d formatType) String() string {
	return [...]string{"FC", "JSON", "Text"}[d]
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
	ctx        context.Context
	config     *config
	fp         *os.File
	lastReport *Report
	worker     *worker
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

// FCEntry describes an entry used by FullCalendar report format
type FCEntry struct {
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Title      string    `json:"title"`
	URL        string    `json:"duration"`
	ClassNames []string  `json:"classNames"`
}

// FCReport describes the format of a FullCalendar-compatible report
type FCReport struct {
	Events []FCEntry `json:"events"`
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

type worker struct {
	cmd            string
	leftShiftDown  bool
	rightShiftDown bool
}

// Add appends the current time and task to your timesheet
func (b *Backend) Add(args []string) error {
	task := strings.Join(args, " ")
	return b.addEntry(task)
}

// Close cleans up before exiting
func (b *Backend) Close() error {
	if b.fp != nil {
		b.fp.Close()
	}
	return nil
}

// Edit opens your current timesheet in your default editor or
// in the editor specified by the EDITOR environment variable
func (b *Backend) Edit() error {
	editor := DefaultEditor
	fileLock := flock.New(b.config.omwFile)
	locked, err := fileLock.TryLock()
	defer fileLock.Unlock()
	if err != nil {
		return err
	}
	if !locked {
		return errors.New("Unable to get file lock")
	}
	if preferred := os.Getenv("EDITOR"); preferred != "" {
		editor = preferred
	}
	if term := os.Getenv("OMW_TERM"); runtime.GOOS != "windows" && term != "" {
		editor = fmt.Sprintf("%s -e %s", term, editor)
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
func (b *Backend) Hello() error {
	return b.addEntry("hello")
}

// Report outputs various report formats to one of the following types:
// Text - command-line default
// JSON - web default
// FC   - web fullcalendar JSON feed URL
// Add 24 hours to the parsed end time so that when a user specifies
// --from 2019-01-01 --to 2019-01-02
// that translates to "report on tasks that occurred between 2019-01-01 00:00
// and "2019-01-03 00:00"
func (b *Backend) Report(start, end string, format string) (output string, err error) {
	fcLayout := "2006-01-02T15:04:05-07:00"
	layout := "2006-1-2" // should support optional leading zeros
	layoutEvent := "2006-1-2 15:4"
	report := Report{}
	report.From, err = time.Parse(layout, start)
	if err != nil {
		report.From, err = time.Parse(fcLayout, start)
	}
	if err != nil {
		return "", errors.Wrap(err, "can't parse start time")
	}

	report.To, err = time.Parse(layout, end)
	if err != nil {
		report.To, err = time.Parse(fcLayout, end)
	}
	if err != nil {
		return "", errors.Wrap(err, "can't parse end time")
	}
	report.To = report.To.Add(24 * time.Hour)
	r, err := os.Open(b.config.omwFile)
	defer r.Close()
	if err != nil {
		return "", err
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
			return "", errors.New("Entry has both break and ignore set to true.  Something's wrong.")
		}
		report.Entries = append(report.Entries, *entry)
	}
	f := FormatText
	if format == "json" {
		f = FormatJSON
	}
	if format == "fc" {
		f = FormatFC
	}
	b.lastReport = &report
	output, err = b.formatReport(report, formatType(f))
	if err != nil {
		return "", err
	}
	return output, nil
}

// Stretch append current timestamp to end of timesheet and copy previous task
// fp is opened in append mode, so seek to beginning of file first
func (b *Backend) Stretch() error {
	r, err := os.Open(b.config.omwFile)
	defer r.Close()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	lastEntry := []string{}
	// Brute force retrieval of last line in timesheet
	for scanner.Scan() {
		lastEntry = strings.Fields(scanner.Text())
	}
	if len(lastEntry) <= 2 {
		return errors.New("Missing task description for stretch")
	}
	err = b.addEntry(strings.Join(lastEntry[2:], " "))
	if err != nil {
		return err
	}
	return nil
}

// addEntry seeks to end of file and appends a formatted string
// will create a new empty file if file is missing
func (b *Backend) addEntry(s string) error {
	fp, err := os.OpenFile(b.config.omwFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	defer fp.Close()
	if err != nil {
		errors.Wrapf(err, "Can't open or create %s", b.config.omwFile)
		return err
	}
	ts := time.Now()
	tsFmt := fmt.Sprintf("%d-%02d-%02d %02d:%02d",
		ts.Year(),
		ts.Month(),
		ts.Day(),
		ts.Hour(),
		ts.Minute(),
	)
	entry := fmt.Sprintf("%s  %s\n", tsFmt, s)
	if s == "hello" {
		entry = fmt.Sprintf("\n%s", entry)
	}
	fileLock := flock.New(b.config.omwFile)
	locked, err := fileLock.TryLock()
	defer fileLock.Unlock()
	if err != nil {
		return err
	}
	if !locked {
		return errors.New("Unable to get file lock")
	}
	fp.WriteString(entry)
	return nil
}

func (b *Backend) formatReport(report Report, format formatType) (string, error) {
	if format == FormatJSON {
		output, err := json.Marshal(report)
		return string(output), err
	}

	entries := []FCEntry{}
	if format == FormatFC {
		for _, entry := range report.Entries {
			classes := []string{}
			if entry.Brk {
				classes = append(classes, "breakEntry")
			}
			if entry.Ignore {
				classes = append(classes, "ignoreEntry")
			}

			entries = append(entries, FCEntry{
				Start:      entry.Start,
				End:        entry.Start.Add(entry.Duration),
				Title:      entry.Task,
				URL:        "",
				ClassNames: classes,
			})
		}
		output, err := json.Marshal(FCReport{
			Events: entries,
		})
		return string(output), err
	}

	// fallback to text format
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

func (b *Backend) parseEntry(s string) (*Entry, error) {
	re := regexp.MustCompile(`(?P<task>[a-zA-Z0-9,._+:@%\/-]+[a-zA-Z0-9,._+:@%\/\-\t ]*) ?(?P<mod>\*\*\*?)*`)
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

// Create an instance of the structures that operate on Omw data
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

// runCommand Executes cmd and handles any output
func runCommand(cmd *exec.Cmd) error {
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

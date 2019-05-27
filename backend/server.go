//go:generate go run -tags generate gen.go

package backend

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"time"

	hook "github.com/robotn/gohook"
	"github.com/zserge/lorca"
)

// Backend represents the context and configuration of every instance of the omw command
// Immediate commands (like omw add, omw report), immediately affect the timesheet
// Long-running commands (like omw server), maintain a context
type Backend struct {
	ctx    context.Context
	config *config
	fp     *os.File
	worker *worker
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
func (b *Backend) Add(task string) {
	b.addEntry(task)
	return
}

// Close cleans up before exiting
func (b *Backend) Close() (err error) {
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
		fp: fp,
		worker: nil,
	}
}

// Edit opens your current timesheet in your default editor
func (b *Backend) Edit() error {
	return nil
}

// Hello appends a newline and then another line to end of timesheet with current time
// and the word "Hello".  Meant to be run at the beginning of a new work day
func (b *Backend) Hello() {
	b.fp.WriteString("\n")
	b.addEntry("hello")
	return
}

// Report outputs various report formats to specified location (for now - just the screen)
func (b *Backend) Report(start, end string) {
	return
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
	ui.Bind("runUtt", b.worker.RunUTT)
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
	buf := []byte{}
	_, err := b.fp.ReadAt(buf, 0)
	if err != nil {
		return err
	}
	_, lastLine, err := bufio.ScanLines(buf, true)
	if err != nil {
		return err
	}
	lastEntry := strings.Fields(string(lastLine))[1:]
	b.addEntry(strings.Join(lastEntry, " "))
	return err
}

func (b *Backend) addEntry(s string) {
	ts := time.Now()
	tsFmt := fmt.Sprintf("%d-%d-%d %d:%d",
		ts.Year(),
		ts.Month(),
		ts.Day(),
		ts.Hour(),
		ts.Minute(),
	)
	entry := fmt.Sprintf("%s\t%s\n", tsFmt, s)
	b.fp.WriteString(entry)
}

// RunUTT Executes 'utt' on the command-line and prints the results
func (c *worker) RunUTT(argv []string) {
	c.Lock()
	defer c.Unlock()
	if len(argv) == 1 {
		cmd := exec.Command("utt", argv[0])
		processOutput(cmd)
	} else if len(argv) > 1 {
		args := append([]string{"utt"}, argv...)
		cmd := exec.Command(args[0], args[1:]...)
		processOutput(cmd)
	} else {
		return
	}
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

// processOutput executes cmd and handles any results
func processOutput(cmd *exec.Cmd) {
	out, err := cmd.Output()
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf(string(out))
	}
	return
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

//go:generate go run -tags generate gen.go

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"

	"github.com/go-vgo/robotgo"
	"github.com/zserge/lorca"
)

// Go types that are bound to the UI must be thread-safe, because each binding
// is executed in its own goroutine. In this simple case we may use atomic
// operations, but for more complex cases one should use proper synchronization.
type worker struct {
	sync.Mutex
	cmd    string
	bounds *lorca.Bounds
	ui     lorca.UI
}

// RunUTT Executes 'utt' on the command-line and prints the results
func (c *worker) RunUTT(argv []string) {
	c.Lock()
	defer c.Unlock()
	if len(argv) == 1 {
		cmd := exec.Command(argv[0])
		processOutput(cmd)
	} else if len(argv) > 1 {
		cmd := exec.Command(argv[0], argv[1:]...)
		processOutput(cmd)
	} else {
		return
	}
}

// Minimize Hides the application window
// Saves the current window lorca.Bounds
func (c *worker) Minimize(s string) {
	c.Lock()
	defer c.Unlock()
	fmt.Printf("RECEIVED STRING MINIMIZE = %s\n\n", s)
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

	(*c.bounds).WindowState = lorca.WindowStateNormal
	err := c.ui.SetBounds(*c.bounds)
	if err != nil {
		log.Println("[ERROR] Restore.SetBounds() WindowStateNormal: ", err)
		return
	}
}

// processOutput handles any results received from exec.Command()
func processOutput(_ *exec.Cmd) {
	return
}

func main() {

	args := []string{}
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
	c := &worker{ui: ui, cmd: ""}
	ui.Bind("runUtt", c.RunUTT)
	ui.Bind("minimize", c.Minimize)
	ui.Bind("restore", c.Restore)

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
	select {
	case <-sigc:
	case <-ui.Done():
	}

	// Setup hotkey loop
	for {
		ok := robotgo.AddEvents("ctrl", "shift", "l")
		if ok {
			log.Println("Launching command window...")
		}
		c.Restore()
	}

	log.Println("exiting...")
}

//go:generate go run -tags generate gen.go

package backend

import (
	"context"
	"os"
	"os/exec"
	"reflect"
	"sync"
	"testing"

	hook "github.com/robotn/gohook"
	"github.com/zserge/lorca"
)

func TestBackend_Add(t *testing.T) {
	type fields struct {
		ctx    context.Context
		config *config
		fp     *os.File
		worker *worker
	}
	type args struct {
		args []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				ctx:    tt.fields.ctx,
				config: tt.fields.config,
				fp:     tt.fields.fp,
				worker: tt.fields.worker,
			}
			b.Add(tt.args.args)
		})
	}
}

func TestBackend_Close(t *testing.T) {
	type fields struct {
		ctx    context.Context
		config *config
		fp     *os.File
		worker *worker
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				ctx:    tt.fields.ctx,
				config: tt.fields.config,
				fp:     tt.fields.fp,
				worker: tt.fields.worker,
			}
			if err := b.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Backend.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		fp      *os.File
		omwDir  string
		omwFile string
	}
	tests := []struct {
		name string
		args args
		want *Backend
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Create(tt.args.fp, tt.args.omwDir, tt.args.omwFile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackend_Edit(t *testing.T) {
	type fields struct {
		ctx    context.Context
		config *config
		fp     *os.File
		worker *worker
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				ctx:    tt.fields.ctx,
				config: tt.fields.config,
				fp:     tt.fields.fp,
				worker: tt.fields.worker,
			}
			if err := b.Edit(); (err != nil) != tt.wantErr {
				t.Errorf("Backend.Edit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBackend_Hello(t *testing.T) {
	type fields struct {
		ctx    context.Context
		config *config
		fp     *os.File
		worker *worker
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				ctx:    tt.fields.ctx,
				config: tt.fields.config,
				fp:     tt.fields.fp,
				worker: tt.fields.worker,
			}
			b.Hello()
		})
	}
}

func TestBackend_Report(t *testing.T) {
	type fields struct {
		ctx    context.Context
		config *config
		fp     *os.File
		worker *worker
	}
	type args struct {
		start string
		end   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				ctx:    tt.fields.ctx,
				config: tt.fields.config,
				fp:     tt.fields.fp,
				worker: tt.fields.worker,
			}
			b.Report(tt.args.start, tt.args.end, FormatText)
		})
	}
}

func TestBackend_Run(t *testing.T) {
	type fields struct {
		ctx    context.Context
		config *config
		fp     *os.File
		worker *worker
	}
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				ctx:    tt.fields.ctx,
				config: tt.fields.config,
				fp:     tt.fields.fp,
				worker: tt.fields.worker,
			}
			if err := b.Run(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("Backend.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBackend_Stretch(t *testing.T) {
	type fields struct {
		ctx    context.Context
		config *config
		fp     *os.File
		worker *worker
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				ctx:    tt.fields.ctx,
				config: tt.fields.config,
				fp:     tt.fields.fp,
				worker: tt.fields.worker,
			}
			if err := b.Stretch(); (err != nil) != tt.wantErr {
				t.Errorf("Backend.Stretch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBackend_addEntry(t *testing.T) {
	type fields struct {
		ctx    context.Context
		config *config
		fp     *os.File
		worker *worker
	}
	type args struct {
		s string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Backend{
				ctx:    tt.fields.ctx,
				config: tt.fields.config,
				fp:     tt.fields.fp,
				worker: tt.fields.worker,
			}
			b.addEntry(tt.args.s)
		})
	}
}

func Test_worker_Minimize(t *testing.T) {
	type fields struct {
		Mutex          sync.Mutex
		cmd            string
		bounds         *lorca.Bounds
		ui             lorca.UI
		leftShiftDown  bool
		rightShiftDown bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &worker{
				Mutex:          tt.fields.Mutex,
				cmd:            tt.fields.cmd,
				bounds:         tt.fields.bounds,
				ui:             tt.fields.ui,
				leftShiftDown:  tt.fields.leftShiftDown,
				rightShiftDown: tt.fields.rightShiftDown,
			}
			c.Minimize()
		})
	}
}

func Test_worker_Restore(t *testing.T) {
	type fields struct {
		Mutex          sync.Mutex
		cmd            string
		bounds         *lorca.Bounds
		ui             lorca.UI
		leftShiftDown  bool
		rightShiftDown bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &worker{
				Mutex:          tt.fields.Mutex,
				cmd:            tt.fields.cmd,
				bounds:         tt.fields.bounds,
				ui:             tt.fields.ui,
				leftShiftDown:  tt.fields.leftShiftDown,
				rightShiftDown: tt.fields.rightShiftDown,
			}
			c.Restore()
		})
	}
}

func Test_processOutput(t *testing.T) {
	type args struct {
		cmd *exec.Cmd
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runCommand(tt.args.cmd)
		})
	}
}

func Test_eventLoop(t *testing.T) {
	type args struct {
		c      *worker
		sigc   *chan os.Signal
		ui     lorca.UI
		hotkey *chan hook.Event
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eventLoop(tt.args.c, tt.args.sigc, tt.args.hotkey)
		})
	}
}

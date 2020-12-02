package zsh

import (
	"fmt"
	"os/exec"
	"strings"
)

// Status executes a shell command and returns the exit status.
func Status(command string) error {
	cmd := cmdPrep(command)
	return cmd.Run()
}

// ShWait executes a shell command and waits for the exit status.
// Start and Wait are used to make this call nonblocking with respect to the
// shell environment, but response may be delayed.
func ShWait(command string) error {
	cmd := cmdPrep(command)
	cmd.Run()
	cmd.Start()
	return cmd.Wait()

	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

// Sh executes a shell command line string and returns the result.
// Any error encountered is returned as an ANSI red string.
func Sh(command string) string {
	cmd := cmdPrep(command)
	stdout, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("%Terror: %v", red, err).Error()
	}

	return string(stdout)
}

// cmdPrep prepares a Cmd struct from a command line string.
// Notes: Cmd struct summary:
/*
type Cmd struct {
	Path            string
	Args            []string
	Env             []string
	Dir             string
	Stdin           io.Reader
	Stdout          io.Writer
	Stderr          io.Writer
	ExtraFiles      []*os.File
	SysProcAttr     *syscall.SysProcAttr
	Process         *os.Process
	ProcessState    *os.ProcessState
	ctx             context.Context // nil means none
	lookPathErr     error           // LookPath error, if any.
	finished        bool            // when Wait was called
	childFiles      []*os.File
	closeAfterStart []io.Closer
	closeAfterWait  []io.Closer
	goroutine       []func() error
	errch           chan error // one send per goroutine
	waitDone        chan struct{}
}
*/
func cmdPrep(command string) *exec.Cmd {
	commandSlice := strings.Split(command, " ")
	app := commandSlice[0]
	args := strings.Join(commandSlice[1:], " ")
	return exec.Command(app, args)
}

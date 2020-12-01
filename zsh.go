package zsh

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	redb, _        = hex.DecodeString("1b5b33316d0a") // byte code for ANSI red
	red     string = string(redb)                     // ANSI red
)

// WriteFile creates the file 'fileName' and writes 'data' to it.
// It returns any error encountered. If the file already exists, it
// will be TRUNCATED and OVERWRITTEN.
func WriteFile(fileName string, data string) error {
	dataFile, err := openTrunc(fileName)
	if err != nil {
		log.Println(err)
		return err
	}
	defer dataFile.Close()

	n, err := dataFile.WriteString(data)
	if err != nil {
		log.Println(err)
		return err
	}
	if n != len(data) {
		log.Printf("incorrect string length written (wanted %d): %d\n", len(data), n)
		return fmt.Errorf("incorrect string length written (wanted %d): %d", len(data), n)
	}
	return nil
}

// openTrunc creates and opens the named file for writing. If successful, methods on
// the returned file can be used for writing; the associated file descriptor has mode
//      O_WRONLY|O_CREATE|O_TRUNC
// If the file does not exist, it is created with mode o644;
//
// If the file already exists, it is TRUNCATED and overwritten
//
// If there is an error, it will be of type *PathError.
func openTrunc(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 644)
}

// Sh executes a shell command line string and returns the result.
func Sh(command string) string {
	cmd := cmdPrep(command)
	stdout, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("%Terror: %v", red, err).Error()
	}

	return string(stdout)
}

// FileExists checks if a file exists and is not a directory
func FileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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

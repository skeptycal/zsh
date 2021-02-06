package zsh

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/skeptycal/util/stringutils/ansi"
)

// version represents the SemVer version of this package.
/* The format is:
  vMAJOR[.MINOR[.PATCH[-PRERELEASE][+BUILD]]]
where square brackets indicate optional parts of the syntax;
MAJOR, MINOR, and PATCH are decimal integers without extra
leading zeros; PRERELEASE and BUILD are each a series of
non-empty dot-separated identifiers using only alphanumeric
characters and hyphens; and all-numeric PRERELEASE identifiers
must not have leading zeros.

".-A..Z"

Reference: https://godoc.org/golang.org/x/mod/semver
*/
type version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
	Build      string
}

var (
    DefaultContext context.Context = context.Background()
    errorColor string = "\033[31m"
    red =  ansi.BuildAnsi( ansi.Black, ansi.RedBackground, ansi.Bold)
    errorWriter io.Writer = os.Stdout// ansi.NewAnsiWriter(os.Stdout)
)

const defaultGetStdinArgs = `Example: tr lowercase to UpperCase`

func GetRecentTagHash() (string, error) {

	// git-rev-list - Lists commit objects in reverse chronological order
	// --tags limits the objects returned to git tags
	// --max-count=1 limit the number of commits to one
	// --pretty=oneline
	out, err := Shell("git rev-list --tags --max-count=1")

	if err != nil {
		log.Error(err)
		return "", err
	}
	return out, nil
}

func GitVersionTag() string {
	tag, err := GetRecentTagHash()
	if Err(err) != nil {
		return ""
	}
	out, err := Shell(fmt.Sprintf("git describe --tags %s", tag))
	if Err(err) != nil {
		return ""
	}
	return out
}

// GetStdin sets the stdin pipe for cmd in order of preference.
/*
1. os.Stdin - This will contain piped data allowing UNIX-style
chaining of CLI commands. e.g.
    ls -lah | grep '.go' | myprogram

2. os.Args - This will contain cli arguments.
    myprogram These are the arguments we want.

3. default - A default is available, but optional.

4. none - GetStdin returns an error. Cmd is unchanged.
*/
func GetStdin(cmd *exec.Cmd) error {
	// Prefer piped data
	// TODO - options in os.Args are not supported

	// if Stdin already contains piped data, add it to cmd
	// (This is the reason cmd is passed as a pointer.)
	if os.Stdin != nil {
		cmd.Stdin = os.Stdin
		return nil
	}

	// if there are command line arguments,
	// add them to cmd.stdin
	if len(os.Args) > 1 {
		args := strings.Join(os.Args[1:], " ")
		cmd.Stdin = strings.NewReader(args)
		return nil
	}

	if len(defaultGetStdinArgs) > 0 {
		cmd.Stdin = strings.NewReader(defaultGetStdinArgs)
		return nil
	}
	return fmt.Errorf("no arguments found for stdin")
}

// CombinedOutput executes a shell command line string and returns
// the result. There is no error or statuscode returned.
//
// There is no programatic error information returned at all.
// This has the advantage of returning a single string variable
// that can easily be used as a function argument. e.g.
//
//      fmt.Printf(CombinedOutput("fmtstring 'temp'"),CombinedOutput("statustemp"))
//
// Any error encountered is returned as an ANSI errorColor
// (default bold black on maroon) string. This is most useful
// for REPL style commands where the error can be seen by the
// user but that can fail without effecting progress, such as:
//
//      ls -lah  # it does not change anything important
//
// If you must have a way to test the return status or error
// returned, the string begins with:
//
//      "/x1b["
//
// This works with:
//
//      s := OutErr("fake_program_that_is_not_real")
//      if s[:2] == `/x1b[` {
//          // handle error
//      }
//
// However, it is probably better to use one of the following:
//
//      Shell(command string) (string, error)
//      Status(command string) error
//
func CombinedOutput(command string) string {
	cmd := CmdPrep(command)
	stdout, err := cmd.Output()

	if err != nil {
		return fmt.Sprintf("%verror: %v", errorColor, err)
	}

	return strings.TrimSpace(string(stdout))
}

func Out(command string) string {
	cmd := CmdPrep(command)
	out, err := cmd.Output()

	if err != nil {
		_ = Err(fmt.Errorf("%verror during command %v: %v", errorColor, command, err))
		return ""
	}

	return strings.TrimSpace(string(out))
}



// Repl executes the command and returns the the result.
// Unlike Shell(), the Repl() process has access to the parent's
// stdin, stdout, and stderr streams.
//
// Stdin will be consumed by any process that is able to use it.
//
// StdOut will be sent to os.Stdout
// StdErr will be sent to os.Stderr
func Repl(command string) (string, error) {
	return shell(command, os.Stdin, os.Stdout, os.Stderr)
}

// PipeIn executes the command and returns the the result.
// Unlike Shell(), the PipeIn() process has access to the parent's stdin.
// This can be used to preload the stdin with the string stdInString.
//
// If stdInString == "", it can still be used to query stdin for
// parameters like 'size':
//      ShellIn("stty size")
//      out: "36 118\n"
//      err: <nil>
//
func PipeIn(command string, stdInString string) (string, error) {
	// return shell(command, os.Stdin, nil, nil)
	cmd := CmdPrep(command)
	stdin, err := cmd.StdinPipe()
	if Err(err) != nil {
		return "", err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, stdInString)
	}()

	out, err := cmd.CombinedOutput()
	if Err(err) != nil {
		return string(out), err
	}

	return string(out), nil
}

// PipeOut executes the command and returns the the result.
// Unlike Shell(), the ShellOut() process has access to the parent's
// stdout and stderr streams.
func PipeOut(command string) (string, error) {
	return shell(command, nil, os.Stdout, os.Stderr)
}

// Shell executes a command line string and returns the result.
func Shell(cmd string) (string, error) {
	out, err := shell(cmd, nil, nil, nil)
	return strings.TrimSpace(out), err
}



// Status executes a shell command and returns the exit status as an error.
func Status(command string) error {
	cmd := CmdPrep(command)
	return cmd.Run()
}


// ShWait executes a shell command and waits for the exit status.
// Start and Wait are used to make this call nonblocking with respect to the
// shell environment, but response may be delayed.
func ShWait(command string) error {
	cmd := CmdPrep(command)
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

// Sh executes a shell command line string and returns the result.
// Any error encountered is returned as an ANSI red string.
func Sh(command string) string {
	cmd := CmdPrep(command)
	stdout, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("%Terror: %v", red, err).Error()
	}

	return string(stdout)
}

// CmdPrep returns a Cmd object ready for execution.
//
// The default Background() context is used to kill the process
// (by calling os.Process.Kill) if the context becomes done before
// the command completes on its own.
//
func CmdPrep(command string) *exec.Cmd {
    return cmdPrep(DefaultContext, command)
}


// cmdPrep prepares a Cmd struct from a command line string.
func cmdPrep( ctx context.Context, command string) *exec.Cmd {
	if ctx == nil {
		ctx = DefaultContext
    }
    app, args := AppArgs(command)
    return exec.CommandContext(ctx, app, args...)
}

// shell executes a prepared command structure and returns the result.
// []byte values are converted to string
// sin, sout, and serr are used to redirect input and output of the command.
func shell(command string, sin io.Reader, sout, serr io.Writer) (string, error) {

    cmd := cmdPrep(DefaultContext, command)
	if sin != nil {
		cmd.Stdin = sin
	}
	if sout != nil {
		cmd.Stdout = sout
	}
	if serr != nil {
		cmd.Stderr = serr
	}

	out, err := cmd.Output()

	if err != nil {
		log.Error(err)
		return string(out), err
	}
	return string(out), err
}

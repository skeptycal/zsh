package zsh

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/skeptycal/util/gofile/redlogger"
)

type TrueMap map[string]bool

func (a TrueMap) String() string {
	sb := strings.Builder{}
	for k, v := range a {
		if v {
			sb.WriteString(k)
			sb.WriteString(",")
		}
	}
	return sb.String()
}

func (a TrueMap) List() (s []string) {
	for k, v := range a {
		if v {
			s = append(s, k)
		}
	}
	return
}

type Stringer interface {
	String() string
}

func Tree(pathname string) []string {
	// TODO stuff
	return nil
}

func Dir(pathname string) (files []os.FileInfo, err error) {
	return files, filepath.Walk(pathname,
		func(path string, info os.FileInfo, err error) error {

			files = append(files, info)
			return nil
		})
}

// Err calls error handling and logging routines
func Err(err error) error {
	return redlogger.DoOrDie(err)
}

// GetEnv returns the environment variable specified by 'key'; if the value is empty or not set, the
// default value is returned.
// If the environment variable is not set, a log event is also
// triggered.
func GetEnv(key string, defValue string) string {
	// value, b := os.LookupEnv(key)
	value, b := syscall.Getenv(key)

	if !b {
		log.Infof("environment variable not set: %s", key)
		return defValue
	}

	if strings.TrimSpace(value) == "" {
		return defValue
	}
	return value
}

// AppArgs returns two strings representing the
// (app,args) arguments for os.Command
//
//  `app` is the first word (space delimited) in command
//  `args` is a string containing the remaining words
//
// Leading and trailing spaces are trimmed using
// strings.TrimSpace() so index 0 cannot contain a space.
//
// If command is a single word then, by definition, there
// are no args, but single letter commands are possible,
// e.g. aliases and short os commands like [
//
func AppArgs(command string) (string, []string) {
	command = strings.TrimSpace(command)
    command = strings.ToValidUTF8(command, "")
    command = strings.ReplaceAll(command, "  ", "_")
	list := strings.Split(command, " ")

	switch len(list) {
	case 0:
		return "", nil
	case 1:
		return list[0], nil
	default:
		return list[0], list[1:]
	}
}

func RemoveDuplicates(s string, c byte) string {

	sb := strings.Builder{}
	defer sb.Reset()

    s = strings.TrimSpace(s)

    buf := bytes.NewBuffer([]byte(s))
    last, err := buf.ReadByte()
    if err != nil {
        return ""
    }

    for {
        b, err := buf.ReadByte();
        if err != nil {
            break
        }
        if b == c && b == last {
            continue
        }
        sb.WriteByte(b)

        last = b

    }
    return sb.String()
}

// SliceUniqMap returns unique elements of a slice
// Reference: https://www.reddit.com/r/golang/comments/5ia523/idiomatic_way_to_remove_duplicates_in_a_slice/
func SliceUniqMap(s []int) []int {
	seen := make(map[int]struct{}, len(s))
	j := 0
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		s[j] = v
		j++
	}
	return s[:j]
}

func unique(intSlice []int) []int {
    keys := make(map[int]bool)
    list := []int{}
    for _, entry := range intSlice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}

// Home returns the current user's home directory.
//
// On Unix, including macOS, it returns the $HOME
// environment variable. On Windows, it returns
// %USERPROFILE%. On Plan 9, it returns the $home
// environment variable.
func Home() string {
	s, err := os.UserHomeDir()
	if Err(err) != nil {
		return "~/"
	}
	return s
}

// NameSplit returns separate name and extension of file.
func NameSplit(filename string) (string, string) {
	s := strings.Split(filepath.Base(filename), ".")
	name := s[0]
	ext := ""
	if len(s) > 1 {
		ext = s[len(s)-1]
	}
	return name, ext
}

// Name returns the name of file without path information.
func Name(filename string) string {
	ns, _ := NameSplit(filename)
	return ns
}

// AbsPath returns the absolute path of file.
func AbsPath(file string) string {
	dir, _ := filepath.Abs(file)
	return dir
}

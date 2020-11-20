package process

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	root = "/proc"
)

type Process struct {
	PID     uint64
	CmdLine string
}

func Processes() ([]*Process, error) {
	return findProc(func(process *Process) bool {
		return true
	})
}

func ProcessesByPattern(pattern string) ([]*Process, error) {
	var expression = regexp.MustCompile(pattern)
	return findProc(func(process *Process) bool {
		return expression.MatchString(process.CmdLine)
	})
}

func findProc(filter func(*Process) bool) ([]*Process, error) {
	fh, err := os.Open(root)
	if err != nil {
		return nil, err
	}

	dirNames, err := fh.Readdirnames(-1)
	fh.Close()
	if err != nil {
		return nil, err
	}

	var processes []*Process
	for _, dirName := range dirNames {
		pid, err := strconv.ParseUint(dirName, 10, 0)
		if err != nil {
			// Not a number, so not a PID subdir.
			continue
		}

		fdBase := filepath.Join(root, dirName, "fd")
		dfh, err := os.Open(fdBase)
		if err != nil {
			// Process is be gone by now, or we don't have access.
			continue
		}
		dfh.Close()

		cmd, err := ioutil.ReadFile(filepath.Join(root, dirName, "/cmdline"))
		if err != nil {
			continue
		}
		cmdLine := extractCmdLine(cmd)

		p := &Process{PID: pid, CmdLine: cmdLine}
		if filter(p) {
			processes = append(processes, p)
		}
	}
	return processes, nil
}

func extractCmdLine(content []byte) string {
	c := 0 // cursor of useful bytes

	for i := 0; i < len(content)-1; i++ {
		// Check if next byte is not a '\0' byte.
		if content[i+1] != byte(0) {
			// Offset must match a '\0' byte.
			c = i + 2
			if content[i] == 0 {
				content[i] = byte(0x20)
			}
		}
	}
	return strings.TrimSpace(string(content[0:c]))
}

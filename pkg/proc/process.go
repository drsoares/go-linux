package proc

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type Process struct {
	PID     int
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
		pid, err := strconv.Atoi(dirName)
		if err != nil {
			// Not a number, so not a PID subdir.
			continue
		}

		procDir := filepath.Join(root, dirName)
		err = checkIfProcExists(procDir)
		if err != nil {
			// Process is be gone by now, or we don't have access.
			continue
		}

		cmdLine, err := extractCmdLine(procDir)
		if err != nil {
			continue
		}

		p := &Process{PID: pid, CmdLine: cmdLine}
		if filter(p) {
			processes = append(processes, p)
		}
	}
	return processes, nil
}

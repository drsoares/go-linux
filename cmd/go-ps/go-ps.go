package main

import (
	"flag"
	"fmt"
	"github.com/drsoares/go-linux/pkg/proc"
	"os"
	"runtime"
)

var pattern string

func main() {
	goos := runtime.GOOS
	if goos != "linux" {
		fmt.Println("only linux distros are supported")
		os.Exit(1)
	}

	flag.StringVar(&pattern, "pattern", "", "filter processes by pattern")
	flag.Parse()

	var processes []*proc.Process
	var err error

	if pattern == "" {
		processes, err = proc.Processes()
	} else {
		processes, err = proc.ProcessesByPattern(pattern)
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	for _, proc := range processes {
		fmt.Println(fmt.Sprintf("PID: %d CmdLine: %s", proc.PID, proc.CmdLine))
	}
}

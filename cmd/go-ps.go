package main

import (
	"flag"
	"fmt"
	"github.com/drsoares/go-ps/process"
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

	flag.Parse()
	flag.StringVar(&pattern, "pattern", "", "filter processes by pattern")

	var processes []*process.Process
	var err error
	if pattern != "" {
		processes, err = process.Processes()
	} else {
		processes, err = process.ProcessesByPattern(pattern)
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	for _, proc := range processes {
		fmt.Println(fmt.Sprintf("PID: %s CmdLine: %s", proc.PID, proc.CmdLine))
	}
}

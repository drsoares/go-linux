package proc

import (
	"os"
	"path/filepath"
)

const (
	root    = "/proc"
	cmdLine = "/cmdline"
	fd      = "fd"
	netTcp  = "/net/tcp"
)

func checkIfProcExists(procDir string) error {
	fdBase := filepath.Join(procDir, fd)
	dfh, err := os.Open(fdBase)
	if err != nil {
		return err
	}
	dfh.Close()
	return nil
}

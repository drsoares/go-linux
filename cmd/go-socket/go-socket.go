package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/drsoares/go-linux/pkg/proc"
	"os"
	"runtime"
)

var states = map[string]proc.SocketState{
	"established": proc.Established,
	"sync_sent":   proc.SynSent,
	"sync_recv":   proc.SynRecv,
	"fin_wait1":   proc.FinWait1,
	"fin_wait2":   proc.FinWait2,
	"time_wait":   proc.TimeWait,
	"close":       proc.Close,
	"close_wait":  proc.CloseWait,
	"last_ack":    proc.LastAck,
	"listen":      proc.Listen,
	"closing":     proc.Closing,
}

func main() {
	goos := runtime.GOOS
	if goos != "linux" {
		fmt.Println("only linux distros are supported")
		os.Exit(1)
	}

	var pid string
	//var state proc.SocketState
	flag.StringVar(&pid, "pid", "", "sockets by pid")
	//flag.UintVar(&state, "state", -1)
	flag.Parse()

	var sockets []*proc.TcpSocket
	var err error

	if pid != "" {
		sockets, err = proc.SocketsByPID(pid)
	} else {
		err = errors.New("not implemented yet")
	}
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	for _, socket := range sockets {
		fmt.Println(fmt.Sprintf("%s\t%d\t%s\t%d\t%d", socket.LocalAddress.IP.String(), socket.LocalAddress.Port,
			socket.RemoteAddress.IP.String(), socket.RemoteAddress.Port, socket.State))
	}
}

package port

import (
	"net"

	"github.com/taluos/Malt/pkg/log"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer func() {
		err := l.Close()
		if err != nil {
			log.Errorf("Cannot close tcp listener: ", err)
			return
		}
	}()
	return l.Addr().(*net.TCPAddr).Port, nil
}

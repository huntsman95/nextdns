package ctl

import (
	"net"

	"github.com/nextdns/nextdns/ctl/internal/winio"
)

func listen(addr string) (net.Listener, error) {
	return winio.ListenPipe(`\\.\pipe\`+addr, &winio.PipeConfig{
		SecurityDescriptor: "D:P(A;;GA;;;WD)",
	})
}

func dial(addr string) (net.Conn, error) {
	return winio.DialPipe(`\\.\pipe\`+addr, nil)
}

//go:build windows
// +build windows

package proxy

import (
	"net"
)

// raises not implemented errors from `golang.org/x/net/ipv4{4,6}` lib
// @TODO: check for implementation at a later time
func setUDPDstOptions(c *net.UDPConn) error { return nil }

func parseDstFromOOB([]byte) net.IP { return nil }

func readUDP(c *net.UDPConn, buf []byte) (n int, lip net.IP, raddr *net.UDPAddr, err error) {
	n, raddr, err = c.ReadFromUDP(buf)
	if err != nil {
		return -1, nil, nil, err
	}
	return n, nil, raddr, nil
}

func writeUDP(c *net.UDPConn, buf []byte, _ net.IP, raddr *net.UDPAddr) error {
	_, err := c.WriteToUDP(buf, raddr)
	return err
}

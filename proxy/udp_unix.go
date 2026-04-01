//go:build aix || darwin || dragonfly || linux || netbsd || openbsd || solaris || freebsd
// +build aix darwin dragonfly linux netbsd openbsd solaris freebsd

package proxy

import (
	"net"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// setUDPDstOptions sets the FlagDst on c to request the destination address as
// part of the oob data.
func setUDPDstOptions(c *net.UDPConn) error {
	// Try setting the flags for both families and ignore the errors unless they
	// both error.
	err6 := ipv6.NewPacketConn(c).SetControlMessage(ipv6.FlagDst|ipv6.FlagInterface, true)
	err4 := ipv4.NewPacketConn(c).SetControlMessage(ipv4.FlagDst|ipv4.FlagInterface, true)
	if err6 != nil && err4 != nil {
		return err4
	}
	return nil
}

// parseDstFromOOB takes oob data and returns the destination IP.
func parseDstFromOOB(oob []byte) net.IP {
	cm6 := &ipv6.ControlMessage{}
	if cm6.Parse(oob) == nil && cm6.Dst != nil {
		return cm6.Dst
	}
	cm4 := &ipv4.ControlMessage{}
	if cm4.Parse(oob) == nil && cm4.Dst != nil {
		return cm4.Dst
	}
	return nil
}

func readUDP(c *net.UDPConn, buf []byte) (n int, lip net.IP, raddr *net.UDPAddr, err error) {
	var oobn int
	oob := make([]byte, udpOOBSize)
	n, oobn, _, raddr, err = c.ReadMsgUDP(buf, oob)
	if err != nil {
		return -1, nil, nil, err
	}
	lip = parseDstFromOOB(oob[:oobn])
	return n, lip, raddr, nil
}

func writeUDP(c *net.UDPConn, buf []byte, lip net.IP, raddr *net.UDPAddr) error {
	_, _, err := c.WriteMsgUDP(buf, oobWithSrc(lip), raddr)
	return err
}

package socks4server

import (
	"io"
	"net"
)

// SOCKS4 Connection
type SOCKS4Conn struct {
	vn      byte   // 1
	cd      byte   // 1
	dstport int16  // 2byte
	dstip   net.IP // 4byte
	userid  string // null terminated string
	dstconn *net.TCPConn
	srcconn *net.TCPConn
}

func (c *SOCKS4Conn) DoRelay() {
	go io.Copy(c.SrcConn(), c.DstConn())
	io.Copy(c.DstConn(), c.SrcConn())
	c.Close()
}

func (c *SOCKS4Conn) DstConn() *net.TCPConn {
	return c.dstconn
}

func (c *SOCKS4Conn) SrcConn() *net.TCPConn {
	return c.srcconn
}

func (c *SOCKS4Conn) Close() {
	c.srcconn.Close()
	c.dstconn.Close()
}

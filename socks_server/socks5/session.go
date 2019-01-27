package socks5

import (
	"fmt"
	"net"
	"study-golang1/socks_server/relay"
)

const (
	VerRequestSocks5 byte = 5

	MethodNoAuthenticationRequired = 0
	MethodGSSAPI                   = 1
	MethodUsernamePassword         = 2
	MethodNoAcceptable             = byte(255)

	CmdConnect      = 1
	CmdBind         = 2
	CmdUdpAssociate = 3

	RsvReserved = 0

	AtypIPV4Address = 1
	AtypDomainName  = 3
	AtypIPV6Address = 4
)

type Session struct {
	ver    byte // 1
	method byte // 1
	relay  *relay.Relay
}

func Negotiate(vn byte, conn *net.TCPConn) (*Session, error) {
	return nil, fmt.Errorf("SOCKS5 not implemented")
}

func (s *Session) Version() string {
	return "5"
}

func (s *Session) Relay() *relay.Relay {
	return s.relay
}

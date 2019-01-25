package socks4server

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
)

// SOCKS4 Listner server
type SOCKS4Server struct {
	listener *net.TCPListener
}

func NewSOCKS4Server(proto string, addr string) (s4srv *SOCKS4Server, err error) {
	ln, err := net.Listen(proto, addr)
	if err != nil {
		return nil, err
	}

	s4srv = &SOCKS4Server{ln.(*net.TCPListener)}
	return s4srv, nil
}

func (s *SOCKS4Server) Listener() (listener *net.TCPListener) {
	return s.listener
}

func (s *SOCKS4Server) Accept() (s4conn *SOCKS4Conn, err error) {

	srcconn, err := s.listener.Accept()
	if err != nil {
		return nil, err
	}

	srcreader := bufio.NewReader(srcconn)

	vn, err := srcreader.ReadByte()
	if err != nil {
		__SendResponse(srcconn, byte(0), byte(91), int16(0), net.IP{})
		return nil, fmt.Errorf("%s < %s", err, "failed to recv vn")
	}

	cd, err := srcreader.ReadByte()
	if err != nil {
		__SendResponse(srcconn, byte(0), byte(91), int16(0), net.IP{})
		return nil, fmt.Errorf("%s < %s", err, "failed to recv cd")
	}

	dstport := int16(0)
	err = binary.Read(srcreader, binary.BigEndian, &dstport)
	if err != nil {
		__SendResponse(srcconn, byte(0), byte(91), int16(0), net.IP{})
		return nil, fmt.Errorf("%s < %s", err, "failed to recv dstport")
	}

	dstip := net.IP{0, 0, 0, 0}
	_, err = srcreader.Read(dstip)
	if err != nil {
		__SendResponse(srcconn, byte(0), byte(91), int16(0), net.IP{})
		return nil, fmt.Errorf("%s < %s", err, "failed to recv dstip")
	}

	userid, err := srcreader.ReadString(0) // NULL TERMINATED
	if err != nil {
		__SendResponse(srcconn, byte(0), byte(91), int16(0), net.IP{})
		return nil, fmt.Errorf("%s < %s", err, "failed to recv userid")
	}

	dstconn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dstip, dstport))
	if err != nil {
		__SendResponse(srcconn, byte(0), byte(91), int16(0), net.IP{})
		return nil, fmt.Errorf("%s < %s", err, "failed to dial")
	}

	__SendResponse(srcconn, byte(0), byte(90), int16(0), net.IP{})
	s4conn = &SOCKS4Conn{vn, cd, dstport, dstip, userid, dstconn.(*net.TCPConn), srcconn.(*net.TCPConn)}
	return s4conn, nil
}

func __SendResponse(conn net.Conn, vn byte, cd byte, dstport int16, dstip net.IP) {
	conn.Write([]byte{vn})
	conn.Write([]byte{cd})
	conn.Write([]byte{0, 0, 0, 0, 0, 0})
}

package socks4

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"study-golang1/socks_server/relay"
)

const (
	VnRequestSocks4  byte = 4
	VnResponseSocks4 byte = 0
	CdConnect        byte = 1
	CdBind           byte = 2
	CdReject         byte = 91
	CdGrant          byte = 90
)

type Session struct {
	vn        byte   // 1
	cd        byte   // 1
	dstport   uint16 // 2byte
	dstip     net.IP // 4byte
	userid    string // null terminated string
	relayConn *relay.Relay
}

func Negotiate(vn byte, conn *net.TCPConn) (*Session, error) {
	_, err := conn.Write([]byte{VnResponseSocks4})
	if nil != err {
		return nil, err
	}

	errorResult := func(err error) (*Session, error) {
		_ = sendResponse(conn, CdReject, uint16(0), net.IP{})
		return nil, fmt.Errorf("%s", err)
	}

	reader := bufio.NewReader(conn)

	cd, err := reader.ReadByte()
	if err != nil {
		return errorResult(fmt.Errorf("%s < %s", err, "failed to recv cd"))
	}

	if cd == CdBind {
		return errorResult(fmt.Errorf("bind is not supported"))
	}

	if cd != CdConnect {
		return errorResult(fmt.Errorf("unknown code (%d)", cd))
	}

	dstport := uint16(0)
	err = binary.Read(reader, binary.BigEndian, &dstport)
	if err != nil {
		return errorResult(fmt.Errorf("%s < %s", err, "failed to recv dstport"))
	}

	if dstport == 0 {
		return errorResult(fmt.Errorf("dstport is zero"))
	}

	dstip := net.IP{0, 0, 0, 0}
	_, err = reader.Read(dstip)
	if err != nil {
		return errorResult(fmt.Errorf("%s < %s", err, "failed to recv dstip"))
	}

	userid, err := reader.ReadString(0) // NULL TERMINATED
	if err != nil {
		return errorResult(fmt.Errorf("%s < %s", err, "failed to recv userid"))
	}

	dstConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", dstip, dstport))
	if err != nil {
		return errorResult(fmt.Errorf("%s < %s", err, "failed to connect"))
	}

	err = sendResponse(conn, CdGrant, uint16(0), []byte{0, 0, 0, 0})
	if err != nil {
		return errorResult(fmt.Errorf("%s < %s", err, "failed to send response"))
	}

	relayConn := relay.NewRelay(conn, dstConn.(*net.TCPConn))
	sess := &Session{vn, cd, dstport, dstip, userid, relayConn}

	return sess, nil
}

func (s *Session) Version() int {
	return int(s.vn)
}

func (s *Session) RelayConn() *relay.Relay {
	return s.relayConn
}

func sendResponse(conn net.Conn, cd byte, dstport uint16, dstip net.IP) error {
	n, err := conn.Write([]byte{cd})
	if nil != err {
		return err
	}
	if n != 1 {
		return fmt.Errorf("unexpected close (%d)", n)
	}

	b := []byte{0, 0}
	binary.BigEndian.PutUint16(b, dstport)
	n, err = conn.Write(b)
	if nil != err {
		return err
	}
	if n != 2 {
		return fmt.Errorf("unexpected close (%d)", n)
	}

	n, err = conn.Write(dstip)
	if nil != err {
		return err
	}
	if n != 4 {
		return fmt.Errorf("unexpected close (%d)", n)
	}

	return nil
}

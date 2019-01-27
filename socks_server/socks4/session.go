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
	vn      byte   // 1
	cd      byte   // 1
	dstport uint16 // 2byte
	dstip   net.IP // 4byte
	userid  string // null terminated string
	domain  string // null terminated string
	relay   *relay.Relay
}

func Negotiate(vn byte, conn *net.TCPConn) (*Session, error) {
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

	// dstip==0.0.0.x : SOCKS4A : 名前解決をsocks serverで行うモードです。
	domain := ""
	if dstip[0] == 0 && dstip[1] == 0 && dstip[2] == 0 && dstip[3] != 0 {
		var err4a error
		domain, err4a = reader.ReadString(0) // NULL TERMINATED
		if err4a != nil {
			return errorResult(fmt.Errorf("%s < %s", err4a, "failed to recv domain"))
		}

		addr, err4a := net.ResolveIPAddr("ip4", domain)
		if err4a != nil {
			return errorResult(fmt.Errorf("%s < %s (%s)", err4a, "failed to resolve domain", domain))
		}
		dstip = addr.IP
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
	sess := &Session{vn, cd, dstport, dstip, userid, domain, relayConn}

	return sess, nil
}

func (s *Session) Version() string {
	if s.domain == "" {
		return "4A"
	} else {
		return "4"
	}
}

func (s *Session) Relay() *relay.Relay {
	return s.relay
}

func sendResponse(conn net.Conn, cd byte, dstport uint16, dstip net.IP) error {
	n, err := conn.Write([]byte{VnResponseSocks4})
	if nil != err {
		return err
	}
	if n != 1 {
		return fmt.Errorf("unexpected close (vn:%d)", n)
	}

	n, err = conn.Write([]byte{cd})
	if nil != err {
		return err
	}
	if n != 1 {
		return fmt.Errorf("unexpected close (cd:%d)", n)
	}

	b := []byte{0, 0}
	binary.BigEndian.PutUint16(b, dstport)
	n, err = conn.Write(b)
	if nil != err {
		return err
	}
	if n != 2 {
		return fmt.Errorf("unexpected close (dstport:%d)", n)
	}

	n, err = conn.Write(dstip)
	if nil != err {
		return err
	}
	if n != 4 {
		return fmt.Errorf("unexpected close (dstip:%d)", n)
	}

	return nil
}

package socks_server

import (
	"fmt"
	"log"
	"net"
	"study-golang1/socks_server/relay"
	"study-golang1/socks_server/socks4"
)

type SocksServer struct {
	proto    string
	addr     string
	listener *net.TCPListener
}

type SocksSession interface {
	Version() int
	RelayConn() *relay.Relay
}

func NewSocksServer(proto string, addr string) *SocksServer {
	return &SocksServer{proto, addr, nil}
}

func (s *SocksServer) Start() error {
	ln, err := net.Listen(s.proto, s.addr)
	if err != nil {
		return err
	}

	s.listener = ln.(*net.TCPListener)

	for {
		sess, err := s.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		log.Printf("client accepted: %s -> %s", sess.RelayConn().SrcConn().RemoteAddr(), sess.RelayConn().DestConn().RemoteAddr())
		go sess.RelayConn().Start()
	}
}

func (s *SocksServer) Accept() (SocksSession, error) {
	srcConn, err := s.listener.Accept()
	if err != nil {
		return nil, err
	}

	socksConn, err := s.negotiate(srcConn.(*net.TCPConn))
	if err != nil {
		return nil, err
	}

	return socksConn, nil
}

func (s *SocksServer) negotiate(conn *net.TCPConn) (SocksSession, error) {
	temp := []byte{0}
	sz, err := conn.Read(temp)
	if err != nil {
		return nil, err
	}

	if sz == 0 {
		return nil, fmt.Errorf("unexpected close")
	}
	vn := temp[0]

	var sess SocksSession
	if vn == socks4.VnRequestSocks4 {
		sess, err = socks4.Negotiate(vn, conn)
		if err != nil {
			return nil, err
		}
	} else if vn == 5 {

	} else {
		return nil, fmt.Errorf("unexpected version (%d)", vn)
	}

	return sess, nil
}

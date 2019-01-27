package socks_server

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"study-golang1/socks_server/relay"
	"study-golang1/socks_server/socks4"
	"study-golang1/socks_server/socks5"
	"sync"
)

type SocksServer struct {
	proto    string
	addr     string
	listener *net.TCPListener
	clientWg sync.WaitGroup
	srcConns chan *net.TCPConn
}

var (
	maxClient = runtime.NumCPU()
)

type ISocksSession interface {
	Version() int
	Relay() *relay.Relay
}

func NewSocksServer(proto string, addr string) *SocksServer {
	return &SocksServer{
		proto:    proto,
		addr:     addr,
		clientWg: sync.WaitGroup{}, // 現実装だとワーカーの出力を集約する必要がないのでこれは不要っぽい
		srcConns: make(chan *net.TCPConn, maxClient),
	}
}

func (s *SocksServer) Start() error {
	ln, err := net.Listen(s.proto, s.addr)
	if err != nil {
		return err
	}
	s.listener = ln.(*net.TCPListener)

	s.startClientHandlers()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			err := s.listener.Close()
			if err != nil {
				log.Println(err)
			}
			close(s.srcConns)
			return err
		}

		s.srcConns <- conn.(*net.TCPConn)
	}
}

func (s *SocksServer) startClientHandlers() {
	go func() {
		for i := 0; i < maxClient; i++ {
			s.clientWg.Add(1)
			go func() {
				for it := range s.srcConns {
					s.handleClient(it)
				}
				s.clientWg.Done()
			}()
		}
		s.clientWg.Wait()
	}()
}

func (s *SocksServer) handleClient(conn *net.TCPConn) {
	sess, err := s.negotiate(conn)
	if err != nil {
		log.Println(err)
		err := conn.Close()
		if err != nil {
			log.Println(err)
		}
		return
	}

	debugPrint()

	log.Printf(
		"socks%d relay: %s -> %s",
		sess.Version(), sess.Relay().Src().RemoteAddr(), sess.Relay().Dest().RemoteAddr(),
	)
	sess.Relay().Start()
}

func (s *SocksServer) negotiate(conn *net.TCPConn) (ISocksSession, error) {
	temp := []byte{0}
	sz, err := conn.Read(temp)
	if err != nil {
		return nil, err
	}

	if sz == 0 {
		return nil, fmt.Errorf("unexpected close")
	}
	ver := temp[0]

	var sess ISocksSession
	if ver == socks4.VnRequestSocks4 {
		sess, err = socks4.Negotiate(ver, conn)
		if err != nil {
			return nil, err
		}
	} else if ver == socks5.VerRequestSocks5 {
		sess, err = socks5.Negotiate(ver, conn)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unexpected version (%d)", ver)
	}

	return sess, nil
}

func debugPrint() {
	log.Printf("NumGoroutine=%d", runtime.NumGoroutine())
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	log.Printf("MemStats.Alloc=%d", stats.Alloc)
}

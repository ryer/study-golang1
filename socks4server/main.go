package socks4server

import (
	"fmt"
	"log"
)

func Main(port int) {
	log.SetFlags(log.Lshortfile)

	addr := fmt.Sprintf(":%d", port)
	s4srv, err := NewSOCKS4Server("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("server started: %s", addr)

	for {
		s4conn, err := s4srv.Accept()
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("client accepted: %s -> %s", s4conn.SrcConn().RemoteAddr(), s4conn.DstConn().RemoteAddr())
		go s4conn.DoRelay()
	}
}

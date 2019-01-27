package socks_server

import (
	"fmt"
	"log"
)

func Main(port int) {
	log.SetFlags(log.Lshortfile)

	addr := fmt.Sprintf(":%d", port)
	server := NewSocksServer("tcp", addr)

	log.Printf("start server: %s", addr)

	err := server.Start()
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("server started: %s", addr)
}

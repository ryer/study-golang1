package echosrv

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func Main(port int) {
	log.SetFlags(log.Lshortfile)

	addr := fmt.Sprintf(":%d", port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	for {
		client, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go HandleClient(client)
	}
}

func HandleClient(client net.Conn) {
	clientIo := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))

	for {
		line, isPrefix, err := clientIo.ReadLine()
		if isPrefix {
			// unreached end of line
			log.Printf("<<%s>>", err)
		}

		if err != nil {
			// handle error
			log.Printf("<<%s>>", err)
			break
		}

		sz, err := clientIo.WriteString(fmt.Sprintf("%s\n", line))
		if err != nil || sz == 0 {
			// handle error
			log.Printf("<<%s>>", err)
			break
		}

		err = clientIo.Flush()
		if err != nil {
			log.Printf("<<%s>>", err)
			break
		}
	}

	err := client.Close()
	if err != nil {
		log.Printf("<<%s>>", err)
		return
	}

	log.Println("client closed.")
}

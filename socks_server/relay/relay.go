package relay

import (
	"io"
	"log"
	"net"
)

type Relay struct {
	src  *net.TCPConn
	dest *net.TCPConn
}

func NewRelay(src *net.TCPConn, dest *net.TCPConn) *Relay {
	return &Relay{src, dest}
}

func (r *Relay) Src() *net.TCPConn {
	return r.src
}

func (r *Relay) Dest() *net.TCPConn {
	return r.dest
}

func (r *Relay) Start() {
	go r.doRelay(r.Src(), r.Dest())
	r.doRelay(r.Dest(), r.Src())
	r.close()
}

func (r *Relay) doRelay(conn1 *net.TCPConn, conn2 *net.TCPConn) {
	_, err := io.Copy(conn1, conn2)
	if nil != err {
		log.Print(err)
	}
}

func (r *Relay) close() {
	_ = r.src.Close()
	_ = r.dest.Close()
}

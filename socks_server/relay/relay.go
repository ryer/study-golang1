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

func (r *Relay) Start() {
	go r.relay(r.SrcConn(), r.DestConn())
	r.relay(r.DestConn(), r.SrcConn())
	r.Close()
}

func (r *Relay) relay(conn1 *net.TCPConn, conn2 *net.TCPConn) {
	_, err := io.Copy(conn1, conn2)
	if nil != err {
		log.Print(err)
	}
}

func (r *Relay) SrcConn() *net.TCPConn {
	return r.src
}

func (r *Relay) DestConn() *net.TCPConn {
	return r.dest
}

func (r *Relay) Close() {
	_ = r.src.Close()
	_ = r.dest.Close()
}

package broker

import (
	"strings"
	"study-golang1/data"
	"sync"
)

type Broker struct {
	wg     sync.WaitGroup
	input  *data.Data
	mtx    sync.Mutex
	cur    int
	Output chan data.Item
	Exit   chan struct{}
}

func NewBroker(input *data.Data) *Broker {
	b := Broker{}
	b.input = input
	b.Output = make(chan data.Item)
	b.Exit = make(chan struct{})

	return &b
}

func (b *Broker) Run() {
	go func() {
		for i := 0; i < 5; i++ {
			b.startWorker()
		}
		b.wg.Wait()
		close(b.Exit)
	}()
}

func (b *Broker) startWorker() {
	b.wg.Add(1)
	go func() {
		for {
			b.mtx.Lock()
			if len(b.input.Items) == b.cur {
				b.mtx.Unlock()
				b.wg.Done()
				return
			}
			it := b.input.Items[b.cur]
			b.cur++
			b.mtx.Unlock()

			it.Url = strings.Replace(it.Url, "https://", "/", -1)

			b.Output <- it
		}
	}()
}

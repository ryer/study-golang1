package broker

import (
	"study-golang1/data"
	"sync"
)

type Broker struct {
	wg     sync.WaitGroup
	target *data.Data
	mtx    sync.Mutex
	cur    int
	Result chan data.Item
	Wait   chan struct{}
}

func NewBroker(target *data.Data) *Broker {
	b := Broker{}
	b.target = target
	b.Result = make(chan data.Item, 5)
	b.Wait = make(chan struct{})

	return &b
}

func (b *Broker) Run() {
	go func() {
		for i := 0; i < 5; i++ {
			b.startWorker()
		}
		b.wg.Wait()
		close(b.Wait)
	}()
}

func (b *Broker) startWorker() {
	b.wg.Add(1)
	go func() {
		for {
			b.mtx.Lock()
			if len(b.target.Items) == b.cur {
				b.mtx.Unlock()
				b.wg.Done()
				return
			}
			it := b.target.Items[b.cur]
			b.cur++
			b.mtx.Unlock()

			b.Result <- it
		}
	}()
}

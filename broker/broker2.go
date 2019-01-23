package broker

import (
	"strings"
	"study-golang1/data"
	"sync"
)

type Broker2 struct {
	wg     sync.WaitGroup
	mtx    sync.Mutex
	cur    int
	Input  chan data.Item
	Output chan data.Item
	Exit   chan struct{}
}

func NewBroker2() *Broker2 {
	b := Broker2{}
	b.Input = make(chan data.Item)
	b.Output = make(chan data.Item)
	b.Exit = make(chan struct{})

	return &b
}

func (b *Broker2) Run() {
	go func() {
		for i := 0; i < 5; i++ {
			b.startWorker()
		}
		b.wg.Wait()
		close(b.Exit)
	}()
}

func (b *Broker2) startWorker() {
	b.wg.Add(1)
	go func() {
		for {
			it, ok := <-b.Input
			if !ok {
				b.wg.Done()
				return
			}

			it.Url = strings.Replace(it.Url, "https://", "/", -1)

			b.Output <- it
		}
	}()
}

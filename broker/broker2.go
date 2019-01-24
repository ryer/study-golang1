package broker

// channelから読むタイプ

import (
	"runtime"
	"study-golang1/data"
	"sync"
)

type Broker2 struct {
	wg     sync.WaitGroup
	work   ItemWork
	input  chan data.Item
	output chan data.Item
}

func NewBroker2(input chan data.Item) *Broker2 {
	b := Broker2{}
	b.input = input
	b.output = make(chan data.Item)
	return &b
}

func (b *Broker2) Invoke(work ItemWork) {
	b.work = work
	b.run()
}

func (b *Broker2) Output() chan data.Item {
	return b.output
}

func (b *Broker2) run() {
	go func() {
		for i := 0; i < runtime.NumCPU(); i++ {
			b.startWorker()
		}
		b.wg.Wait()
		close(b.output)
	}()
}

func (b *Broker2) startWorker() {
	b.wg.Add(1)
	go func() {
		for it := range b.input {
			result, err := b.work(it)
			if nil != err {
				break
			}
			b.output <- result
		}
		b.wg.Done()
	}()
}

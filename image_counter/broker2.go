package image_counter

// channelから読むタイプ

import (
	"runtime"
	"sync"
)

type Broker2 struct {
	wg     sync.WaitGroup
	work   ItemWork
	input  chan Item
	output chan Item
}

func NewBroker2(input chan Item) *Broker2 {
	b := Broker2{}
	b.input = input
	b.output = make(chan Item)
	return &b
}

func (b *Broker2) Invoke(work ItemWork) {
	b.work = work
	b.run()
}

func (b *Broker2) Output() chan Item {
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
				panic(err)
			}
			b.output <- result
		}
		b.wg.Done()
	}()
}

package image_counter

// Mutexで守りつつ1つずつ消費していくタイプ

import (
	"runtime"
	"sync"
)

type Broker1 struct {
	wg     sync.WaitGroup
	work   ItemWork
	mtx    sync.Mutex
	cur    int
	input  Data
	output chan Item
}

func NewBroker1(input Data) *Broker1 {
	b := Broker1{}
	b.input = input
	b.output = make(chan Item)
	return &b
}

func (b *Broker1) Invoke(work ItemWork) {
	b.work = work
	b.run()
}

func (b *Broker1) Output() chan Item {
	return b.output
}

func (b *Broker1) run() {
	go func() {
		for i := 0; i < runtime.NumCPU(); i++ {
			b.startWorker()
		}
		b.wg.Wait()
		close(b.output)
	}()
}

func (b *Broker1) startWorker() {
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

			result, err := b.work(it)
			if nil != err {
				b.wg.Done()
				return
			}

			b.output <- result
		}
	}()
}

package wgchan

import "sync"

type WgChan struct {
	wg *sync.WaitGroup
}

func NewWgChan() WgChan {
	return WgChan{
		wg: &sync.WaitGroup{},
	}
}

func (wgc WgChan) Add(delta int) {
	wgc.wg.Add(delta)
}

func (wgc WgChan) AddChan() chan<- struct{} {
	wgc.wg.Add(1)
	c := make(chan struct{})
	go func() {
		<-c
		wgc.wg.Done()
	}()
	return c
}

func (wgc WgChan) Done() {
	wgc.wg.Done()
}

func (wgc *WgChan) Wait() {
	wgc.wg.Wait()
}

func (wgc *WgChan) WaitChan() <-chan struct{} {
	c := make(chan struct{})
	go func() {
		wgc.wg.Wait()
		close(c)
	}()
	return c
}

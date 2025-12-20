package utils

import (
	"sync"
)

type IRunner interface {
	Close() error
	ExecLocked(fct func())
	Start()
	Wait()
}

type runner struct {
	mx             sync.Mutex
	chWaitChannels []chan bool
	isClosed       bool
	isStarted      bool

	onNext  func() (canContinue bool)
	onClose func() error
}

func NewRunner(
	onNext func() (canContinue bool),
	onClose func() error,
) IRunner {
	return &runner{
		onNext:  onNext,
		onClose: onClose,
	}
}

func (p *runner) ExecLocked(fct func()) {
	p.mx.Lock()
	defer p.mx.Unlock()
	if p.isClosed {
		return
	}
	fct()
}

func (p *runner) Close() error {
	isClosed := p.closeLocked()
	if !isClosed && p.onClose != nil {
		return p.onClose()
	}
	return nil
}

func (p *runner) Start() {
	ExecLocked(&p.mx, func() {
		p.isStarted = true
	})
	if p.onNext != nil {
		go p.loop()
	}
}

func (p *runner) Wait() {
	ExecLocked(&p.mx, func() {
		if !p.isStarted {
			panic("runner is not started")
		}
	})
	ch := p.addWaitChannel()
	_ = <-ch
}

func (p *runner) addWaitChannel() chan bool {
	p.mx.Lock()
	defer p.mx.Unlock()
	ch := make(chan bool, 1)
	p.chWaitChannels = append(p.chWaitChannels, ch)
	return ch
}

func (p *runner) checkClosed() bool {
	p.mx.Lock()
	defer p.mx.Unlock()
	return p.isClosed
}

func (p *runner) closeLocked() bool {
	p.mx.Lock()
	defer p.mx.Unlock()
	isClosed := p.isClosed
	for _, ch := range p.chWaitChannels {
		ch <- true
		close(ch)
	}
	p.chWaitChannels = nil
	p.isClosed = true
	return isClosed
}

func (p *runner) loop() {
	for {
		if p.checkClosed() {
			return
		}
		if canContinue := p.onNext(); !canContinue {
			p.Close()
			return
		}
	}
}

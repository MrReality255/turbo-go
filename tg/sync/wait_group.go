package sync

import (
	"sync"

	"github.com/MrReality255/turbo-go/tg/utils"
)

type IWaitGroup interface {
	Go(fct any)
	Wait() error
}

type waitGroup struct {
	cl      IConcurrentLauncher
	counter int
	mx      sync.Mutex
	err     error
	chDone  chan error
}

func NewWaitGroup(maxCount int, threadCount int) IWaitGroup {
	threadCount = utils.Coalesce(threadCount, maxCount)
	chDone := make(chan error, 1)

	return &waitGroup{
		cl:     NewConcurrentLauncher(maxCount, threadCount),
		chDone: chDone,
	}
}

func (wg *waitGroup) Go(fct any) {
	wg.mx.Lock()
	defer wg.mx.Unlock()
	wg.counter++
	go wg.cl.Go(func() error {
		err := execFct(fct)
		wg.mx.Lock()
		defer wg.mx.Unlock()
		if wg.err == nil {
			wg.err = err
		}
		wg.counter--
		if wg.counter == 0 {
			wg.chDone <- wg.err
		}
		return nil

	})
}

func (wg *waitGroup) Wait() error {
	return <-wg.chDone
}

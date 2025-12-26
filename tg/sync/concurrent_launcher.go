package sync

import (
	"sync"

	"github.com/MrReality255/turbo-go/tg/utils"
)

type IConcurrentLauncher interface {
	Go(fct interface{})
	Locked(fct interface{}) error
	Wait() error
}

type ConcurrentLauncher struct {
	wg      sync.WaitGroup
	errList *utils.ErrorList
	ch      chan int
	mx      sync.Mutex
}

func NewConcurrentLauncher(maxCount int, limit int) IConcurrentLauncher {
	cl := &ConcurrentLauncher{
		errList: utils.NewErrorList(maxCount),
	}

	if limit != 0 {
		cl.ch = make(chan int, limit)
	}

	return cl
}

func (cl *ConcurrentLauncher) Go(fct interface{}) {
	cl.wg.Add(1)
	if cl.ch != nil {
		cl.ch <- 1
	}
	go func() {
		defer func() {
			if cl.ch != nil {
				<-cl.ch
			}
			cl.wg.Done()
		}()
		cl.errList.Add(execFct(fct))
	}()
}

func (cl *ConcurrentLauncher) Wait() error {
	cl.wg.Wait()
	return cl.errList.Err()
}

func (cl *ConcurrentLauncher) Locked(fct interface{}) error {
	cl.mx.Lock()
	defer cl.mx.Unlock()
	return execFct(fct)
}

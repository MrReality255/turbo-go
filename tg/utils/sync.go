package utils

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type IConcurrentLauncher interface {
	Go(fct interface{})
	Locked(fct interface{}) error
	Wait() error
}

type ConcurrentLauncher struct {
	wg      sync.WaitGroup
	errList *ErrorList
	ch      chan int
	mx      sync.Mutex
}

type ErrorList struct {
	mx   sync.Mutex
	errs []error
}

func NewConcurrentLauncher(maxCount int, limit int) IConcurrentLauncher {
	cl := &ConcurrentLauncher{
		errList: NewErrorList(maxCount),
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

func (l *ErrorList) Add(err error) {
	if err != nil {
		l.errs = append(l.errs, err)
	}
}

func (l *ErrorList) Err() error {
	list := make([]string, 0, len(l.errs))
	for _, err := range l.errs {
		list = append(list, err.Error())
	}
	if len(list) == 0 {
		return nil
	}
	return errors.New(strings.Join(list, "; "))
}

func NewErrorList(maxCount int) *ErrorList {
	return &ErrorList{
		errs: make([]error, 0, maxCount),
	}
}

func execFct(fct interface{}) error {
	switch fct := fct.(type) {
	case func():
		fct()
		return nil
	case func() error:
		return fct()
	case nil:
		return nil
	default:
		panic(fmt.Errorf("unable to execute %T: %v", fct, fct))
	}
}

func ExecLocked(ptr *sync.Mutex, fct func()) {
	ptr.Lock()
	defer ptr.Unlock()
	fct()
}

func ExecLockedErr(pt *sync.Mutex, fct func() error) error {
	pt.Lock()
	defer pt.Unlock()
	return fct()
}

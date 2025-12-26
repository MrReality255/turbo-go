package utils

import "sync"

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

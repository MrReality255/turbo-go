package utils

import (
	"errors"
	"strings"
	"sync"
)

type ErrorList struct {
	mx   sync.Mutex
	errs []error
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

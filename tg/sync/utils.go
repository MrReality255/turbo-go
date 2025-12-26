package sync

import (
	"fmt"
)

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

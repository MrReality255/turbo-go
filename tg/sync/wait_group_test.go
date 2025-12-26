package sync

import (
	"testing"

	"github.com/MrReality255/turbo-go/tg/utils"
)

func TestWaitGroup(t *testing.T) {

	var (
		wg = NewWaitGroup(10, 2)
		l  = utils.NewStringList(1024, false)
	)
	for i := 0; i < 10; i++ {
		i := i
		wg.Go(func() error {
			l.Addf("item: %v", i)
			return nil
		})
	}
	utils.TestAsString(t, 0, "err", `<nil>`, wg.Wait())
	utils.TestAsString(t, 0, "err", `item: 0,item: 1,item: 2,item: 3,item: 4,item: 5,item: 6,item: 7,item: 8,item: 9`, l.SortJoin(","))
}

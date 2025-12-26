package sync

import (
	"testing"
	"time"

	"github.com/MrReality255/turbo-go/tg/utils"
)

func TestBatchProcessorOverflow(t *testing.T) {
	var (
		buffer = make([]int, 0, 100)
		log    = utils.NewStringList(1024, false)
	)

	bp := NewBatchProcessor(
		batchProcessorConfig{
			FlushThreads: 1,
			Timeout:      time.Millisecond * 50,
			Length:       5,
		}, func(reason BatchFlushReason) func() error {
			buffer = nil
			return func() error {
				log.Addf("flush %v", reason)
				return nil
			}
		})

	wg := NewWaitGroup(1000, 1000)
	for i := 0; i < 22; i++ {
		wg.Go(func() {
			time.Sleep(time.Millisecond)
			utils.MustSucceed(bp.Add(func() int {
				buffer = append(buffer, i)
				return 1
			}))
		})
	}

	utils.TestAsString(t, 0, "err", "<nil>", wg.Wait())
	utils.TestAsString(t, 0, "err", "<nil>", bp.Flush())
	utils.TestAsString(t, 0, "err", `[flush overflow flush overflow flush overflow flush overflow flush explicit]`, log.Content())
}

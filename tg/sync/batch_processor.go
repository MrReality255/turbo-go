package sync

import (
	"fmt"
	"sync"
	"time"

	"github.com/MrReality255/turbo-go/tg/utils"
)

type BatchFlushReason uint8

const (
	BatchFlushTimeout  BatchFlushReason = 1
	BatchFlushOverflow BatchFlushReason = 2
	BatchFlushExplicit BatchFlushReason = 3
)

type IBatchProcessor interface {
	Add(adderFct func() int) error
	Flush() error
}

type BatchProcessorConfig struct {
	FlushThreads int
	Length       int
	Timeout      time.Duration
}

type batchRec struct {
	currentLength int
}

type batchProcessor struct {
	config BatchProcessorConfig
	onSwap func(reason BatchFlushReason) func() error

	mx sync.Mutex

	batch     *batchRec
	chThreads chan int
	lastErr   error
}

func NewBatchProcessor(
	config BatchProcessorConfig,
	onSwap func(reason BatchFlushReason) func() error,
) IBatchProcessor {
	return &batchProcessor{config: config, onSwap: onSwap, chThreads: make(chan int, utils.Coalesce(config.FlushThreads, 1))}
}

func (b *batchProcessor) Add(adderFct func() int) error {
	flusher, err := b.doAdd(adderFct)
	if err != nil || flusher == nil {
		return err
	}
	return b.doFlush(flusher)
}

func (b *batchProcessor) Flush() error {
	b.mx.Lock()
	defer b.mx.Unlock()
	// wait until all threads are done
	for i := 0; i < b.config.FlushThreads; i++ {
		b.chThreads <- 1
	}

	// clean up the channel
	defer func() {
		for i := 0; i < b.config.FlushThreads; i++ {
			<-b.chThreads
		}
	}()

	flusher := b.onSwap(BatchFlushExplicit)
	b.batch = nil
	if err := flusher(); err != nil {
		return err
	}
	defer func() {
		b.lastErr = nil
	}()
	return b.lastErr
}

func (b *batchProcessor) checkTimeout(batch *batchRec, timeout time.Duration) {
	time.Sleep(timeout)
	b.mx.Lock()
	defer b.mx.Unlock()
	if batch != b.batch {
		return
	}
	flusher := b.doSwapLocked(BatchFlushTimeout)
	go b.flushTimeout(flusher)
}

func (b *batchProcessor) doAdd(adderFct func() int) (func() error, error) {
	b.mx.Lock()
	defer b.mx.Unlock()
	if err := b.lastErr; err != nil {
		b.lastErr = nil
		return nil, err
	}

	if b.batch == nil {
		b.batch = &batchRec{currentLength: 0}
		go b.checkTimeout(b.batch, b.config.Timeout)
	}

	b.batch.currentLength += adderFct()
	if b.batch.currentLength >= b.config.Length {
		return b.doSwapLocked(BatchFlushOverflow), nil
	}
	return nil, nil
}

func (b *batchProcessor) doFlush(flusher func() error) error {
	defer func() {
		<-b.chThreads
	}()
	return flusher()
}

func (b *batchProcessor) doSwapLocked(reason BatchFlushReason) func() error {
	b.chThreads <- 1
	defer func() {
		b.batch = nil
	}()
	return b.onSwap(reason)
}

func (b *batchProcessor) flushTimeout(flusher func() error) {
	err := b.doFlush(flusher)
	if err == nil {
		return
	}

	b.mx.Lock()
	defer b.mx.Unlock()
	b.lastErr = err
}

func (r BatchFlushReason) String() string {
	switch r {
	case BatchFlushTimeout:
		return "timeout"
	case BatchFlushOverflow:
		return "overflow"
	case BatchFlushExplicit:
		return "explicit"
	default:
		return fmt.Sprintf("unknown: %v", int(r))
	}
}

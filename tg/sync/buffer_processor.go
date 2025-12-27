package sync

type IBufferProcessor[T any] interface {
	Add(items ...T) error
	Flush() error
}

type bufferProcessor[T any] struct {
	buffer []T
	cfg    BatchProcessorConfig
	p      IBatchProcessor
}

func NewBufferProcessor[T any](
	cfg BatchProcessorConfig,
	onFlush func(data []T, reason BatchFlushReason) error,
) IBufferProcessor[T] {
	c := &bufferProcessor[T]{cfg: cfg}
	c.p = NewBatchProcessor(cfg, func(reason BatchFlushReason) func() error {
		procBuffer := c.buffer
		c.buffer = nil
		return func() error {
			return onFlush(procBuffer, reason)
		}
	})
	return c
}

func (b *bufferProcessor[T]) Add(items ...T) error {
	return b.p.Add(func() int {
		if b.buffer == nil {
			b.buffer = make([]T, 0, b.cfg.Length)
		}
		b.buffer = append(b.buffer, items...)
		return len(items)
	})
}

func (b *bufferProcessor[T]) Flush() error {
	return b.p.Flush()
}

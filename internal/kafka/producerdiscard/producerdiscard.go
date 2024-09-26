/*
Abstract producer for tests
*/
package producerdiscard

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type DiscardProducer struct{}

func NewDiscardProducer() *DiscardProducer {
	return &DiscardProducer{}
}

func (w *DiscardProducer) Close() error {
	return nil
}

func (w *DiscardProducer) Stats() kafka.WriterStats {
	return kafka.WriterStats{}
}

func (w *DiscardProducer) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	return nil
}
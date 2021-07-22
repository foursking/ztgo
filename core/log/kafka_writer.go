package log

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type kafkaWriterOptions struct {
	Topic   string   `toml:"topic" json:"topic"`
	Brokers []string `toml:"brokers" json:"brokers"`
}

type kafkaWriter struct {
	writer *kafka.Writer
}

// Write writes log to kafka
func (kw *kafkaWriter) Write(bs []byte) (n int, err error) {
	// copy是因为zap有自己的buffer，此处bs参数地址中的值可能随时被更改
	v := make([]byte, len(bs))
	copy(v, bs)
	err = kw.writer.WriteMessages(context.TODO(), kafka.Message{Value: v})
	return
}

// Sync shuts down the producer and waits for any buffered messages to be flushed
func (kw *kafkaWriter) Sync() error {
	return kw.writer.Close()
}

func newKafkaWriter(opt *kafkaWriterOptions) (*kafkaWriter, error) {
	kw := kafkaWriter{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:       opt.Brokers,
			Topic:         opt.Topic,
			QueueCapacity: 10240,
			Async:         true,
		}),
	}
	return &kw, nil
}

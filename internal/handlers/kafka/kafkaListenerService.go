package kafka_handler

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func listenKafkaHanlderFunc(topic string, groupID string, maximumRetry int, fn func(kafka.Message) error) {
	kr := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		Topic:       topic,
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.FirstOffset,
		GroupID:     groupID,
	})

	kw := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Balancer: &kafka.LeastBytes{},
		// Topic:    "heyhey",
	}

	go func() {
		for {
			ctx := context.Background()
			msg, err := kr.FetchMessage(ctx)

			if err != nil {
				fmt.Println("Failed to consume message", err)
			} else {
				withRetry(ctx, kr, kw, msg, maximumRetry, fn)
			}
		}
	}()
}

func withRetry(ctx context.Context, kr *kafka.Reader, kw *kafka.Writer, msg kafka.Message, maximumRetry int, fn func(kafka.Message) error) {
	count := 0
	for {
		err := fn(msg)
		fmt.Println(err)
		if err == nil {
			break
		}

		if count > maximumRetry {
			sendToDLQ(context.Background(), kw, msg)
			break
		}
		count++
	}
	kr.CommitMessages(ctx, msg)
}

func sendToDLQ(ctx context.Context, kw *kafka.Writer, msg kafka.Message) {
	msg.Topic = fmt.Sprintf("%s.DLQ", msg.Topic)
	err := kw.WriteMessages(ctx, msg)
	if err != nil {
		fmt.Println(err)
	}
}

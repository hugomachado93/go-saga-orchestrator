package payment_service

import (
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type PaymentService struct {
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

func (ps *PaymentService) AuthorizePayment(msg kafka.Message) error {
	fmt.Println("Authorizing payment...")
	time.Sleep(5 * time.Second)

	fmt.Println(msg.Value)

	// a := &saga.Response{SagaUUID: "123456", Payload: "mamas", Status: saga.SUCCESS}
	// ja, _ := json.Marshal(a)

	// w := &kafka.Writer{
	// 	Addr:     kafka.TCP("localhost:9092"),
	// 	Topic:    "teste3",
	// 	Balancer: &kafka.LeastBytes{},
	// }

	// _ = w.WriteMessages(context.Background(),
	// 	kafka.Message{
	// 		Value: ja,
	// 	},
	// )
	return nil
}

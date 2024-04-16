package services

import (
	"encoding/json"
	"fmt"
	"main/internal/domain/saga"
	"main/internal/infrastructure/kfk"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
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

	var s saga.Response

	json.Unmarshal(msg.Value, &s)

	h := protocol.Header{Key: "x-api-key", Value: []byte("644e032b-3ee0-441e-a10c-d265a986ca2c")}

	headers := make([]protocol.Header, 0)
	headers = append(headers, h)

	rjson, _ := json.Marshal(saga.Response{SagaUUID: s.SagaUUID, SagaName: s.SagaName, Payload: "payload2", Event: "PAYMENT_AUTHORIZED"})

	err := kfk.SendMessage(string(rjson), "APP_ORCHESTRATOR", headers, "")
	if err != nil {
		return fmt.Errorf("error %w", err)
	}

	return nil
}

func (ps *PaymentService) CapturePayment(msg kafka.Message) error {
	fmt.Println("Capturing payment...")
	time.Sleep(5 * time.Second)

	fmt.Println(msg.Value)

	var s saga.Response

	json.Unmarshal(msg.Value, &s)

	h := protocol.Header{Key: "x-api-key", Value: []byte("644e032b-3ee0-441e-a10c-d265a986ca2c")}

	headers := make([]protocol.Header, 0)
	headers = append(headers, h)

	rjson, _ := json.Marshal(saga.Response{SagaUUID: s.SagaUUID, SagaName: s.SagaName, Payload: "payload3", Event: "PAYMENT_CAPTURED"})

	err := kfk.SendMessage(string(rjson), "APP_ORCHESTRATOR", headers, "")
	if err != nil {
		return fmt.Errorf("error %w", err)
	}

	return nil
}

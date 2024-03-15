package orchestrator_service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"main/internal/domain/saga"
	"main/internal/repository"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
)

type OrchestratorService struct {
	transction *repository.Transaction
	sagaRepo   *repository.SagaRepository
}

func NewOrchestratorService(transction *repository.Transaction,
	sagaRepo *repository.SagaRepository) *OrchestratorService {
	return &OrchestratorService{transction: transction, sagaRepo: sagaRepo}
}

func (o *OrchestratorService) OrchestrateSaga(msg kafka.Message) error {
	ctx := context.Background()

	err := o.transction.WithTransaction(ctx, func(ctx context.Context) error {
		var response saga.Response
		var s *saga.Saga
		var err error
		isNewSaga := false

		json.Unmarshal(msg.Value, &response)
		s, err = o.sagaRepo.FindSagaByUUID(ctx, response.SagaUUID)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s = saga.NewSaga()
				isNewSaga = true
			} else {
				fmt.Println("Failed to findSaga")
				return err
			}
		}
		fmt.Println("Prepering saga")
		s.PrepareSaga(response)

		if isNewSaga {
			err = o.sagaRepo.InsertSaga(ctx, s)
		} else {
			err = o.sagaRepo.UpdateSaga(ctx, s)
		}

		if err != nil {
			return err
		}

		fmt.Println("send saga command")

		w := &kafka.Writer{
			Addr:     kafka.TCP("localhost:9092"),
			Balancer: &kafka.LeastBytes{},
		}

		fmt.Println(string(s.CurrentCommand))

		headers := make([]protocol.Header, 0)
		h := protocol.Header{Key: "saga_uuid", Value: []byte(s.SagaUUID)}
		headers = append(headers, h)

		err = w.WriteMessages(context.Background(),
			kafka.Message{
				Value:   []byte(response.Payload),
				Topic:   string(s.CurrentCommand),
				Headers: headers,
				Key:     []byte(s.SagaUUID),
			},
		)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil

}

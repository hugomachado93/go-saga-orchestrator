package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"main/internal/domain/saga"
	"main/internal/err"
	kafka_adapter "main/internal/kafka"
	"main/internal/repository"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type OrchestratorService struct {
	transaction *repository.Transaction
	sagaRepo    *repository.SagaRepository
	stmRepo     *repository.StateMachineRepository
}

func NewOrchestratorService(transction *repository.Transaction,
	sagaRepo *repository.SagaRepository, stmRepo *repository.StateMachineRepository) *OrchestratorService {
	return &OrchestratorService{transaction: transction, sagaRepo: sagaRepo, stmRepo: stmRepo}
}

func (o *OrchestratorService) OrchestrateSaga(msg kafka.Message) error {
	ctx := context.Background()

	return o.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		var response saga.Response

		apkey, err := getApiKey(msg)
		if err != nil {
			return err
		}

		json.Unmarshal(msg.Value, &response)

		s, err := o.sagaRepo.FindSagaByUUID(ctx, response.SagaUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				s = saga.NewSaga(apkey, uuid.New().String(), response.SagaName)
			} else {
				fmt.Println("Failed to findSaga")
				return err
			}
		}

		s.ReceiveResponse()

		if err := o.sagaRepo.UpdateSaga(ctx, s); err != nil {
			return fmt.Errorf("failed to update saga %w", err)
		}

		stm, err := o.stmRepo.FindSettingsByName(ctx, response.SagaName)
		if err != nil {
			return err
		}

		fmt.Println("Preparing saga")
		nextState, err := stm.FindNextStep(s.CurrentState, response.Event)
		if err != nil {
			return err
		}

		s.UpdateSaga(response.Payload, nextState)

		if err := o.sagaRepo.InsertSaga(ctx, s); err != nil {
			return err
		}

		return nil
	})
}

func getApiKey(msg kafka.Message) (string, error) {
	for _, v := range msg.Headers {
		if v.Key == "x-api-key" {
			return string(v.Value), nil
		}
	}

	return "", err.ErroFailedFindApiKey
}

func (o *OrchestratorService) SendCommand() {
	o.transaction.WithTransaction(context.Background(), func(ctx context.Context) error {
		sagas, err := o.sagaRepo.FindSagasByStatus(ctx, saga.WAITING_PROCESS)
		if err != nil {
			return err
		}

		for _, saga := range sagas {
			if err := kafka_adapter.SendMessage(saga.Payload, saga.CurrentState, nil, saga.SagaUUID); err == nil {
				saga.RequestSaga()
				o.sagaRepo.UpdateSaga(ctx, &saga)
			} else {
				fmt.Errorf("failed to send message %w", err)
			}
		}

		return nil
	})
}

func (o *OrchestratorService) CheckTimeout() {
	o.transaction.WithTransaction(context.Background(), func(ctx context.Context) error {
		sagas, err := o.sagaRepo.FindSagasByStatus(ctx, saga.REQUEST_SENT)
		if err != nil {
			return err
		}

		for _, saga := range sagas {
			if !saga.HasSagaTimedout() {
				continue
			}

			if err := kafka_adapter.SendMessage(saga.Payload, saga.CurrentState, nil, saga.SagaUUID); err == nil {
				saga.RequestSaga()
				o.sagaRepo.UpdateSaga(ctx, &saga)
			} else {
				fmt.Println("failed to send message %w", err)
			}
		}

		return nil
	})
}

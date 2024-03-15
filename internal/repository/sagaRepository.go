package repository

import (
	"context"
	"main/internal/domain/saga"
)

type SagaRepository struct {
}

func NewSagaRepositor() *SagaRepository {
	return &SagaRepository{}
}

func (sr *SagaRepository) InsertSaga(ctx context.Context, saga *saga.Saga) error {
	tx := extractTx(ctx)
	sql := "INSERT INTO saga (saga_name, saga_uuid, saga_state, saga_index, payload, last_update) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := tx.Exec(sql, saga.SagaName, saga.Id, saga.Status, saga.Idx, saga.Payload, saga.LastUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (sr *SagaRepository) UpdateSaga(ctx context.Context, saga *saga.Saga) error {
	tx := extractTx(ctx)
	sql := "UPDATE saga SET saga_name = $1, saga_uuid = $2, saga_index = $3, payload = $4, saga_state = $5, last_update = $6"
	_, err := tx.Exec(sql, saga.SagaName, saga.SagaUUID, saga.Idx, saga.Payload, saga.Status, saga.LastUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SagaRepository) FindSagaByUUID(ctx context.Context, UUID string) (*saga.Saga, error) {
	tx := extractTx(ctx)
	sql := "SELECT saga_uuid, saga_name, saga_uuid, saga_index, payload, saga_state, last_update FROM saga WHERE saga_uuid = $1"
	r := tx.QueryRow(sql, UUID)
	if err := r.Err(); err != nil {
		return nil, err
	}

	var saga saga.Saga

	err := r.Scan(&saga.Id, &saga.SagaName, &saga.Idx, &saga.Payload, &saga.Status, &saga.LastUpdate)
	if err != nil {
		return nil, err
	}

	return &saga, nil
}

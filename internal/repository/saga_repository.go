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
	sql := "INSERT INTO saga (api_key, name, uuid, state, status, payload, timeout, created_at, last_update) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	_, err := tx.Exec(sql, saga.ApiKey, saga.SagaName, saga.SagaUUID, saga.CurrentState, saga.Status, saga.Payload, saga.Timeout, saga.CreatedAt, saga.LastUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (sr *SagaRepository) UpdateSaga(ctx context.Context, saga *saga.Saga) error {
	tx := extractTx(ctx)
	sql := "UPDATE saga SET api_key = $1, name = $2, uuid = $3, state = $4, status = $5, payload = $6, timeout = $7, created_at = $8, last_update =$9 where id = $10"
	_, err := tx.Exec(sql, saga.ApiKey, saga.SagaName, saga.SagaUUID, saga.CurrentState, saga.Status, saga.Payload, saga.Timeout, saga.CreatedAt, saga.LastUpdate, saga.Id)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SagaRepository) FindSagaByUUID(ctx context.Context, UUID string) (*saga.Saga, error) {
	tx := extractTx(ctx)
	sql := "SELECT id, api_key, name, uuid, payload, state, status, timeout, created_at, last_update FROM saga WHERE uuid = $1 order by created_at desc limit 1"
	r := tx.QueryRow(sql, UUID)
	if err := r.Err(); err != nil {
		return nil, err
	}

	var saga saga.Saga

	err := r.Scan(&saga.Id, &saga.ApiKey, &saga.SagaName, &saga.SagaUUID, &saga.Payload, &saga.CurrentState, &saga.Status, &saga.Timeout, &saga.CreatedAt, &saga.LastUpdate)
	if err != nil {
		return nil, err
	}

	return &saga, nil
}

func (sr *SagaRepository) FindSagasByStatus(ctx context.Context, status saga.Status) ([]saga.Saga, error) {
	tx := extractTx(ctx)
	sql := "SELECT id, api_key, name, uuid, payload, state, status, timeout, created_at, last_update FROM saga WHERE status = $1 order by created_at desc limit 1"
	r, err := tx.Query(sql, status)
	if err != nil {
		return nil, err
	}

	var sagas []saga.Saga

	for r.Next() {
		var saga saga.Saga
		err := r.Scan(&saga.Id, &saga.ApiKey, &saga.SagaName, &saga.SagaUUID, &saga.Payload, &saga.CurrentState, &saga.Status, &saga.Timeout, &saga.CreatedAt, &saga.LastUpdate)
		if err != nil {
			return nil, err
		}
		sagas = append(sagas, saga)
	}

	return sagas, nil
}

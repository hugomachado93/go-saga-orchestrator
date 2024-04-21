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
	sql := "INSERT INTO saga_command (api_key, name, uuid, state, status, payload, timeout, created_at, last_update, delayed_message) VALUES (:api_key, :name, :uuid, :state, :status, :payload, :timeout, :created_at, :last_update, :delayed_message)"
	_, err := tx.NamedExec(sql, saga)
	if err != nil {
		return err
	}

	return nil
}

func (sr *SagaRepository) UpdateSaga(ctx context.Context, saga *saga.Saga) error {
	tx := extractTx(ctx)
	sql := "UPDATE saga_command SET api_key = :api_key, name = :name, uuid = :uuid, state = :state, status = :status, payload = :paylaod, timeout = :timeout, created_at = :created_at, last_update = :last_update, delayed_message = :delayed_message where id = :id"
	_, err := tx.NamedExec(sql, saga)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SagaRepository) FindSagaByUUID(ctx context.Context, UUID string) (*saga.Saga, error) {
	tx := extractTx(ctx)
	sql := "SELECT id, api_key, name, uuid, payload, state, status, timeout, created_at, last_update, delayed_message FROM saga_command WHERE uuid = ? order by created_at desc limit 1"
	r := tx.QueryRowx(sql, UUID)
	if err := r.Err(); err != nil {
		return nil, err
	}

	var saga saga.Saga

	err := r.StructScan(&saga)
	if err != nil {
		return nil, err
	}

	return &saga, nil
}

func (sr *SagaRepository) FindSagasByStatus(ctx context.Context, status saga.Status) ([]saga.Saga, error) {
	tx := extractTx(ctx)
	sql := "SELECT id, api_key, name, uuid, payload, state, status, timeout, created_at, last_update, delayed_message FROM saga_command WHERE status = ? and delayed_message is null  order by created_at desc limit 1"
	r, err := tx.Queryx(sql, status)
	if err != nil {
		return nil, err
	}

	var sagas []saga.Saga

	for r.Next() {
		var saga saga.Saga
		err := r.StructScan(&saga)
		if err != nil {
			return nil, err
		}
		sagas = append(sagas, saga)
	}

	return sagas, nil
}

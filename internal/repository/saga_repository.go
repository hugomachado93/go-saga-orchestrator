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

func (r *SagaRepository) SaveSaga(ctx context.Context, saga *saga.Saga) error {
	tx := extractTx(ctx)
	if saga.Id != nil {
		_, err := tx.NamedExec("UPDATE saga SET api_key = :api_key, uuid = :uuid, name = :name, state = :state, status = :status, last_update = :last_update WHERE id = id", saga)
		if err != nil {
			return err
		}
	} else {
		res, err := tx.PrepareNamed("INSERT INTO saga (api_key, uuid, name, state, status, created_at, last_update) VALUES (:api_key, :uuid, :name, :state, :status, :created_at, :last_update) RETURNING id")
		if err != nil {
			return err
		}
		var id int64
		err = res.Get(&id, saga)
		if err != nil {
			return err
		}
		saga.Id = &id
	}

	// Insert associated saga commands
	for _, command := range saga.SagaCommands {
		command.SagaId = saga.Id
		_, err := tx.NamedExec("INSERT INTO saga_command (payload, state, status, timeout, created_at, last_update, saga_id) VALUES (:payload, :state, :status, :timeout, :created_at, :last_update, :saga_id)", command)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetSaga retrieves a saga record by ID.
func (r *SagaRepository) GetSagaByUUID(ctx context.Context, uuid string) (*saga.Saga, error) {
	tx := extractTx(ctx)

	var saga saga.Saga
	err := tx.Get(&saga, "SELECT * FROM saga WHERE uuid = $1", uuid)
	if err != nil {
		return nil, err
	}

	// Retrieve associated saga commands
	err = tx.Select(&saga.SagaCommands, "SELECT * FROM saga_command WHERE saga_id = $1", saga.Id)
	if err != nil {
		return nil, err
	}

	return &saga, nil
}

// UpdateSaga updates a saga record in the database.
func (r *SagaRepository) UpdateSaga(ctx context.Context, saga *saga.Saga) error {
	tx := extractTx(ctx)

	// Update saga record
	_, err := tx.Exec("UPDATE saga SET api_key = $1, uuid = $2, name = $3, state = $4, status = $5, last_update = $6 WHERE id = $7", saga.ApiKey, saga.UUID, saga.Name, saga.State, saga.Status, saga.LastUpdate, saga.Id)
	if err != nil {
		return err
	}

	// Insert updated associated saga commands
	for _, command := range saga.SagaCommands {
		_, err := tx.Exec("UPDATE saga_command SET payload = $1, state = $2, status = $3, timeout = $4, created_at = $5, last_update = $6, saga_id = $7 where id = $8",
			command.Payload, command.CurrentState, command.Status, command.Timeout, command.CreatedAt, command.LastUpdate, saga.Id, command.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

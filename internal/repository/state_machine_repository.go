package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/domain/statemachine"
	"time"
)

type StateMachineRepository struct {
}

func NewStatemachineRepository() *StateMachineRepository {
	return &StateMachineRepository{}
}

func (sr *StateMachineRepository) InsertSettings(ctx context.Context, stm statemachine.IStatemachineDefintion) error {
	tx := extractTx(ctx)

	sql := "INSERT INTO statemachine_settings (client_api_key, name, context, created_at) VALUES ($1, $2, $3, $4)"
	_, err := tx.Exec(sql, stm.XApiKey(), stm.GetName(), stm.ContextToJson(), time.Now())
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (sr *StateMachineRepository) FindSettingsByName(ctx context.Context, name string) (*statemachine.StatemachineDefinition, error) {
	tx := extractTx(ctx)

	var stm statemachine.StatemachineDefinition
	var context string
	var teste string

	sql := "SELECT id, client_api_key, name, context, created_at from statemachine_settings WHERE name = $1 limit 1"

	r := tx.QueryRow(sql, name)
	if err := r.Err(); err != nil {
		return nil, fmt.Errorf("failed %w", err)
	}

	err := r.Scan(&stm.Id, &teste, &stm.Name, &context, &stm.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed %w", err)
	}

	json.Unmarshal([]byte(context), &stm.Context)

	return &stm, nil
}

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
	_, err := tx.Exec(sql, stm.GetName(), stm.GetName(), stm.ContextToJson(), time.Now())
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

	sql := "select * from statemachine_settings where name = $1"

	r := tx.QueryRow(sql, name)
	if r.Err() != nil {
		return nil, fmt.Errorf("error here: %w", r.Err())
	}

	r.Scan(&stm.Id, &stm.ClientApiKey, &stm.Name, &context, stm.CreatedAt)

	json.Unmarshal([]byte(context), &stm.Context)

	return &stm, nil
}

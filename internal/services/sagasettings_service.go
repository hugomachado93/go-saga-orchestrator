package services

import (
	"context"
	"main/internal/domain/statemachine"
	"main/internal/repository"
)

type SagaSettingsService struct {
	cr *repository.StateMachineRepository
	t  *repository.Transaction
}

func NewStateMachine(cr *repository.StateMachineRepository, t *repository.Transaction) *SagaSettingsService {
	return &SagaSettingsService{cr: cr, t: t}
}

func (c *SagaSettingsService) InsertStateMachineSettings(sm *statemachine.Statemachine, xApiKey string) error {
	return c.t.WithTransaction(context.Background(), func(ctx context.Context) error {
		stm, err := statemachine.NewDefinition(xApiKey, sm)
		if err != nil {
			return err
		}
		c.cr.InsertSettings(ctx, stm)
		return nil
	})
}

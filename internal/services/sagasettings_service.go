package services

import (
	"context"
	"main/internal/api/requests"
	"main/internal/repository"
)

type SagaSettingsService struct {
	cr *repository.StateMachineRepository
	t  *repository.Transaction
}

func NewStateMachine(cr *repository.StateMachineRepository, t *repository.Transaction) *SagaSettingsService {
	return &SagaSettingsService{cr: cr, t: t}
}

func (c *SagaSettingsService) InsertStateMachineSettings(stm *requests.Statemachine, xApiKey string) {
	c.t.WithTransaction(context.Background(), func(ctx context.Context) error {

		c.cr.InsertSettings(ctx, stm, xApiKey)
		return nil
	})
}

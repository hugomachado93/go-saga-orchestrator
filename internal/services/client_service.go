package services

import (
	"context"
	"main/internal/repository"
)

type ClientService struct {
	cr *repository.ClientRepository
	tr *repository.Transaction
}

func NewClientService(cr *repository.ClientRepository, tr *repository.Transaction) *ClientService {
	return &ClientService{cr: cr, tr: tr}
}

func (cs *ClientService) IsClientAuthorized(api_key string) bool {
	var isClientAuthorized bool
	cs.tr.WithTransaction(context.Background(), func(ctx context.Context) error {
		c, err := cs.cr.FindClientByApiKey(ctx, api_key)
		if err != nil {
			return err
		}

		if c != nil {
			isClientAuthorized = true
		} else {
			isClientAuthorized = false
		}

		return nil
	})

	return isClientAuthorized
}

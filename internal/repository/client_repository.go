package repository

import (
	"context"
	"fmt"
	"main/internal/domain/client"
)

type ClientRepository struct{}

func NewClientRepository() *ClientRepository {
	return &ClientRepository{}
}

func (cr *ClientRepository) FindClientByApiKey(ctx context.Context, api_key string) (*client.Client, error) {
	tx := extractTx(ctx)
	sql := "SELECT id, name, api_key, last_update FROM clients WHERE api_key = $1"
	r := tx.QueryRow(sql, api_key)
	if err := r.Err(); err != nil {
		fmt.Print(err)
		return nil, err
	}

	var c client.Client

	err := r.Scan(&c.Id, &c.Name, &c.ApiKey, &c.LastUpdate)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

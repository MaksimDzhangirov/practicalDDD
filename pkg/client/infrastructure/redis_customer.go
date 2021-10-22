package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/client/domain/model"
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/client/infrastructure/dto"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

// Репозиторий для Redis

type CustomerRedisRepository struct {
	client *redis.Client
}

func (r *CustomerRedisRepository) GetCustomer(ctx context.Context, ID uuid.UUID) (*model.Customer, error) {
	data, err := r.client.Get(ctx, fmt.Sprintf("user-%s", ID.String())).Result()
	if err != nil {
		return nil, err
	}

	var row dto.CustomerJSON
	err = json.Unmarshal([]byte(data), &row)
	if err != nil {
		return nil, err
	}

	customer, err := row.ToEntity()
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

package infrastructure

import (
	"context"
	"encoding/json"
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/client/domain/model"
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/client/infrastructure/dto"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"path"
)

// API

type CustomerRedisAPIRepository struct {
	client  *http.Client
	baseUrl string
}

func (r *CustomerRedisAPIRepository) GetCustomer(ctx context.Context, ID uuid.UUID) (*model.Customer, error) {
	resp, err := r.client.Get(path.Join(r.baseUrl, "users", ID.String()))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var row dto.CustomerJSON
	err = json.Unmarshal(data, &row)
	if err != nil {
		return nil, err
	}

	customer, err := row.ToEntity()
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

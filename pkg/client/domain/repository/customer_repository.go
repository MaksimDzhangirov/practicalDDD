package repository

import (
	"context"
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/client/domain/model"
	"github.com/google/uuid"
)

type CustomerSpecification struct {
}

type Customers []model.Customer

type CustomerRepository interface {
	GetCustomer(ctx context.Context, ID uuid.UUID) (*model.Customer, error)
	Search(ctx context.Context, specification CustomerSpecification) (Customers, int, error)
	SaveCustomer(ctx context.Context, customer model.Customer) (*model.Customer, error)
	UpdateCustomer(ctx context.Context, customer model.Customer) (*model.Customer, error)
	DeleteCustomer(ctx context.Context, ID uuid.UUID) (*model.Customer, error)
}

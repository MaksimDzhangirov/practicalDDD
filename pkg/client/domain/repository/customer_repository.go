package repository

import (
	"context"
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/client/domain/model"
)

type CustomerSpecification struct {

}

type CustomerRepository interface {
	Search(ctx context.Context, specification CustomerSpecification) ([]model.Customer, error)
	Create(ctx context.Context, customer model.Customer) (*model.Customer, error)
	UpdatePerson(ctx context.Context, customer model.Customer) (*model.Customer, error)
	UpdateCompany(ctx context.Context, customer model.Customer) (*model.Customer, error)
	//
	// и много других методов
	//
}

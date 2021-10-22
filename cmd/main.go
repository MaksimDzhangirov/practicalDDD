package main

import (
	"fmt"
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/product/domain/model"
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/product/infrastructure/create"
	"github.com/google/uuid"
)

func main() {
	spec := create.NewAndSpecification(
		create.NewHasAtLeast(10),
		create.FunctionSpecification(create.IsPlastic),
		create.FunctionSpecification(create.IsDeliverable),
	)

	fmt.Printf("%+v", spec.Create(model.Product{
		ID: uuid.New(),
	}))
	// выводит: {ID:befaf2b9-73cd-44cf-95f1-5fba087e46d9 Material:plastic IsDeliverable:true Quantity:10}
}

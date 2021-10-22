package create

import "github.com/MaksimDzhangirov/PracticalDDD/pkg/product/domain/model"

type ProductSpecification interface {
	Create(product model.Product) model.Product
}

type AndSpecification struct {
	specifications []ProductSpecification
}

func NewAndSpecification(specifications ...ProductSpecification) ProductSpecification {
	return AndSpecification{
		specifications: specifications,
	}
}

func (s AndSpecification) Create(product model.Product) model.Product {
	for _, specification := range s.specifications {
		product = specification.Create(product)
	}
	return product
}

type HasAtLeast struct {
	pieces int
}

func NewHasAtLeast(pieces int) ProductSpecification {
	return HasAtLeast{
		pieces: pieces,
	}
}

func (h HasAtLeast) Create(product model.Product) model.Product {
	product.Quantity = h.pieces
	return product
}

func IsPlastic(product model.Product) model.Product {
	product.Material = model.Plastic
	return product
}

func IsDeliverable(product model.Product) model.Product {
	product.IsDeliverable = true
	return product
}

type FunctionSpecification func(product model.Product) model.Product

func (fs FunctionSpecification) Create(product model.Product) model.Product {
	return fs(product)
}
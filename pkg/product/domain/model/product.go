package model

import "github.com/google/uuid"

type MaterialType = string

const Plastic = "plastic"

type Product struct {
	ID            uuid.UUID
	Material      MaterialType
	IsDeliverable bool
	Quantity      int
}

type ProductSpecification interface {
	IsValid(product Product) bool
}

type AndSpecification struct {
	specifications []ProductSpecification
}

func NewAndSpecification(specifications ...ProductSpecification) ProductSpecification {
	return AndSpecification{
		specifications: specifications,
	}
}

func (s AndSpecification) IsValid(product Product) bool {
	for _, specification := range s.specifications {
		if !specification.IsValid(product) {
			return false
		}
	}

	return true
}

type OrSpecification struct {
	specifications []ProductSpecification
}

func NewOrSpecification(specifications ...ProductSpecification) ProductSpecification {
	return OrSpecification{
		specifications: specifications,
	}
}

func (s OrSpecification) IsValid(product Product) bool {
	for _, specification := range s.specifications {
		if specification.IsValid(product) {
			return true
		}
	}

	return false
}

type NotSpecification struct {
	specification ProductSpecification
}

func NewNotSpecification(specification ProductSpecification) ProductSpecification {
	return NotSpecification{
		specification: specification,
	}
}

func (s NotSpecification) IsValid(product Product) bool {
	return !s.specification.IsValid(product)
}

type HasAtLeast struct {
	pieces int
}

func NewHasAtLeast(pieces int) ProductSpecification {
	return HasAtLeast{
		pieces: pieces,
	}
}

func (h HasAtLeast) IsValid(product Product) bool {
	return product.Quantity >= h.pieces
}

func IsPlastic(product Product) bool {
	return product.Material == Plastic
}

func IsDeliverable(product Product) bool {
	return product.IsDeliverable
}

type FunctionSpecification func(product Product) bool

func (fs FunctionSpecification) IsValid(product Product) bool {
	return fs(product)
}

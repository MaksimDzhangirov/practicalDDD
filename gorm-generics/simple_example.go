package main

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/ompluscator/gorm-generics"
	// какие-то импорты
)

// Product - сущность предметной области
type Product struct {
	// какие-то поля
}

// ProductGorm - это DTO для сопоставления сущности Product с базой данных
type ProductGorm struct {
	// какие-то поля
}

// ToEntity соответствует интерфейсу gorm_generics.GormModel
func (g ProductGorm) ToEntity() Product {
	return Product{
		// какие-то поля
	}
}

// FromEntity соответствует интерфейсу gorm_generics.GormModel
func (g ProductGorm) FromEntity(product Product) interface{} {
	return ProductGorm{
		// какие-то поля
	}
}

func main() {
	db, err := gorm.Open( /* строка подключения к БД */ )
	// обработка ошибки

	err = db.AutoMigrate(ProductGorm{})
	// обработка ошибки

	// инициализируем новый репозиторий, передавая
	// GORM модель и сущность как тип
	repository := gorm_generics.NewRepository[ProductGorm, Product](db)

	ctx := context.Background()

	// создаём новую сущность
	product := Product{
		// какие-то поля
	}

	// посылаем новую сущность в репозиторий для сохранения
	err = repository.Insert(ctx, &product)
	// обработка ошибки

	fmt.Println(product)
	// Выводит:
	// {1 product1 100 true}

	single, err := repository.FindByID(ctx, product.ID)
	// обработка ошибки

	fmt.Println(single)
	// Выводит:
	// {1 product1 100 true}
}

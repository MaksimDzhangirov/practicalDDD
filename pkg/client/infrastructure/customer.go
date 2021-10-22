package infrastructure

import (
	"context"
	"errors"
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/client/domain/model"
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/client/infrastructure/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// инфраструктурный уровень

type CustomerRepository struct {
	connection *gorm.DB
}

func (r *CustomerRepository) GetCustomer(ctx context.Context, ID uuid.UUID) (*model.Customer, error) {
	var row dto.CustomerGorm
	err := r.connection.WithContext(ctx).Where("uuid = ?", ID).First(&row).Error
	if err != nil {
		return nil, err
	}

	customer, err := row.ToEntity()
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerRepository) SaveCustomer(ctx context.Context, customer model.Customer) (*model.Customer, error) {
	row := NewRow(customer)
	err := r.connection.WithContext(ctx).Save(&row).Error
	if err != nil {
		return nil, err
	}

	customer, err = row.ToEntity()
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerRepository) CreateCustomer(ctx context.Context, customer model.Customer) (*model.Customer, error) {
	tx := r.connection.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	//
	// какой-то код
	//

	var total int64
	var err error
	if customer.Person != nil {
		err = tx.Model(dto.PersonGorm{}).Where("ssn = ?", customer.Person.SSN).Count(&total).Error
	} else if customer.Person != nil {
		err = tx.Model(dto.CompanyGorm{}).Where("registration_number = ?", customer.Company.RegistrationNumber).Count(&total).Error
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	} else if total > 0 {
		tx.Rollback()
		return nil, errors.New("there is already such record in DB")
	}

	//
	// какой-то код
	//
	row := NewRow(customer)
	err = tx.Save(&row).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	customer, err = row.ToEntity()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &customer, nil
}

//
// другие методы
//

func NewRow(customer model.Customer) dto.CustomerGorm {
	var person *dto.PersonGorm
	if customer.Person != nil {
		person = &dto.PersonGorm{
			SSN:       customer.Person.SSN,
			FirstName: customer.Person.FirstName,
			LastName:  customer.Person.LastName,
			Birthday:  time.Time(customer.Person.Birthday),
		}
	}

	var company *dto.CompanyGorm
	if customer.Company != nil {
		company = &dto.CompanyGorm{
			Name:               customer.Company.Name,
			RegistrationNumber: customer.Company.RegistrationNumber,
			RegistrationDate:   customer.Company.RegistrationDate,
		}
	}

	return dto.CustomerGorm{
		UUID:     uuid.NewString(),
		Person:   person,
		Company:  company,
		Street:   customer.Address.Street,
		Number:   customer.Address.Number,
		Postcode: customer.Address.Postcode,
		City:     customer.Address.City,
	}
}

package fake

import access_model "github.com/MaksimDzhangirov/PracticalDDD/pkg/access/domain/model"

type UserFakeRepository struct {
	//
	// какие-то поля
	//
}

func (r *UserFakeRepository) Create(user access_model.User) error {
	//
	// какие-то код
	//
	return nil
}
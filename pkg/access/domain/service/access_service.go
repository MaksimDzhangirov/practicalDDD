package access_service

import access_model "github.com/MaksimDzhangirov/PracticalDDD/pkg/access/domain/model"

type UserService interface {
	Create(user access_model.User) error
	//
	// какие-то методы
	//
}

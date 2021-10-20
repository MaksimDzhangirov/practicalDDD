package access_repository

import access_model "github.com/MaksimDzhangirov/PracticalDDD/pkg/access/domain/model"

type UserRepository interface {
	Create(user access_model.User) error
}

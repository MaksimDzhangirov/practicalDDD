package access

import (
	access_repository "github.com/MaksimDzhangirov/PracticalDDD/pkg/access/domain/repository"
	access_service "github.com/MaksimDzhangirov/PracticalDDD/pkg/access/domain/service"
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/access/infrastructure/database"
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/access/infrastructure/fake"
)

type AccessModule struct {
	repository access_repository.UserRepository
	service    access_service.UserService
}

func NewAccessModule(useDatabase bool) *AccessModule {
	var repository access_repository.UserRepository
	if useDatabase {
		repository = &database.UserDBRepository{}
	} else {
		repository = &fake.UserFakeRepository{}
	}
	var service access_service.UserService
	//
	// какой-то код
	//
	return &AccessModule{
		repository: repository,
		service:    service,
	}
}

func (m *AccessModule) GetRepository() access_repository.UserRepository {
	return m.repository
}

func (m *AccessModule) GetService() access_service.UserService {
	return m.service
}
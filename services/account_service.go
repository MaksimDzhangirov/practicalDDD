package services

import (
	"fmt"
	"github.com/MaksimDzhangirov/PracticalDDD/domain/services"
	domain "github.com/MaksimDzhangirov/PracticalDDD/domain/userAccount/entity"
	"net/http"
)

// инфраструктурный уровень
type AccountAPIService struct {
	client *http.Client
}

func NewAccountService(client *http.Client) services.AccountService {
	return &AccountAPIService{
		client: client,
	}
}

func (s *AccountAPIService) Update(account domain.Account) error {
	var request *http.Request
	//
	// какой-то код
	//
	response, err := s.client.Do(request)
	if err != nil {
		return err
	}
	//
	// какой-то код
	//
	fmt.Printf("Response code: %d", response.StatusCode)
	return nil
}

package entity

import domain "github.com/MaksimDzhangirov/PracticalDDD/domain/userAccount/entity"

type Bonus struct {
	Name string
}

func (b *Bonus) Apply(account *domain.Account) error {
	//
	// какой-то код
	//

	return nil
}

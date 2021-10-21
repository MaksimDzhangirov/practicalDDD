package infrastructure

import (
	"github.com/MaksimDzhangirov/PracticalDDD/infrastructure/bankAccount/dto"
	"github.com/MaksimDzhangirov/PracticalDDD/pkg/investment/domain/model"
	value_objects "github.com/MaksimDzhangirov/PracticalDDD/value-objects"
	"github.com/google/uuid"
)

type CryptoCurrencyGorm struct {
	UUID string `gorm:"column:uuid"`
}

// инфраструктурный уровень
type CryptoInvestmentGorm struct {
	ID                 int                 `gorm:"primaryKey;column:id"`
	UUID               string              `gorm:"column:uuid"`
	CryptoCurrencyID   int                 `gorm:"column:crypto_currency_id"`
	CryptoCurrency     CryptoCurrencyGorm  `gorm:"foreignKey:CryptoCurrencyID"`
	InvestedAmount     int                 `gorm:"column:amount"`
	InvestedCurrencyID int                 `gorm:"column:currency_id"`
	Currency           dto.CurrencyGorm    `gorm:"foreignKey:InvestedCurrencyID"`
	BankAccountID      int                 `gorm:"column:bank_account_id"`
	BankAccount        dto.BankAccountGorm `gorm:"foreignKey:BankAccountID"`
}

type CryptoInvestmentDBFactory struct {
}

func (f *CryptoInvestmentDBFactory) ToEntity(dto CryptoInvestmentGorm) (model.CryptoInvestment, error) {
	id, err := uuid.Parse(dto.UUID)
	if err != nil {
		return model.CryptoInvestment{}, err
	}

	cryptoId, err := uuid.Parse(dto.CryptoCurrency.UUID)
	if err != nil {
		return model.CryptoInvestment{}, err
	}

	currencyId, err := uuid.Parse(dto.Currency.UUID)
	if err != nil {
		return model.CryptoInvestment{}, err
	}

	accountId, err := uuid.Parse(dto.BankAccount.UUID)
	if err != nil {
		return model.CryptoInvestment{}, err
	}

	return model.CryptoInvestment{
		ID:               id,
		CryptoCurrencyID: cryptoId,
		InvestedMoney:    value_objects.NewMoney(dto.InvestedAmount, currencyId),
		BankAccountID:    accountId,
	}, nil
}

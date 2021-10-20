package billing

import "flamingo.me/dingo"

type TransactionLog struct {

}

type DatabaseTransactionLog struct {

}

type CreditCardProcessor struct {

}

type PaypalCreditCardProcessor struct {

}

type BillingModule struct {}

func (module *BillingModule) Configure(injector *dingo.Injector) {
	// Эта команда сообщает Dingo, что всякий раз, когда она видит зависимость от TransactionLog
	// она должна удовлетворить её, используя DatabaseTransactionLog.
	injector.Bind(new(TransactionLog)).To(DatabaseTransactionLog{})

	// По аналогии такая запись сообщает Dingo, что когда используется CreditCardProcessor в зависимости,
	// она должна быть удовлетворена с помощью PaypalCreditCardProcessor.
	injector.Bind(new(CreditCardProcessor)).To(PaypalCreditCardProcessor{})
}

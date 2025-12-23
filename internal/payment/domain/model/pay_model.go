package model

type PaymentMethod int

const (
	PaymentMethodUnspecified PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

var PMname = map[PaymentMethod]string{
	PaymentMethodUnspecified:   "PaymentMethodUnspecified",
	PaymentMethodCard:          "PaymentMethodCard",
	PaymentMethodSBP:           "PaymentMethodSBP",
	PaymentMethodCreditCard:    "PaymentMethodCreditCard",
	PaymentMethodInvestorMoney: "PaymentMethodInvestorMoney",
}

// test for enums

type PayOrder struct {
	OrderUuid     string
	UserUuid      string
	PaymentMethod PaymentMethod
}
type PayOrderResponse struct {
	TransactionUuid string
}

// vo
// aggregate

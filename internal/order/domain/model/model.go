package modelOrder

type Order struct {
	orderUuid       string
	userUuid        string
	partUuids       string
	totalPrice      float64
	transactionUuid *string
	paymentMethod   *string
	status          string
}

type PaymentMethod int

const (
	UNKNOWN = iota
	CARD
	SBP
	CREDIT_CARD
	INVESTOR_MONEY
)

func (p PaymentMethod) String() string {
	switch p {
	case UNKNOWN:
		return "UNKNOWN"
	case CARD:
		return "CARD"
	case SBP:
		return "SBP"
	case CREDIT_CARD:
		return "CREDIT_CARD"
	case INVESTOR_MONEY:
		return "INVESTOR_MONEY"
	default:
		return "UNKNOWN"
	}
}

type CreateOrderRequest struct {
	UserUuid string
	PartUUid []string
}

type GetOrderResponse struct {
	OrderUuid  string
	UserUuid   string
	PartUuid   []string
	TotalPrice float64
}
type GetOrderParams struct {
	OrderUuid string
}

type CancelOrderParams struct {
	OrderUuid string
}

type CancelOrderRes struct {
}

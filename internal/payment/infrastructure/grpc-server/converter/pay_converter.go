package converter

import (
	"gitlab.com/godevs2/micro/internal/payment/domain/model"
	paymentV1 "gitlab.com/godevs2/micro/shared/pkg/proto/payment/v1"
)

func convertPaymentMethod(method paymentV1.PaymentMethod) model.PaymentMethod {
	switch method {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED:
		return model.PaymentMethodUnspecified
	default:
		return model.PaymentMethodUnspecified
	}
}

// PayOrderRequestToModel конвертирует gRPC запрос в доменную модель
func PayOrderRequestToModel(req *paymentV1.PayOrderRequest) *model.PayOrder {
	if req == nil {
		return &model.PayOrder{}
	}
	return &model.PayOrder{
		OrderUuid:     req.OrderUuid,
		UserUuid:      req.UserUuid,
		PaymentMethod: convertPaymentMethod(req.PaymentMethod),
	}
}

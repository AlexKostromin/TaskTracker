/*
package domain

import (

	"context"

	"gitlab.com/godevs2/micro/internal/payment/domain/model"
)

func validateUser(ctx context.Context, user) {

		// check perms

		// can do payment
	}

package domain

import (
	"context"
	"errors"

	"gitlab.com/godevs2/micro/internal/payment/domain/model"
)

func ValidateUser(ctx context.Context, userUUID string) error {
	// check perms
	// can do payment
	if userUUID == "" {
		return errors.New("user UUID is required")
	}
	return nil
}

func ValidatePaymentRequest(order *model.PayOrder) error {
	if order.OrderUuid == "" {
		return errors.New("order UUID is required")
	}
	if order.UserUuid == "" {
		return errors.New("user UUID is required")
	}
	return nil
}
*/
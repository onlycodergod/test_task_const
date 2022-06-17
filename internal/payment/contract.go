package payment

import (
	"context"
)

// PaymentRepository — это интерфейс с 5 методами: CreatePayment, UpdateStatus, GetStatus, GetPayments
// и CancelPayment.
// @property CreatePayment - Этот метод используется для создания нового платежа.
// @property UpdateStatus - Это используется для обновления статуса платежа.
// @property GetStatus - Этот метод используется для получения статуса платежа.
// @property GetPayments - Этот метод используется для получения всех платежей, сделанных
// пользователем.
// @property CancelPayment - Используется для отмены платежа.
type PaymentRepository interface {
	CreatePayment(ctx context.Context, input PaymentInput) (int64, error)
	UpdateStatus(ctx context.Context, input PaymentStatus) (int64, error)
	GetStatus(ctx context.Context, PaymentID int64) (string, error)
	GetPayments(ctx context.Context, input PaymentUser) ([]payment, error)
	CancelPayment(ctx context.Context, PaymentID int64) (int64, error)
}

// PaymentUseCase — это интерфейс с 5 методами: CreatePayment, UpdateStatus, GetStatus, GetPayments и
// CancelPayment.
// @property CreatePayment - Эта функция используется для создания платежа.
// @property {error} UpdateStatus - Это используется для обновления статуса платежа.
// @property GetStatus - Используется для получения статуса платежа.
// @property GetPayments - Это используется для получения всех платежей пользователя.
// @property {error} CancelPayment - Это функция, которая будет использоваться для отмены платежа.
type PaymentUseCase interface {
	CreatePayment(ctx context.Context, input PaymentInput) (int64, error)
	UpdateStatus(ctx context.Context, input PaymentStatus) error
	GetStatus(ctx context.Context, PaymentID int64) (string, error)
	GetPayments(ctx context.Context, input PaymentUser) ([]payment, error)
	CancelPayment(ctx context.Context, PaymentID int64) error
}

package payment

import (
	"context"
	"fmt"
	"sync"
)

// usecase — это структура с полем repo типа PaymentRepository.
// @property {PaymentRepository} repo - Это репозиторий, который будет использоваться для хранения
// платежа.
type usecase struct {
	repo PaymentRepository
}

// > Эта функция создает новый экземпляр структуры usecase и возвращает указатель на нее.
func NewPaymentusecase(repo PaymentRepository) *usecase {
	return &usecase{
		repo: repo,
	}
}

// Эта функция используется для создания нового платежа.
func (u *usecase) CreatePayment(ctx context.Context, input PaymentInput) (int64, error) {
	PaymentID, err := u.repo.CreatePayment(
		ctx,
		input,
	)
	if err != nil {
		wg := &sync.WaitGroup{}

		wg.Add(1)
		go func() {
			_ = u.UpdateStatus(
				ctx,
				PaymentStatus{
					ID:     PaymentID,
					Status: StatusError,
				},
			)

			wg.Done()
		}()
		wg.Wait()

		return 0, err
	}

	return PaymentID, nil
}

// Эта функция используется для обновления статуса платежа.
func (u *usecase) UpdateStatus(ctx context.Context, input PaymentStatus) error {
	ErrorExeption := make(chan error)
	uCtx, cancel := context.WithCancel(ctx)

	go func() {
		status, err := u.repo.GetStatus(
			uCtx,
			input.ID,
		)
		if err != nil {
			return
		}

		if status == StatusSuccess || status == StatusFailure {
			ErrorExeption <- fmt.Errorf("payment-usecase-UpdateStatus-GetStatus, terminal status %s", status)
		}
	}()

	go func() {
		r, err := u.repo.UpdateStatus(
			uCtx,
			input,
		)
		if err != nil {
			ErrorExeption <- err
			return
		}

		if err := checkTerminalStatusRow(r); err != nil {
			ErrorExeption <- fmt.Errorf("payment-usecase-UpdateStatus, %s", err.Error())
		}

		cancel()
	}()

	select {
	case <-uCtx.Done():
		return nil
	case err := <-ErrorExeption:
		return err
	}
}

// Эта функция используется для получения статуса платежа.
func (u *usecase) GetStatus(ctx context.Context, PaymentID int64) (string, error) {
	return u.repo.GetStatus(
		ctx,
		PaymentID,
	)
}

// Эта функция используется для получения всех платежей для пользователя.
func (u *usecase) GetPayments(ctx context.Context, input PaymentUser) ([]payment, error) {
	return u.repo.GetPayments(
		ctx,
		input,
	)
}

// Эта функция используется для отмены платежа.
func (u *usecase) CancelPayment(ctx context.Context, PaymentID int64) error {
	ErrorExeption := make(chan error)
	dCtx, cancel := context.WithCancel(ctx)

	go func() {
		status, err := u.repo.GetStatus(
			dCtx,
			PaymentID,
		)
		if err != nil {
			return
		}

		if status == StatusSuccess || status == StatusFailure {
			ErrorExeption <- fmt.Errorf("payment-usecase-CancelPayment-GetStatus, terminal status %s", status)
		}
	}()

	go func() {
		r, err := u.repo.CancelPayment(
			dCtx,
			PaymentID,
		)
		if err != nil {
			ErrorExeption <- err
			return
		}

		if err := checkTerminalStatusRow(r); err != nil {
			ErrorExeption <- fmt.Errorf("payment-usecase-deletePayment, %s", err.Error())
		}

		cancel()
	}()

	select {
	case <-dCtx.Done():
		return nil
	case err := <-ErrorExeption:
		return err
	}
}

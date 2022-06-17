package payment

import (
	"context"
	"database/sql"
	"fmt"
)

const payments = "payments"

// Репозиторий — это структура с полем db, которое является указателем на sql.DB.
// @property db - Это указатель на объект sql.DB, который будет использоваться для подключения к базе
// данных.
type repository struct {
	db *sql.DB
}

// Он создает новый экземпляр структуры репозитория и возвращает указатель на него.
func NewPaymentRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}

// Создание нового платежа.
func (r *repository) CreatePayment(ctx context.Context, input PaymentInput) (int64, error) {
	const format = `INSERT INTO %s (user_id, user_email, amount, currency)
						VALUES ($1, $2, $3, $4)
					RETURNING id`

	query := fmt.Sprintf(
		format,
		payments,
	)

	row := r.db.QueryRowContext(
		ctx,
		query,
		input.UserID,
		input.UserEmail,
		input.Amount,
		input.Currency,
	)

	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("payment-repository-CreatePayment, %s", err.Error())
	}

	var id int64
	if err := row.Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("payment-repository-CreatePayment, %s", "no result")
		}

		return 0, fmt.Errorf("payment-repository-CreatePayment, %s", err.Error())
	}

	return id, nil
}

// Обновление статуса платежа.
func (r *repository) UpdateStatus(ctx context.Context, input PaymentStatus) (int64, error) {
	const format = `UPDATE %s SET status = $1
						WHERE id = $2
						AND status NOT IN ($3, $4)`

	query := fmt.Sprintf(
		format,
		payments,
	)

	rows, err := r.db.ExecContext(
		ctx,
		query,
		input.Status,
		input.ID,
		StatusSuccess,
		StatusFailure,
	)
	if err != nil {
		return 0, fmt.Errorf("payment-reposiroty-UpdateStatus, %s", err.Error())
	}

	return rows.RowsAffected()
}

// Получение статуса платежа.
func (r *repository) GetStatus(ctx context.Context, PaymentID int64) (string, error) {
	const format = `SELECT status from %s
						WHERE id = $1`

	query := fmt.Sprintf(
		format,
		payments,
	)

	rows := r.db.QueryRowContext(
		ctx,
		query,
		PaymentID,
	)

	var status string
	if err := rows.Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("payment-repository-CreatePayment, %s", "no result")
		}

		return "", fmt.Errorf("payment-reposiroty-GetStatus, %s", err.Error())
	}

	return status, nil
}

// Функция, которая возвращает срез платежей.
func (r *repository) GetPayments(ctx context.Context, input PaymentUser) ([]payment, error) {
	var arg string
	var value interface{}

	if input.UserEmail != "" && input.UserID == 0 {
		arg = "user_email"
		value = input.UserEmail
	}

	if input.UserID != 0 && input.UserEmail == "" {
		arg = "user_id"
		value = input.UserID
	}

	const format = `SELECT
						id,
						user_id,
						user_email,
						currency,
						amount,
						created_at,
						updated_at,
						status
					from %s
						WHERE %s = $1`

	query := fmt.Sprintf(
		format,
		payments,
		arg,
	)

	rows, err := r.db.QueryContext(
		ctx,
		query,
		value,
	)
	if err != nil {
		return []payment{}, fmt.Errorf("payment-reposiroty-GetPayments, %s", err.Error())
	}

	defer rows.Close()

	output := make([]payment, 0)
	for rows.Next() {
		value := payment{}

		err := rows.Scan(
			&value.ID,
			&value.UserID,
			&value.UserEmail,
			&value.Currency,
			&value.Amount,
			&value.CreatedAt,
			&value.UpdatedAt,
			&value.Status,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return []payment{}, fmt.Errorf("payment-repository-CreatePayment, %s", "no result")
			}

			return []payment{}, fmt.Errorf("payment-reposiroty-GetPayments, %s", err.Error())
		}

		output = append(output, value)
	}

	if err := rows.Err(); err != nil {
		return []payment{}, fmt.Errorf("payment-reposiroty-GetPayments, %s", err.Error())
	}

	return output, nil
}

// Обновление статуса "Отмены" платежа.
func (r *repository) CancelPayment(ctx context.Context, PaymentID int64) (int64, error) {
	const format = `UPDATE %s SET status = $1
						WHERE id = $2
						AND status NOT IN ($3, $4)`

	query := fmt.Sprintf(
		format,
		payments,
	)

	rows, err := r.db.ExecContext(
		ctx,
		query,
		StatusCanceled,
		PaymentID,
		StatusSuccess,
		StatusFailure,
	)
	if err != nil {
		return 0, fmt.Errorf("payment-repository-CancelPayment, %s", err.Error())
	}

	return rows.RowsAffected()
}

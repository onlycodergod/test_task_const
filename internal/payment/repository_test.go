package payment

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

// Он создает новое соединение с базой данных и макет для этого соединения.
func TestCreatePayment(t *testing.T) {
	t.Parallel()

	db, dbMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewPaymentRepository(db)

	tests := []struct {
		name   string
		mock   func()
		input  PaymentInput
		expect int64
		err    error
	}{
		{
			name: "Create payment",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)

				dbMock.ExpectQuery("INSERT INTO payments").
					WithArgs(1, "user_email", 10.5, "currency").
					WillReturnRows(rows)
			},
			input: PaymentInput{
				UserID:    1,
				UserEmail: "user_email",
				Amount:    10.5,
				Currency:  "currency",
			},
			expect: 1,
			err:    nil,
		},
		{
			name: "Fail",
			mock: func() {
				dbMock.ExpectQuery("INSERT INTO payments").
					WithArgs(0, "", 0.0, "").
					WillReturnError(errors.New("insert error"))
			},
			input: PaymentInput{
				UserID:    0,
				UserEmail: "",
				Amount:    0.0,
				Currency:  "",
			},
			err: errors.New("insert error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			// Вызов функции CreatePayment с входными данными.
			got, err := r.CreatePayment(
				context.TODO(),
				PaymentInput{
					UserID:    tt.input.UserID,
					UserEmail: tt.input.UserEmail,
					Amount:    tt.input.Amount,
					Currency:  tt.input.Currency,
				})

			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, got)
			}

			assert.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

// Он обновляет статус данного ресурса «Развертывание».
func TestUpdateStatus(t *testing.T) {
	t.Parallel()

	db, dbMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// Он создает фиктивное соединение с базой данных.
	defer db.Close()

	r := NewPaymentRepository(db)

	tests := []struct {
		name   string
		mock   func()
		input  PaymentStatus
		expect int64
		err    error
	}{
		{
			name: "Update status to success",
			mock: func() {
				dbMock.ExpectExec("UPDATE payments").
					WithArgs("success", 1, StatusSuccess, StatusFailure).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: PaymentStatus{
				ID:     1,
				Status: "success",
			},
			expect: 1,
			err:    nil,
		},
		{
			name: "Update status to failure",
			mock: func() {
				dbMock.ExpectExec("UPDATE payments").
					WithArgs("failure", 1, StatusSuccess, StatusFailure).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: PaymentStatus{
				ID:     1,
				Status: "failure",
			},
			expect: 1,
			err:    nil,
		},
		{
			name: "Fail",
			mock: func() {
				dbMock.ExpectExec("UPDATE payments").
					WithArgs("failure", 1, StatusSuccess, StatusFailure).
					WillReturnError(errors.New("update error"))
			},
			input: PaymentStatus{
				ID:     1,
				Status: StatusFailure,
			},
			err: errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.UpdateStatus(
				context.TODO(),
				PaymentStatus{
					ID:     tt.input.ID,
					Status: tt.input.Status,
				})

			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, got)
			}

			assert.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

// Он создает фиктивное соединение с базой данных
func TestGetStatus(t *testing.T) {
	t.Parallel()

	db, dbMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	r := NewPaymentRepository(db)

	tests := []struct {
		name   string
		mock   func()
		input  int64
		expect string
		err    error
	}{
		{
			name: "Get payment status",
			mock: func() {
				row := sqlmock.NewRows([]string{"status"}).AddRow("new")

				dbMock.ExpectQuery("SELECT").
					WithArgs(1).
					WillReturnRows(row)
			},
			input:  1,
			expect: "new",
			err:    nil,
		},
		{
			name: "Fail",
			mock: func() {
				dbMock.ExpectQuery("SELECT").
					WithArgs(1).
					WillReturnError(errors.New("not found"))
			},
			input: 1,
			err:   errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetStatus(
				context.TODO(),
				tt.input,
			)

			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, got)
			}

			assert.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

// Он создает фиктивное соединение с базой данных.
func TestGetPayments(t *testing.T) {
	t.Parallel()

	db, dbMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// Создание структуры с именем теста, фиктивной функцией, вводом, ожидаемым выводом и ошибкой.
	defer db.Close()

	r := NewPaymentRepository(db)

	tests := []struct {
		name   string
		mock   func()
		input  PaymentUser
		expect []payment
		err    error
	}{
		{
			name: "Get user payments by ID",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "user_email", "currency", "amount", "created_at", "updated_at", "status"}).
					AddRow(1, 1, "user_email", "currency", 10.5, "created_at", "updated_at", "status").
					AddRow(2, 2, "user_email", "currency", 10.5, "created_at", "updated_at", "status")

				dbMock.ExpectQuery("SELECT").
					WithArgs(1).
					WillReturnRows(rows)
			},
			input: PaymentUser{
				UserID: 1,
			},
			expect: []payment{
				{1, 1, 10.5, "user_email", "currency", "created_at", "updated_at", "status"},
				{2, 2, 10.5, "user_email", "currency", "created_at", "updated_at", "status"},
			},
			err: nil,
		},
		{
			name: "Get user payments by Email",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "user_email", "currency", "amount", "created_at", "updated_at", "status"}).
					AddRow(1, 1, "user_email", "currency", 10.5, "created_at", "updated_at", "status").
					AddRow(2, 2, "user_email", "currency", 10.5, "created_at", "updated_at", "status")

				dbMock.ExpectQuery("SELECT").
					WithArgs("email").
					WillReturnRows(rows)
			},
			input: PaymentUser{
				UserEmail: "email",
			},
			expect: []payment{
				{1, 1, 10.5, "user_email", "currency", "created_at", "updated_at", "status"},
				{2, 2, 10.5, "user_email", "currency", "created_at", "updated_at", "status"},
			},
			err: nil,
		},
		{
			name: "Fail",
			mock: func() {
				dbMock.ExpectQuery("SELECT").
					WithArgs(1).
					WillReturnError(errors.New("not found"))
			},
			input: PaymentUser{
				UserID: 1,
			},
			err: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.GetPayments(
				context.TODO(),
				PaymentUser{
					UserID:    tt.input.UserID,
					UserEmail: tt.input.UserEmail,
				},
			)

			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, got)
			}

			assert.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

// Он создает фиктивное соединение с базой данных
func TestCancelPayment(t *testing.T) {
	t.Parallel()

	db, dbMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// Создание фейкового подключения к базе данных.
	defer db.Close()

	r := NewPaymentRepository(db)

	tests := []struct {
		name   string
		mock   func()
		input  int64
		expect int64
		err    error
	}{
		{
			name: "Cancel payment",
			mock: func() {
				dbMock.ExpectExec("UPDATE payments").
					WithArgs("canceled", 1, StatusSuccess, StatusFailure).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input:  1,
			expect: 1,
			err:    nil,
		},
		{
			name: "Fail",
			mock: func() {
				dbMock.ExpectExec("UPDATE payments").
					WithArgs("canceled", 1, StatusSuccess, StatusFailure).
					WillReturnError(errors.New("update error"))
			},
			input: 1,
			err:   errors.New("update error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.CancelPayment(
				context.TODO(),
				tt.input,
			)

			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, got)
			}

			assert.NoError(t, dbMock.ExpectationsWereMet())
		})
	}
}

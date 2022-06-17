package payment

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/onlycodergod/payment-api-emulator/pkg/loggin"
)

// > Тип контроллера — это структура с интерфейсом PaymentUseCase и интерфейсом регистратора.
// @property {PaymentUseCase} UseCase - Это интерфейс, который определяет методы, которые контроллер
// будет использовать для взаимодействия с вариантом использования.
// @property logger - Это регистратор, который будет использоваться для регистрации ошибок и другой
// информации.
type controller struct {
	UseCase PaymentUseCase
	logger  loggin.ILogger
}

// > Эта функция создает новый экземпляр структуры контроллера и возвращает указатель на него
func NewPaymentController(l loggin.ILogger, u PaymentUseCase) *controller {
	return &controller{
		logger:  l,
		UseCase: u,
	}
}

// Это константа, определяющая маршрут.
const (
	CreatePayment          = "/payment"
	UpdateStatusByID       = "/payments/{id}/status"
	GetStatusByID          = "/payments/{id}/status"
	GetPaymentsByUserEmail = "/payments/user" // query /payments/user?email=email
	GetPaymentsByUserID    = "/payments/user/{id}"
	CancelPaymentByID      = "/payments/{id}"
)

// Эта функция представляет собой обработчик, который будет вызываться при запросе маршрута.
// // метод `/платежа` `POST`.
func (c *controller) Register(router *mux.Router) *mux.Router {
	router.HandleFunc(CreatePayment, c.CreatePayment).Methods(http.MethodPost)
	router.HandleFunc(UpdateStatusByID, c.UpdateStatus).Methods(http.MethodPut)
	router.HandleFunc(GetStatusByID, c.GetStatus).Methods(http.MethodGet)
	router.HandleFunc(GetPaymentsByUserEmail, c.GetPaymentsByUserEmail).Methods(http.MethodGet)
	router.HandleFunc(GetPaymentsByUserID, c.GetPaymentsByUserID).Methods(http.MethodGet)
	router.HandleFunc(CancelPaymentByID, c.CancelPayment).Methods(http.MethodPut)
	return router
}

// Эта функция представляет собой обработчик, который будет вызываться при выполнении запроса к
// маршруту.
// // `/платеж` методом `POST`.
func (c *controller) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var input PaymentInput

	// Это проверка правильности данных в теле запроса.
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, InvalidBodyData, http.StatusBadRequest)
		return
	}

	// Это проверка правильности адреса электронной почты в теле запроса.
	if ok := isEmail(input.UserEmail); !ok {
		http.Error(w, InvalidBodyEmail, http.StatusBadRequest)
		return
	}

	// Это вызов метода варианта использования, который создает платеж.
	id, err := c.UseCase.CreatePayment(
		r.Context(),
		input,
	)
	if err != nil {
		c.logger.Error(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	// Это ответ клиенту.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(
		PaymentStatus{
			ID: id,
		},
	)
}

// Эта функция представляет собой обработчик, который будет вызываться при выполнении запроса к
// маршруту.
// // `/payments/{id}/status` методом `PUT`.
func (c *controller) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	PaymentID, err := GetQueryId(r)
	if err != nil {
		http.Error(w, InvalidQueryID, http.StatusBadRequest)
		return
	}

	var input PaymentStatus
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, InvalidBodyData, http.StatusBadRequest)
		return
	}

	input.ID = PaymentID

	err = c.UseCase.UpdateStatus(
		r.Context(),
		input,
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Эта функция представляет собой обработчик, который будет вызываться при выполнении запроса к
// маршруту.
// // `/payments/{id}/status` методом `GET`.
func (c *controller) GetStatus(w http.ResponseWriter, r *http.Request) {
	PaymentID, err := GetQueryId(r)
	if err != nil {
		http.Error(w, InvalidQueryID, http.StatusBadRequest)
		return
	}

	status, err := c.UseCase.GetStatus(
		r.Context(),
		PaymentID,
	)
	if err != nil {
		c.logger.Error(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		PaymentStatus{
			Status: status,
		},
	)
}

// Эта функция представляет собой обработчик, который будет вызываться при выполнении запроса к
// маршруту.
// // `/payments/user` методом `GET`.
func (c *controller) GetPaymentsByUserEmail(w http.ResponseWriter, r *http.Request) {
	UserEmail := r.URL.Query().Get("email")
	if ok := isEmail(UserEmail); !ok {
		http.Error(w, InvalidQueryEmail, http.StatusBadRequest)
		return
	}

	data, err := c.UseCase.GetPayments(
		r.Context(),
		PaymentUser{
			UserEmail: UserEmail,
		},
	)
	if err != nil {
		c.logger.Error(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		PaymentsData{
			Data: data,
		},
	)
}

// Эта функция представляет собой обработчик, который будет вызываться при выполнении запроса к
// маршруту.
// // `/payments/user/{id}` методом `GET`.
func (c *controller) GetPaymentsByUserID(w http.ResponseWriter, r *http.Request) {
	userID, err := GetQueryId(r)
	if err != nil {
		http.Error(w, InvalidQueryID, http.StatusBadRequest)
		return
	}

	data, err := c.UseCase.GetPayments(
		r.Context(),
		PaymentUser{
			UserID: userID,
		},
	)
	if err != nil {
		c.logger.Error(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		PaymentsData{
			Data: data,
		},
	)
}

// Обработчик, который будет вызываться при запросе маршрута
// // `/payments/{id}` методом `PUT`.
func (c *controller) CancelPayment(w http.ResponseWriter, r *http.Request) {
	PaymentID, err := GetQueryId(r)
	if err != nil {
		http.Error(w, InvalidQueryID, http.StatusBadRequest)
		return
	}

	err = c.UseCase.CancelPayment(
		r.Context(),
		PaymentID,
	)

	if err != nil {
		c.logger.Error(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

Техническое задание: создать эмулятор платежного сервиса.

Стиль написания
```
CamelCase
```
 Что нужно для запуска проекта:
    1. Docker
    2. Docker-compose (version by - support 3.9 verion file)

 Запуск проекта:
 ```sh
    1. sudo make compose-up
```

 API 
   1. "/payment", Method: POST - создает транзакцию, request body params: {"user_id": type int, "amount": type decimal, "user_email": type varchar, "currency": type varchar}
```go
func (c *controller) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var input PaymentInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, InvalidBodyData, http.StatusBadRequest)
		return
	}

	if ok := isEmail(input.UserEmail); !ok {
		http.Error(w, InvalidBodyEmail, http.StatusBadRequest)
		return
	}

	id, err := c.UseCase.CreatePayment(
		r.Context(),
		input,
	)
	if err != nil {
		c.logger.Error(err)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}
```
```go
   2. "/payments/{id}/status", Method: PUT - обновляет статус транзакции по ее id
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
```
    3. "/payments/{id}/status", Method: GET - возвращает статус транзакции по ее id
```go
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
```
    4. "/payments/user/{id}", Method: GET - возвращает транзакции пользователя по его id
```go
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
```
    5. "/payments/user?email=...", Method: GET - возвращает транзакции пользователя по его email
```go
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
```
    6. "/payments/{id}", Method: PUT - отменяет транзакцию транзакцию по ее id
```go
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
```

 Архитектура проекта из соображений:
 ```git
    1. https://github.com/golang-standards/project-layout
    2. https://github.com/moby/moby
    3. https://github.com/digitalocean
 ```
 
 Дизайн архитектуры и работы транзакций
```
 <a href="https://ibb.co/w7PHgQ4" target="_blank"><img src="https://i.ibb.co/xYBvLg6/photo-diagram-project.jpg" target="_blank"></a> 
```


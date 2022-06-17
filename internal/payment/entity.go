package payment

// Это структура, содержащая поля, используемые для представления платежа.
//
// Поля снабжены тегами, которые используются пакетом sqlx для сопоставления полей со столбцами в базе
// данных.
//
// Поля также снабжены тегами, которые используются пакетом json для сопоставления полей с ключами в
// объекте JSON.
//
// Поля также снабжены тегами, которые используются пакетом go-swagger для сопоставления полей с
// ключами в объекте JSON.
//
// Поля также аннотируются тегами, которые используются в go-swagger.
// @property {int64} ID - Уникальный идентификатор платежа.
// @property {int64} UserID - ID пользователя, совершившего платеж.
// @property {float64} Amount - Сумма платежа.
// @property {string} UserEmail - Электронная почта пользователя, совершившего платеж
// @property {string} Currency - Валюта платежа.
// @property {string} CreatedAt - Дата и время создания платежа.
// @property {string} UpdatedAt - Дата и время последнего обновления платежа.
// @property {string} Status - Статус платежа. Он может быть «ожидающим», «завершенным» или
// «неудачным».
type payment struct {
	ID        int64   `json:"id" db:"id"`
	UserID    int64   `json:"user_id" db:"user_id"`
	Amount    float64 `json:"amount" db:"amount"`
	UserEmail string  `json:"user_email" db:"user_email"`
	Currency  string  `json:"currency" db:"currency"`
	CreatedAt string  `json:"created_at" db:"created_at"`
	UpdatedAt string  `json:"updated_at" db:"updated_at"`
	Status    string  `json:"status" db:"status"`
}

// «PaymentInput» — это структура с четырьмя полями: «UserID», «Amount», «UserEmail» и «Currency».
//
// Первое поле, `UserID`, представляет собой `int64` (64-битное целое число). Второе поле, «Сумма»,
// представляет собой «float64» (64-битное число с плавающей запятой). Третье поле, `UserEmail`,
// представляет собой `строку` (строку символов). Четвертое поле «Валюта» также является строкой.
//
// Теги `json` в каждом поле сообщают компилятору Go
// @property {int64} UserID - Идентификатор пользователя, который осуществляет платеж.
// @property {float64} Amount - Сумма к оплате.
// @property {string} UserEmail - Адрес электронной почты пользователя, который осуществляет платеж.
// @property {string} Currency - Валюта платежа.
type PaymentInput struct {
	UserID    int64   `json:"user_id"`
	Amount    float64 `json:"amount"`
	UserEmail string  `json:"user_email"`
	Currency  string  `json:"currency"`
}

// PaymentUser — это структура, содержащая идентификатор пользователя и адрес электронной почты.
// @property {int64} UserID - ID пользователя в вашей системе.
// @property {string} UserEmail - Электронный адрес пользователя.
type PaymentUser struct {
	UserID    int64  `json:"user_id"`
	UserEmail string `json:"user_email"`
}

// PaymentsData — это структура, содержащая фрагмент платежных структур.
// @property {[]payment} Data - Это массив платежей, которые мы будем возвращать.
type PaymentsData struct {
	Data []payment `json:"data"`
}

// PaymentStatus — это структура, содержащая идентификатор и статус.
// @property {int64} ID - Идентификатор платежа.
// @property {string} Status - Статус платежа. Возможные значения:
type PaymentStatus struct {
	ID     int64  `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
}

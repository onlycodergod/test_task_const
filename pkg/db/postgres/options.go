package postgres

import "fmt"

// DBOptions — это структура с 6 полями, все из которых являются строками.
// @property {string} User - Имя пользователя для подключения к базе данных.
// @property {string} Password - Пароль для пользователя.
// @property {string} Host - Имя хоста сервера базы данных.
// @property {string} Port - Номер порта сервера базы данных.
// @property {string} DB - Имя базы данных для подключения.
// @property {string} SSLmode - Это режим SSL для использования. По умолчанию «требуется».
type DBOptions struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       string
	SSLmode  string
}

// Он принимает структуру DBOptions и возвращает строку, которую можно использовать для подключения к
// базе данных Postgres.
func getDSN(options DBOptions) string {
	const format = "postgres://%s:%s@%s:%s/%s?sslmode=%s"

	dsn := fmt.Sprintf(
		format,
		options.User,
		options.Password,
		options.Host,
		options.Port,
		options.DB,
		options.SSLmode,
	)

	return dsn
}

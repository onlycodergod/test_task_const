package postgres

import "fmt"

type DBOptions struct {
	User     string
	Password string
	Host     string
	Port     string
	DB       string
	SSLmode  string
}

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

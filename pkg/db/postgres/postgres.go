package postgres

import (
	"database/sql"
)

// `postgres` — это тип, который имеет поле, называемое `options` типа `DBOptions`.
// @property {DBOptions} options - Это параметры, которые передаются в базу данных.
type postgres struct {
	options DBOptions
}

// Он создает новый объект postgres с переданными параметрами.
func NewPostgres(options DBOptions) *postgres {
	return &postgres{
		options: options,
	}
}

// Метод коннекта к нашей database postgres.
func (p *postgres) Connect() (*sql.DB, error) {
	dsn := getDSN(p.options)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}

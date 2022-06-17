package postgres

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/onlycodergod/payment-api-emulator/pkg/loggin"
)

// Он создает новый экземпляр миграции, а затем запускает миграцию
func InitMigrate(logger loggin.ILogger, options DBOptions) {
	dsn := getDSN(options)

	if len(dsn) == 0 {
		logger.Fatal("migrate: environment variable not declared")
	}

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		logger.Fatal(err)
	}

	if err := m.Up(); err != nil {
		logger.Debug(err)
	}

	m.Close()
}

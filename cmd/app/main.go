package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/onlycodergod/payment-api-emulator/config"
	"github.com/onlycodergod/payment-api-emulator/internal/payment"
	"github.com/onlycodergod/payment-api-emulator/pkg/db/postgres"
	"github.com/onlycodergod/payment-api-emulator/pkg/http/server"
	"github.com/onlycodergod/payment-api-emulator/pkg/loggin"
)

// Он создает новый объект конфигурации, и в случае сбоя он регистрирует ошибку и выходит из программы.
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config initialization error: %s", err.Error())
	}

	// Logger
	logger := loggin.NewLogger(cfg.Logger.Debug)

	// Database
	dbOptions := postgres.DBOptions{
		User:     cfg.Postgres.User,
		Password: cfg.Postgres.Password,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		DB:       cfg.Postgres.DB,
		SSLmode:  cfg.Postgres.SSLMode,
	}

	pg, err := postgres.NewPostgres(dbOptions).Connect()
	if err != nil {
		logger.Fatalf("postgres connection failed, %s", err.Error())
	}
	defer pg.Close()

	// Миграция базы данных.
	postgres.InitMigrate(
		logger,
		dbOptions,
	)

	// Создание нового репозитория платежей, варианта использования и контроллера.
	rep := payment.NewPaymentRepository(pg)
	usc := payment.NewPaymentUseCase(rep)
	con := payment.NewPaymentController(
		logger,
		usc,
	)

	// Создание нового маршрутизатора и http-сервера.
	router := mux.NewRouter()

	httpServer := server.NewHttpServer(
		con.Register(router),
		cfg.HTTP.Port,
		time.Duration(cfg.HTTP.ReadTimeout),
		time.Duration(cfg.HTTP.WriteTimeout),
		time.Duration(cfg.HTTP.ShutdownTimeout),
	)

	logger.Infof("http server created and started at http://localhost:%s", cfg.HTTP.Port)

	// Обработчик сигнала.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Infof("app - run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		logger.Errorf("app - run - httpServer.Notify: %s", err.Error())
	}

	// Грамотное завершение работы сервера. (Shutdown)
	err = httpServer.Shutdown()
	if err != nil {
		logger.Errorf("app - run - httpServer.Shutdown: %s", err.Error())
	}
}

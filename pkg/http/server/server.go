package server

import (
	"context"
	"net/http"
	"time"
)

// Это структура, содержащая указатель на http.Server, канал ошибок и time.Duration.
// @property server - Это фактический HTTP-сервер, который будет прослушивать запросы.
// @property notify - Это канал, который будет использоваться для уведомления основной горутины об
// остановке сервера.
// @property shutdownTimeout - Время ожидания завершения работы сервера перед возвратом ошибки.
type server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// > Эта функция создает новый HTTP-сервер с заданным обработчиком, портом, тайм-аутом чтения,
// тайм-аутом записи и временем завершения работы.
func NewHttpServer(handler http.Handler, port string, readTimeout, writeTimeout, shutdownTime time.Duration) *server {
	httpServer := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  readTimeout * time.Second,
		WriteTimeout: writeTimeout * time.Second,
	}

	serv := &server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: shutdownTime * time.Second,
	}

	serv.start()

	return serv
}

// Он запускает сервер в горутине.
func (s *server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Возврат канала только для чтения.
func (s *server) Notify() <-chan error {
	return s.notify
}

// Выключение сервера.
func (s *server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

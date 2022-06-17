package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// «Logger» — это структура с одним полем «Debug», которое является логическим значением.
// @property {bool} Debug - Если true, регистратор будет печатать отладочные сообщения.
type Logger struct {
	Debug bool `yaml:"debug"`
}

// HTTP — это структура, содержащая строку, int64, int64 и int64.
// @property {string} Port - Порт, на котором будет прослушиваться HTTP-сервер.
// @property {int64} WriteTimeout - Максимальная продолжительность до истечения времени ожидания записи
// ответа.
// @property {int64} ReadTimeout - Максимальная продолжительность чтения всего запроса, включая тело.
// @property {int64} ShutdownTimeout - Время ожидания выключения сервера перед его уничтожением.
type HTTP struct {
	Port            string `yaml:"port" env:"HTTP_PORT" env-required:"true"`
	WriteTimeout    int64  `yaml:"writeTimeout" env:"HTTP_WRITE_TIMEOUT" env-required:"true"`
	ReadTimeout     int64  `yaml:"readTimeout" env:"HTTP_READ_TIMEOUT" env-required:"true"`
	ShutdownTimeout int64  `yaml:"shutdownTimeout" env:"HTTP_SHUT_DOWN_TIMEOUT" env-required:"true"`
}

// Это структура с полями, которые являются строками, и каждое поле имеет тег, который сообщает пакету
// env, как заполнять поле.
//
// Пакет env заполнит поля типа Postgres значениями из среды. Пакет env будет искать переменные среды,
// соответствующие значению тега. Например, поле «Пользователь» будет заполнено значением переменной
// среды POSTGRES_USER.
//
// Пакет env заполнит поля типа Postgres значениями из среды. Пакет env будет искать переменные среды,
// соответствующие значению тега. Например,
// @property {string} User - Имя пользователя, используемое при подключении к базе данных.
// @property {string} Password - Пароль пользователя Postgres.
// @property {string} Host - Имя хоста сервера Postgres.
// @property {string} Port - Порт, на котором работает Postgres.
// @property {string} DB - Имя базы данных для подключения.
// @property {string} SSLMode - Это режим SSL для использования. По умолчанию установлено значение
// «требовать», что означает, что соединение не будет установлено, если SSL не используется.
type Postgres struct {
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
	DB       string `env:"POSTGRES_DB"`
	SSLMode  string `env:"POSTGRES_SSLMODE"`
}

// «Config» — это структура, которая содержит структуру «Logger», структуру «HTTP» и структуру
// «Postgres».
// @property {Logger}  - Регистратор: это конфигурация регистратора.
// @property {HTTP}  - Регистратор: это конфигурация регистратора.
// @property {Postgres}  - Регистратор: это конфигурация регистратора.
type Config struct {
	Logger `yaml:"logger"`
	HTTP   `yaml:"http"`
	Postgres
}

// Он читает файл конфигурации, затем читает переменные среды, а затем возвращает конфигурацию.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("config/development.yml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

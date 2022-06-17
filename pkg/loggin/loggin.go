package loggin

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ILogger — это интерфейс, определяющий методы ведения журнала.
// @property Debugf - Отладка с форматом
// @property Debug - Записывает сообщение на уровне отладки.
// @property Infof - Печатает информационное сообщение
// @property Info - Это наиболее распространенный уровень журнала. Он используется для регистрации
// информации о нормальной работе приложения.
// @property Warnf - Регистрирует сообщение уровня Warn в стандартном регистраторе.
// @property Warn - Регистрирует сообщение уровня Warn в стандартном регистраторе.
// @property Errorf - Записывает сообщение на уровне ошибки.
// @property Error - Записывает сообщение на уровне ошибки.
// @property Fatalf - Fatal эквивалентен l.Critical(fmt.Sprint()), за которым следует вызов os.Exit(1).
// @property Fatal - Fatal регистрирует сообщение уровня Fatal в стандартном регистраторе.
type ILogger interface {
	Debugf(message string, args ...interface{})
	Debug(args ...interface{})
	Infof(message string, args ...interface{})
	Info(args ...interface{})
	Warnf(message string, args ...interface{})
	Warn(args ...interface{})
	Errorf(message string, args ...interface{})
	Error(args ...interface{})
	Fatalf(message string, args ...interface{})
	Fatal(args ...interface{})
}

// Тип регистратора — это структура, которая имеет одно поле с именем loggin типа ILogger.
// @property {ILogger} loggin - Это интерфейс, который мы будем использовать для регистрации сообщений.
type logger struct {
	loggin ILogger
}

// `NewLogger` возвращает указатель на структуру `logger`, которая содержит структуру `logging`
func NewLogger(debug bool) *logger {
	return &logger{
		loggin: InitZap(debug),
	}
}

// Он создает новый каталог с именем logs.
func InitZap(debug bool) ILogger {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "time"
	config.MessageKey = "message"
	config.EncodeTime = zapcore.TimeEncoderOfLayout("02-01-2006,15:04:05")
	config.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConsole := zapcore.NewConsoleEncoder(config)

	err := os.MkdirAll("logs", 0o755)
	if err != nil {
		log.Fatalln(err.Error())
	}

	var defaultLogLevel zapcore.Level
	if debug {
		defaultLogLevel = zapcore.DebugLevel
	} else {
		defaultLogLevel = zapcore.InfoLevel
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoderConsole, zapcore.AddSync(os.Stdout), defaultLogLevel),
	)

	return zap.New(core, zap.AddCallerSkip(1), zap.AddCaller()).Sugar()
}

// Метод, определенный в структуре logger. Он принимает сообщение и переменную
// количество аргументов. Затем он вызывает метод `Debugf` для поля `loggin` структуры `logger`.
func (l *logger) Debugf(message string, args ...interface{}) {
	l.loggin.Debugf(message, args...)
}

// Метод, определенный в структуре logger. Он принимает переменное количество аргументов.
// Затем он вызывает метод `Debug` для поля `loggin` структуры `logger`.
func (l *logger) Debug(args ...interface{}) {
	l.loggin.Debug(args...)
}

// Метод, определенный в структуре logger. Он принимает сообщение и переменную
// количество аргументов. Затем он вызывает метод `Infof` для поля `loggin` структуры `logger`.
func (l *logger) Infof(message string, args ...interface{}) {
	l.loggin.Infof(message, args...)
}

// Метод, определенный в структуре logger. Он принимает переменное количество аргументов.
// Затем он вызывает метод `Info` для поля `loggin` структуры `logger`.
func (l *logger) Info(args ...interface{}) {
	l.loggin.Info(args...)
}

// Метод, определенный в структуре logger. Он принимает сообщение и переменную
// количество аргументов. Затем он вызывает метод `Warnf` для поля `loggin` структуры `logger`.
func (l *logger) Warnf(message string, args ...interface{}) {
	l.loggin.Warnf(message, args...)
}

// Метод, определенный в структуре logger. Он принимает переменное количество аргументов.
// Затем он вызывает метод «Warn» для поля «loggin» структуры «logger».
func (l *logger) Warn(args ...interface{}) {
	l.loggin.Warn(args...)
}

// Метод, определенный в структуре logger. Он принимает сообщение и переменную
// количество аргументов. Затем он вызывает метод Errorf для поля loggin структуры logger.
func (l *logger) Errorf(message string, args ...interface{}) {
	l.loggin.Errorf(message, args...)
}

// Метод, определенный в структуре logger. Он принимает переменное количество аргументов.
// Затем он вызывает метод `Error` для поля `loggin` структуры `logger`.
func (l *logger) Error(args ...interface{}) {
	l.loggin.Error(args...)
}

// Метод, определенный в структуре logger. Он принимает переменное количество аргументов.
// Затем он вызывает метод Fatalf для поля loggin структуры logger.
func (l *logger) Fatalf(message string, args ...interface{}) {
	l.loggin.Fatalf(message, args...)
}

// Метод, определенный в структуре регистратора. Он принимает переменное количество аргументов. Затем
// он вызывает метод Fatal для поля loggin структуры регистратора.
func (l *logger) Fatal(args ...interface{}) {
	l.loggin.Fatal(args...)
}

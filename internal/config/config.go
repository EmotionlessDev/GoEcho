// Package config собирает конфигурацию платформенных сервисов из переменных
package config

import (
	"fmt"
	"log/slog"
	"time"
)

// Окружения, влияющие на формат логов и подробность ошибок.
const (
	EnvDev  = "dev"
	EnvProd = "prod"
)

// Форматы вывода логов.
const (
	LogFormatJSON = "json"
	LogFormatText = "text"
)

// App — общая часть конфигурации, одинаковая у всех бинарей платформы.
type App struct {
	// Env влияет на дефолтный формат логов и на то, отдаём ли мы наружу детали ошибок.
	Env string
	// ServiceName попадает в каждую строку лога — в compose их несколько.
	ServiceName string
	LogLevel    slog.Level
	LogFormat   string
	// ShutdownTimeout ограничивает graceful shutdown: сколько ждём завершения
	// начатых запросов, прежде чем закрыть соединения принудительно.
	ShutdownTimeout time.Duration
}

func (a App) IsDev() bool { return a.Env == EnvDev }

type Postgres struct {
	DSN             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	// ConnectTimeout ограничивает установку соединения: без него старт сервиса
	// с недоступной базой висит до таймаута ядра.
	ConnectTimeout time.Duration
}

// HTTP — параметры http.Server. Таймауты заданы явно: сервер без них
// держит соединение бесконечно и утекает под медленным клиентом.
type HTTP struct {
	Addr              string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

type JWT struct {
	Secret string
	TTL    time.Duration
}

type GoEchoAPI struct {
	App
	Postgres Postgres
	HTTP     HTTP
	JWT      JWT
}

func (l *loader) app(service string) App {
	env := l.oneOf("APP_ENV", EnvDev, EnvDev, EnvProd)

	// В dev человекочитаемый текст, в prod — JSON для сборщика логов.
	defaultFormat := LogFormatJSON
	if env == EnvDev {
		defaultFormat = LogFormatText
	}

	return App{
		Env:             env,
		ServiceName:     service,
		LogLevel:        l.logLevel("LOG_LEVEL", slog.LevelInfo),
		LogFormat:       l.oneOf("LOG_FORMAT", defaultFormat, LogFormatJSON, LogFormatText),
		ShutdownTimeout: l.duration("SHUTDOWN_TIMEOUT", 15*time.Second),
	}
}

func (l *loader) postgres() Postgres {
	return Postgres{
		// Дефолт указывает на порт стенда с хоста (docker-compose.yml),
		// чтобы `go run ./cmd/...` работал без единой заданной переменной.
		DSN:             l.str("POSTGRES_DSN", "postgres://goecho:goecho@localhost:55432/goecho?sslmode=disable"),
		MaxConns:        l.int32("POSTGRES_MAX_CONNS", 10),
		MinConns:        l.int32("POSTGRES_MIN_CONNS", 2),
		MaxConnLifetime: l.duration("POSTGRES_MAX_CONN_LIFETIME", time.Hour),
		MaxConnIdleTime: l.duration("POSTGRES_MAX_CONN_IDLE_TIME", 30*time.Minute),
		ConnectTimeout:  l.duration("POSTGRES_CONNECT_TIMEOUT", 5*time.Second),
	}
}

func (l *loader) http(defaultAddr string) HTTP {
	return HTTP{
		Addr:              l.str("HTTP_ADDR", defaultAddr),
		ReadHeaderTimeout: l.duration("HTTP_READ_HEADER_TIMEOUT", 5*time.Second),
		ReadTimeout:       l.duration("HTTP_READ_TIMEOUT", 15*time.Second),
		WriteTimeout:      l.duration("HTTP_WRITE_TIMEOUT", 30*time.Second),
		IdleTimeout:       l.duration("HTTP_IDLE_TIMEOUT", 60*time.Second),
	}
}

func (l *loader) jwt() JWT {
	return JWT{
		Secret: l.str("JWT_SECRET", ""),
		TTL:    l.duration("JWT_TTL", 15*time.Minute),
	}
}

func LoadGoEchoAPI() (GoEchoAPI, error) {
	l := &loader{}
	cfg := GoEchoAPI{
		App:      l.app("goecho-api"),
		HTTP:     l.http(":8080"),
		Postgres: l.postgres(),
		JWT:      l.jwt(),
	}

	if err := l.err(); err != nil {
		return GoEchoAPI{}, fmt.Errorf("конфигурация goecho-api: %w", err)
	}

	return cfg, nil
}

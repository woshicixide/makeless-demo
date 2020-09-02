package main

import (
	"github.com/go-saas/go-saas/mailer"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-saas/go-saas"
	"github.com/go-saas/go-saas/authenticator/basic"
	"github.com/go-saas/go-saas/config/basic"
	"github.com/go-saas/go-saas/database/basic"
	"github.com/go-saas/go-saas/event/basic"
	"github.com/go-saas/go-saas/http"
	"github.com/go-saas/go-saas/http/basic"
	"github.com/go-saas/go-saas/logger/basic"
	"github.com/go-saas/go-saas/mailer/basic"
	"github.com/go-saas/go-saas/security/basic"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// logger
	logger := new(go_saas_logger_basic.Logger)

	// config
	config := &go_saas_config_basic.Config{
		RWMutex: new(sync.RWMutex),
	}

	// mailer
	mailer := &go_saas_mailer_basic.Mailer{
		Handlers: make(map[string]func(data map[string]interface{}) (go_saas_mailer.Mail, error)),
		Host:     os.Getenv("MAILER_HOST"),
		Port:     os.Getenv("MAILER_PORT"),
		Identity: os.Getenv("MAILER_IDENTITY"),
		Username: os.Getenv("MAILER_USERNAME"),
		Password: os.Getenv("MAILER_PASSWORD"),
		RWMutex:  new(sync.RWMutex),
	}

	// database
	database := &go_saas_database_basic.Database{
		Dialect:  "mysql",
		Host:     os.Getenv("DB_HOST"),
		Database: os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		RWMutex:  new(sync.RWMutex),
	}

	// security
	security := &go_saas_security_basic.Security{
		Database: database,
		RWMutex:  new(sync.RWMutex),
	}

	// event hub
	hub := &go_saas_event_basic.Hub{
		List:    new(sync.Map),
		RWMutex: new(sync.RWMutex),
	}

	// event
	event := &go_saas_event_basic.Event{
		Hub:     hub,
		Error:   make(chan error),
		RWMutex: new(sync.RWMutex),
	}

	// jwt authenticator
	authenticator := &go_saas_authenticator_basic.Authenticator{
		Security:    security,
		Realm:       "auth",
		Key:         os.Getenv("JWT_KEY"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "id",
		RWMutex:     new(sync.RWMutex),
	}

	// http
	http := &go_saas_http_basic.Http{
		Router:        gin.Default(),
		Handlers:      make(map[string]func(http go_saas_http.Http) error),
		Logger:        logger,
		Event:         event,
		Authenticator: authenticator,
		Security:      security,
		Database:      database,
		Mailer:        mailer,
		Tls:           nil,
		Origins:       strings.Split(os.Getenv("ORIGINS"), ","),
		Headers:       []string{"Team"},
		Port:          os.Getenv("API_PORT"),
		Mode:          os.Getenv("API_MODE"),
		RWMutex:       new(sync.RWMutex),
	}

	saas := &go_saas.Saas{
		Config:   config,
		Logger:   logger,
		Mailer:   mailer,
		Database: database,
		Http:     http,
		RWMutex:  new(sync.RWMutex),
	}

	if err := saas.Init("./go-saas.json"); err != nil {
		saas.GetLogger().Fatal(err)
	}

	if err := saas.Run(); err != nil {
		saas.GetLogger().Fatal(err)
	}
}

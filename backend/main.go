package main

import (
	"os"
	"strings"
	"sync"

	"github.com/go-saas/go-saas"
	"github.com/go-saas/go-saas/api"
	"github.com/go-saas/go-saas/database/basic"
	"github.com/go-saas/go-saas/event/basic"
	"github.com/go-saas/go-saas/jwt/basic"
	"github.com/go-saas/go-saas/logger/basic"
	"github.com/go-saas/go-saas/security/basic"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// logger
	logger := new(go_saas_basic_logger.Logger)

	// database
	database := &go_saas_basic_database.Database{
		Dialect:  "mysql",
		Host:     os.Getenv("DB_HOST"),
		Database: os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		RWMutex:  new(sync.RWMutex),
	}

	// security
	security := &go_saas_basic_security.Security{
		Database: database,
		RWMutex:  new(sync.RWMutex),
	}

	// jwt
	jwt := &go_saas_basic_jwt.Jwt{
		Key:     os.Getenv("JWT_KEY"),
		RWMutex: new(sync.RWMutex),
	}

	// event
	event := &go_saas_basic_event.Event{
		Hub:     new(go_saas_basic_event.Hub).Init(),
		RWMutex: new(sync.RWMutex),
	}

	// api
	api := &saas_api.Api{
		Logger:   logger,
		Event:    event,
		Security: security,
		Database: database,
		Jwt:      jwt,
		Tls:      nil,
		Origins:  strings.Split(os.Getenv("ORIGINS"), ","),
		Port:     os.Getenv("API_PORT"),
		Mode:     os.Getenv("API_MODE"),
		RWMutex:  new(sync.RWMutex),
	}

	saas := &go_saas.Saas{
		Logger:   logger,
		Database: database,
		Api:      api,
		RWMutex:  new(sync.RWMutex),
	}

	if err := saas.Run(); err != nil {
		saas.GetLogger().Fatal(err)
	}
}

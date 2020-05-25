package main

import (
	"os"
	"strings"
	"sync"

	go_saas "github.com/go-saas/go-saas"
	saas_api "github.com/go-saas/go-saas/api"
	saas_database "github.com/go-saas/go-saas/database"
	saas_event_basic "github.com/go-saas/go-saas/event/basic"
	saas_logger_stdio "github.com/go-saas/go-saas/logger/stdio"
	saas_security_basic "github.com/go-saas/go-saas/security/basic"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// logger
	logger := new(saas_logger_stdio.Stdio)

	// database
	database := &saas_database.Database{
		Dialect:  "mysql",
		Host:     os.Getenv("DB_HOST"),
		Database: os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		RWMutex:  new(sync.RWMutex),
	}

	// security
	security := &saas_security_basic.Basic{
		Database: database,
		RWMutex:  new(sync.RWMutex),
	}

	// jwt
	jwt := &saas_api.Jwt{
		Key:     os.Getenv("JWT_KEY"),
		RWMutex: new(sync.RWMutex),
	}

	// event
	event := &saas_event_basic.Event{
		Hub:     new(saas_event_basic.Hub).Init(),
		RWMutex: new(sync.RWMutex),
	}

	// api
	api := &saas_api.Api{
		Logger:   logger,
		Event:    event,
		Security: security,
		Database: database,
		Origins:  strings.Split(os.Getenv("ORIGINS"), ","),
		Jwt:      jwt,
		Tls:      nil,
		Port:     os.Getenv("API_PORT"),
		Mode:     os.Getenv("API_MODE"),
		RWMutex:  new(sync.RWMutex),
	}

	saas := &go_saas.Saas{
		License:  "abc",
		Logger:   logger,
		Database: database,
		Api:      api,
		RWMutex:  new(sync.RWMutex),
	}

	if err := saas.Run(); err != nil {
		saas.GetLogger().Fatal(err)
	}
}

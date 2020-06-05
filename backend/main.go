package main

import (
	"os"
	"strings"
	"sync"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/go-saas/go-saas"
	"github.com/go-saas/go-saas/authenticator/basic"
	"github.com/go-saas/go-saas/database/basic"
	"github.com/go-saas/go-saas/event/basic"
	"github.com/go-saas/go-saas/http"
	"github.com/go-saas/go-saas/http/basic"
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

	// event hub
	hub := &go_saas_basic_event.Hub{
		List:    make(map[uint]map[uint]chan sse.Event),
		RWMutex: new(sync.RWMutex),
	}

	// event
	event := &go_saas_basic_event.Event{
		Hub:     hub,
		RWMutex: new(sync.RWMutex),
	}

	// jwt authenticator
	authenticator := &go_saas_basic_authenticator.Authenticator{
		RWMutex: new(sync.RWMutex),
	}

	// http
	http := &go_saas_basic_http.Http{
		Router:        gin.Default(),
		Handlers:      make(map[string]func(http go_saas_http.Http) error),
		Logger:        logger,
		Event:         event,
		Authenticator: authenticator,
		Security:      security,
		Database:      database,
		Jwt:           jwt,
		Tls:           nil,
		Origins:       strings.Split(os.Getenv("ORIGINS"), ","),
		Port:          os.Getenv("API_PORT"),
		Mode:          os.Getenv("API_MODE"),
		RWMutex:       new(sync.RWMutex),
	}

	saas := &go_saas.Saas{
		Logger:   logger,
		Database: database,
		Http:     http,
		RWMutex:  new(sync.RWMutex),
	}

	if err := saas.Init(); err != nil {
		saas.GetLogger().Fatal(err)
	}

	if err := saas.Run(); err != nil {
		saas.GetLogger().Fatal(err)
	}
}

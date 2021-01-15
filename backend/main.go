package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/makeless/makeless-go"
	"github.com/makeless/makeless-go/authenticator/basic"
	"github.com/makeless/makeless-go/config/basic"
	"github.com/makeless/makeless-go/database/basic"
	"github.com/makeless/makeless-go/event/basic"
	"github.com/makeless/makeless-go/http"
	"github.com/makeless/makeless-go/http/basic"
	"github.com/makeless/makeless-go/logger"
	"github.com/makeless/makeless-go/logger/basic"
	"github.com/makeless/makeless-go/mailer"
	"github.com/makeless/makeless-go/mailer/basic"
	"github.com/makeless/makeless-go/queue"
	"github.com/makeless/makeless-go/queue/basic"
	"github.com/makeless/makeless-go/security/basic"
	"gorm.io/driver/mysql"
)

func main() {
	// logger
	logger := new(makeless_go_logger_basic.Logger)

	// config
	config := &makeless_go_config_basic.Config{
		RWMutex: new(sync.RWMutex),
	}

	// queue
	mailQueue := &makeless_go_queue_basic.Queue{
		Context: context.Background(),
		RWMutex: new(sync.RWMutex),
	}

	// mailer
	mailer := &makeless_go_mailer_basic.Mailer{
		Handlers: make(map[string]func(data map[string]interface{}) (makeless_go_mailer.Mail, error)),
		Queue:    mailQueue,
		Host:     os.Getenv("MAILER_HOST"),
		Port:     os.Getenv("MAILER_PORT"),
		Identity: os.Getenv("MAILER_IDENTITY"),
		Username: os.Getenv("MAILER_USERNAME"),
		Password: os.Getenv("MAILER_PASSWORD"),
		RWMutex:  new(sync.RWMutex),
	}

	// database
	database := &makeless_go_database_basic.Database{
		Host:     os.Getenv("DB_HOST"),
		Database: os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		RWMutex:  new(sync.RWMutex),
	}

	// security
	security := &makeless_go_security_basic.Security{
		Database: database,
		RWMutex:  new(sync.RWMutex),
	}

	// event hub
	hub := &makeless_go_event_basic.Hub{
		List:    new(sync.Map),
		RWMutex: new(sync.RWMutex),
	}

	// event
	event := &makeless_go_event_basic.Event{
		Hub:     hub,
		Error:   make(chan error),
		RWMutex: new(sync.RWMutex),
	}

	// jwt authenticator
	authenticator := &makeless_go_authenticator_basic.Authenticator{
		Security:    security,
		Realm:       "auth",
		Key:         os.Getenv("JWT_KEY"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "id",
		RWMutex:     new(sync.RWMutex),
	}

	// router
	router := &makeless_go_http_basic.Router{
		RWMutex: new(sync.RWMutex),
	}

	// http
	http := &makeless_go_http_basic.Http{
		Router:        router,
		Handlers:      make(map[string]func(http makeless_go_http.Http) error),
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

	makeless := &makeless_go.Makeless{
		Config:   config,
		Logger:   logger,
		Mailer:   mailer,
		Database: database,
		Http:     http,
		RWMutex:  new(sync.RWMutex),
	}

	if err := makeless.Init(mysql.Open(database.GetConnectionString()), "./makeless.json"); err != nil {
		makeless.GetLogger().Fatal(err)
	}

	// async mail
	go func(mailer makeless_go_mailer.Mailer, logger makeless_go_logger.Logger) {
		for {
			select {
			case <-mailer.GetQueue().GetContext().Done():
				return
			case <-time.After(1 * time.Second):
				var err error
				var empty bool
				var node makeless_go_queue.Node

				for {
					if empty, err = mailer.GetQueue().Empty(); err != nil {
						logger.Fatal(err)
					}

					if empty {
						break
					}

					if node, err = mailer.GetQueue().Remove(); err != nil {
						logger.Fatal(err)
					}

					var mail = &makeless_go_mailer_basic.Mail{
						RWMutex: new(sync.RWMutex),
					}
					
					if err := json.Unmarshal(node.GetData(), mail); err != nil {
						logger.Fatal(err)
					}

					if err = mailer.Send(context.Background(), mail); err != nil {
						logger.Fatal(err)
					}
				}
			}
		}
	}(mailer, logger)

	if err := makeless.Run(); err != nil {
		makeless.GetLogger().Fatal(err)
	}
}

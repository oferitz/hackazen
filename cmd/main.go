package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/oferitz/hackazen/internal/data"
	"github.com/oferitz/hackazen/internal/db"
	"github.com/oferitz/hackazen/internal/mailer"
	sessionStore "github.com/oferitz/hackazen/internal/session"
	"go.uber.org/zap"
	"log"
	"sync"
)

type application struct {
	config       *koanf.Koanf
	server       *fiber.App
	sessionStore *session.Store
	logger       *zap.SugaredLogger
	models       data.Models
	mailer       mailer.Mailer
	wg           sync.WaitGroup
}

func main() {
	var cfg = koanf.New(".")
	if err := cfg.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	dbConn, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}

	defer dbConn.Close()

	store, err := sessionStore.New(cfg)
	if err != nil {
		log.Fatalf(err.Error())
	}

	app := &application{
		config:       cfg,
		server:       fiber.New(),
		sessionStore: store,
		logger:       logger.Sugar(),
		models:       data.NewModels(dbConn),
		mailer:       mailer.New(cfg),
	}

	app.initRoutes()
	err = app.serve()
	if err != nil {
		log.Fatalf(err.Error())
	}
}

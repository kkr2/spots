package main

import (
	"log"
	"os"

	"github.com/kkr2/spots/internal/config"
	"github.com/kkr2/spots/internal/server"
	"github.com/kkr2/spots/pkg/db/postgres"
	"github.com/kkr2/spots/pkg/logger"
	"github.com/kkr2/spots/pkg/utils"
)

// @title Spot Query
// @version 1.0
// @description Example Api for querying GIS data
// @contact.name Krist Kokali
// @contact.url https://github.com/kkr2
// @contact.email kristkokali21@gmail.com
// @BasePath /api/v1
func main() {
	log.Println("Starting api server")

	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode)

	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgres connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	s := server.NewServer(cfg, psqlDB, appLogger)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}

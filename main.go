package main

import (
	"os"
	"os/signal"
	"syscall"

	"transport/lib/errors"
	"transport/lib/graceful"
	"transport/lib/utils/config"
	"transport/lib/utils/logger"

	"telegram-bot/app"
	v1 "telegram-bot/app/router/v1"
	"telegram-bot/pkg/utils"
)

// @title Event Management System
// @version 1.0
// @description Event Management System API Swagger documents

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	//load config
	config.LoadConfig("config")

	//Init errors map
	errors.SetErrorMap(utils.DefaultErrorMap)

	// Build DIG container
	container := app.BuildContainer()

	serv := v1.Initialize(container)

	// graceful service
	s := graceful.Register(serv)
	defer s.Close()
	go s.StartServer(serv)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	close(quit)
	logger.Info("Shutting down")

}

package main

import (
	"os"
	"os/signal"
	"syscall"

	"transport/lib/errors"
	"transport/lib/graceful"
	"transport/lib/utils/config"
	"transport/lib/utils/logger"
	"transport/lib/utils/tservice"

	"github.com/spf13/viper"

	"telegram-bot/app"
	v1 "telegram-bot/app/router/v1"
	"telegram-bot/app/services"
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
	config.LoadConfig("trip-config")

	//Init errors map
	errors.Initialize(tservice.EMS)
	// errors.SetErrorMap(utils.DefaultErrorMap)

	// Build DIG container
	container := app.BuildContainer()

	mode := viper.GetInt("ts_service.mode")
	logger.Info("Server run mode: ", mode)
	// Start by mode
	// mode = 0: run all: consumer, cron, publisher and api gateway
	// mode = 1: run publisher and api gateway
	// mode = 2: run consumer
	// mode = 3: run cron gateway
	if mode == 0 || mode == 1 {
		//Init serv
		serv := v1.Initialize(container)

		// graceful service
		s := graceful.Register(serv)
		defer s.Close()
		go s.StartServer(serv)
	}

	if mode == 0 || mode == 2 {
		container.Invoke(func(
			inService services.ActionService,
		) {
			go inService.Consume()
		})
	}

	if mode == 3 {
		//Init serv
		serv := v1.InitializeCron(container)

		// graceful service
		s := graceful.Register(serv)
		defer s.Close()
		go s.StartServer(serv)
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	close(quit)
	logger.Info("Shutting down")

}

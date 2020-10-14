package api

import (
	"github.com/gin-gonic/gin"
	"transport/lib/errors"
	"transport/lib/thttp"
	"transport/lib/utils/logger"

	"telegram-bot/app/services"
)

type Cron struct {
	outService services.OutService
	inService  services.ActionService
}

func NewCron(outService services.OutService, inService services.ActionService) *Cron {
	return &Cron{
		outService: outService,
		inService:  inService,
	}
}

func (s *Cron) Resend(c *gin.Context) thttp.Response {
	logger.Info("Start cronjob resend wait messages")
	go s.outService.CronResend(c)
	return thttp.Response{
		Error: errors.Success.New(),
	}
}

func (s *Cron) Retry(c *gin.Context) thttp.Response {
	logger.Info("Start cronjob retry wait messages")
	go s.inService.CronRetry(c)
	return thttp.Response{
		Error: errors.Success.New(),
	}
}

func (s *Cron) RetryPrevious(c *gin.Context) thttp.Response {
	logger.Info("Start cronjob resend wait previous messages")
	go s.inService.CronRetryPrevious(c)
	return thttp.Response{
		Error: errors.Success.New(),
	}
}

func (s *Cron) ArchivedMessages(c *gin.Context) thttp.Response {
	logger.Info("Start cronjob archive messages")
	go s.outService.CronArchivedMessages(c)
	go s.inService.CronArchivedMessages(c)
	return thttp.Response{
		Error: errors.Success.New(),
	}
}

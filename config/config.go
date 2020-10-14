package config

import "github.com/spf13/viper"

const (
	DefaultPageSize  = 25
	DefaultSortField = "-_id"

	DefaultMaxRetryTimes   = 3
	DefaultConsumerThreads = 10

	CronRetrySize             = 1000
	CronResendOutMessageLimit = 100
	CronArchivedMessageLimit  = 1000

	ArchivedMessageDays = 7
)

func GetConsumerThreads() int {
	retryTimes := viper.GetInt("consumer.threads")
	if retryTimes <= 0 {
		retryTimes = DefaultConsumerThreads
	}

	return retryTimes
}

func GetArchivedDays() int {
	days := viper.GetInt("cron.archived_days")
	if days <= 0 {
		days = ArchivedMessageDays
	}

	return days
}

func GetRetrySize() int {
	size := viper.GetInt("cron.retry_size")
	if size <= 0 {
		size = CronRetrySize
	}

	return size
}

func GetArchivedSize() int {
	size := viper.GetInt("cron.archived_size")
	if size <= 0 {
		size = CronArchivedMessageLimit
	}

	return size
}

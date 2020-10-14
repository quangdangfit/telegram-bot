package dbs

import (
	"github.com/spf13/viper"
	"transport/lib/database"
	"transport/lib/utils/logger"
)

func NewDatabase() database.MongoDB {
	dbConfig := database.Config{
		Host:         viper.GetString("ts_mongodb.host"),
		Replica:      viper.GetString("ts_mongodb.replica"),
		UserName:     viper.GetString("ts_mongodb.user"),
		Password:     viper.GetString("ts_mongodb.pass"),
		AuthDatabase: viper.GetString("ts_mongodb.authdb"),
		Database:     viper.GetString("ts_mongodb.db"),
	}
	db, err := database.New(dbConfig)
	if err != nil {
		logger.Panic(err)
	}

	logger.Infof("Connect database %s successfully!", dbConfig.Database)
	return db
}

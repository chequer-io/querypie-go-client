package config

import (
	_ "github.com/mattn/go-sqlite3" // Import for side effects to register the SQLite3 driver
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

var LocalDatabase *gorm.DB

func getLogLevel(level string) logger.LogLevel {
	switch level {
	case "debug":
		return logger.Info
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	case "fatal":
		return logger.Error
	case "panic":
		return logger.Error
	default:
		return logger.Warn
	}
}

func initLocalDatabase(dataSourceName string, logLevel string) {
	var gormLogLevel = getLogLevel(logLevel)

	// Initialize the GORM connection with SQLite
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		logrus.Fatal(err)
	}

	LocalDatabase = db.
		// Enable full save associations by default.
		// https://gorm.io/docs/associations.html#Updating-Associations-with-FullSaveAssociations
		Session(&gorm.Session{FullSaveAssociations: true})
}

func InitConfigForLocalDatabase(viper *viper.Viper, logLevel string) {
	var dataSourceName string
	if err := viper.UnmarshalKey("sqlite3-data-source", &dataSourceName); err != nil {
		logrus.Fatalf("Unable to unmarshal sqlite3-data-source: %v", err)
		os.Exit(1)
	}

	if dataSourceName == "" {
		logrus.Fatalf("sqlite3-data-source is not set in the configuration")
		os.Exit(1)
	}

	logrus.Infof("SQLite Data Source: %s", dataSourceName)
	initLocalDatabase(dataSourceName, logLevel)
}

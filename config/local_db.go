package config

import (
	_ "github.com/mattn/go-sqlite3" // Import for side effects to register the SQLite3 driver
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var LocalDatabase *gorm.DB

func initLocalDatabase(dataSourceName string) {
	// Initialize the GORM connection with SQLite
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}
	LocalDatabase = db
}

func InitConfigForLocalDatabase(viper *viper.Viper) {
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
	initLocalDatabase(dataSourceName)
}

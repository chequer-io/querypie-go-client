package local_db

import (
	_ "github.com/mattn/go-sqlite3" // Import for side effects to register the SQLite3 driver
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"qpc/models"
)

var LocalDatabase *gorm.DB

func initLocalDatabase(dataSourceName string) {
	// Initialize the GORM connection with SQLite
	db, err := gorm.Open(sqlite.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}
	LocalDatabase = db

	err1 := db.AutoMigrate(
		&models.UserV2{},
		&models.AdminRole{},
		&models.UserV1{},
		&models.UserRole{},
		&models.Role{},
		&models.SummarizedConnectionV2{},
	)
	if err1 != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("AutoMigrate has done successfully!")
}

func InitConfigForResource(viper *viper.Viper) {
	var dataSourceName string
	if err := viper.UnmarshalKey("sqlite3_data_source", &dataSourceName); err != nil {
		logrus.Fatalf("Unable to unmarshal sqlite3_data_source: %v", err)
		os.Exit(1)
	}

	if dataSourceName == "" {
		logrus.Fatalf("sqlite3_data_source is not set in the configuration")
		os.Exit(1)
	}

	logrus.Infof("SQLite Data Source: %s", dataSourceName)
	initLocalDatabase(dataSourceName)
}

package local_db

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import for side effects to register the SQLite3 driver
)

var localDatabase *sql.DB

func initLocalDatabase(dataSourceName string) {
	// Open a connection to the SQLite database
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		logrus.Fatal(err)
	}
	localDatabase = db
	defer db.Close()

	// Create a table
	//goland:noinspection SqlNoDataSourceInspection
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"name" TEXT,
		"age" INTEGER
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("Table created successfully!")
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

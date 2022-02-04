package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ockibagusp/golang-website-example/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var connection *gorm.DB
var db *sql.DB
var err error

// Init: new database
func Init(env string) *gorm.DB {
	var connectString string
	configuration := config.GetConfig()

	if env == "PROD" {
		connectString = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			configuration.PROD.DB_USERNAME,
			configuration.PROD.DB_PASSWORD,
			configuration.PROD.DB_HOST,
			configuration.PROD.DB_PORT,
			configuration.PROD.DB_NAME,
		)
	} else if env == "DEV" {
		connectString = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			configuration.DEV.DB_USERNAME,
			configuration.DEV.DB_PASSWORD,
			configuration.DEV.DB_HOST,
			configuration.DEV.DB_PORT,
			configuration.DEV.DB_NAME,
		)
	} else {
		panic("failed to env: PROD or DEV")
	}

	connection, err = gorm.Open(mysql.Open(connectString), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db, _ = connection.DB()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(10 * time.Minute)

	return connection
}

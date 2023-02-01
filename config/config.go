package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	Config struct {
		// app
		DBDriver string
		// mysql
		DBMySQLHost     string
		DBMySQLPort     string
		DBMySQLUser     string
		DBMySQLPassword string
		DBMySQLName     string
		// sqlite
		DBSQLiteName string
	}
)

func GetAPPConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	return &Config{
		// app
		DBDriver: os.Getenv("APP_DB_DRIVER"),
		// mysql
		DBMySQLHost:     os.Getenv("DB_MYSQL_HOST"),
		DBMySQLPort:     os.Getenv("DB_MYSQL_PORT"),
		DBMySQLUser:     os.Getenv("DB_MYSQL_USER"),
		DBMySQLPassword: os.Getenv("DB_MYSQL_PASSWORD"),
		DBMySQLName:     os.Getenv("DB_MYSQL_NAME"),
		// sqlite
		DBSQLiteName: os.Getenv("DB_SQLITE_NAME"),
	}
}

func (config *Config) GetDatabaseConnection() *gorm.DB {
	if config.DBDriver == "mysql" {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.DBMySQLUser,
			config.DBMySQLPassword,
			config.DBMySQLHost,
			config.DBMySQLPort,
			config.DBMySQLName,
		)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newDBLogger()})
		if err != nil {
			log.Fatal(err)
		}

		return db.Debug()
	}

	if config.DBDriver == "sqlite" {
		db, err := gorm.Open(sqlite.Open(config.DBSQLiteName), &gorm.Config{Logger: newDBLogger()})
		if err != nil {
			log.Fatal(err)
		}

		return db.Debug()
	}

	log.Println(config.DBDriver)
	log.Fatal("unsupported driver")

	return nil
}

func newDBLogger() logger.Interface {
	return logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold:             30 * time.Second, // Slow SQL threshold
			LogLevel:                  logger.Silent,    // Log level
			IgnoreRecordNotFoundError: false,            // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,            // Enable Color
		},
	)
}

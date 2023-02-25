package config

import (
	"fmt"
	"log"
	"os"
	"regexp"
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
		AppDBDriver string
		// mysql
		DBMySQLHost     string
		DBMySQLPort     string
		DBMySQLUser     string
		DBMySQLPassword string
		DBMySQLName     string
		// sqlite
		DBSQLiteName string
		// secure cookie
		SessionsCookieStore string
	}
)

const projectDirName = "golang-website-example"

func loadEnv() {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func GetAPPConfig() *Config {
	loadEnv()

	return &Config{
		// app
		AppDBDriver: os.Getenv("APP_DB_DRIVER"),
		// mysql
		DBMySQLHost:     os.Getenv("DB_MYSQL_HOST"),
		DBMySQLPort:     os.Getenv("DB_MYSQL_PORT"),
		DBMySQLUser:     os.Getenv("DB_MYSQL_USER"),
		DBMySQLPassword: os.Getenv("DB_MYSQL_PASSWORD"),
		DBMySQLName:     os.Getenv("DB_MYSQL_NAME"),
		// sqlite
		DBSQLiteName: os.Getenv("DB_SQLITE_NAME"),
		// secure cookie
		SessionsCookieStore: os.Getenv("DB_SQLITE_NAME"),
	}
}

func (config *Config) GetDatabaseConnection() *gorm.DB {
	if config.AppDBDriver == "mysql" {
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

	if config.AppDBDriver == "sqlite" {
		db, err := gorm.Open(sqlite.Open(config.DBSQLiteName), &gorm.Config{Logger: newDBLogger()})
		if err != nil {
			log.Fatal(err)
		}

		return db.Debug()
	}

	log.Println(config.AppDBDriver)
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

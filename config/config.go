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
		// session test
		SessionTest string
		// debug
		Debug string
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
		SessionsCookieStore: os.Getenv("SESSIONS_COOKIE_STORE"),
		// session test
		SessionTest: os.Getenv("SESSION_TEST"),
		// debug
		Debug: os.Getenv("DEBUG"),
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

func (config *Config) GetSessionTest() {
	if config.SessionTest == "true" || config.SessionTest == "1" {
		os.Setenv("SESSION_TEST", "1")
	} else if config.SessionTest == "" || config.SessionTest == "0" {
		os.Setenv("SESSION_TEST", "0")
	} else {
		log.Fatal("unsupported session test")
	}
}

func (config *Config) GetDebug() {
	if config.Debug == "true" || config.Debug == "1" {
		os.Setenv("DEBUG", "1")
	} else if config.Debug == "false" || config.Debug == "0" {
		os.Setenv("DEBUG", "0")
	} else {
		log.Fatal("unsupported debug")
	}
}

func (config *Config) GetDebugAsTrue(debug []bool) bool {
	if (len(debug) == 1 && debug[0] == true) || os.Getenv("DEBUG") == "1" {
		return true
	}
	panic("func GetDebugAsTrue: (debug [1]: true or false) or no debug")
}

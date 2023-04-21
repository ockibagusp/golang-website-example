package config

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/ockibagusp/golang-website-example/app/main/controller/mock/method"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	Config struct {
		// app
		AppDBDriver    string
		AppJWTAuthSign string
		// mysql
		DBMySQLHost     string
		DBMySQLPort     string
		DBMySQLUser     string
		DBMySQLPassword string
		DBMySQLName     string
		// secure cookie
		SessionsCookieStore string
		// session test
		SessionTest string
		// debug
		Debug string
	}
)

const projectDirName = "golang-website-example"

func fullProjetDir() string {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	return string(rootPath)
}

func loadEnv() {
	err := godotenv.Load(fullProjetDir() + `/.env`)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func GetAPPConfig() *Config {
	loadEnv()

	return &Config{
		// app
		AppDBDriver:    os.Getenv("APP_DB_DRIVER"),
		AppJWTAuthSign: os.Getenv("APP_JWT_AUTH_SIGN"),
		// mysql
		DBMySQLHost:     os.Getenv("DB_MYSQL_HOST"),
		DBMySQLPort:     os.Getenv("DB_MYSQL_PORT"),
		DBMySQLUser:     os.Getenv("DB_MYSQL_USER"),
		DBMySQLPassword: os.Getenv("DB_MYSQL_PASSWORD"),
		DBMySQLName:     os.Getenv("DB_MYSQL_NAME"),
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

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			log.Fatal(err)
		}

		return db.Debug()
	}

	log.Fatal("unsupported driver")
	return nil
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

func (config *Config) SetSessionToFalse() bool {
	return os.Getenv("SESSION_TEST") == "1" && method.SetSession == false
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
	} else if len(debug) > 1 {
		panic("func (*Config) GetDebugAsTrue: (debug [1]: true or false) or no debug")
	}

	return false
}

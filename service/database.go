package service

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Here we initializes the database
	"github.com/pressly/goose"
	"github.com/sirupsen/logrus"
)

const (
	postgresHostKey        = "POSTGRES_HOST"
	postgresPasswordKey    = "POSTGRES_PASSWORD"
	postgresUserKey        = "POSTGRES_USER"
	postgresSSLModeKey     = "POSTGRES_SSL_MODE"
	serviceDataBaseNameKey = "SERVICE_DATABASE_NAME"
)

func (a *App) initializeDB() {
	if !a.appUsesDB() {
		logrus.Info("No DB connection string found. Skipping DB ")
		return
	}
	err := a.connectToDBServer()
	defer a.serverDB.Close()

	if err != nil {
		logrus.Fatalf("Can't ping database!! About: %v", err)
	}

	a.assertDBExists()
	err = a.connectToDB()

	if err != nil {
		logrus.Fatalf("can't reach database: %v", err)
	}

	a.checkAndDoDBUpdates()
}

func (a *App) appUsesDB() bool {

	return envExist(postgresHostKey) &&
		envExist(postgresPasswordKey) &&
		envExist(postgresUserKey) &&
		envExist(postgresSSLModeKey) &&
		envExist(serviceDataBaseNameKey)
}

func (a *App) connectToDBServer() error {
	var err error
	connectionParams := getServerConnectionString()
	for i := 0; i < 5; i++ {
		a.serverDB, err = sqlx.Open("postgres", connectionParams) // gorm checks Ping on Open
		a.serverDB.Ping()
		if err == nil {
			return nil
		}
		time.Sleep(5 * time.Second)
	}

	logrus.Warningf("Can't stablish connection to postgres server: %v", err)
	return err
}

func (a *App) assertDBExists() error {
	if a.dbExists() {
		return nil
	}
	a.createDB()

	return nil
}

func (a *App) dbExists() bool {
	dbName := os.Getenv(serviceDataBaseNameKey)
	query := "SELECT true from pg_database WHERE datname=$1"
	var exists bool
	err := a.serverDB.Get(&exists, query, dbName)
	return err == nil && exists
}

func (a *App) createDB() {
	dbName := os.Getenv(serviceDataBaseNameKey)
	query := fmt.Sprintf("CREATE DATABASE %v;", dbName)
	a.serverDB.MustExec(query)
}

func (a *App) checkAndDoDBUpdates() {
	migrations := getMigrationsPath()
	datastoreName := getDBConnectionString()
	var gooseDB *sql.DB

	gooseDB, _ = sql.Open("postgres", datastoreName)
	err := gooseDB.Ping()

	if err != nil {
		logrus.Fatalf("Error Opening gooseDB: %v\n", err)
	}
	err = goose.Up(gooseDB, migrations)
	if err != nil {
		logrus.Fatalf("Error on DB migrations: %v", err)
	}
	gooseDB.Close()
}

func (a *App) connectToDB() error {
	connectionParams := getDBConnectionString()
	db, err := sqlx.Open("postgres", connectionParams)
	err = db.Ping()
	if err != nil {
		logrus.Warningf("Error Opening mainDB: %v\n", err)
		return err
	}
	a.db = db
	return nil
}

func getServerConnectionString() string {
	host := os.Getenv(postgresHostKey)
	password := os.Getenv(postgresPasswordKey)
	user := os.Getenv(postgresUserKey)
	sslMode := os.Getenv(postgresSSLModeKey)

	return fmt.Sprintf(
		"user=%v password=%v sslmode=%v host=%v",
		user,
		password,
		sslMode,
		host,
	)
}

func getDBConnectionString() string {
	dbName := os.Getenv(serviceDataBaseNameKey)
	return fmt.Sprintf("dbname=%v %v", dbName, getServerConnectionString())
}

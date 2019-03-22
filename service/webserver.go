package service

import (
	"database/sql"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose"
	"github.com/sirupsen/logrus"
	"google.golang.org/appengine"
)

const appName = "service-backbone"
const migrationDir = "../migrations"

// App models the main app
type App struct {
	router *gin.Engine
	db     *sqlx.DB
}

// Initialize start the DB and the router
func (a *App) Initialize() {
	a.initializeDB()
	a.router = gin.Default()
	a.setCors()
	a.initializeRoutes()
}

// Run starts listening and serving HHTP requests
func (a *App) Run() {
	switch true {
	case appengine.IsAppEngine():
		http.Handle("/", a.router)
		appengine.Main()

	default:
		a.router.Run(os.Getenv("PORT"))
	}
}

func (a *App) initializeDB() {
	if !a.appUsesDB() {
		logrus.Info("No DB connection string found. Skipping DB ")
		return
	}
	a.checkAndDoDBUpdates()
	a.connectToDB()
	a.assertDBConnection()
}

func (a *App) setCors() {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "PUT", "POST", "HEAD", "DELETE", "PATCH", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "authorization"}
	a.router.Use(cors.New(config))
}

func (a *App) appUsesDB() bool {
	datastoreName := os.Getenv("POSTGRES_CONNECTION")
	return datastoreName != ""
}

func (a *App) checkAndDoDBUpdates() {
	migrations := getMigrationsPath()
	datastoreName := os.Getenv("POSTGRES_CONNECTION")
	gooseDB, err := sql.Open("postgres", datastoreName)
	if err != nil {
		logrus.Fatalf("Error Opening gooseDB: %v\n", err)
	}
	err = goose.Up(gooseDB, migrations)
	if err != nil {
		logrus.Fatalf("Error on DB migrations: %v", err)
	}
	gooseDB.Close()
}

func (a *App) connectToDB() {
	datastoreName := os.Getenv("POSTGRES_CONNECTION")
	db, err := sqlx.Open("postgres", datastoreName)
	if err != nil {
		logrus.Fatalf("Error Opening mainDB: %v\n", err)
	}
	a.db = db
}

func (a *App) assertDBConnection() {
	err := a.db.Ping()
	if err != nil {
		logrus.Fatalf("Error conecting with mainDB: %v\n", err)
	}
}

func getMigrationsPath() string {
	migrations := getDir(migrationDir)
	return migrations
}

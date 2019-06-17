package service

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"

	firebase "firebase.google.com/go"

	"cloud.google.com/go/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

const (
	appName            = "admin_service"
	migrationDir       = "../migrations"
	firebaseFileName   = "../firebase_credential.json"
	production         = "prod"
	local              = "local"
	postgresConnection = "POSTGRES_CONNECTION"
)

// App models the main app
type App struct {
	router   *gin.Engine
	db       *sqlx.DB
	serverDB *sqlx.DB
	authApp  *firebase.App
}

// Initialize start the DB and the router
func (a *App) Initialize() {
	a.initializeDB()
	a.initializeRouter()
	a.setCors()
	a.initializeLogger()
	a.initializeAuthSystem()
	a.initializeRoutes()
	a.serRouterMode()
}

// Run starts listening and serving HHTP requests
func (a *App) Run() {
	env := os.Getenv("ENV")
	switch true {
	case env == local:
		logrus.Infof("Local env")
		a.router.Run(os.Getenv("PORT"))

	default:
		logrus.Infof("Running on appengine env")
		http.Handle("/", a.router)
		appengine.Main()
	}
}

func (a *App) initializeRouter() {
	a.router = gin.Default()
}

func (a *App) setCors() {
	if !isPublicAPI() {
		return
	}
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "PUT", "POST", "HEAD", "DELETE", "PATCH", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "authorization"}
	a.router.Use(cors.New(config))
}

func (a *App) initializeLogger() {
	profile := os.Getenv("ENV")

	if profile == "local" {
		logrus.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000",
			FullTimestamp:   true,
		})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}

func (a *App) initializeAuthSystem() {
	a.fetchCredentials()
	a.initAuthApp()
}

func (a *App) fetchCredentials() {
	if !mustFetchFirebaseCredentials() {
		return
	}
	credentials := getCredentialFile()
	saveCredentialFile(credentials)
}

func (a *App) initAuthApp() {
	ctx := context.Background()

	firebaseFile := getDir(firebaseFileName)
	opt := option.WithCredentialsFile(firebaseFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		logrus.Fatalf("error initializing app: %v\n", err)
	}
	a.authApp = app
}

func (a *App) serRouterMode() {
	env := os.Getenv("ENV")
	if env == production {
		gin.SetMode(gin.ReleaseMode)
	}
}

func getCredentialFile() *storage.Reader {
	ctx := context.Background()
	bucket := os.Getenv("GCLOUD_STORAGE_BUCKET")

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		logrus.Fatalf("Can't get storage conection: %v", err)
	}
	sw, err := storageClient.Bucket(bucket).Object("firebase_credential.json").NewReader(ctx)
	if err != nil {
		logrus.Fatalf("Can't stablish reader connection: %v", err)
	}
	return sw
}

func saveCredentialFile(credentials *storage.Reader) {
	buff, err := ioutil.ReadAll(credentials)
	if err != nil {
		logrus.Fatalf("Can't read credentials: %v", err)
	}
	err = ioutil.WriteFile("firebase_credential.json", buff, 0644)
	if err != nil {
		logrus.Fatalf("Can't write firebase_credential.json: %v", err)
	}
}

func getMigrationsPath() string {
	migrations := getDir(migrationDir)
	return migrations
}

func mustFetchFirebaseCredentials() bool {
	bucket := os.Getenv("GCLOUD_STORAGE_BUCKET")
	return bucket != ""
}

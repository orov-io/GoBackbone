package service

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"

	"firebase.google.com/go/auth"
	_ "github.com/lib/pq" // Here we initializes the database
)

// ServeTestHTTP allows serve http responses from tests.
func (a *App) ServeTestHTTP(resp *httptest.ResponseRecorder, req *http.Request) {
	a.router.ServeHTTP(resp, req)
}

func isUserType(userType string, token *auth.Token) bool {
	if thisUser, ok := token.Claims[userType]; !(ok && thisUser.(bool)) {
		return false
	}
	return true
}

func (a *App) getAppName() string {
	prefix := ""
	if os.Getenv("ENV") != "prod" {
		prefix = os.Getenv("ENV")
	}
	return prefix + "-" + appName
}

func getFirebaseCredentials() string {
	return "firebase_credential.json"
}

func getDir(dir string) string {
	_, file, _, _ := runtime.Caller(1)
	FileDir := filepath.Dir(file)
	assets := filepath.Join(FileDir, dir)
	return assets
}

func dirExist(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}

func dirIsNotEmpty(dir string) bool {
	files, err := ioutil.ReadDir(dir)
	return err != nil && len(files) > noFiles
}

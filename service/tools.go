package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"firebase.google.com/go/auth"
)

const authorizedUserRole = "admin"
const (
	verifyCustomTokenURL = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyCustomToken?key=%s"
	verifyPasswordURL    = "https://www.googleapis.com/identitytoolkit/v3/relyingparty/verifyPassword?key=%s"
)

// ServeTestHTTP allows serve http responses for tests.
func (a *App) ServeTestHTTP(resp *httptest.ResponseRecorder, req *http.Request) {
	a.router.ServeHTTP(resp, req)
}

func (a *App) getAppName() string {
	prefix := ""
	if os.Getenv("ENV") != "prod" {
		prefix = os.Getenv("ENV")
	}
	return prefix + "-" + appName
}

func (a *App) getAuthClient() *auth.Client {
	ctx := context.Background()
	client, _ := a.authApp.Auth(ctx)
	return client
}

func (a *App) patrol() gin.HandlerFunc {
	return func(c *gin.Context) {
		userToken, err := getUserTokenFromHeader(c)
		if err != nil {
			sendUnauthorizedAccess(c, err)
			return
		}

		authClient := a.getAuthClient()
		token, err := authClient.VerifyIDToken(c, userToken)
		if err != nil {
			sendUnauthorizedAccess(c, err)
			return
		}

		if !isSuperAdmin(token) {
			sendUnauthorizedAccess(c, nil)
			return
		}

		logAuthorizedAccess(token)
	}
}

func getUserTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return "", fmt.Errorf("Header must start with Bearer")
	}
	return splitToken[1], nil
}

func logAuthorizedAccess(user *auth.Token) {
	logrus.Infof("Authorized user: %v", user.UID)
}

func isSuperAdmin(token *auth.Token) bool {
	superAdminFlag, ok := token.Claims["SuperAdmin"]
	return ok && superAdminFlag.(bool)
}

func getDir(dir string) string {
	_, file, _, _ := runtime.Caller(1)
	FileDir := filepath.Dir(file)
	directory := filepath.Join(FileDir, dir)
	return directory
}

func dirExist(dir string) bool {
	_, err := os.Stat(dir)
	return !os.IsNotExist(err)
}

func dirIsNotEmpty(dir string) bool {
	files, err := ioutil.ReadDir(dir)
	return err != nil && len(files) > noFiles
}

func isPublicAPI() bool {
	isPublicAPI, err := strconv.ParseBool(os.Getenv("PUBLIC_API"))
	if err != nil {
		return false
	}
	return isPublicAPI
}

// SignInWithCustomToken provides custom tokens functionality
func SignInWithCustomToken(token string) (string, error) {
	req, err := json.Marshal(map[string]interface{}{
		"token":             token,
		"returnSecureToken": true,
	})
	if err != nil {
		return "", err
	}

	apiKey := os.Getenv("GOOGLE_IDENTITY_API_KEY")
	if err != nil {
		return "", err
	}
	resp, err := postRequest(fmt.Sprintf(verifyCustomTokenURL, apiKey), req)
	if err != nil {
		logrus.Infof("1: %v", err)
		return "", err
	}
	var respBody struct {
		IDToken string `json:"idToken"`
	}
	if err := json.Unmarshal(resp, &respBody); err != nil {
		return "", err
	}
	return respBody.IDToken, err
}

func postRequest(url string, req []byte) ([]byte, error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected http status code: %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func inferErrorCode(theError error) int {
	r, _ := regexp.Compile("googleapi: Error [0-9]+")
	matchedError := r.FindString(theError.Error())
	if matchedError == "" {
		return http.StatusInternalServerError
	}
	matchedCode := matchedError[len(matchedError)-3:]
	code64, err := strconv.ParseInt(matchedCode, 10, 64)
	code := int(code64)
	if err != nil || !isInHTTPResponseRange(code) {
		return http.StatusInternalServerError
	}
	return code
}

func isInHTTPResponseRange(code int) bool {
	return code >= 100 && code <= 505
}

func envExist(envVar string) bool {
	return os.Getenv(envVar) != ""
}

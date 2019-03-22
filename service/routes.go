package service

import (
	"github.com/gin-gonic/gin"

	"github.com/sirupsen/logrus"
)

const (
	templatesPattern = "/**/*.gohtml"
	templatesDir     = "../templates"
	assetsDir        = "../assets"
	assetsEndpoint   = "/assets"
	noFiles          = 0
)

func (a *App) initializeRoutes() {
	a.serveAssets()
	a.parseTemplates()
	a.addRoutes()
}

func (a *App) addRoutes() {
	a.addPingRoute()
	a.addPrivateGroupExample()
}

func (a *App) serveAssets() {
	assets := getAssetsDir()
	if !mustServeDir(assets) {
		logrus.Info("No assets to serve")
		return
	}
	a.router.Static(assetsEndpoint, assets)
}

func (a *App) parseTemplates() {
	templates := getTemplatesDir()
	if !mustServeDir(templates) {
		logrus.Info("No templates to serve")
		return
	}
	a.router.LoadHTMLGlob(templates + templatesPattern)
}

func (a *App) addPingRoute() {
	a.router.GET("/ping", a.pong)
}

func (a *App) addPrivateGroupExample() {
	v1 := a.router.Group("/v1")
	v1.Use(accessControl())
	{
		v1.GET("/ping", a.pong)
	}
}

func accessControl() gin.HandlerFunc {
	return accessControlHandlerFunc
}

func accessControlHandlerFunc(c *gin.Context) {
	logrus.Info("In a private route")
}

func mustServeDir(dir string) bool {
	return dirExist(dir) && dirIsNotEmpty(dir)
}

func getTemplatesDir() string {
	return getDir(templatesDir)
}

func getAssetsDir() string {
	return getDir(assetsDir)
}

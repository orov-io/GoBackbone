package main

import (
	"github.com/orov.io/GoBackbone/service"
)

func main() {
	app := service.App{}
	app.Initialize()
	app.Run()
}

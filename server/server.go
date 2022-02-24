package server

import (
	"github.com/jerryan999/CryptoAlert/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() {

	// echo is a web framework for Go
	e := echo.New()
	e.Debug = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// controller
	ctl := NewAlertController()

	// CRUD operations for alerts
	e.POST("/add-alert", ctl.AddAlert)
	e.POST("/remove-alert", ctl.RemoveAlert)
	e.GET("/get-alerts", ctl.GetAlerts)
	e.POST("/update-alert", ctl.UpdateAlert)

	// Start the server
	e.Logger.Fatal(e.Start(":9000"))
}

func init() {
	database.Migrate(db)
}

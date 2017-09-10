package main

import (
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.GET("/healthz", func(c echo.Context) error {
		return c.String(200, "ok")
	})
	e.HideBanner = true

	e.Logger.Fatal(e.Start(":5309"))
}

package main

import (
	// "github.com/labstack/echo/v4"
	"github.com/KKGo-Software-engineering/assessment-tax/postgres"
	// "net/http"
)

func main() {
	// Create a new Postgres instance
	_, err := postgres.New();
	if err != nil {
		panic(err)
	}
	// e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, Go Bootcamp!")
	// })
	// e.Logger.Fatal(e.Start(":1323"))
}

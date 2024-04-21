package main

import (
	"github.com/labstack/echo/v4"
	"github.com/KKGo-Software-engineering/assessment-tax/postgres"
	"github.com/KKGo-Software-engineering/assessment-tax/calculateTax"
	// "net/http"
)

func main() {

	// Create a new Postgres instance
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	handler := calculateTax.New(p)
	e := echo.New()
	e.POST("tax/calculation", handler.HandleIncomeData)
	e.Logger.Fatal(e.Start(":8080"))
}

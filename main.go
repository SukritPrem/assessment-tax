package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/KKGo-Software-engineering/assessment-tax/postgres"
	"github.com/KKGo-Software-engineering/assessment-tax/calculateTax"
	// "net/http"
)

func AuthMiddleware(username, password string, c echo.Context) (bool, error) {
	if username == "adminTax" || password == "admin!" {
		return true, nil
	}
	return false, nil
}

func main() {

	// Create a new Postgres instance
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	handler := calculateTax.New(p)
	e := echo.New()
	e.POST("tax/calculation", handler.HandleIncomeData)
	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(AuthMiddleware))
	g.POST("/deductions/personal", handler.DeductionsPersonal)
	g.POST("/deductions/k-receipt", handler.DeductionsKReceipt)
	e.Logger.Fatal(e.Start(":8080"))
}

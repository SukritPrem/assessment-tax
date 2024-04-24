package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/KKGo-Software-engineering/assessment-tax/postgres"
	"github.com/KKGo-Software-engineering/assessment-tax/calculateTax"
	"os"
	"os/signal"
	"context"
	"time"
	"net/http"
	"fmt"
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
	e.POST("tax/calculation", handler.HandleCalculateTaxData)
	e.POST("tax/calculations/upload-csv", handler.HandleIncomeDataCSV)
	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(AuthMiddleware))
	g.POST("/deductions/personal", handler.DeductionsPersonal)
	g.POST("/deductions/k-receipt", handler.DeductionsKReceipt)

	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed { // Start server
			e.Logger.Fatal("shutting down the server")
		}
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Println("bye bye")
}

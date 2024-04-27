package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/SukritPrem/assessment-tax/postgres"
	"github.com/SukritPrem/assessment-tax/calculateTax"
	"os"
	"os/signal"
	"context"
	"time"
	"net/http"
	"fmt"
)

func AuthMiddleware(username, password string, c echo.Context) (bool, error) {
	if username == os.Getenv("ADMIN_USERNAME") || password == os.Getenv("ADMIN_PASSWORD") {
		return true, nil
	}
	return false, nil
}

func main() {

	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	handler := calculateTax.New(p)
	e := echo.New()
	e_client := e.Group("/")
	e_client.POST("tax/calculations", handler.HandleCalculateTaxData)
	e_client.POST("tax/calculations/upload-csv", handler.HandleIncomeDataCSV)
	g := e.Group("/admin")
	g.Use(middleware.BasicAuth(AuthMiddleware))
	g.POST("/deductions/personal", handler.DeductionsPersonal)
	g.POST("/deductions/k-receipt", handler.DeductionsKReceipt)

	go func() {
		if err := e.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed { // Start server
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

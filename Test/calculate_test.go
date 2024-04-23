package Test

import (
	"testing"
	"github.com/labstack/echo/v4"
	"net/http/httptest"
	"net/http"
	"github.com/KKGo-Software-engineering/assessment-tax/postgres"
	"github.com/KKGo-Software-engineering/assessment-tax/calculateTax"
	// "io/ioutil"
	// "strings"
	"bytes"
	"encoding/json"
)

type IncomeData struct {
  TotalIncome float64 `json:"TotalIncome"`
  Wht        float64 `json:"wht"`
  Allowances []struct {
    AllowanceType string  `json:"allowanceType"`
    Amount        float64 `json:"amount"`
  } `json:"allowances"`
}


func TestCalculate(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	reqBody := IncomeData{
		TotalIncome: 500000.0,
		Wht: 0.0,
		Allowances: []struct {
			AllowanceType string  `json:"allowanceType"`
			Amount        float64 `json:"amount"`
		}{
			{
				AllowanceType: "donation",
				Amount: 0.0,
			},
		},
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/tax/calculation", bytes.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// req.Body = ioutil.NopCloser(strings.NewReader(`
	// 		{
	// 			"totalIncome": 500000.0,
	// 			"wht": 0.0,
	// 			"allowances": [
	// 					{
	// 						"allowanceType": "donation",
	// 						"amount": 0.0
	// 					}
	// 				]
	// 		}`
	// 	)
	// )
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tax/calculation")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	expected := `{"tax":29000,"taxLevel":[{"level":"0-150,000","tax":0},{"level":"150,001-500,000","tax":29000},{"level":"500,001-1,000,000","tax":0},{"level":"1,000,001-2,000,000","tax":0},{"level":"2,000,001 ขึ้นไป","tax":0}]}`

	handler := calculateTax.New(p)
	err = handler.HandleIncomeData(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", rec.Code)
	}

	if rec.Body.String() != expected {
		t.Errorf("Expected %v, got %v", expected, rec.Body.String())
	}
}
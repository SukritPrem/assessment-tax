package Test

import (
	"testing"
	"github.com/labstack/echo/v4"
	"net/http/httptest"
	"net/http"
	"github.com/KKGo-Software-engineering/assessment-tax/postgres"
	"github.com/KKGo-Software-engineering/assessment-tax/calculateTax"
	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	// "io/ioutil"
	// // "strings"
	"bytes"
	"encoding/json"
	// "mime/multipart"
)

func TestCalculateExp01_TaxRefund(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	reqBody := IncomeData{
		TotalIncome: 220001,
		Wht: 7500,
		Allowances: []struct {
			AllowanceType string  `json:"allowanceType"`
			Amount        float64 `json:"amount"`
		}{
			{
				AllowanceType: "k-receipt",
				Amount: 0,
			},
		},
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/tax/calculation", bytes.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tax/calculation")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}

	expected := `{
		"tax": 0,
		"taxLevel": [
			{
				"level": "0-150,000",
				"tax": 0
			},
			{
				"level": "150,001-500,000",
				"tax": 0
			},
			{
				"level": "500,001-1,000,000",
				"tax": 0
			},
			{
				"level": "1,000,001-2,000,000",
				"tax": 0
			},
			{
				"level": "2,000,001 ขึ้นไป",
				"tax": 0
			}
		],
		"taxRefund": 6500
	}`

	handler := calculateTax.New(p)
	err = handler.HandleCalculateTaxData(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}
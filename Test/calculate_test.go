package Test

import (
	"testing"
	"github.com/labstack/echo/v4"
	"net/http/httptest"
	"net/http"
	"github.com/SukritPrem/assessment-tax/postgres"
	"github.com/SukritPrem/assessment-tax/calculateTax"
	// "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestCalculateExp02(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	reqBody := IncomeData{
		TotalIncome: 500000.0,
		Wht: 25000.0,
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
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tax/calculation")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	expected := `{
		"tax": 4000,
		"taxLevel": [
			{
				"level": "0-150,000",
				"tax": 0
			},
			{
				"level": "150,001-500,000",
				"tax": 29000
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
		"taxRefund": 0
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

func TestCalculateExp03_04(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	reqBody := IncomeData{
		TotalIncome: 500000.0,
		Wht: 0,
		Allowances: []struct {
			AllowanceType string  `json:"allowanceType"`
			Amount        float64 `json:"amount"`
		}{
			{
				AllowanceType: "donation",
				Amount: 200000.0,
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
		"tax": 19000,
		"taxLevel": [
			{
				"level": "0-150,000",
				"tax": 0
			},
			{
				"level": "150,001-500,000",
				"tax": 19000
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
		"taxRefund": 0
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

func TestCalculateExp07(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	reqBody := IncomeData{
		TotalIncome: 500000.0,
		Wht: 0,
		Allowances: []struct {
			AllowanceType string  `json:"allowanceType"`
			Amount        float64 `json:"amount"`
		}{
			{
				AllowanceType: "k-receipt",
				Amount: 200000.0,
			},
			{
				AllowanceType: "donation",
				Amount: 100000.0,
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
		"tax": 14000,
		"taxLevel": [
			{
				"level": "0-150,000",
				"tax": 0
			},
			{
				"level": "150,001-500,000",
				"tax": 14000
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
		"taxRefund": 0
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


type Request_amount struct {
  Amount float64 `json:"amount"`
}

func TestCalculateExp05(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	reqBody := Request_amount{
		Amount: 60000.0,
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bytes.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/deductions/personal")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	expected := `
				{		
    				"personalDeduction": 60000
				}
				`
	handler := calculateTax.New(p)
	err = handler.DeductionsPersonal(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}

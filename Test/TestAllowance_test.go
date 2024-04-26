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
	// "encoding/json"
)

func Test_DuplicateDonation(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 10,
		"wht": 0.0,
		"allowances": [
			{
			"allowanceType": "donation",
			"amount": 0.0
			},
			{
			"allowanceType": "donation",
			"amount": 0.0
			}
		]
	}`)
	req := httptest.NewRequest(http.MethodPost, "/tax/calculation", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tax/calculation")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}

	expected := `"Duplicate Donation or K-Receipt"`

	handler := calculateTax.New(p)
	err = handler.HandleCalculateTaxData(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}

func Test_DonationRandomIncomeThenMax(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 750000,
		"wht": 0.0,
		"allowances": [
			{
			"allowanceType": "donation",
			"amount": 500000
			}
		]
	}`)
	req := httptest.NewRequest(http.MethodPost, "/tax/calculation", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tax/calculation")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}

	expected := `{
		"tax": 48500,
		"taxLevel": [
			{
				"level": "0-150,000",
				"tax": 0
			},
			{
				"level": "150,001-500,000",
				"tax": 35000
			},
			{
				"level": "500,001-1,000,000",
				"tax": 13500
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

func Test_KreceiveRandomIncomeThenMax(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 1500000,
		"wht": 0.0,
		"allowances": [
			{
			"allowanceType": "k-receipt",
			"amount": 500000
			}
		]
	}`)
	req := httptest.NewRequest(http.MethodPost, "/tax/calculation", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tax/calculation")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}

	expected := `{
		"tax": 188000,
		"taxLevel": [
			{
				"level": "0-150,000",
				"tax": 0
			},
			{
				"level": "150,001-500,000",
				"tax": 35000
			},
			{
				"level": "500,001-1,000,000",
				"tax": 75000
			},
			{
				"level": "1,000,001-2,000,000",
				"tax": 78000
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


func Test_KreceiveAndDonationRandomAmount(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	// reqBody := IncomeData{
	// 	TotalIncome: 1500000,
	// 	Wht: 0.0,
	// 	Allowances: []struct {
	// 		AllowanceType string  `json:"allowanceType"`
	// 		Amount        float64 `json:"amount"`
	// 	}{
	// 		{
	// 			AllowanceType: "k-receipt",
	// 			Amount: 20000,
	// 		},
	// 		{
	// 			AllowanceType: "donation",
	// 			Amount: 20000,
	// 		},
	// 	},
	// }
	jsonBytes := []byte(`{
		"totalIncome": 1500000,
		"wht": 0.0,
		"allowances": [
			{
			"allowanceType": "k-receipt",
			"amount": 20000
			},
			{
			"allowanceType": "donation",
			"amount": 20000
			}
		]
	}`)
	// reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/tax/calculation", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tax/calculation")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}

	expected := `{
		"tax": 190000,
		"taxLevel": [
			{
				"level": "0-150,000",
				"tax": 0
			},
			{
				"level": "150,001-500,000",
				"tax": 35000
			},
			{
				"level": "500,001-1,000,000",
				"tax": 75000
			},
			{
				"level": "1,000,001-2,000,000",
				"tax": 80000
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
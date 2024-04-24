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
	// "strings"
	"bytes"
	"encoding/json"
)

func Test_WhtThenMax(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	reqBody := IncomeData{
		TotalIncome: 500000,
		Wht: 25000.1,
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

	expected := `"Wht is greater than TotalIncome * 0.05"`

	handler := calculateTax.New(p)
	err = handler.HandleCalculateTaxData(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest{
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}

func Test_WhtIsNeg(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	reqBody := IncomeData{
		TotalIncome: 500000,
		Wht: -1,
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

	expected := `"Wht Is Negative"`

	handler := calculateTax.New(p)
	err = handler.HandleCalculateTaxData(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest{
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}
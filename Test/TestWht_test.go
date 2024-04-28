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
)

func Test_WhtThenMax(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 500000.0,
		"wht": 500001,
		"allowances": [
			{
			"allowanceType": "donation",
			"amount": 0
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

	expected := `"Wht is greater than TotalIncome"`

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
	jsonBytes := []byte(`{
		"totalIncome": 500000.0,
		"wht": -25000.1,
		"allowances": [
			{
			"allowanceType": "donation",
			"amount": 0
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

	expected := `"Error: Wht checkValuefloat"`

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
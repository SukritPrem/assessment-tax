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


func Test_ValidatorReqHaveExtraField(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 500000.0,
		"wht": 25000.1,
		"allowances": [
			{
			"allowanceType": "donation",
			"amount": 0
			}
		],
		"Hi":1
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

	expected := `"Invalid JSON data Hi"`

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

func Test_ValidatorReqNotHaveWhtField(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 500000.0,
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

	expected := `"Error: Wht required"`

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

func Test_ValidatorReqNotHaveTotalIncomeField(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
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

	expected := `"Error: TotalIncome required"`

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

func Test_ValidatorReqNotHaveAllowancesField(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 220001,
		"wht": 7500
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

	expected := `"Error: Allowances required"`

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

func Test_ValidatorReqNotHaveAllowancesFieldButNotHaveAllowanceType(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 220001,
		"wht": 7500,
		"allowances": [
			{
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

	expected := `"Error: AllowanceType required"`

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

func Test_ValidatorReqNotHaveAllowancesFieldButNotHaveAllowanceAmount(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 220001,
		"wht": 7500,
		"allowances": [
			{
				"allowanceType": "donation"
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

	expected := `"Error: Amount required"`

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

func Test_ValidatorReqIsWantFloatButGotStringKeyWht(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 220001,
		"wht": "a",
		"allowances": [
			{
				"allowanceType": "donation"
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

	expected := `"json: cannot unmarshal string into Go struct field IncomeDataValidate.wht of type float64"`

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

func Test_ValidatorReqIsFloatMax(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome":  1.7976931348623197e+308,
		"wht": 1,
		"allowances": [
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

	expected := `"json: cannot unmarshal number 1.7976931348623197e+308 into Go value of type float64"`

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

func Test_ValidatorReqIsFloatMin(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome":  -1.7976931348623157e+308,
		"wht": 1,
		"allowances": [
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

	expected := `"Error: TotalIncome checkValuefloat"`

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
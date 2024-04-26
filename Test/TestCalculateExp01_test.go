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

func TestCalculateExp01(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 500000.0,
		"wht": 0.0,
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
	expected := `{"tax":29000,"taxLevel":[{"level":"0-150,000","tax":0},{"level":"150,001-500,000","tax":29000},{"level":"500,001-1,000,000","tax":0},{"level":"1,000,001-2,000,000","tax":0},{"level":"2,000,001 ขึ้นไป","tax":0}],
	"taxRefund": 0}`

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

func TestCalculateExp01_donationIsNeg(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 500000.0,
		"wht": 0.0,
		"allowances": [
			{
			"allowanceType": "donation",
			"amount": -1
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
	expected := `"Error: Amount checkValuefloat"`

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

func TestCalculateExp01_KreceiveIsNeg(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 500000.0,
		"wht": 0.0,
		"allowances": [
			{
			"allowanceType": "k-receipt",
			"amount": -1
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
	expected := `"Error: Amount checkValuefloat"`

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

func TestCalculateExp01_TotalIncomeIsZero(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 0,
		"wht": 0.0,
		"allowances": [
			{
			"allowanceType": "k-receipt",
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
		"taxRefund": 0
	}`

	handler := calculateTax.New(p)
	err = handler.HandleCalculateTaxData(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}

func TestCalculateExp01_TotalIncomeIsOneMilion(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 1000000.0,
		"wht": 0.0,
		"allowances": [
			{
			"allowanceType": "k-receipt",
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

	expected := `{
		"tax": 101000,
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
				"tax": 66000
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

func TestCalculateExp01_TotalIncomeIsTwoMilion(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 2000000.0,
		"wht": 0.0,
		"allowances": [
			{
			"allowanceType": "k-receipt",
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

	expected := `{
		"tax": 298000,
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
				"tax": 188000
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

func TestCalculateExp01_TotalIncomeIsThreeMilion(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 3000000.0,
		"wht": 0.0,
		"allowances": [
			{
			"allowanceType": "k-receipt",
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

	expected := `{
		"tax": 639000,
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
				"tax": 200000
			},
			{
				"level": "2,000,001 ขึ้นไป",
				"tax": 329000
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

func TestCalculateExp01_notHaveBody(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(``)
	req := httptest.NewRequest(http.MethodPost, "/tax/calculation", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/tax/calculation")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}

	expected := `"unexpected end of JSON input"`

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

func TestInvalidType_allowance(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
		"totalIncome": 3000000.0,
		"wht": 0.0,
		"allowances": [
			{
			"allowanceType": "don",
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

	expected := `"Error: AllowanceType checkType"`

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
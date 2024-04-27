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

func TestCalculateExp05_adminUpdateValueMinDeductions_personal(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	reqBody := Request_amount{
		Amount: 9999,
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
				"Amount is not in range"
				`
	handler := calculateTax.New(p)
	err = handler.DeductionsPersonal(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}

func TestCalculateExp05_adminUpdateValueThenMaxValueDeductions_personal(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	reqBody := Request_amount{
		Amount: 100001,
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
				"Amount is not in range"
				`
	handler := calculateTax.New(p)
	err = handler.DeductionsPersonal(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}


func TestCalculateExp08_adminUpdateValueMinDeductionsKrecipt(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	reqBody := Request_amount{
		Amount: 0,
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", bytes.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/deductions/k-receipt")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	expected := `
				"Amount is not in range"
				`
	handler := calculateTax.New(p)
	err = handler.DeductionsKReceipt(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}

func TestCalculateExp08_adminUpdateValueThenMaxValueDeductions_personal(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	reqBody := Request_amount{
		Amount: 100001,
	}
	reqJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", bytes.NewReader(reqJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/deductions/k-receipt")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	expected := `
				"Amount is not in range"
				`
	handler := calculateTax.New(p)
	err = handler.DeductionsKReceipt(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}

func TestCalculateExp08_adminUpdateDuplicateKeyAmount_kReceipt(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
			"amount": 60000.0,
			"amount": 60000.0
	}`)
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/deductions/k-receipt")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	expected := `
				"Duplicate key found"
				`
	handler := calculateTax.New(p)
	err = handler.DeductionsKReceipt(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}

func TestCalculateExp08_HaveAnotherkey_kReceipt(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
			"Hi": 1
	}`)
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/deductions/k-receipt")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	expected := `
				"Invalid JSON data Hi"
				`
	handler := calculateTax.New(p)
	err = handler.DeductionsKReceipt(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}


func TestCalculateExp08_ValueAmountIsNeg_kReceipt(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
			"amount": -1
	}`)
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/deductions/k-receipt")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	expected := `"Error:Amount validateValueFloat"`
	handler := calculateTax.New(p)
	err = handler.DeductionsKReceipt(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}

func TestCalculateExp08_adminUpdateDuplicateKeyAmount_personal(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
			"amount": 60000.0,
			"amount": 60000.0
	}`)
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/deductions/k-receipt")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	expected := `
				"Duplicate key found"
				`
	handler := calculateTax.New(p)
	err = handler.DeductionsKReceipt(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}

func TestCalculateExp08_HaveAnotherkey_personal(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
			"Hi": 1
	}`)
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/deductions/k-receipt")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	expected := `
				"Invalid JSON data Hi"
				`
	handler := calculateTax.New(p)
	err = handler.DeductionsKReceipt(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}


func TestCalculateExp08_ValueAmountIsNeg_personal(t *testing.T) {
	// Create a new Postgres instance
	e := echo.New()
	jsonBytes := []byte(`{
			"amount": -1
	}`)
	req := httptest.NewRequest(http.MethodPost, "/admin/deductions/k-receipt", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/deductions/k-receipt")
	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	expected := `"Error:Amount validateValueFloat"`
	handler := calculateTax.New(p)
	err = handler.DeductionsKReceipt(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
}
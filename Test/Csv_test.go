
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
	"io/ioutil"
	// "strings"
	"bytes"
	"mime/multipart"
)

func TestHandleIncomeDataCSV_errorWhtThenMax(t *testing.T) {
	fileContent, err := ioutil.ReadFile("./fileCsv/WhtThenMax.csv")
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("file")

	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)
	part, err := multipartWriter.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	_, err = part.Write(fileContent)
	if err != nil {
		t.Fatal(err)
	}
	multipartWriter.Close()

	req.Header.Set(echo.HeaderContentType, multipartWriter.FormDataContentType())
	req.Body = ioutil.NopCloser(body)

	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	expected := `"Wht is greater than TotalIncome"`
	handler := calculateTax.New(p)
	err = handler.HandleIncomeDataCSV(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest{
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
	// Here you can check the response body as well
}

func TestHandleIncomeDataCSV_TrueFileCsv(t *testing.T) {
	fileContent, err := ioutil.ReadFile("./fileCsv/True.csv")
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("file")

	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)
	part, err := multipartWriter.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	_, err = part.Write(fileContent)
	if err != nil {
		t.Fatal(err)
	}
	multipartWriter.Close()

	req.Header.Set(echo.HeaderContentType, multipartWriter.FormDataContentType())
	req.Body = ioutil.NopCloser(body)

	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	expected := `{
		"taxes": [
			{
				"totalIncome": 500000,
				"tax": 29000
			},
			{
				"totalIncome": 600000,
				"tax": 8000
			},
			{
				"totalIncome": 750000,
				"tax": 23750
			}
		]
	}`
	handler := calculateTax.New(p)
	err = handler.HandleIncomeDataCSV(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusOK{
		t.Errorf("Expected status 200, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
	// Here you can check the response body as well
}

func TestHandleIncomeDataCSV_Empty(t *testing.T) {
	fileContent, err := ioutil.ReadFile("./fileCsv/Empty.csv")
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("file")

	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)
	part, err := multipartWriter.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	_, err = part.Write(fileContent)
	if err != nil {
		t.Fatal(err)
	}
	multipartWriter.Close()

	req.Header.Set(echo.HeaderContentType, multipartWriter.FormDataContentType())
	req.Body = ioutil.NopCloser(body)

	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	expected := `"Error no data"`
	handler := calculateTax.New(p)
	err = handler.HandleIncomeDataCSV(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest{
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
	// Here you can check the response body as well
}

func TestHandleIncomeDataCSV_ErrorFormatHeader(t *testing.T) {
	fileContent, err := ioutil.ReadFile("./fileCsv/ErrorFormatHeader.csv")
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("file")

	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)
	part, err := multipartWriter.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	_, err = part.Write(fileContent)
	if err != nil {
		t.Fatal(err)
	}
	multipartWriter.Close()

	req.Header.Set(echo.HeaderContentType, multipartWriter.FormDataContentType())
	req.Body = ioutil.NopCloser(body)

	p, err := postgres.New();
	if err != nil {
		panic(err)
	}

	expected := `"Error format header: totalIncome,wht,donation"`
	handler := calculateTax.New(p)
	err = handler.HandleIncomeDataCSV(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest{
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
	// Here you can check the response body as well
}

func TestHandleIncomeDataCSV_ErrorParseDonate(t *testing.T) {
	fileContent, err := ioutil.ReadFile("./fileCsv/ErrorParseDonate.csv")
	if err != nil {
		t.Fatal(err)
	}
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("file")

	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)
	part, err := multipartWriter.CreateFormFile("file", "test.csv")
	if err != nil {
		t.Fatal(err)
	}
	_, err = part.Write(fileContent)
	if err != nil {
		t.Fatal(err)
	}
	multipartWriter.Close()

	req.Header.Set(echo.HeaderContentType, multipartWriter.FormDataContentType())
	req.Body = ioutil.NopCloser(body)

	p, err := postgres.New();
	if err != nil {
		panic(err)
	}
	
	expected := `"Error can't parse donation amount"`
	handler := calculateTax.New(p)
	err = handler.HandleIncomeDataCSV(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusBadRequest{
		t.Errorf("Expected status 400, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
	// Here you can check the response body as well
}
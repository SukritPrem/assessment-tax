
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
	"io/ioutil"
	// "strings"
	"bytes"
	"mime/multipart"
)

func TestHandleIncomeDataCSV(t *testing.T) {
	fileContent, err := ioutil.ReadFile("../file.csv")
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
		"taxs": [
			{
				"totalIncome": 500000,
				"tax": 29000
			},
			{
				"totalIncome": 600000,
				"tax": 38000
			},
			{
				"totalIncome": 750000,
				"tax": 61250
			}
		]
	}`
	handler := calculateTax.New(p)
	err = handler.HandleIncomeDataCSV(c)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %v", rec.Code)
	}

	require.JSONEq(t, expected, rec.Body.String())
	// Here you can check the response body as well
}

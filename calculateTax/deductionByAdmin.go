package calculateTax

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"fmt"
  "io/ioutil"
  "strings"
  "encoding/json"
  "github.com/go-playground/validator/v10"
)

type Request_amount struct {
  Amount float64 `json:"amount" `
}

type Request_amount_new struct {
  Amount *float64 `json:"amount" validate:"required,maxfloat"`
}

type Response_amount_personalDeduction struct {
  Amount float64 `json:"personalDeduction"`
}

type Response_amount_kReceipt struct {
  Amount float64 `json:"k-receipt"`
}

func (h *Handler) DeductionsPersonal(c echo.Context) error {
  a := new(Request_amount)
  err := c.Bind(&a)
  if err != nil {
    return c.JSON(http.StatusBadRequest, "Invalid JSON data")
  }
  fmt.Println(a.Amount)
  if (a.Amount <= 10000 || a.Amount > 100000){
    return c.JSON(http.StatusBadRequest, "Amount is not in range")
  }
  _, err = h.store.UpdateAmountByTaxType("personalDeduction",a.Amount)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, "Not found")
  }
  r := &Response_amount_personalDeduction{
    Amount: a.Amount,
  }
  return c.JSON(http.StatusOK, r)
}

func (h *Handler) DeductionsKReceipt(c echo.Context) error {
  a := new(Request_amount)
  err := c.Bind(&a)
  if err != nil {
    return c.JSON(http.StatusBadRequest, "Invalid JSON data")
  }
  if (a.Amount <= 0 || a.Amount > 100000){
    return c.JSON(http.StatusBadRequest, "Amount is not in range")
  }
  _, err = h.store.UpdateAmountByTaxType("k-receipt",a.Amount)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, "Not found")
  }

  r := &Response_amount_kReceipt{
    Amount: a.Amount,
  }
  return c.JSON(http.StatusOK, r)
}

func (h *Handler) Deductions(c echo.Context) error {
  // taxType := c.Param("taxType")
  body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
  }
	defer c.Request().Body.Close()
  var jsonBytes []byte
  jsonBytes = body
  err = check(json.NewDecoder(strings.NewReader(string(jsonBytes))), nil, dupErr)
	if err != nil {
		fmt.Println("Error:", err)
	}
  err = validateKeyReqAdmin(body)
  if err != nil {
    return c.JSON(http.StatusBadRequest, err.Error())
  }
  amount := new(Request_amount_new)
  validate := validator.New()
  validate.RegisterValidation("maxfloat", validateValueFloat)
	err = validate.Struct(amount)
  	if err != nil {
		errors := err.(validator.ValidationErrors)
		allErrors := errors.Error()
		for _, e := range errors {
			allErrors = allErrors + e.Field() + " " + e.Tag() + "\n"
			fmt.Println(e.Field(), e.Tag())
		}
	}
	fmt.Println(*amount.Amount)
  
  return c.JSON(http.StatusOK, "Deductions")
}
package calculateTax

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"fmt"
)
type Request_amount struct {
  Amount float64 `json:"amount"`
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
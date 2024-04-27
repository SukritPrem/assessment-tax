package calculateTax

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"fmt"
	"io/ioutil"
)
type Reponse_csv struct {
  Taxs []TotalIncomeAndTax `json:"taxs"`
}

type TotalIncomeAndTax struct {
  TotalIncome float64 `json:"totalIncome"`
  Tax float64 `json:"tax"`
}

func (h *Handler) HandleIncomeDataCSV(c echo.Context) error {

  data,err := OpenfileAndGetData(c)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, err.Error())
  }
  
  personalDeduction, k_receipt, err := GetValuepersonalAndKreceipt(h)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, err.Error())
  }

  err, r := validateCSV(data, personalDeduction, k_receipt)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, err.Error()) 
  }
  return c.JSON(http.StatusOK, r)
}

func OpenfileAndGetData(c echo.Context) ([]byte,error){
  file, err := c.FormFile("file")
  if err != nil {
    return nil,err
  }
  src, err := file.Open()
  if err != nil {
    return nil,fmt.Errorf("Error opening file")
  }
  defer src.Close()
  data, err := ioutil.ReadAll(src)
  if err != nil {
    return nil,fmt.Errorf("Error opening file")
  }
  return data,nil
}

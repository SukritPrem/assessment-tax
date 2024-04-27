package calculateTax


import (
	"github.com/labstack/echo/v4"
	"fmt"
  "math"
	"net/http"
  "io/ioutil"
	// "strconv"
	// "strings"
  // "encoding/json"
  "github.com/go-playground/validator/v10"
)
type Handler struct {
	store Storer
}

type Storer interface {
	GetAmountByTaxType(v string) (float64, error)
  UpdateAmountByTaxType(v string,a float64) (float64, error)
}


func New(db Storer) *Handler {
	return &Handler{store: db}
}

type IncomeData struct {
  TotalIncome float64 `json:"totalIncome"`
  Wht        float64 `json:"wht"`
  Allowances []Allowance `json:"allowances"`
}

type Allowance struct {
  AllowanceType string `json:"allowanceType"`
  Amount        float64 `json:"amount"`
}

type taxlevel struct {
  level int
  tax  float64
  rate_min float64
  rate_max float64
  pay float64
}

type Response_tax struct {
  Tax_sum float64 `json:"tax"`
  Tax_level []LevelWithTax  `json:"taxLevel"`
  Tax_refund float64 `json:"taxRefund"`
}

type LevelWithTax struct {
  Level string `json:"level"`
  Tax float64 `json:"tax"`
}

func (h *Handler) HandleCalculateTaxData(c echo.Context) error {

  body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
  }
	defer c.Request().Body.Close()

  incomeData, err := ValidateReqClientTax(body)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, err.Error())
  }

  personalDeduction, k_receipt, err := GetValuepersonalAndKreceipt(h)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, err.Error())
  }


  sum_tax,taxlevels,err := CalculateAndGetSumTax(&incomeData,personalDeduction,k_receipt)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, err.Error())
  }

  taxRefund := 0.0
  if(sum_tax < 0){
    taxRefund = math.Abs(sum_tax)
    sum_tax = 0
  }

  return c.JSON(http.StatusOK, ReponseSumTaxWithTaxLevel(taxlevels,sum_tax,taxRefund))
}


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


//handler CalculateTaxData
func GetValuepersonalAndKreceipt(h *Handler) (float64,float64,error){
  personalDeduction, err := h.store.GetAmountByTaxType("personalDeduction")
  if(err != nil){
    return 0,0,err
  }
  k_receipt, err := h.store.GetAmountByTaxType("k-receipt")
  if(err != nil){
    return 0,0,err
  }
  return personalDeduction,k_receipt,nil
}


func ValidateReqClientTax(body []byte) (IncomeData,error){
  var incomeData IncomeData

  err := validateKey(body)
  if err != nil {
    return incomeData,err
  }

  incomeData, err = validateValueByStuct(body)
  if err != nil {
    if errors, ok := err.(validator.ValidationErrors); ok {
      allErrors := "Error: "
      for _, e := range errors {
        allErrors += e.Field() + " "+ e.Tag()
        return incomeData,fmt.Errorf(allErrors)
      }
    }
    return incomeData,err
  }

  err = checkErrorIncomeData(&incomeData)
  if(err != nil){
    return incomeData,err
  }

  return incomeData,nil
}

func CalculateAndGetSumTax(incomeData *IncomeData, personalDeduction float64, k_receipt float64) (float64,[]taxlevel,error){
  incomeData.TotalIncome = incomeData.TotalIncome - personalDeduction
  err := IncomeDataDecrease(incomeData,k_receipt)
  if(err != nil){
    return 0,nil,err
  }

  taxlevels := CalculateTaxLevelWithNetIncomeData(incomeData)
  sum_tax := sumAllTaxLevel(taxlevels)
  sum_tax = sum_tax - incomeData.Wht
  return sum_tax,taxlevels,nil
}
///////////////////////////
//handler CSV

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

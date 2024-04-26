package calculateTax


import (
	"github.com/labstack/echo/v4"
	// "fmt"
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
  var incomeData IncomeData
  body, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
  }
	defer c.Request().Body.Close()
  err = validateKey(body)
  if err != nil {
    return c.JSON(http.StatusBadRequest, err.Error())
  }
  incomeData, err = validateValueByStuct(body)
  if err != nil {
    errors := err.(validator.ValidationErrors)
    allErrors := "Error: "
    for _, e := range errors {
      allErrors += e.Field() + " "+ e.Tag()
			return c.JSON(http.StatusBadRequest, allErrors)
		}
  }
  personalDeduction, err := h.store.GetAmountByTaxType("personalDeduction")
  if(err != nil){
    return c.JSON(http.StatusBadRequest, "Not found")
  }
  k_receipt, err := h.store.GetAmountByTaxType("k-receipt")
  if(err != nil){
    return c.JSON(http.StatusBadRequest, "Not found")
  }

  err = checkErrorIncomeData(&incomeData)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, err.Error())
  }

  incomeData.TotalIncome = incomeData.TotalIncome - personalDeduction
  err = IncomeDataDecrease(&incomeData,k_receipt)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, err.Error())
  }
  taxlevels := CalculateTaxLevelWithNetIncomeData(&incomeData)
  sum_tax := sumAllTaxLevel(taxlevels)
  sum_tax = sum_tax - incomeData.Wht
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
  file, err := c.FormFile("file")
  // r := &Reponse_csv{}
  req := []IncomeData{}
  if err != nil {
    return c.JSON(http.StatusBadRequest, "Error opening file")
  }
  src, err := file.Open()
  if err != nil {
    return c.JSON(http.StatusBadRequest, "Error opening file")
  }
  defer src.Close()

  personalDeduction, err := h.store.GetAmountByTaxType("personalDeduction")
  if(err != nil){
    return c.JSON(http.StatusBadRequest, "Not found")
  }
  k_receipt, err := h.store.GetAmountByTaxType("k-receipt")
  if(err != nil){
    return c.JSON(http.StatusBadRequest, "Not found")
  }
  data, err := ioutil.ReadAll(src)
  if err != nil {
    return c.JSON(http.StatusBadRequest, "Error opening file")
  }
  // fmt.Println(string(data))
  err, r := validateCSV(data, &req, personalDeduction, k_receipt)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, err.Error()) 
  }
  return c.JSON(http.StatusOK, r)
}
package calculateTax


import (
	"github.com/labstack/echo/v4"
	// "fmt"
  "math"
	"net/http"
  "io/ioutil"
	// "strconv"
	// "strings"
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
  TotalIncome float64 `json:"TotalIncome"`
  Wht        float64 `json:"wht"`
  Allowances []struct {
    AllowanceType string  `json:"allowanceType"`
    Amount        float64 `json:"amount"`
  } `json:"allowances"`
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
  err := c.Bind(&incomeData)
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON data")
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
  // fmt.Println(incomeData.TotalIncome)
  err = IncomeDataDecrease(&incomeData,k_receipt)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, err.Error())
  }
  // fmt.Printf("k_receipt: %f\n",k_receipt)
  // fmt.Printf("totalincome: %f",incomeData.TotalIncome)
  // fmt.Println(incomeData.TotalIncome)
  taxlevels := CalculateTaxLevelWithNetIncomeData(&incomeData)
  sum_tax := sumAllTaxLevel(taxlevels)
  // fmt.Println(sum_tax)

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
  // Read the file
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
package calculateTax


import (
	"github.com/labstack/echo/v4"
	"fmt"
  "math"
	"net/http"
  "io/ioutil"
	"strconv"
	"strings"
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
}

type LevelWithTax struct {
  Level string `json:"level"`
  Tax float64 `json:"tax"`
}

func roundFloat(val float64, precision uint) float64 {
    ratio := math.Pow(10, float64(precision))
    return math.Round(val*ratio) / ratio
}

func IncomeDataDecrease(incomeData *IncomeData,k_receipt float64) error {
  for _, allowance := range incomeData.Allowances {
      if(allowance.AllowanceType == "donation") {
        if(allowance.Amount < 100000 && allowance.Amount > 0) {
          incomeData.TotalIncome = incomeData.TotalIncome - allowance.Amount  
        } else if allowance.Amount >= 100000 {
          incomeData.TotalIncome = incomeData.TotalIncome - 100000 
        }
      } else if(allowance.AllowanceType == "k-receipt") {
        if(allowance.Amount > k_receipt) {
          incomeData.TotalIncome = incomeData.TotalIncome - k_receipt
        } else if(allowance.Amount <= k_receipt && allowance.Amount > 0) {
          incomeData.TotalIncome = incomeData.TotalIncome - allowance.Amount
        } 
      }
  }
  return nil
}

func CalculateTaxLevelWithNetIncomeData(incomeData *IncomeData) []taxlevel{
  taxlevels := []taxlevel{
    {1, 0, 0, 150000,0},
    {2, 0.1, 150001, 500000,0},
    {3, 0.15,500001, 1000000,0},
    {4, 0.2,1000001, 2000000,0},
    {5, 0.35,2000001, 2000001,0},
  }

  for i := 0; i < len(taxlevels); i++ {
    if incomeData.TotalIncome >= taxlevels[i].rate_min && incomeData.TotalIncome <= taxlevels[i].rate_max && i != 4 {
      taxlevels[i].pay = roundFloat((incomeData.TotalIncome - taxlevels[i].rate_min) * taxlevels[i].tax,0)
    }
    if i == 4 && incomeData.TotalIncome >= taxlevels[i].rate_min {
      taxlevels[i].pay = roundFloat((incomeData.TotalIncome - taxlevels[i].rate_min) * taxlevels[i].tax,0)
    }
  }
  return taxlevels
}

func sumAllTaxLevel(taxlevels []taxlevel) float64 {
  sum_tax := 0.0
  for i := 0; i < len(taxlevels); i++ {
      sum_tax = sum_tax + taxlevels[i].pay
  }
  return sum_tax
}

func ReponseSumTaxWithTaxLevel(taxlevels []taxlevel,sum_tax float64) Response_tax {
  r := &Response_tax{
    Tax_sum: sum_tax,
    Tax_level: []LevelWithTax{
      {
        Level: "0-150,000",
        Tax: taxlevels[0].pay,
      },
      {
        Level: "150,001-500,000",
        Tax: taxlevels[1].pay,
      },
      {
        Level: "500,001-1,000,000",
        Tax: taxlevels[2].pay,
      },
      {
        Level: "1,000,001-2,000,000",
        Tax: taxlevels[3].pay,
      },
      {
        Level: "2,000,001 ขึ้นไป",
        Tax: taxlevels[4].pay,
      },
    },
  }
  return *r
}
func (h *Handler) HandleIncomeData(c echo.Context) error {
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
  if(incomeData.TotalIncome < 0){
    return c.JSON(http.StatusBadRequest, "TotalIncome Is Negative")
  }
  incomeData.TotalIncome = incomeData.TotalIncome - personalDeduction
  IncomeDataDecrease(&incomeData,k_receipt)

  taxlevels := CalculateTaxLevelWithNetIncomeData(&incomeData)
  sum_tax := sumAllTaxLevel(taxlevels)
  if(incomeData.Wht > 0){
    sum_tax = sum_tax - incomeData.Wht
  }
  return c.JSON(http.StatusOK, ReponseSumTaxWithTaxLevel(taxlevels,sum_tax))
}

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

func validateCSV(data []byte) bool {

    lines := strings.Split(string(data), "\n")

    // Check for header row
    if lines[0] != "totalIncome,wht,donation" {
        fmt.Println("Warning: Header row not found. Assuming data starts from the first line.")
    }

    // Validate each line
    for _, line := range lines[1:] {
        values := strings.Split(line, ",")

        // Check for correct number of values
        if len(values) != 3 {
            fmt.Println("Error: Invalid line format: {line}")
            return false
        }

        // Check if all values are numbers
        for _, value := range values {
            if _, err := strconv.ParseFloat(value, 64); err != nil {
                fmt.Println("Error: Invalid value in line: {line}")
                return false
            }
        }
    }

    return true
}

type reponse_csv struct {
  Taxs []totalIncomeAndTax `json:"taxs"`
}

type totalIncomeAndTax struct {
  TotalIncome float64 `json:"totalIncome"`
  Tax float64 `json:"tax"`
}
func HandleIncomeDataCSV(c echo.Context) error {
  file, err := c.FormFile("file")
  if err != nil {
    return c.JSON(http.StatusBadRequest, "Error opening file")
  }
  src, err := file.Open()
  if err != nil {
    return c.JSON(http.StatusBadRequest, "Error opening file")
  }
  defer src.Close()
  // Read the file

  data, err := ioutil.ReadAll(src)
  if err != nil {
    return c.JSON(http.StatusBadRequest, "Error opening file")
  }
  fmt.Println(string(data))
  if validateCSV(data) {
      fmt.Println("CSV format is valid.")
  } else {
      fmt.Println("CSV format is invalid.")
  }
  return c.JSON(http.StatusOK, "CSV")
}
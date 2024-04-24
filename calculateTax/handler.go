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
  if(incomeData.TotalIncome < 0){
    return c.JSON(http.StatusBadRequest, "TotalIncome Is Negative")
  }
  if(incomeData.Wht < 0){
    return c.JSON(http.StatusBadRequest, "Wht Is Negative")
  }
  if(incomeData.Wht > incomeData.TotalIncome * 0.05){
    return c.JSON(http.StatusBadRequest, "Wht is greater than TotalIncome * 0.05")
  }

  incomeData.TotalIncome = incomeData.TotalIncome - personalDeduction
  // fmt.Println(incomeData.TotalIncome)
  err = IncomeDataDecrease(&incomeData,k_receipt)
  if(err != nil){
    return c.JSON(http.StatusBadRequest, err.Error())
  }
  fmt.Println(incomeData.TotalIncome)
  taxlevels := CalculateTaxLevelWithNetIncomeData(&incomeData)
  sum_tax := sumAllTaxLevel(taxlevels)
  fmt.Println(sum_tax)

  sum_tax = sum_tax - incomeData.Wht
  taxRefund := 0.0
  if(sum_tax < 0){
    taxRefund = math.Abs(sum_tax)
  }
  
  return c.JSON(http.StatusOK, ReponseSumTaxWithTaxLevel(taxlevels,sum_tax,taxRefund))
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


type Reponse_csv struct {
  Taxs []TotalIncomeAndTax `json:"taxs"`
}

type TotalIncomeAndTax struct {
  TotalIncome float64 `json:"totalIncome"`
  Tax float64 `json:"tax"`
}

func validateCSV(data []byte, req *[]IncomeData, personalDeduction float64, k_receipt float64) (bool, Reponse_csv) {
    r := Reponse_csv{}
    // Reponse_csv.Taxs = append(Reponse_csv.Taxs, TotalIncomeAndTax{})
    result := TotalIncomeAndTax{}
    dataOneLine := IncomeData{}
    data = []byte(strings.Replace(string(data), "\r", "", -1))
    lines := strings.Split(string(data), "\n")

    lines[0] = strings.Replace(lines[0], "\r", "", -1)
    if(strings.Compare(lines[0],"totalIncome,wht,donation") != 0){
      fmt.Println("Header row not found. Assuming data starts from the first line.")
    }

    // Validate each line
    var err error
    var num float64
    for _, line := range lines[1:] {
        fmt.Println(line)
        values := strings.Split(line, ",")
        for i := 0; i < len(values); i++ {
          _, err = strconv.ParseFloat(values[i], 64)
          if(err != nil){
            fmt.Printf("Error in line %d: %s\n", i, values[i])
          } 
          if(i == 0){
            dataOneLine.TotalIncome, err = strconv.ParseFloat(values[0], 64)
            if(err != nil){
              fmt.Printf("Error in line %d: %s\n", i, values[i])
            }
          } else if(i == 1){
            dataOneLine.Wht, err = strconv.ParseFloat(values[1], 64)
            if(err != nil){
              fmt.Printf("Error in line %d: %s\n", i, values[i])
            }
          } else if(i == 2){
            num, err = strconv.ParseFloat(values[2], 64)
            if(err != nil){
              fmt.Printf("Error in line %d: %s\n", i, values[i])
            }
            dataOneLine.Allowances = append(dataOneLine.Allowances, struct {
              AllowanceType string  `json:"allowanceType"`
              Amount        float64 `json:"amount"`
            }{"donation", num})
          }
        }
        result.TotalIncome = dataOneLine.TotalIncome
        dataOneLine.TotalIncome = dataOneLine.TotalIncome - personalDeduction
        fmt.Printf("totalincome: %f\n",dataOneLine.TotalIncome)
        fmt.Printf("personalDeduction: %f\n",personalDeduction)
        IncomeDataDecrease(&dataOneLine,k_receipt)
        fmt.Printf("totalincome: %f",dataOneLine.TotalIncome)
        taxlevels := CalculateTaxLevelWithNetIncomeData(&dataOneLine)
        sum_tax := sumAllTaxLevel(taxlevels)
        result.Tax = sum_tax
        r.Taxs = append(r.Taxs, result)
        
    }

    return true, r
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
  check, r := validateCSV(data, &req, personalDeduction, k_receipt)
  if(check == false){
    return c.JSON(http.StatusBadRequest, "Error in CSV file")
  }
  return c.JSON(http.StatusOK, r)
}
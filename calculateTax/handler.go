package calculateTax


import (
	"github.com/labstack/echo/v4"
	"fmt"
	"net/http"
)
type Handler struct {
	store Storer
}

type Storer interface {
	GetPersonalDeduction() (int, error)
}


func New(db Storer) *Handler {
	return &Handler{store: db}
}

type IncomeData struct {
  TotalIncome float64 `json:"totalIncome"`
  Wht        float64 `json:"wht"`
  Allowances []struct {
    AllowanceType string  `json:"allowanceType"`
    Amount        float64 `json:"amount"`
  } `json:"allowances"`
}

func (h *Handler) HandleIncomeData(c echo.Context) error {
  var incomeData IncomeData
  // Bind the JSON request body to the IncomeData struct
  err := c.Bind(&incomeData)
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON data")
  }
  fmt.Println(h.store.GetPersonalDeduction())
  // Process the income data (logic omitted for brevity)
//   fmt.Println("Received Income Data:")
//   fmt.Printf("  Total Income: %.2f\n", incomeData.TotalIncome)
//   fmt.Printf("  Withholding Tax: %.2f\n", incomeData.Wht)
//   fmt.Println("  Allowances:")
//   for _, allowance := range incomeData.Allowances {
//     fmt.Printf("    - Type: %s, Amount: %.2f\n", allowance.AllowanceType, allowance.Amount)
//   }

  return c.JSON(http.StatusOK, incomeData)
}
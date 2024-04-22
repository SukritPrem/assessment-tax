package calculateTax


import (
	"github.com/labstack/echo/v4"
	// "fmt"
	"net/http"
)
type Handler struct {
	store Storer
}

type Storer interface {
	GetPersonalDeduction() (float64, error)
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

type reposne_tax struct {
  tax_sum float64 `json:"tax"`
  tax_level []LevelWithTax `josn:"taxLevel"`
}

type LevelWithTax struct {
  level string
  tax float64
}

func (h *Handler) HandleIncomeData(c echo.Context) error {
  var incomeData IncomeData
  taxlevels := []taxlevel{}
  taxlevels = append(taxlevels, taxlevel{1, 0, 0, 150000,0})
  taxlevels = append(taxlevels, taxlevel{2, 0.1, 150001, 500000,0})
  taxlevels = append(taxlevels, taxlevel{3, 0.15,500001, 1000000,0})
  taxlevels = append(taxlevels, taxlevel{4, 0.2,1000001, 2000000,0})
  taxlevels = append(taxlevels, taxlevel{5, 0.35,2000001, 2000001,0})
  // Bind the JSON request body to the IncomeData struct
  err := c.Bind(&incomeData)
  if err != nil {
    return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON data")
  }
  //when connect database
  // fmt.Println(h.store.GetPersonalDeduction())
  personalDeduction := 60000.0
  k_receipt := 50000.0
  // Process the income data (logic omitted for brevity)
  //   fmt.Println("Received Income Data:")
  //   fmt.Printf("  Total Income: %.2f\n", incomeData.TotalIncome)
  //   fmt.Printf("  Withholding Tax: %.2f\n", incomeData.Wht)
  //   fmt.Println("  Allowances:")
  //   for _, allowance := range incomeData.Allowances {
  //     fmt.Printf("    - Type: %s, Amount: %.2f\n", allowance.AllowanceType, allowance.Amount)
  //   }

    incomeData.TotalIncome = incomeData.TotalIncome - personalDeduction

    for _, allowance := range incomeData.Allowances {
      if(allowance.AllowanceType == "donation") {
        if(allowance.Amount < 100000 ) {
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

    for i := 0; i < len(taxlevels); i++ {
      if incomeData.TotalIncome >= taxlevels[i].rate_min && incomeData.TotalIncome <= taxlevels[i].rate_max && i != 4 {
        taxlevels[i].pay = incomeData.TotalIncome * taxlevels[i].tax
      }
      if i == 4 && incomeData.TotalIncome >= taxlevels[i].rate_min {
        taxlevels[i].pay = incomeData.TotalIncome * taxlevels[i].tax
      }
    }
    
    sum_tax := 0.0
    if(incomeData.Wht > 0){
      for i := 0; i < len(taxlevels); i++ {
        sum_tax = sum_tax + taxlevels[i].pay
      }
      sum_tax = sum_tax - incomeData.Wht
    }
    response_tax := reposne_tax{}

    response_tax.tax_sum = sum_tax
    for i := 0; i < len(taxlevels); i++ {
		switch i {
		case 0:
			response_tax.tax_level = append(response_tax.tax_level, LevelWithTax{"0-150,000", taxlevels[i].pay})
		case 1:
			response_tax.tax_level = append(response_tax.tax_level, LevelWithTax{"150,001-500,000", taxlevels[i].pay})
		case 2:
			response_tax.tax_level = append(response_tax.tax_level, LevelWithTax{"500,001-1,000,000", taxlevels[i].pay})
		case 3:
			response_tax.tax_level = append(response_tax.tax_level, LevelWithTax{"1,000,001-2,000,000", taxlevels[i].pay})
		case 4:
			response_tax.tax_level = append(response_tax.tax_level, LevelWithTax{"2,000,001 ขึ้นไป", taxlevels[i].pay})
		}
	}


  return c.JSON(http.StatusOK, response_tax)
}
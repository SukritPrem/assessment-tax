package calculateTax

import (
	  "math"
    "errors"
)

func roundFloat(val float64, precision uint) float64 {
    ratio := math.Pow(10, float64(precision))
    return math.Round(val*ratio) / ratio
}

func IncomeDataDecrease(incomeData *IncomeData,k_receipt float64) error {
  for _, allowance := range incomeData.Allowances {
      if(allowance.AllowanceType == "donation") {
        if(allowance.Amount < 100000 && allowance.Amount >= 0) {
          incomeData.TotalIncome = incomeData.TotalIncome - allowance.Amount  
        } else if allowance.Amount >= 100000 {
          incomeData.TotalIncome = incomeData.TotalIncome - 100000 
        } else {
          return errors.New("Donation is negative")
        }
      } else if(allowance.AllowanceType == "k-receipt") {
        if(allowance.Amount >= k_receipt) {
          incomeData.TotalIncome = incomeData.TotalIncome - k_receipt
        } else if(allowance.Amount < k_receipt && allowance.Amount >= 0) {
          incomeData.TotalIncome = incomeData.TotalIncome - allowance.Amount
        } else {
          return errors.New("k-receipt is negative")
        }
      }
  }
  return nil
}

func CalculateTaxLevelWithNetIncomeData(incomeData *IncomeData) []taxlevel{
  taxlevels := []taxlevel{
    {1, 0, 0, 150000,0},
    {2, 0.1, 150000, 500000,0},
    {3, 0.15,500000, 1000000,0},
    {4, 0.2,1000000, 2000000,0},
    {5, 0.35,2000000, 2000000,0},
  }

  for i := 0; i < len(taxlevels); i++ {
    if incomeData.TotalIncome >= taxlevels[i].rate_min + 1 && incomeData.TotalIncome <= taxlevels[i].rate_max && i != 4 {
        taxlevels[i].pay = roundFloat((incomeData.TotalIncome - taxlevels[i].rate_min) * taxlevels[i].tax,0)
        break;
    } else {
        taxlevels[i].pay = roundFloat((taxlevels[i].rate_max - taxlevels[i].rate_min) * taxlevels[i].tax,0)
    }
    if i == 4 && incomeData.TotalIncome >= taxlevels[i].rate_min + 1 {
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

func ReponseSumTaxWithTaxLevel(taxlevels []taxlevel,sum_tax float64,taxRefund float64) Response_tax {
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
    Tax_refund: taxRefund,
  }
  return *r
}

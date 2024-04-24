package calculateTax

import (
	  "math"
)

func roundFloat(val float64, precision uint) float64 {
    ratio := math.Pow(10, float64(precision))
    return math.Round(val*ratio) / ratio
}

func IncomeDataDecrease(incomeData *IncomeData,k_receipt float64) error {
  for _, allowance := range incomeData.Allowances {
      if(allowance.AllowanceType == "donation") {
        if(allowance.Amount < 100000 && allowance.Amount > 10000) {
          incomeData.TotalIncome = incomeData.TotalIncome - allowance.Amount  
        } else if allowance.Amount >= 100000 {
          incomeData.TotalIncome = incomeData.TotalIncome - 100000 
        }
      } else if(allowance.AllowanceType == "k-receipt") {
        if(allowance.Amount >= k_receipt) {
          incomeData.TotalIncome = incomeData.TotalIncome - k_receipt
        } else if(allowance.Amount < k_receipt && allowance.Amount > 0) {
          incomeData.TotalIncome = incomeData.TotalIncome - allowance.Amount
        } 
      }
  }
  return nil
}

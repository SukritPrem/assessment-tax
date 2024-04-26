package calculateTax

import (
	"encoding/json"
	// "reflect"
	// "errors"
	"github.com/go-playground/validator/v10"
	"math"
	// "fmt"
)

type IncomeDataValidate struct {
  TotalIncome *float64 `json:"totalIncome" validate:"required,maxfloat"`
  Wht        *float64 `json:"wht" validate:"required,maxfloat"`
  Allowances *[]struct {
    AllowanceType *string  `json:"allowanceType" validate:"required,checkType"`
    Amount        *float64 `json:"amount" validate:"required,maxfloat"`
  } `json:"allowances" validate:"required,dive,required"`
}

func validateValueByStuct(jsonBytes []byte) (IncomeData, error) {
	var incomeData IncomeDataValidate
	d := IncomeData{}
	if err := json.Unmarshal(jsonBytes, &incomeData); err != nil {
		return d, err
	}
	// Custom validation
	validate := validator.New()
	validate.RegisterValidation("maxfloat", validateFloatMax)
	validate.RegisterValidation("checkType", validateAllowanceType)
	if err := validate.Struct(incomeData); err != nil {
		return d, err
	}
	if(incomeData.TotalIncome != nil) {
		d.TotalIncome = *incomeData.TotalIncome
	}
	if(incomeData.Wht != nil) {
		d.Wht = *incomeData.Wht
	}
	allowance := Allowance{}
	for _, allowanceValidate := range *incomeData.Allowances {
		if(allowanceValidate.AllowanceType != nil) {
			allowance.AllowanceType = *allowanceValidate.AllowanceType
		}
		if(allowanceValidate.Amount != nil) {
			allowance.Amount = *allowanceValidate.Amount
		}
		d.Allowances = append(d.Allowances, allowance)
	}
	// fmt.Printf("Data: %v\n", d)
	return d, nil
}

func validateFloatMax(fl validator.FieldLevel) bool {
	data := fl.Field().Interface().(float64)
	if data > float64(math.MaxFloat64) || data < 0	{
		return false
	}
	return true
}

func validateAllowanceType(fl validator.FieldLevel) bool {
	data := fl.Field().Interface().(string)
	if !(data == "donation" || data == "k-receipt") {
		return false
	}
	return true
}
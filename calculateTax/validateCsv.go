
package calculateTax

import (
	"fmt"
	"strconv"
	"strings"
	"errors"
)

func pasreCommaToIncomeData(values []string) (IncomeData, error){
	dataOneLine := IncomeData{}
	var err error
	var num float64

	dataOneLine.TotalIncome, err = strconv.ParseFloat(values[0], 64)
	if(err != nil){
		return dataOneLine, fmt.Errorf("Error can't parse TotalIncome")
	}
	dataOneLine.Wht, err = strconv.ParseFloat(values[1], 64)
	if(err != nil){
		return dataOneLine, fmt.Errorf("Error can't parse Wht")
	}	
	num, err = strconv.ParseFloat(values[2], 64)
	if(err != nil){
		return dataOneLine, fmt.Errorf("Error can't parse donation amount")
	}
	dataOneLine.Allowances = append(dataOneLine.Allowances, struct {
		AllowanceType string  `json:"allowanceType"`
		Amount        float64 `json:"amount"`
	}{"donation", num})

	return dataOneLine, nil
}

func validateCSV(data []byte, req *[]IncomeData, personalDeduction float64, k_receipt float64) (error, Reponse_csv) {
    r := Reponse_csv{}
    result := TotalIncomeAndTax{}
    data = []byte(strings.Replace(string(data), "\r", "", -1))
    lines := strings.Split(string(data), "\n")
	if(len(lines) < 2){
		return errors.New("Error no data"), r
	}
    // lines[0] = strings.Replace(lines[0], "\r", "", -1)
    if(strings.Compare(lines[0],"totalIncome,wht,donation") != 0){
      return errors.New("Error format header: totalIncome,wht,donation"), r	
    }

    for _, line := range lines[1:] {
        values := strings.Split(line, ",")
		if(len(values) != 3){
			return errors.New("Error format data need equal 3"), r
		}
		dataOneLine, err := pasreCommaToIncomeData(values)
		if(err != nil){
			return err, r
		}
		err = checkErrorIncomeData(&dataOneLine)
		if(err != nil){
			return err, r
		}
        result.TotalIncome = dataOneLine.TotalIncome
	
        dataOneLine.TotalIncome = dataOneLine.TotalIncome - personalDeduction
        err = IncomeDataDecrease(&dataOneLine,k_receipt)
		if(err != nil){
			return err, r
		}
        taxlevels := CalculateTaxLevelWithNetIncomeData(&dataOneLine)
        sum_tax := sumAllTaxLevel(taxlevels)
        result.Tax = sum_tax
        r.Taxs = append(r.Taxs, result)
        
    }

    return nil, r
}

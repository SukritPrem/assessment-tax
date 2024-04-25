
package calculateTax

import (
	"fmt"
	"strconv"
	"strings"
)

func pasreCommaToIncomeData(values []string) (IncomeData, error){
	dataOneLine := IncomeData{}
	// fmt.Printf("values: %v\n", values)
	var err error
	var num float64

	if(len(values) != 3){
		fmt.Println("Error in line: ", values)
	}
	dataOneLine.TotalIncome, err = strconv.ParseFloat(values[0], 64)
	if(err != nil){
		fmt.Printf("Error")
	}
	dataOneLine.Wht, err = strconv.ParseFloat(values[1], 64)
	if(err != nil){
		fmt.Printf("Error")
	}	
	num, err = strconv.ParseFloat(values[2], 64)
	if(err != nil){
		fmt.Printf("Error")
	}
	dataOneLine.Allowances = append(dataOneLine.Allowances, struct {
		AllowanceType string  `json:"allowanceType"`
		Amount        float64 `json:"amount"`
	}{"donation", num})
	// fmt.Printf("dataOneLine: %v\n", dataOneLine)
	return dataOneLine, nil
}

func validateCSV(data []byte, req *[]IncomeData, personalDeduction float64, k_receipt float64) (error, Reponse_csv) {
    r := Reponse_csv{}
    // Reponse_csv.Taxs = append(Reponse_csv.Taxs, TotalIncomeAndTax{})
    result := TotalIncomeAndTax{}
    // dataOneLine := IncomeData{}
    data = []byte(strings.Replace(string(data), "\r", "", -1))
    lines := strings.Split(string(data), "\n")

    lines[0] = strings.Replace(lines[0], "\r", "", -1)
    if(strings.Compare(lines[0],"totalIncome,wht,donation") != 0){
      fmt.Println("Header row not found. Assuming data starts from the first line.")
    }

    // Validate each line
    // var err error
    // var num float64
    for _, line := range lines[1:] {
        fmt.Println(line)
        values := strings.Split(line, ",")
		// fmt.Printf("values: %v\n", values)
		// pasreCommaToIncomeData(values)
		dataOneLine, err := pasreCommaToIncomeData(values)
		if(err != nil){
			return err, r
		}
        result.TotalIncome = dataOneLine.TotalIncome
        dataOneLine.TotalIncome = dataOneLine.TotalIncome - personalDeduction
        IncomeDataDecrease(&dataOneLine,k_receipt)
        taxlevels := CalculateTaxLevelWithNetIncomeData(&dataOneLine)
        sum_tax := sumAllTaxLevel(taxlevels)
        result.Tax = sum_tax
        r.Taxs = append(r.Taxs, result)
        
    }

    return nil, r
}

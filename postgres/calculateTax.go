package postgres

import (
	"fmt"
)	
type tax struct {
	id 	 int
	amount float64
	taxType string
}

func (p *Postgres) GetAmountByTaxType(v string) (float64, error) {
	rows, err := p.Db.Query("SELECT * FROM ktaxes WHERE taxtype=$1", v)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	fmt.Println("rows")
	for rows.Next() {
		var t tax
		err := rows.Scan(&t.id, &t.amount, &t.taxType)
		if err != nil {
			return 0, err
		}
		fmt.Println(t.amount)
		return t.amount, nil
	}
	return 0, nil
}

func (p *Postgres) UpdateAmountByTaxType(v string,a float64) (float64, error) {
	_, err := p.Db.Exec("UPDATE ktaxes SET amount=$1 WHERE taxtype=$2", a,v)
	if err != nil {
		return 0, err
	}
	return 0, nil
}



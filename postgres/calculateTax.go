package postgres

import (
	"fmt"
)	
type tax struct {
	id 	 int
	amount int
	taxType string
}

func (p *Postgres) GetPersonalDeduction() (int, error) {
	rows, err := p.Db.Query("SELECT * FROM ktaxes WHERE taxtype=$1", "personalDeduction")
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

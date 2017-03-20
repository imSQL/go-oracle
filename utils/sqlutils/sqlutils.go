package sqlutils

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"os"
)

type Result map[int]map[string]string

func (res *Result) GetMetric(db *sql.DB, query_text string) {

	rows, err := db.Query(query_text)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rows.Close()

	cols, _ := rows.Columns()

	values := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))

	for i := range values {
		scans[i] = &values[i]
	}

	result := make(Result)
	i := 0

	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		row := make(map[string]string)

		for k, v := range values {
			key := cols[k]
			row[key] = string(v)
		}

		result[i] = row
		i++
	}
	*res = result

}

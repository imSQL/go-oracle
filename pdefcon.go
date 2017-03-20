package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"os"
	//	"utils/controlfiles"
	//	"utils/database"
	//	"utils/datafile"
	//	"utils/eventmetrics"
	//	"utils/instance"
	//	"utils/onlinelogs"
	//	"utils/parameters"
	//	"utils/sga"
	"pdefcon-for-oracle/utils/systemload"
	//	"utils/tablespace"
	//	"utils/users"
	//	"utils/version"
)

type ID string

func (id ID) Scan(src interface{}) error {
	fmt.Println(src)
	return nil
}

func getDSN() string {
	var dsn string
	if len(os.Args) > 1 {
		dsn = os.Args[1]
		if dsn != "" {
			return dsn
		}
	}
	dsn = os.Getenv("GO_OCI8_CONNECT_STRING")
	if dsn != "" {
		return dsn
	}
	fmt.Fprintln(os.Stderr, `Please specifiy connection parameter in GO_OCI8_CONNECT_STRING environment variable,
or as the first argument! (The format is user/name@host:port/sid)`)
	return "scott/tiger@XE"
}

func main() {
	os.Setenv("NLS_LANG", "")

	db, err := sql.Open("oci8", getDSN())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	var wait_class string
	rows, err := db.Query(systemload.ViewSystemLoad)
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

	result := make(map[int]map[string]string)
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

	fmt.Println(result)

	fmt.Println(wait_class)
}

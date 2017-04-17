package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"os"
	//	"utils/controlfiles"
	"pdefcon-for-oracle/utils/database"
	"pdefcon-for-oracle/utils/datafile"
	//	"utils/eventmetrics"
	"pdefcon-for-oracle/utils/instance"
	//	"utils/onlinelogs"
	//	"utils/parameters"
	//	"utils/sga"
	"pdefcon-for-oracle/utils/filestat"
	"pdefcon-for-oracle/utils/systemload"
	"pdefcon-for-oracle/utils/table"
	"pdefcon-for-oracle/utils/tablespace"
	"pdefcon-for-oracle/utils/users"
	"pdefcon-for-oracle/utils/version"

	"pdefcon-for-oracle/awr/loadprofile"
	"pdefcon-for-oracle/awr/topsql"
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

	result1 := new(systemload.SystemLoad)
	result1.DbHandler = db

	result1.GetMetrics()
	result1.PrintMetrics()

	result2 := new(version.Version)
	result2.DbHandler = db

	result2.GetMetrics()
	result2.PrintMetrics()

	result3 := new(users.Users)
	result3.DbHandler = db
	result3.GetMetrics()
	result3.PrintMetrics()
	result3.GetTotalCounterByUsername()
	result3.GetTotalCounterByStatus()
	result3.GetTotalCounter()

	result4 := new(database.Database)
	result4.DbHandler = db
	result4.GetMetrics()
	result4.PrintMetrics()

	result5 := new(instance.Instances)
	result5.DbHandler = db
	result5.GetMetrics()
	result5.PrintMetrics()

	result6 := new(tablespace.Tablespace)
	result6.DbHandler = db
	result6.GetMetrics()
	result6.PrintMetrics()

	result7 := new(datafile.Datafile)
	result7.DbHandler = db
	result7.GetMetrics()
	result7.PrintMetrics()

	result8 := new(filestat.Filestat)
	result8.DbHandler = db
	result8.GetMetrics()
	result8.PrintMetrics()

	result9 := new(table.Table)
	result9.DbHandler = db
	result9.GetMetrics()
	result9.PrintMetrics()

	result10 := new(loadprofile.Cursor)
	result10.DbHandler = db
	result10.GetMetrics()
	result10.PrintMetrics()

	result11 := new(topsql.Cursor)
	result11.DbHandler = db
	result11.GetMetrics()
	//result11.PrintMetrics()

	//fmt.Println(wait_class)

}

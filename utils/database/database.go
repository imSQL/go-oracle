package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"os"
	"pdefcon-for-oracle/utils/sqlutils"
	//	"sort"
	"strconv"
	"strings"
)

type Database struct {
	Odb       sqlutils.Result
	DbHandler *sql.DB
}

const ViewDatabase = `SELECT * FROM V$DATABASE`

func (odbs *Database) GetMetrics() {
	odbs.Odb.GetMetric(odbs.DbHandler, ViewDatabase)
}

func (odbs *Database) PrintMetrics() {
	current_hostname, _ := os.Hostname()

	for _, v := range odbs.Odb {
		fmt.Fprintf(os.Stdout, "OracleDatabases,host=%s,databasename=%s ", current_hostname, v["NAME"])
		length_v := len(v)
		counter := 0
		for ak, av := range v {
			if counter == length_v-1 {
				if _, ok := strconv.ParseInt(av, 10, 64); ok != nil {
					fmt.Fprintf(os.Stdout, "%s=%q", strings.Replace(strings.ToLower(ak), "#", "", -1), strings.ToLower(av))
				} else {
					fmt.Fprintf(os.Stdout, "%s=%s", strings.Replace(strings.ToLower(ak), "#", "", -1), strings.ToLower(av))
				}
			} else {
				if _, ok := strconv.ParseInt(av, 10, 64); ok != nil {
					fmt.Fprintf(os.Stdout, "%s=%q,", strings.Replace(strings.ToLower(ak), "#", "", -1), strings.ToLower(av))
				} else {
					fmt.Fprintf(os.Stdout, "%s=%s,", strings.Replace(strings.ToLower(ak), "#", "", -1), strings.ToLower(av))
				}

			}
			counter++
		}
		fmt.Fprintf(os.Stdout, "\n")
	}
}

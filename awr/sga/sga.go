package sga

import (
	"database/sql"
	"fmt"
	"os"
	"pdefcon-for-oracle/utils/sqlutils"
	"strconv"
	"strings"
)

type OracleSGA struct {
	poolname      string
	poolsize      int64
	componentname string
	componentsize int64
	rate          float64
}

const (
	SGA = `
SELECT 
    cmp.component poolname,
    cmp.current_size poolsize,
    ss.name componentname,
    ss.bytes componentsize,
    (ss.bytes/cmp.current_size*100) rate
FROM
    v$sga_dynamic_components cmp,
    v$sgastat ss
WHERE
	cmp.component = ss.pool
`
)

type Cursor struct {
	sga       []OracleSGA
	cursor    sqlutils.Result
	DbHandler *sql.DB
}

func (cs *Cursor) GetMetrics() {
	cs.cursor.GetMetric(cs.DbHandler, SGA)
	for _, val := range cs.cursor {
		tmp := new(OracleSGA)
		for ak, av := range val {
			switch ak {
			case "POOLNAME":
				tmp.poolname = strings.Replace(av, " ", "_", -1)
			case "POOLSIZE":
				tmp.poolsize, _ = strconv.ParseInt(av, 10, 64)
			case "COMPONENTNAME":
				tmp.componentname = strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(av, ",", "", -1), "#", "", -1), "/", "", -1), ":", "", -1), " ", "_", -1)
			case "COMPONENTSIZE":
				tmp.componentsize, _ = strconv.ParseInt(av, 10, 64)
			case "RATE":
				tmp.rate, _ = strconv.ParseFloat(av, 64)
			default:
				fmt.Println("Nothing")
			}
		}
		cs.sga = append(cs.sga, *tmp)
	}

}

func (cs *Cursor) PrintMetrics() {
	current_hostname, _ := os.Hostname()
	for _, val := range cs.sga {
		fmt.Fprintf(os.Stdout, "OracleSGA,host=%s,poolname=%s,componentname=%s poolname=%q,poolsize=%d,componentname=%q,componentsize=%d,rate=%.6f\n",
			current_hostname,
			val.poolname,
			val.componentname,
			val.poolname,
			val.poolsize,
			val.componentname,
			val.componentsize,
			val.rate,
		)
	}
}

package table

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"os"
	"pdefcon-for-oracle/utils/sqlutils"
	"strconv"
	"strings"
)

type Table struct {
	tablestat sqlutils.Result
	DbHandler *sql.DB
}

const (
	TableStat = `
SELECT
    dt.OWNER,
    dt.TABLESPACE_NAME,
    dt.TABLE_NAME,
    dt.CLUSTER_NAME,
    dt.STATUS,
    dt.LOGGING,
    dt.BACKED_UP,
    dt.NUM_ROWS,
    dt.AVG_ROW_LEN,
    dt.CACHE,
    dt.TABLE_LOCK,
    dt.TEMPORARY,
    dt.BUFFER_POOL,
    dt.FLASH_CACHE,
    dt.CELL_FLASH_CACHE,
    dt.ROW_MOVEMENT,
    dt.GLOBAL_STATS,
    dt.USER_STATS,
    dt.MONITORING,
    dt.COMPRESSION,
    dt.COMPRESS_FOR,
    dt.DROPPED,
    dt.READ_ONLY,
    dt.PARTITIONED,
    dtm.PARTITION_NAME,
    dtm.SUBPARTITION_NAME,
    dtm.inserts,
    dtm.deletes,
    dtm.updates,
    dtm.TRUNCATED,
    dtm.DROP_SEGMENTS
FROM
    DBA_TABLES dt,
    SYS.DBA_TAB_MODIFICATIONS dtm
WHERE
    dt.table_name = dtm.table_name
AND dt.OWNER not in ('SYS','SYSTEM')
`
)

func (tb *Table) GetMetrics() {
	tb.tablestat.GetMetric(tb.DbHandler, TableStat)
}

func (tb *Table) PrintMetrics() {
	current_hostname, _ := os.Hostname()

	for _, v := range tb.tablestat {
		fmt.Fprintf(os.Stdout, "OracleTableStat,host=%s,tablename=%s ", current_hostname, v["TABLE_NAME"])
		length_v := len(v)
		counter := 0
		for ak, av := range v {
			if av == "" {
				av = "Nil"
			}
			if counter == length_v-1 {
				if _, ok := strconv.ParseInt(av, 10, 64); ok != nil {
					fmt.Fprintf(os.Stdout, "%s=%q", strings.Replace(strings.ToLower(ak), "#", "", -1), strings.ToLower(strings.TrimSpace(av)))
				} else {
					fmt.Fprintf(os.Stdout, "%s=%q", strings.Replace(strings.ToLower(ak), "#", "", -1), strings.ToLower(strings.TrimSpace(av)))

				}
			} else {
				if _, ok := strconv.ParseInt(av, 10, 64); ok != nil {
					fmt.Fprintf(os.Stdout, "%s=%q,", strings.Replace(strings.ToLower(ak), "#", "", -1), strings.ToLower(strings.TrimSpace(av)))
				} else {
					fmt.Fprintf(os.Stdout, "%s=%q,", strings.Replace(strings.ToLower(ak), "#", "", -1), strings.ToLower(strings.TrimSpace(av)))

				}
			}
			counter++
		}
		fmt.Fprintf(os.Stdout, "\n")
	}
}

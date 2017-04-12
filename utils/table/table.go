package table

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	//"os"
	"pdefcon-for-oracle/utils/sqlutils"
	//"strconv"
	//"strings"
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
    dt.PCT_USED,
    dt.PCT_INCREASE,
    dt.PCT_USED,
    dt.LOGGING,
    dt.BACKED_UP,
    dt.NUM_ROWS,
    dt.BLOCKS,
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
    DBA_TAB_MODIFICATIONS dtm
WHERE
    dt.table_name = dtm.table_name
`
)

func (tb *Table) GetMetrics() {
	tb.tablestat.GetMetric(tb.DbHandler, TableStat)
}

func (tb *Table) PrintMetrics() {
	fmt.Println(tb)
}

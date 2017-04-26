package bufferpool

import (
	"database/sql"
	"fmt"
	"os"
	"pdefcon-for-oracle/utils/sqlutils"
	"strconv"
)

type OracleBufferPool struct {
	name            string
	physical_reads  int64
	db_block_gets   int64
	consistent_gets int64
	hit             float64
}

type Cursor struct {
	OracleBufferPool
	cursor    sqlutils.Result
	DbHandler *sql.DB
}

const (
	BufferPool = `
SELECT
    NAME,
    PHYSICAL_READS,
    DB_BLOCK_GETS,
    CONSISTENT_GETS,
    SUBSTR(1-(PHYSICAL_READS/(DB_BLOCK_GETS+CONSISTENT_GETS)),1,5) HIT
FROM
    V$BUFFER_POOL_STATISTICS
    `
)

func (cs *Cursor) GetMetrics() {
	cs.cursor.GetMetric(cs.DbHandler, BufferPool)
	for ak, av := range cs.cursor[0] {
		switch ak {
		case "NAME":
			cs.OracleBufferPool.name = av
		case "PHYSICAL_READS":
			cs.OracleBufferPool.physical_reads, _ = strconv.ParseInt(av, 10, 64)
		case "DB_BLOCK_GETS":
			cs.OracleBufferPool.db_block_gets, _ = strconv.ParseInt(av, 10, 64)
		case "CONSISTENT_GETS":
			cs.OracleBufferPool.consistent_gets, _ = strconv.ParseInt(av, 10, 64)
		case "HIT":
			cs.OracleBufferPool.hit, _ = strconv.ParseFloat(av, 64)
		default:
			fmt.Println("Nothing")
		}
	}
}

func (cs *Cursor) PrintMetrics() {
	current_host, _ := os.Hostname()
	fmt.Fprintf(os.Stdout, "OracleBufferPool,host=%s,bufferpoolname=%s name=%q,PHYSICAL_READS=%d,db_block_gets=%d,consistent_gets=%d,hit=%.4f\n",
		current_host,
		cs.OracleBufferPool.name,
		cs.OracleBufferPool.name,
		cs.OracleBufferPool.physical_reads,
		cs.OracleBufferPool.db_block_gets,
		cs.OracleBufferPool.consistent_gets,
		cs.OracleBufferPool.hit,
	)
}

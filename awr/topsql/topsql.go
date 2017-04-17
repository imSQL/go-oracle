package topsql

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	//	"os"
	"pdefcon-for-oracle/utils/sqlutils"
	"strconv"
	"time"
)

type OracleTopSQL struct {
	interval_end   string
	snap_id        int64
	sql_id         int64
	execs          int64
	buff_gets      int64
	phys_reads     int64
	cpus           float64
	elap_s         float64
	rows_ret       int64
	io_waits       float64
	exec_rank      int64
	buff_get_rank  int64
	phys_read_rank int64
	cpu_rank       int64
	elap_rank      int64
	io_wait_rank   int64
	concur_wait_s  int64
	parses         int64
	sorts          int64
	plsqls         float64
	app_wait_s     float64
	gets_per_exec  int64
	reads_per_exec int64
	rows_per_exec  int64
	sql_text       string
}

type Cursor struct {
	OracleTopSQL
	cursor    sqlutils.Result
	DbHandler *sql.DB
}

const (
	TopSQL = `
SELECT
    c.end_interval_time AS "INTERVAL_END",
    c.snap_id           AS SNAP_ID,
    c.sql_id            AS "SQL_ID",
    c.execs             AS EXECS,
    c.buff_gets         AS "BUFF_GETS",
    c.phys_reads        AS "PHYS_READS",
    c.cpu_sec           AS "CPUS",
    c.elap_sec          AS "ELAP_S",
    c.rows_ret          AS "ROWS_RET",
    c.io_wait_sec       AS "IO_WAITS",
    c.exec_rank         AS "EXEC_RANK",
    c.buff_get_rank     AS "BUFF_GET_RANK",
    c.phys_read_rank    AS "PHYS_READ_RANK",
    c.cpu_rank          AS "CPU_RANK",
    c.elap_rank         AS "ELAP_RANK",
    c.io_wait_rank      AS "IO_WAIT_RANK",
    c.concur_wait_sec   AS "CONCUR_WAIT_S",
    c.parses            AS PARSES,
    c.sorts             AS SORTS,
    c.plsql_sec         AS "PLSQLS",
    c.app_wait_sec      AS "APP_WAIT_S",
    c.gets_per_exec     AS "GETS_PER_EXEC",
    c.reads_per_exec    AS "READS_PER_EXEC",
    c.rows_per_exec     AS "ROWS_PER_EXEC",
    DBMS_LOB.SUBSTR(DECODE ('Y', 'Y', dhst.sql_text, 'T', SUBSTR (
    dhst.sql_text, 1, 1000), NULL ) , 4000, 1)AS SQL_TEXT
FROM
    (
        SELECT
            b.*
        FROM
            (
                SELECT
                    TO_CHAR(a.end_interval_time,'MM/DD/YYYY HH24:MI:SS') AS
                    end_interval_time,
                    a.snap_id,
                    a.sql_id,
                    a.execs                     AS execs,
                    a.buff_gets                 AS BUFF_GETS,
                    a.phys_reads                AS PHYS_READS,
                    ROUND(a.cpu_time /1000000,2) AS CPU_SEC,
                    ROUND(a.elap_time/1000000,2) AS ELAP_SEC,
                    a.rows_ret                   AS rows_ret,
                    ROUND(a.io_wait/1000000,2)   AS IO_WAIT_SEC,
                    RANK() OVER (PARTITION BY a.snap_id ORDER BY a.execs DESC,
                    a.elap_time DESC) AS EXEC_RANK,
                    RANK() OVER (PARTITION BY a.snap_id ORDER BY a.buff_gets
                    DESC, a.elap_time DESC) AS BUFF_GET_RANK,
                    RANK() OVER (PARTITION BY a.snap_id ORDER BY a.phys_reads
                    DESC, a.elap_time DESC) AS PHYS_READ_RANK,
                    RANK() OVER (PARTITION BY a.snap_id ORDER BY a.cpu_time
                    DESC, a.elap_time DESC) AS CPU_RANK,
                    rank() OVER (PARTITION BY a.snap_id ORDER BY a.elap_time
                    DESC, a.elap_time DESC) AS ELAP_RANK,
                    RANK() OVER (PARTITION BY a.snap_id ORDER BY a.io_wait DESC
                    , a.elap_time DESC)            AS IO_WAIT_RANK,
                    ROUND(a.concur_wait/1000000,2) AS CONCUR_WAIT_SEC,
                    a.parses                       AS parses,
                    a.sorts                        AS sorts,
                    ROUND(a.plsql_time/1000000,2)  AS PLSQL_SEC,
                    ROUND(a.app_wait  /1000000,2)  AS APP_WAIT_SEC,
                    ROUND(a.buff_gets /a.execs,0)  AS GETS_PER_EXEC,
                    ROUND(a.phys_reads/a.execs,0)  AS READS_PER_EXEC,
                    ROUND(a.rows_ret  /a.execs,0)  AS ROWS_PER_EXEC
                FROM
                    (
                        SELECT
                            dhs.end_interval_time,
                            dhss.snap_id,
                            dhss.sql_id,
                            dhss.executions_delta     AS execs,
                            dhss.buffer_gets_delta    AS buff_gets,
                            dhss.sorts_delta          AS sorts,
                            dhss.loads_delta          AS loads,
                            dhss.invalidations_delta  AS invalds,
                            dhss.parse_calls_delta    AS parses,
                            dhss.disk_reads_delta     AS phys_reads,
                            dhss.rows_processed_delta AS rows_ret,
                            dhss.cpu_time_delta       AS cpu_time,
                            dhss.elapsed_time_delta   AS elap_time,
                            dhss.iowait_delta         AS io_wait,
                            dhss.ccwait_delta         AS concur_wait,
                            dhss.plsexec_time_delta   AS plsql_time,
                            dhss.apwait_delta         AS app_wait
                        FROM
                            dba_hist_sqlstat dhss,
                            dba_hist_snapshot dhs
                        WHERE
                            dhss.snap_id = dhs.snap_id
                        AND dhs.snap_id IN
                            (
                                SELECT
                                    snap_id
                                FROM
                                    dba_hist_snapshot
                                WHERE
                                    end_interval_time BETWEEN TO_DATE(
                                    '%s','YYYY-MM-DD HH24:MI:SS') AND
                                    TO_DATE('%s','YYYY-MM-DD HH24:MI:SS')
                            )
                        AND dhss.executions_delta >0
                    )
                    a
                ORDER BY
                    snap_id,
                    buff_get_rank
            )
            b
        WHERE
            b.buff_get_rank  <= %d
         OR b.phys_read_rank <= %d
         OR b.exec_rank      <= %d
         OR b.cpu_rank       <= %d
         OR b.elap_rank      <= %d
         OR b.io_wait_rank   <= %d
    )
    c,
    dba_hist_sqltext dhst
WHERE
    dhst.sql_id = c.sql_id
ORDER BY
    c.snap_id,
    c.buff_get_rank`
)

func (ts *Cursor) GetMetrics() {
	current_date_start := fmt.Sprintf("%s00:00", time.Now().Format("2006-01-02 15:"))
	current_date_end := fmt.Sprintf("%s59:59", time.Now().Format("2006-01-02 15:"))
	query_text := fmt.Sprintf(TopSQL, current_date_start, current_date_end, 10, 10, 10, 10, 10, 10)

	fmt.Println(current_date_start, current_date_end)

	ts.cursor.GetMetric(ts.DbHandler, query_text)

	for ak, av := range ts.cursor[0] {
		switch ak {
		case "INTERVAL_END":
			ts.OracleTopSQL.interval_end = av
		case "SNAP_ID":
			ts.OracleTopSQL.snap_id, _ = strconv.ParseInt(av, 10, 64)
		case "SQL_ID":
			ts.OracleTopSQL.sql_id, _ = strconv.ParseInt(av, 10, 64)
		case "EXECS":
			ts.OracleTopSQL.execs, _ = strconv.ParseInt(av, 10, 64)
		case "BUFF_GETS":
			ts.OracleTopSQL.buff_gets, _ = strconv.ParseInt(av, 10, 64)
		case "PHYS_READS":
			ts.OracleTopSQL.phys_reads, _ = strconv.ParseInt(av, 10, 64)
		case "CPUS":
			ts.OracleTopSQL.cpus, _ = strconv.ParseFloat(av, 64)
		case "ELAP_S":
			ts.OracleTopSQL.elap_s, _ = strconv.ParseFloat(av, 64)
		case "ROWS_RET":
			ts.OracleTopSQL.rows_ret, _ = strconv.ParseInt(av, 10, 64)
		case "IO_WAITS":
			ts.OracleTopSQL.io_waits, _ = strconv.ParseFloat(av, 64)
		case "EXEC_RANK":
			ts.OracleTopSQL.exec_rank, _ = strconv.ParseInt(av, 10, 64)
		case "BUFF_GET_RANK":
			ts.OracleTopSQL.buff_get_rank, _ = strconv.ParseInt(av, 10, 64)
		case "PHYS_READ_RANK":
			ts.OracleTopSQL.phys_read_rank, _ = strconv.ParseInt(av, 10, 64)
		case "CPU_RANK":
			ts.OracleTopSQL.cpu_rank, _ = strconv.ParseInt(av, 10, 64)
		case "ELAP_RANK":
			ts.OracleTopSQL.elap_rank, _ = strconv.ParseInt(av, 10, 64)
		case "IO_WAIT_RANK":
			ts.OracleTopSQL.io_wait_rank, _ = strconv.ParseInt(av, 10, 64)
		case "CONCUR_WAIT_S":
			ts.OracleTopSQL.concur_wait_s, _ = strconv.ParseInt(av, 10, 64)
		case "PARSES":
			ts.OracleTopSQL.parses, _ = strconv.ParseInt(av, 10, 64)
		case "SORTS":
			ts.OracleTopSQL.sorts, _ = strconv.ParseInt(av, 10, 64)
		case "PLSQLS":
			ts.OracleTopSQL.plsqls, _ = strconv.ParseFloat(av, 64)
		case "APP_WAIT_S":
			ts.OracleTopSQL.app_wait_s, _ = strconv.ParseFloat(av, 64)
		case "GETS_PER_EXEC":
			ts.OracleTopSQL.gets_per_exec, _ = strconv.ParseInt(av, 10, 64)
		case "READS_PER_EXEC":
			ts.OracleTopSQL.reads_per_exec, _ = strconv.ParseInt(av, 10, 64)
		case "ROWS_PER_EXEC":
			ts.OracleTopSQL.rows_per_exec, _ = strconv.ParseInt(av, 10, 64)
		case "SQL_TEXT":
			ts.OracleTopSQL.sql_text = av
		default:
			fmt.Println("Nothing")
		}
	}
}

func (ts *Cursor) PrintMetrics() {
	fmt.Println(ts.OracleTopSQL)
}

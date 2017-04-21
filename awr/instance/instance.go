package instance

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"pdefcon-for-oracle/utils/sqlutils"
	"strconv"
	//	"strings"
	"os"
	"time"
)

type OracleInstance struct {
	snap_id                    int64
	interval_end               string
	snap_seconds               float64
	execs                      int64
	execs_s                    float64
	logical_reads              int64
	logical_reads_s            float64
	phys_reads                 int64
	phys_reads_s               float64
	phys_writes                int64
	phys_writes_s              float64
	hard_parses                int64
	hard_parses_s              float64
	total_parses               int64
	total_parses_s             float64
	parses_cpu_s               float64
	parses_cpu_s_s             float64
	parse_elap_s               float64
	parse_elap_s_s             float64
	sorts_disk                 int64
	sorts_disk_s               float64
	sorts_mem                  int64
	sorts_mem_s                float64
	sorts_rows                 int64
	sorts_rows_s               float64
	tbl_scan_rows              int64
	tbl_scan_rows_s            float64
	tbl_fetch_rowid            int64
	tbl_fetch_rowid_s          float64
	io_wait_s                  float64
	io_wait_s_s                float64
	concur_wait_s              float64
	concur_wait_s_s            float64
	user_calls                 int64
	user_calls_s               float64
	user_commits               int64
	user_commits_s             float64
	user_rollbacks             int64
	user_rollbacks_s           float64
	sqlnet_rt_client           int64
	sqlnet_rt_client_s         float64
	sqlnet_rt_dblink           int64
	sqlnet_rt_dblink_s         float64
	sqlnet_bytes_from_client   int64
	sqlnet_bytes_from_client_s float64
	sqlnet_bytes_from_dblink   int64
	sqlnet_bytes_from_dblink_s float64
	sqlnet_bytes_to_client     int64
	sqlnet_bytes_to_client_s   float64
	sqlnet_bytes_to_dblink     int64
	sqlnet_bytes_to_dblink_s   float64
	wa_execs_multi             int64
	wa_execs_multi_s           float64
	wa_execs_onepass           int64
	wa_execs_onepass_s         float64
	wa_execs_optimal           int64
	wa_execs_optimal_s         float64
}

type Cursor struct {
	OracleInstance
	cursor    sqlutils.Result
	DbHandler *sql.DB
}

const (
	INSTANCE_STATUS = `
SELECT
    snap_id                                            AS SNAP,
    TO_CHAR(end_interval_time,'YYYY-MM-DD HH24:MI:SS') AS "INTERVAL_END",
    snap_seconds                                       AS "SNAP_SECONDS",
    execs AS EXECS,
    ROUND(execs/snap_seconds,2)                    AS "EXECS_S",
    logical_reads                                  AS "LOGICAL_READS",
    ROUND(logical_reads/snap_seconds,2)            AS "LOGICAL_READS_S",
    phys_reads                                     AS "PHYS_READS",
    ROUND(phys_reads/snap_seconds,2)               AS "PHYS_READS_S",
    phys_writes                                    AS "PHYS_WRITES",
    ROUND(phys_writes/snap_seconds,2)              AS "PHYS_WRITES_S",
    hard_parses                                    AS "HARD_PARSES",
    ROUND(hard_parses/snap_seconds,2)              AS "HARD_PARSES_S",
    total_parses                                   AS "TOTAL_PARSES",
    ROUND(total_parses/snap_seconds,2)             AS "TOTAL_PARSES_S",
    parse_cpu_s                                    AS "PARSE_CPU_S",
    ROUND(parse_cpu_s/snap_seconds,2)              AS "PARSE_CPU_S_S",
    parse_elap_s                                   AS "PARSE_ELAP_S",
    ROUND(parse_elap_s/snap_seconds,2)             AS "PARSE_ELAP_S_S",
    sorts_disk                                     AS "SORTS_DISK",
    ROUND(sorts_disk/snap_seconds,2)               AS "SORTS_DISK_S",
    sorts_mem                                      AS "SORTS_MEM",
    ROUND(sorts_mem/snap_seconds,2)                AS "SORTS_MEM_S",
    sorts_rows                                     AS "SORTS_ROWS",
    ROUND(sorts_rows/snap_seconds,2)               AS "SORTS_ROWS_S",
    tbl_fetch_rowid                                AS "TBL_FETCH_ROWID",
    ROUND(tbl_fetch_rowid/snap_seconds,2)          AS "TBL_FETCH_ROWID_S",
    tbl_scan_rows                                  AS "TBL_SCAN_ROWS",
    ROUND(tbl_scan_rows/snap_seconds,2)            AS "TBL_SCAN_ROWS_S",
    io_wait_s                                      AS "IO_WAIT_S",
    ROUND(io_wait_s/snap_seconds,2)                AS "IO_WAIT_S_S",
    concur_wait_s                                  AS "CONCUR_WAIT_S",
    ROUND(concur_wait_s/snap_seconds,2)            AS "CONCUR_WAIT_S_S",
    user_calls                                     AS "USER_CALLS",
    ROUND(user_calls/snap_seconds,2)               AS "USER_CALLS_S",
    user_commits                                   AS "USER_COMMITS",
    ROUND(user_commits/snap_seconds,2)             AS "USER_COMMITS_S",
    user_rollbacks                                 AS "USER_ROLLBACKS",
    ROUND(user_rollbacks/snap_seconds,2)           AS "USER_ROLLBACKS_S",
    sqlnet_rt_client                               AS "SQLNET_RT_CLIENT",
    ROUND(sqlnet_rt_client/snap_seconds,2)         AS "SQLNET_RT_CLIENT_S",
    sqlnet_rt_dblink                               AS "SQLNET_RT_DBLINK",
    ROUND(sqlnet_rt_dblink/snap_seconds,2)         AS "SQLNET_RT_DBLINK_S",
    sqlnet_bytes_from_client                       AS "SQLNET_BYTES_FROM_CLIENT",
    ROUND(sqlnet_bytes_from_client/snap_seconds,2) AS "SQLNET_BYTES_FROM_CLIENT_S",
    sqlnet_bytes_from_dblink                       AS "SQLNET_BYTES_FROM_DBLINK",
    ROUND(sqlnet_bytes_from_dblink/snap_seconds,2) AS "SQLNET_BYTES_FROM_DBLINK_S",
    sqlnet_bytes_to_client                       AS "SQLNET_BYTES_TO_CLIENT",
    ROUND(sqlnet_bytes_to_client/snap_seconds,2) AS "SQLNET_BYTES_TO_CLIENT_S",
    sqlnet_bytes_to_dblink                       AS "SQLNET_BYTES_TO_DBLINK",
    ROUND(sqlnet_bytes_to_dblink/snap_seconds,2) AS "SQLNET_BYTES_TO_DBLINK_S",
    wa_execs_multi                               AS "WA_EXECS_MULTI",
    ROUND(wa_execs_multi/snap_seconds,2)         AS "WA_EXECS_MULTI_S",
    wa_execs_onepass                             AS "WA_EXECS_ONEPASS",
    ROUND(wa_execs_onepass/snap_seconds,2)       AS "WA_EXECS_ONEPASS_S",
    wa_execs_optimal                             AS "WA_EXECS_OPTIMAL",
    ROUND(wa_execs_optimal/snap_seconds,2)       AS "WA_EXECS_OPTIMAL_S"
FROM
    (
        SELECT
            c.snap_id,
            c.end_interval_time,
            c.snap_seconds,
            MAX (DECODE (c.stat_name, 'SQL*Net roundtrips to/from client',
            c.VALUE)) AS SQLNET_RT_CLIENT,
            MAX (DECODE (c.stat_name, 'SQL*Net roundtrips to/from dblink',
            c.VALUE)) AS SQLNET_RT_DBLINK,
            MAX (DECODE (c.stat_name, 'bytes received via SQL*Net from client',
            c.VALUE)) AS SQLNet_BYTES_FROM_client,
            MAX (DECODE (c.stat_name, 'bytes received via SQL*Net from dblink',
            c.VALUE)) AS SQLNet_BYTES_from_dblink,
            MAX (DECODE (c.stat_name, 'bytes sent via SQL*Net to client',
            c.VALUE)) AS SQLNet_BYTES_to_client,
            MAX (DECODE (c.stat_name, 'bytes sent via SQL*Net to dblink',
            c.VALUE)) AS SQLNet_BYTES_to_dblink,
            MAX (DECODE (c.stat_name, 'concurrency wait time',ROUND(c.VALUE/100
            ,2)))                                                    AS concur_wait_s,
            MAX (DECODE (c.stat_name, 'execute count',c.VALUE))      AS execs,
            MAX (DECODE (c.stat_name, 'parse count (hard)',c.VALUE)) AS
            hard_parses,
            MAX (DECODE (c.stat_name, 'parse count (total)',c.VALUE)) AS
            total_parses,
            MAX (DECODE (c.stat_name, 'parse time cpu',ROUND(c.VALUE/100,2)))
            AS parse_cpu_s,
            MAX (DECODE (c.stat_name, 'parse time elapsed',ROUND(c.VALUE/100,2)
            ))                                                    AS parse_elap_s,
            MAX (DECODE (c.stat_name, 'physical reads',c.VALUE))  AS phys_reads,
            MAX (DECODE (c.stat_name, 'physical writes',c.VALUE)) AS
            phys_writes,
            MAX (DECODE (c.stat_name, 'session logical reads',c.VALUE)) AS
            logical_reads,
            MAX (DECODE (c.stat_name, 'sorts (disk)',c.VALUE))         AS sorts_disk,
            MAX (DECODE (c.stat_name, 'sorts (memory)',c.VALUE))       AS sorts_mem,
            MAX (DECODE (c.stat_name, 'sorts (rows)',c.VALUE))         AS sorts_rows,
            MAX (DECODE (c.stat_name, 'table fetch by rowid',c.VALUE)) AS
            tbl_fetch_rowid,
            MAX (DECODE (c.stat_name, 'table scan rows gotten',c.VALUE)) AS
            tbl_scan_rows,
            MAX (DECODE (c.stat_name, 'user I/O wait time',ROUND(c.VALUE/100,2)
            ))                                                   AS io_wait_s,
            MAX (DECODE (c.stat_name, 'user calls',c.VALUE))     AS user_calls,
            MAX (DECODE (c.stat_name, 'user commits',c.VALUE))   AS user_commits,
            MAX (DECODE (c.stat_name, 'user rollbacks',c.VALUE)) AS
            user_rollbacks,
            MAX (DECODE (c.stat_name, 'workarea executions - multipass',c.VALUE
            )) AS wa_execs_multi,
            MAX (DECODE (c.stat_name, 'workarea executions - onepass',c.VALUE))
            AS wa_execs_onepass,
            MAX (DECODE (c.stat_name, 'workarea executions - optimal',c.VALUE))
            AS wa_execs_optimal
        FROM
            (
                SELECT
                    dhss.snap_id,
                    dhs.end_interval_time,
                    dhss.stat_name,
                    dhss.VALUE - LAG (dhss.VALUE, 1) OVER (PARTITION BY
                    dhss.stat_name ORDER BY dhs.snap_id) AS VALUE,
                    DECODE ( dhs.startup_time - LAG (dhs.startup_time) OVER (
                    PARTITION BY dhss.stat_name ORDER BY dhs.snap_id), INTERVAL
                    '0' SECOND, 'N', 'Y' ) AS instance_restart,
                    EXTRACT (DAY FROM dhs.end_interval_time - LAG (
                    dhs.end_interval_time, 1) OVER (PARTITION BY dhss.stat_name
                    ORDER BY dhs.snap_id))                                   * 86400 + EXTRACT (HOUR FROM
                    dhs.end_interval_time                                    - LAG (dhs.end_interval_time, 1)
                    OVER (PARTITION BY dhss.stat_name ORDER BY dhs.snap_id)) *
                    3600                                                     +
                    EXTRACT (MINUTE FROM dhs.end_interval_time               -
                    LAG (dhs.end_interval_time, 1) OVER (PARTITION BY
                    dhss.stat_name ORDER BY dhs.snap_id)) * 60 + EXTRACT (
                    SECOND FROM dhs.end_interval_time     - LAG (
                    dhs.end_interval_time, 1) OVER (PARTITION BY dhss.stat_name
                    ORDER BY dhs.snap_id)) AS snap_seconds
                FROM
                    dba_hist_snapshot dhs,
                    dba_hist_sysstat dhss
                WHERE
                    dhs.snap_id     = dhss.snap_id
                AND dhss.stat_name IN ('SQL*Net roundtrips to/from client',
                    'SQL*Net roundtrips to/from dblink',
                    'bytes received via SQL*Net from client',
                    'bytes received via SQL*Net from dblink',
                    'bytes sent via SQL*Net to client',
                    'bytes sent via SQL*Net to dblink', 'concurrency wait time'
                    , 'execute count', 'parse count (hard)',
                    'parse count (total)', 'parse time cpu',
                    'parse time elapsed', 'physical reads', 'physical writes',
                    'session logical reads', 'sorts (disk)', 'sorts (memory)',
                    'sorts (rows)', 'table fetch by rowid',
                    'table scan rows gotten', 'user I/O wait time',
                    'user calls', 'user commits', 'user rollbacks',
                    'workarea executions - multipass',
                    'workarea executions - onepass',
                    'workarea executions - optimal')
                AND dhs.snap_id IN
                    (
                        SELECT
                            snap_id
                        FROM
                            dba_hist_snapshot
                        WHERE
                            end_interval_time BETWEEN TO_DATE('%s',
                            'YYYY-MM-DD HH24:MI:SS') AND TO_DATE('%s',
                            'YYYY-MM-DD HH24:MI:SS')
                    )
            )
            c
        WHERE
            instance_restart = 'N'
        AND extract(hour FROM end_interval_time) BETWEEN %d AND %d
        GROUP BY
            snap_id,
            end_interval_time,
            snap_seconds
        ORDER BY
            snap_id
    )
`
)

func (cs *Cursor) GetMetrics() {
	current_date_start := fmt.Sprintf("%s 00:00:00", time.Now().Format("2006-01-02"))
	current_date_end := fmt.Sprintf("%s 23:59:59", time.Now().Format("2006-01-02"))
	curr_hour := time.Now().Hour()

	query_text := fmt.Sprintf(INSTANCE_STATUS, current_date_start, current_date_end, curr_hour, curr_hour)
	cs.cursor.GetMetric(cs.DbHandler, query_text)

	for ak, av := range cs.cursor[0] {
		switch ak {

		case "SNAP":
			cs.OracleInstance.snap_id, _ = strconv.ParseInt(av, 10, 64)
		case "INTERVAL_END":
			cs.OracleInstance.interval_end = av
		case "SNAP_SECONDS":
			cs.OracleInstance.snap_seconds, _ = strconv.ParseFloat(av, 64)
		case "EXECS":
			cs.OracleInstance.execs, _ = strconv.ParseInt(av, 10, 64)
		case "EXECS_S":
			cs.OracleInstance.execs_s, _ = strconv.ParseFloat(av, 64)
		case "LOGICAL_READS":
			cs.OracleInstance.logical_reads, _ = strconv.ParseInt(av, 10, 64)
		case "LOGICAL_READS_S":
			cs.OracleInstance.logical_reads_s, _ = strconv.ParseFloat(av, 64)
		case "PHYS_READS":
			cs.OracleInstance.phys_reads, _ = strconv.ParseInt(av, 10, 64)
		case "PHYS_READS_S":
			cs.OracleInstance.phys_reads_s, _ = strconv.ParseFloat(av, 64)
		case "PHYS_WRITES":
			cs.OracleInstance.phys_writes, _ = strconv.ParseInt(av, 10, 64)
		case "PHYS_WRITES_S":
			cs.OracleInstance.phys_writes_s, _ = strconv.ParseFloat(av, 64)
		case "HARD_PARSES":
			cs.OracleInstance.hard_parses, _ = strconv.ParseInt(av, 10, 64)
		case "HARD_PARSES_S":
			cs.OracleInstance.hard_parses_s, _ = strconv.ParseFloat(av, 64)
		case "TOTAL_PARSES":
			cs.OracleInstance.total_parses, _ = strconv.ParseInt(av, 10, 64)
		case "TOTAL_PARSES_S":
			cs.OracleInstance.total_parses_s, _ = strconv.ParseFloat(av, 64)
		case "PARSE_CPU_S":
			cs.OracleInstance.parses_cpu_s, _ = strconv.ParseFloat(av, 64)
		case "PARSE_CPU_S_S":
			cs.OracleInstance.parses_cpu_s_s, _ = strconv.ParseFloat(av, 64)
		case "PARSE_ELAP_S":
			cs.OracleInstance.parse_elap_s, _ = strconv.ParseFloat(av, 64)
		case "PARSE_ELAP_S_S":
			cs.OracleInstance.parse_elap_s_s, _ = strconv.ParseFloat(av, 64)
		case "SORTS_DISK":
			cs.OracleInstance.sorts_disk, _ = strconv.ParseInt(av, 10, 64)
		case "SORTS_DISK_S":
			cs.OracleInstance.sorts_disk_s, _ = strconv.ParseFloat(av, 64)
		case "SORTS_MEM":
			cs.OracleInstance.sorts_mem, _ = strconv.ParseInt(av, 10, 64)
		case "SORTS_MEM_S":
			cs.OracleInstance.sorts_mem_s, _ = strconv.ParseFloat(av, 64)
		case "SORTS_ROWS":
			cs.OracleInstance.sorts_rows, _ = strconv.ParseInt(av, 10, 64)
		case "SORTS_ROWS_S":
			cs.OracleInstance.sorts_rows_s, _ = strconv.ParseFloat(av, 64)
		case "TBL_FETCH_ROWID":
			cs.OracleInstance.tbl_fetch_rowid, _ = strconv.ParseInt(av, 10, 64)
		case "TBL_FETCH_ROWID_S":
			cs.OracleInstance.tbl_fetch_rowid_s, _ = strconv.ParseFloat(av, 64)
		case "TBL_SCAN_ROWS":
			cs.OracleInstance.tbl_scan_rows, _ = strconv.ParseInt(av, 10, 64)
		case "TBL_SCAN_ROWS_S":
			cs.OracleInstance.tbl_scan_rows_s, _ = strconv.ParseFloat(av, 64)
		case "IO_WAIT_S":
			cs.OracleInstance.io_wait_s, _ = strconv.ParseFloat(av, 64)
		case "IO_WAIT_S_S":
			cs.OracleInstance.io_wait_s_s, _ = strconv.ParseFloat(av, 64)
		case "CONCUR_WAIT_S":
			cs.OracleInstance.concur_wait_s, _ = strconv.ParseFloat(av, 64)
		case "CONCUR_WAIT_S_S":
			cs.OracleInstance.concur_wait_s_s, _ = strconv.ParseFloat(av, 64)
		case "USER_CALLS":
			cs.OracleInstance.user_calls, _ = strconv.ParseInt(av, 10, 64)
		case "USER_CALLS_S":
			cs.OracleInstance.user_calls_s, _ = strconv.ParseFloat(av, 64)
		case "USER_COMMITS":
			cs.OracleInstance.user_commits, _ = strconv.ParseInt(av, 10, 64)
		case "USER_COMMITS_S":
			cs.OracleInstance.user_commits_s, _ = strconv.ParseFloat(av, 64)
		case "USER_ROLLBACKS":
			cs.OracleInstance.user_rollbacks, _ = strconv.ParseInt(av, 10, 64)
		case "USER_ROLLBACKS_S":
			cs.OracleInstance.user_rollbacks_s, _ = strconv.ParseFloat(av, 64)
		case "SQLNET_RT_CLIENT":
			cs.OracleInstance.sqlnet_rt_client, _ = strconv.ParseInt(av, 10, 64)
		case "SQLNET_RT_CLIENT_S":
			cs.OracleInstance.sqlnet_rt_client_s, _ = strconv.ParseFloat(av, 64)
		case "SQLNET_RT_DBLINK":
			cs.OracleInstance.sqlnet_rt_dblink, _ = strconv.ParseInt(av, 10, 64)
		case "SQLNET_RT_DBLINK_S":
			cs.OracleInstance.sqlnet_rt_dblink_s, _ = strconv.ParseFloat(av, 64)
		case "SQLNET_BYTES_FROM_CLIENT":
			cs.OracleInstance.sqlnet_bytes_from_client, _ = strconv.ParseInt(av, 10, 64)
		case "SQLNET_BYTES_FROM_CLIENT_S":
			cs.OracleInstance.sqlnet_bytes_from_client_s, _ = strconv.ParseFloat(av, 64)
		case "SQLNET_BYTES_FROM_DBLINK":
			cs.OracleInstance.sqlnet_bytes_from_dblink, _ = strconv.ParseInt(av, 10, 64)
		case "SQLNET_BYTES_FROM_DBLINK_S":
			cs.OracleInstance.sqlnet_bytes_from_dblink_s, _ = strconv.ParseFloat(av, 64)
		case "SQLNET_BYTES_TO_CLIENT":
			cs.OracleInstance.sqlnet_bytes_to_client, _ = strconv.ParseInt(av, 10, 64)
		case "SQLNET_BYTES_TO_CLIENT_S":
			cs.OracleInstance.sqlnet_bytes_to_client_s, _ = strconv.ParseFloat(av, 64)
		case "SQLNET_BYTES_TO_DBLINK":
			cs.OracleInstance.sqlnet_bytes_to_dblink, _ = strconv.ParseInt(av, 10, 64)
		case "SQLNET_BYTES_TO_DBLINK_S":
			cs.OracleInstance.sqlnet_bytes_to_dblink_s, _ = strconv.ParseFloat(av, 64)
		case "WA_EXECS_MULTI":
			cs.OracleInstance.wa_execs_multi, _ = strconv.ParseInt(av, 10, 64)
		case "WA_EXECS_MULTI_S":
			cs.OracleInstance.wa_execs_multi_s, _ = strconv.ParseFloat(av, 64)
		case "WA_EXECS_ONEPASS":
			cs.OracleInstance.wa_execs_onepass, _ = strconv.ParseInt(av, 10, 64)
		case "WA_EXECS_ONEPASS_S":
			cs.OracleInstance.wa_execs_onepass_s, _ = strconv.ParseFloat(av, 64)
		case "WA_EXECS_OPTIMAL":
			cs.OracleInstance.wa_execs_optimal, _ = strconv.ParseInt(av, 10, 64)
		case "WA_EXECS_OPTIMAL_S":
			cs.OracleInstance.wa_execs_optimal_s, _ = strconv.ParseFloat(av, 64)
		default:
			fmt.Println("Nothing")

		}
	}
}

func (cs *Cursor) PrintMetrics() {
	current_hostname, _ := os.Hostname()
	fmt.Fprintf(os.Stdout, "OracleInstanceActivity,host=%s,snap_id=%d snap_id=%d,interval_end=%q,snap_seconds=%.2f,execs=%d,execs_s=%.2f,logical_reads=%d,logical_reads_s=%.2f,phys_reads=%d,phys_reads_s=%.2f,phys_writes=%d,phys_writes_s=%.2f,hard_parses=%d,hard_parses_s=%.2f,total_parses=%d,total_parses_s=%.2f,parses_cpu_s=%.2f,parses_cpu_s_s=%.2f,parse_elap_s=%.2f,parse_elap_s_s=%.2f,sorts_disk=%d,sorts_disk_s=%.2f,sorts_mem=%d,sorts_mem_s=%.2f,sorts_rows=%d,sorts_rows_s=%.2f,tbl_scan_rows=%d,tbl_scan_rows_s=%.2f,tbl_fetch_rowid=%d,tbl_fetch_rowid_s=%.2f,io_wait_s=%.2f,io_wait_s_s=%.2f,concur_wait_s=%.2f,concur_wait_s_s=%.2f,user_calls=%d,user_calls_s=%.2f,user_commits=%d,user_commits_s=%.2f,user_rollbacks=%d,user_rollbacks_s=%.2f,sqlnet_rt_client=%d,sqlnet_rt_client_s=%.2f,sqlnet_rt_dblink=%d,sqlnet_rt_dblink_s=%.2f,sqlnet_bytes_from_client=%d,sqlnet_bytes_from_client_s=%.2f,sqlnet_bytes_from_dblink=%d,sqlnet_bytes_from_dblink_s=%.2f,sqlnet_bytes_to_client=%d,sqlnet_bytes_to_client_s=%.2f,sqlnet_bytes_to_dblink=%d,sqlnet_bytes_to_dblink_s=%.2f,wa_execs_multi=%d,wa_execs_multi_s=%.2f,wa_execs_onepass=%d,wa_execs_onepass_s=%.2f,wa_execs_optimal=%d,wa_execs_optimal_s=%.2f\n",
		current_hostname,
		cs.OracleInstance.snap_id,
		cs.OracleInstance.snap_id,
		cs.OracleInstance.interval_end,
		cs.OracleInstance.snap_seconds,
		cs.OracleInstance.execs,
		cs.OracleInstance.execs_s,
		cs.OracleInstance.logical_reads,
		cs.OracleInstance.logical_reads_s,
		cs.OracleInstance.phys_reads,
		cs.OracleInstance.phys_reads_s,
		cs.OracleInstance.phys_writes,
		cs.OracleInstance.phys_writes_s,
		cs.OracleInstance.hard_parses,
		cs.OracleInstance.hard_parses_s,
		cs.OracleInstance.total_parses,
		cs.OracleInstance.total_parses_s,
		cs.OracleInstance.parses_cpu_s,
		cs.OracleInstance.parses_cpu_s_s,
		cs.OracleInstance.parse_elap_s,
		cs.OracleInstance.parse_elap_s_s,
		cs.OracleInstance.sorts_disk,
		cs.OracleInstance.sorts_disk_s,
		cs.OracleInstance.sorts_mem,
		cs.OracleInstance.sorts_mem_s,
		cs.OracleInstance.sorts_rows,
		cs.OracleInstance.sorts_rows_s,
		cs.OracleInstance.tbl_scan_rows,
		cs.OracleInstance.tbl_scan_rows_s,
		cs.OracleInstance.tbl_fetch_rowid,
		cs.OracleInstance.tbl_fetch_rowid_s,
		cs.OracleInstance.io_wait_s,
		cs.OracleInstance.io_wait_s_s,
		cs.OracleInstance.concur_wait_s,
		cs.OracleInstance.concur_wait_s_s,
		cs.OracleInstance.user_calls,
		cs.OracleInstance.user_calls_s,
		cs.OracleInstance.user_commits,
		cs.OracleInstance.user_commits_s,
		cs.OracleInstance.user_rollbacks,
		cs.OracleInstance.user_rollbacks_s,
		cs.OracleInstance.sqlnet_rt_client,
		cs.OracleInstance.sqlnet_rt_client_s,
		cs.OracleInstance.sqlnet_rt_dblink,
		cs.OracleInstance.sqlnet_rt_dblink_s,
		cs.OracleInstance.sqlnet_bytes_from_client,
		cs.OracleInstance.sqlnet_bytes_from_client_s,
		cs.OracleInstance.sqlnet_bytes_from_dblink,
		cs.OracleInstance.sqlnet_bytes_from_dblink_s,
		cs.OracleInstance.sqlnet_bytes_to_client,
		cs.OracleInstance.sqlnet_bytes_to_client_s,
		cs.OracleInstance.sqlnet_bytes_to_dblink,
		cs.OracleInstance.sqlnet_bytes_to_dblink_s,
		cs.OracleInstance.wa_execs_multi,
		cs.OracleInstance.wa_execs_multi_s,
		cs.OracleInstance.wa_execs_onepass,
		cs.OracleInstance.wa_execs_onepass_s,
		cs.OracleInstance.wa_execs_optimal,
		cs.OracleInstance.wa_execs_optimal_s,
	)
}

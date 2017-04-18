package loadprofile

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"os"
	"pdefcon-for-oracle/utils/sqlutils"
	"strconv"
	//	"strings"
	"time"
)

type OracleLoadProfile struct {
	snap_id                 int64
	interval_end            string
	elap_time               float64
	redo_size_persecond     float64
	logical_reads_persecond float64
	block_changes_persecond float64
	phys_reads_persecond    float64
	phys_wrts_persecond     float64
	user_calls_persecond    float64
	parses_persecond        float64
	hard_parses_persecond   float64
	sorts_persecond         float64
	execs_persecond         float64
	trnxs_persecond         float64
}

type Cursor struct {
	OracleLoadProfile
	cursor    sqlutils.Result
	DbHandler *sql.DB
}

const (
	LoadProfile = `
SELECT
    snap_id                                            AS SNAP,
    TO_CHAR(end_interval_time,'YYYY-MM-DD HH24:MI:SS') AS "INTERVAL_END",
    snap_seconds                                       AS "ELAP_TIME",
    ROUND(redo_size    /snap_seconds,2)                    AS "REDO_SIZE/S",
    ROUND(logical_reads/snap_seconds,2)                    AS "LOGICAL_READS/S",
    ROUND(block_changes  /snap_seconds,2)               AS "BLOCK_CHANGES/S",
    ROUND(physical_reads /snap_seconds,2)               AS "PHYS_READS/S",
    ROUND(physical_writes/snap_seconds,2)               AS "PHYS_WRTS/S",
    ROUND(user_calls     /snap_seconds,2)               AS "USER_CALLS/S",
    ROUND(parses         /snap_seconds,2)               AS "PARSES/S",
    ROUND(hard_parses    /snap_seconds,2)               AS "HARD_PARSES/S",
    ROUND((sorts_disk    +sorts_memory)/snap_seconds,2) AS "SORTS/S",
    ROUND(executions     /snap_seconds,2)               AS "EXECS/S",
    ROUND((commits       + rollbacks)/snap_seconds,2)   AS "TRNXS/S"
FROM
    (
        SELECT
            c.snap_id,
            c.end_interval_time,
            c.snap_seconds,
            MAX (DECODE (c.stat_name, 'redo size', c.VALUE))             AS redo_size,
            MAX (DECODE (c.stat_name, 'session logical reads', c.VALUE)) AS
            logical_reads,
            MAX (DECODE (c.stat_name, 'db block changes', c.VALUE)) AS
            block_changes,
            MAX (DECODE (c.stat_name, 'physical reads', c.VALUE)) AS
            physical_reads,
            MAX (DECODE (c.stat_name, 'physical writes', c.VALUE)) AS
            physical_writes,
            MAX (DECODE (c.stat_name, 'user calls', c.VALUE))          AS user_calls,
            MAX (DECODE (c.stat_name, 'parse count (total)', c.VALUE)) AS
            parses,
            MAX (DECODE (c.stat_name, 'parse count (hard)', c.VALUE)) AS
            hard_parses,
            MAX (DECODE (c.stat_name, 'sorts (disk)', c.VALUE))   AS sorts_disk,
            MAX (DECODE (c.stat_name, 'sorts (memory)', c.VALUE)) AS
            sorts_memory,
            MAX (DECODE (c.stat_name, 'execute count', c.VALUE))  AS executions,
            MAX (DECODE (c.stat_name, 'user commits', c.VALUE))   AS commits,
            MAX (DECODE (c.stat_name, 'user rollbacks', c.VALUE)) AS rollbacks
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
                AND dhss.stat_name IN ('redo size', 'session logical reads',
                    'db block changes', 'physical reads', 'physical writes',
                    'user calls', 'parse count (total)', 'parse count (hard)',
                    'sorts (disk)', 'sorts (memory)', 'execute count',
                    'user commits', 'user rollbacks')
                AND dhs.snap_id IN
                    (
                        SELECT
                            snap_id
                        FROM
                            dba_hist_snapshot
                        WHERE
			 end_interval_time BETWEEN TO_DATE('%s','YYYY-MM-DD HH24:MI:SS') AND TO_DATE('%s','YYYY-MM-DD HH24:MI:SS')
                    )
            )
            c
        WHERE
            instance_restart = 'N'
        AND extract(hour FROM end_interval_time) BETWEEN %d AND
            %d
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
	current_date_start := fmt.Sprintf("%s00:00", time.Now().Format("2006-01-02 00:"))
	current_date_end := fmt.Sprintf("%s23:59:59", time.Now().Format("2006-01-02 "))
	curr_hour := time.Now().Hour()

	query_text := fmt.Sprintf(LoadProfile, current_date_start, current_date_end, curr_hour, curr_hour)
	fmt.Println(query_text)

	for ak, av := range cs.cursor[0] {
		switch ak {
		case "HARD_PARSES/S":
			cs.OracleLoadProfile.hard_parses_persecond, _ = strconv.ParseFloat(av, 64)
		case "INTERVAL_END":
			cs.OracleLoadProfile.interval_end = av
		case "BLOCK_CHANGES/S":
			cs.OracleLoadProfile.block_changes_persecond, _ = strconv.ParseFloat(av, 64)
		case "PARSES/S":
			cs.OracleLoadProfile.parses_persecond, _ = strconv.ParseFloat(av, 64)
		case "ELAP_TIME":
			cs.OracleLoadProfile.elap_time, _ = strconv.ParseFloat(av, 64)
		case "LOGICAL_READS/S":
			cs.OracleLoadProfile.logical_reads_persecond, _ = strconv.ParseFloat(av, 64)
		case "SORTS/S":
			cs.OracleLoadProfile.sorts_persecond, _ = strconv.ParseFloat(av, 64)
		case "EXECS/S":
			cs.OracleLoadProfile.execs_persecond, _ = strconv.ParseFloat(av, 64)
		case "REDO_SIZE/S":
			cs.OracleLoadProfile.redo_size_persecond, _ = strconv.ParseFloat(av, 64)
		case "PHYS_READS/S":
			cs.OracleLoadProfile.phys_reads_persecond, _ = strconv.ParseFloat(av, 64)
		case "USER_CALLS/S":
			cs.OracleLoadProfile.user_calls_persecond, _ = strconv.ParseFloat(av, 64)
		case "TRNXS/S":
			cs.OracleLoadProfile.trnxs_persecond, _ = strconv.ParseFloat(av, 64)
		case "SNAP":
			cs.OracleLoadProfile.snap_id, _ = strconv.ParseInt(av, 10, 64)
		case "PHYS_WRTS/S":
			cs.OracleLoadProfile.phys_wrts_persecond, _ = strconv.ParseFloat(av, 64)
		default:
			fmt.Println("Nothing")
		}
	}
}

func (cs *Cursor) PrintMetrics() {
	current_host, _ := os.Hostname()
	fmt.Fprintf(os.Stdout, "OracleLoadProfile,host=%s,region=LoadProfile hard_parses_persecond=%.2f,interval_end=%q,block_changes_persecond=%.2f,parses_persecond=%.2f,elap_time=%.2f,logical_reads_persecond=%.2f,sorts_persecond=%.2f,execs_persecond=%.2f,redo_size_persecond=%.2f,phys_reads_persecond=%.2f,phys_wrts_persecond=%.2f,user_calls_persecon=%.2f,trnxs_persecon=%.2f,snap_id=%d\n",
		current_host,
		cs.OracleLoadProfile.hard_parses_persecond,
		cs.OracleLoadProfile.interval_end,
		cs.OracleLoadProfile.block_changes_persecond,
		cs.OracleLoadProfile.parses_persecond,
		cs.OracleLoadProfile.elap_time,
		cs.OracleLoadProfile.logical_reads_persecond,
		cs.OracleLoadProfile.sorts_persecond,
		cs.OracleLoadProfile.execs_persecond,
		cs.OracleLoadProfile.redo_size_persecond,
		cs.OracleLoadProfile.phys_reads_persecond,
		cs.OracleLoadProfile.phys_wrts_persecond,
		cs.OracleLoadProfile.user_calls_persecond,
		cs.OracleLoadProfile.trnxs_persecond,
		cs.OracleLoadProfile.snap_id)
}

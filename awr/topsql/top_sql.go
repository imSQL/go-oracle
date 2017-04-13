SELECT
    c.end_interval_time AS "INTERVAL END",
    c.snap_id           AS SNAP,
    c.sql_id            AS "SQL ID",
    c.execs             AS EXECS,
    c.buff_gets         AS "BUFF GETS",
    c.phys_reads        AS "PHYS READS",
    c.cpu_sec           AS "CPU S",
    c.elap_sec          AS "ELAP S",
    c.rows_ret          AS "ROWS RET",
    c.io_wait_sec       AS "IO WAIT S",
    c.exec_rank         AS "EXEC RANK",
    c.buff_get_rank     AS "BUFF GET RANK",
    c.phys_read_rank    AS "PHYS READ RANK",
    c.cpu_rank          AS "CPU RANK",
    c.elap_rank         AS "ELAP RANK",
    c.io_wait_rank      AS "IO WAIT RANK",
    c.concur_wait_sec   AS "CONCUR WAIT S",
    c.parses            AS PARSES,
    c.sorts             AS SORTS,
    c.plsql_sec         AS "PLSQL S",
    c.app_wait_sec      AS "APP WAIT S",
    c.gets_per_exec     AS "GETS PER EXEC",
    c.reads_per_exec    AS "READS PER EXEC",
    c.rows_per_exec     AS "ROWS PER EXEC",
    DBMS_LOB.SUBSTR(DECODE (:show_sql_y_n_t, 'Y', dhst.sql_text, 'T', SUBSTR (
    dhst.sql_text, 1, 1000), NULL ) , 4000, 1)AS sql_text
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
                                    :start_date,'MM/DD/YYYY HH24:MI:SS') AND
                                    TO_DATE(:end_date,'MM/DD/YYYY HH24:MI:SS')
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
            b.buff_get_rank  <= :rank
         OR b.phys_read_rank <= :rank
         OR b.exec_rank      <= :rank
         OR b.cpu_rank       <= :rank
         OR b.elap_rank      <= :rank
         OR b.io_wait_rank   <= :rank
    )
    c,
    dba_hist_sqltext dhst
WHERE
    dhst.sql_id = c.sql_id
ORDER BY
    c.snap_id,
    c.buff_get_rank;
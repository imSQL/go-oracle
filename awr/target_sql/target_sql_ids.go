SELECT
    b.*,
    DECODE (:show_sql_y_n_t, 'Y', dhst.sql_text, 'T', SUBSTR (dhst.sql_text, 1,
    1000), NULL ) AS sql_text
FROM
    (
        SELECT
            TO_CHAR (a.end_interval_time, 'MM/DD/YYYY HH24:MI:SS' ) AS
            "INTERVAL END",
            a.sql_id                                      AS "SQL ID",
            a.plan_hash                                   AS "PLAN HASH",
            a.sql_profile                                 AS "SQL PROFILE",
            a.execs                                       AS execs,
            a.buff_gets                                   AS "BUFF GETS",
            ROUND (a.buff_gets / a.execs, 0)              AS "GETS PER EXE",
            a.phys_reads                                  AS "PHYS READS",
            ROUND (a.phys_reads / a.execs, 0)             AS "READS PER EXE",
            ROUND (a.cpu_time   / 1000000, 2)             AS "CPU S",
            ROUND ((a.cpu_time  / 1000000) / a.execs, 3 ) AS "CPU S PER EXE",
            ROUND (a.elap_time  / 1000000, 2)             AS "ELAP S",
            ROUND ((a.elap_time / 1000000) / a.execs, 3 ) AS "ELAP S PER EXE",
            a.rows_ret                                    AS "ROWS RET",
            ROUND (a.io_wait     / 1000000, 2)                AS "IO WAIT S",
            ROUND (a.concur_wait / 1000000, 2)                AS
            "CONCUR WAIT S",
            a.parses                          AS parses,
            a.sorts                           AS sorts,
            ROUND (a.plsql_time / 1000000, 2) AS "PLSQL S",
            ROUND (a.app_wait   / 1000000, 2) AS "APP WAIT S",
            ROUND (a.rows_ret   / a.execs, 0) AS "ROWS PER EXE"
        FROM
            (
                SELECT
                    dhs.end_interval_time,
                    dhss.sql_id,
                    dhss.plan_hash_value      AS plan_hash,
                    dhss.sql_profile          AS sql_profile,
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
                    dhss.snap_id          = dhs.snap_id
                AND dhss.sql_id           = :sql_id
                AND dhss.executions_delta > 0
            )
            a
        ORDER BY
            1
    )
    b,
    dba_hist_sqltext dhst
WHERE
    dhst.sql_id = b."SQL ID"
ORDER BY
    b."SQL ID",
    to_date(b."INTERVAL END",'MM/DD/YYYY HH24:MI:SS');
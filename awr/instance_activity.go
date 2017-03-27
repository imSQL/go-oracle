SELECT
    snap_id                                            AS SNAP,
    TO_CHAR(end_interval_time,'MM/DD/YYYY HH24:MI:SS') AS "INTERVAL END",
    snap_seconds                                       AS "SNAP SECONDS",
    execs,
    ROUND(execs/snap_seconds,2)                    AS "EXECS/S",
    logical_reads                                  AS "LOGICAL READS",
    ROUND(logical_reads/snap_seconds,2)            AS "LOGICAL READS/S",
    phys_reads                                     AS "PHYS READS",
    ROUND(phys_reads/snap_seconds,2)               AS "PHYS READS/S",
    phys_writes                                    AS "PHYS WRITES",
    ROUND(phys_writes/snap_seconds,2)              AS "PHYS WRITES/S",
    hard_parses                                    AS "HARD PARSES",
    ROUND(hard_parses/snap_seconds,2)              AS "HARD PARSES/S",
    total_parses                                   AS "TOTAL PARSES",
    ROUND(total_parses/snap_seconds,2)             AS "TOTAL PARSES/S",
    parse_cpu_s                                    AS "PARSE CPU S",
    ROUND(parse_cpu_s/snap_seconds,2)              AS "PARSE CPU S/S",
    parse_elap_s                                   AS "PARSE ELAP S",
    ROUND(parse_elap_s/snap_seconds,2)             AS "PARSE ELAP S/S",
    sorts_disk                                     AS "SORTS DISK",
    ROUND(sorts_disk/snap_seconds,2)               AS "SORTS DISK/S",
    sorts_mem                                      AS "SORTS MEM",
    ROUND(sorts_mem/snap_seconds,2)                AS "SORTS MEM/S",
    sorts_rows                                     AS "SORTS ROWS",
    ROUND(sorts_rows/snap_seconds,2)               AS "SORTS ROWS/S",
    tbl_fetch_rowid                                AS "TBL FETCH ROWID",
    ROUND(tbl_fetch_rowid/snap_seconds,2)          AS "TBL FETCH ROWID/S",
    tbl_scan_rows                                  AS "TBL SCAN ROWS",
    ROUND(tbl_scan_rows/snap_seconds,2)            AS "TBL SCAN ROWS/S",
    io_wait_s                                      AS "IO WAIT S",
    ROUND(io_wait_s/snap_seconds,2)                AS "IO WAIT S/S",
    concur_wait_s                                  AS "CONCUR WAIT S",
    ROUND(concur_wait_s/snap_seconds,2)            AS "CONCUR WAIT S/S",
    user_calls                                     AS "USER CALLS",
    ROUND(user_calls/snap_seconds,2)               AS "USER CALLS/S",
    user_commits                                   AS "USER COMMITS",
    ROUND(user_commits/snap_seconds,2)             AS "USER COMMITS/S",
    user_rollbacks                                 AS "USER ROLLBACKS",
    ROUND(user_rollbacks/snap_seconds,2)           AS "USER ROLLBACKS/S",
    sqlnet_rt_client                               AS "SQLNET RT CLIENT",
    ROUND(sqlnet_rt_client/snap_seconds,2)         AS "SQLNET RT CLIENT/S",
    sqlnet_rt_dblink                               AS "SQLNET RT DBLINK",
    ROUND(sqlnet_rt_dblink/snap_seconds,2)         AS "SQLNET RT DBLINK/S",
    sqlnet_bytes_from_client                       AS "SQLNET BYTES FROM CLIENT",
    ROUND(sqlnet_bytes_from_client/snap_seconds,2) AS
    "SQLNET BYTES FROM CLIENT/S",
    sqlnet_bytes_from_dblink                       AS "SQLNET BYTES FROM DBLINK",
    ROUND(sqlnet_bytes_from_dblink/snap_seconds,2) AS
    "SQLNET_BYTES FROM DBLINK/S",
    sqlnet_bytes_to_client                       AS "SQLNET BYTES TO CLIENT",
    ROUND(sqlnet_bytes_to_client/snap_seconds,2) AS "SQLNET_BYTES TO CLIENT/S",
    sqlnet_bytes_to_dblink                       AS "SQLNET BYTES TO DBLINK",
    ROUND(sqlnet_bytes_to_dblink/snap_seconds,2) AS "SQLNET BYTES TO DBLINK/S",
    wa_execs_multi                               AS "WA EXECS MULTI",
    ROUND(wa_execs_multi/snap_seconds,2)         AS "WA EXECS MULTI/S",
    wa_execs_onepass                             AS "WA EXECS ONEPASS",
    ROUND(wa_execs_onepass/snap_seconds,2)       AS "WA EXECS ONEPASS/S",
    wa_execs_optimal                             AS "WA EXECS OPTIMAL",
    ROUND(wa_execs_optimal/snap_seconds,2)       AS "WA EXECS OPTIMAL/S"
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
                            end_interval_time BETWEEN TO_DATE(:start_date,
                            'MM/DD/YYYY HH24:MI:SS') AND TO_DATE(:end_date,
                            'MM/DD/YYYY HH24:MI:SS')
                    )
            )
            c
        WHERE
            instance_restart = 'N'
        AND extract(hour FROM end_interval_time) BETWEEN :beginhour AND
            :endhour
        GROUP BY
            snap_id,
            end_interval_time,
            snap_seconds
        ORDER BY
            snap_id
    );
SELECT
    snap_id                                            AS SNAP,
    TO_CHAR(end_interval_time,'MM/DD/YYYY HH24:MI:SS') AS "INTERVAL END",
    snap_seconds                                       AS "ELAP TIME",
    ROUND(redo_size    /snap_seconds,2)                    AS "REDO SIZE/S",
    ROUND(logical_reads/snap_seconds,2)                    AS "LOGICAL READS/S"
    ,
    ROUND(block_changes  /snap_seconds,2)               AS "BLOCK CHANGES/S",
    ROUND(physical_reads /snap_seconds,2)               AS "PHYS READS/S",
    ROUND(physical_writes/snap_seconds,2)               AS "PHYS WRTS/S",
    ROUND(user_calls     /snap_seconds,2)               AS "USER CALLS/S",
    ROUND(parses         /snap_seconds,2)               AS "PARSES/S",
    ROUND(hard_parses    /snap_seconds,2)               AS "HARD PARSES/S",
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
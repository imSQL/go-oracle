SELECT
    c.snap_id                                                       AS SNAP,
    TO_CHAR(c.end_interval_time,'MM/DD/YYYY HH24:MI:SS')            AS "INTERVAL END",
    c.snap_seconds                                                  AS "SNAP SECONDS",
    MAX (DECODE (c.stat_name, 'DB CPU', c.VALUE))                   AS "DB CPU S",
    MAX (DECODE (c.stat_name, 'DB time', c.VALUE))                  AS "DB TIME S",
    MAX (DECODE (c.stat_name, 'hard parse elapsed time', c.VALUE) ) AS
    "HARD PARSE ELAP S",
    MAX (DECODE (c.stat_name, 'parse time elapsed', c.VALUE) ) AS
    "PARSE ELAP S",
    MAX (DECODE (c.stat_name, 'sql execute elapsed time', c.VALUE) ) AS
    "SQL EXEC ELAP S"
FROM
    (
        SELECT
            dhstm.snap_id,
            dhs.end_interval_time,
            dhstm.stat_name,
            ROUND ( ( dhstm.VALUE                   - LAG (dhstm.VALUE, 1) OVER (PARTITION BY
            dhstm.stat_name ORDER BY dhs.snap_id) ) / 1000000, 2 ) AS VALUE,
            DECODE ( dhs.startup_time               - LAG (dhs.startup_time)
            OVER (PARTITION BY dhstm.stat_name ORDER BY dhs.snap_id), INTERVAL
            '0' SECOND, 'N', 'Y' ) AS instance_restart,
            ( CAST (dhs.end_interval_time AS DATE) - CAST (LAG (
            dhs.end_interval_time, 1) OVER (PARTITION BY dhstm.stat_name
            ORDER BY dhs.snap_id) AS DATE ) ) * 86400 AS snap_seconds
        FROM
            dba_hist_snapshot dhs,
            dba_hist_sys_time_model dhstm
        WHERE
            dhs.snap_id      = dhstm.snap_id
        AND dhstm.stat_name IN ('DB time', 'DB CPU', 'hard parse elapsed time',
            'parse time elapsed', 'sql execute elapsed time')
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
AND extract(hour FROM end_interval_time) BETWEEN :beginhour AND :endhour
GROUP BY
    snap_id,
    end_interval_time,
    snap_seconds
ORDER BY
    snap_id;
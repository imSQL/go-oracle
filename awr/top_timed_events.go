WITH
    snaps AS
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
SELECT
    d.snap_id                                              AS snap,
    TO_CHAR (d.end_interval_time, 'MM/DD/YYYY HH24:MI:SS') AS "INTERVAL END",
    d.event                                                AS event,
    d.waits                                                AS waits,
    ROUND (d.seconds, 0)                                   AS SECONDS,
    ROUND ((d.seconds / e.db_time) * 100, 2)               AS "PCT DB TIME",
    d.wait_class                                           AS CLASS,
    d.time_rank                                            AS RANK
FROM
    (
        SELECT
            c.*,
            RANK () OVER (PARTITION BY c.snap_id ORDER BY c.seconds DESC) AS
            time_rank
        FROM
            (
                SELECT
                    dhse.snap_id,
                    dhs.end_interval_time,
                    dhse.event_name AS event,
                    dhse.wait_class,
                    dhse.total_waits-LAG (dhse.total_waits, 1) OVER (PARTITION
                    BY dhse.event_name ORDER BY dhs.snap_id) AS waits,
                    ROUND((dhse.time_waited_micro-LAG (dhse.time_waited_micro,
                    1) OVER (PARTITION BY dhse.event_name ORDER BY dhs.snap_id)
                    )                      / 1000000,2) AS seconds,
                    DECODE(dhs.startup_time-LAG(dhs.startup_time) OVER (
                    PARTITION BY dhse.event_name ORDER BY dhs.snap_id),
                    INTERVAL '0' SECOND,'N','Y') AS instance_restart
                FROM
                    dba_hist_snapshot dhs,
                    dba_hist_system_event dhse
                WHERE
                    dhs.snap_id  = dhse.snap_id
                AND dhs.snap_id IN
                    (
                        SELECT
                            snap_id
                        FROM
                            snaps
                    )
                AND dhse.wait_class NOT IN ('Idle')
                --                     AND dhse.event_name not in ('Queue
                -- Monitor Task Wait')
                UNION
                SELECT
                    dhstm.snap_id,
                    dhs.end_interval_time,
                    dhstm.stat_name AS event,
                    'CPU'           AS CLASS,
                    NULL            AS waits,
                    ROUND((dhstm.VALUE                    -LAG(dhstm.VALUE, 1) OVER (PARTITION BY
                    dhstm.stat_name ORDER BY dhs.snap_id))/ 1000000,2) AS
                    seconds,
                    DECODE(dhs.startup_time-LAG(dhs.startup_time) OVER (
                    PARTITION BY dhstm.stat_name ORDER BY dhs.snap_id),
                    INTERVAL '0' SECOND, 'N','Y') AS instance_restart
                FROM
                    dba_hist_snapshot dhs,
                    dba_hist_sys_time_model dhstm
                WHERE
                    dhstm.stat_name IN ('DB CPU')
                AND dhs.snap_id      = dhstm.snap_id
                AND dhs.snap_id     IN
                    (
                        SELECT
                            snap_id
                        FROM
                            snaps
                    )
            )
            c
        WHERE
            c.seconds          > 0
        AND c.instance_restart = 'N'
    )
    d,
    (
        SELECT
            snap_id,
            ROUND((dhstm.VALUE                      -LAG(dhstm.VALUE, 1) OVER (PARTITION BY
            dhstm.stat_name ORDER BY dhstm.snap_id))/ 1000000,2) AS db_time
        FROM
            dba_hist_sys_time_model dhstm
        WHERE
            dhstm.stat_name IN ('DB time')
        AND dhstm.snap_id   IN
            (
                SELECT
                    snap_id
                FROM
                    snaps
            )
    )
    e
WHERE
    e.snap_id  = d.snap_id
AND time_rank <= :RANK
ORDER BY
    d.snap_id,
    seconds DESC;
SELECT
    TO_CHAR(end_interval_time,'MM/DD/YYYY HH24:MI') AS "INTERVAL END",
    SUBSTR(eq_type
    || '-'
    || TO_CHAR(NVL(name,' '))
    || DECODE( upper(req_reason) , 'CONTENTION', NULL , '-', NULL , ' ('
    ||req_reason
    ||')'), 1, 78) AS "ENQ TYPE (REASON)",
    waits,
    ROUND(mseconds/1000,2) AS "WAIT S"
FROM
    (
        SELECT
            dhs.end_interval_time,
            dhes.eq_type,
            dhes.req_reason,
            dhes.total_wait#-LAG (dhes.total_wait#, 1) OVER (PARTITION BY
            dhes.eq_type,dhes.req_reason ORDER BY dhs.snap_id) AS waits,
            ROUND((dhes.cum_wait_time-LAG (dhes.cum_wait_time, 1) OVER (
            PARTITION BY dhes.eq_type,dhes.req_reason ORDER BY dhs.snap_id)),2)
            AS mseconds,
            DECODE(dhs.startup_time-LAG(dhs.startup_time) OVER (PARTITION BY
            dhes.eq_type,dhes.req_reason ORDER BY dhs.snap_id), INTERVAL '0'
            SECOND,'N','Y') AS instance_restart
        FROM
            dba_hist_snapshot dhs,
            dba_hist_enqueue_stat dhes
        WHERE
            dhs.snap_id  = dhes.snap_id
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
    a,
    v$lock_type l
WHERE
    a.waits            > 0
AND a.instance_restart = 'N'
AND a.eq_type          = l.type --and name like 'Segment High Water Mark' and
    -- eq_type like 'HW'
AND extract(hour FROM end_interval_time) BETWEEN :beginhour AND :endhour
ORDER BY
    a.end_interval_time,
    a.mseconds DESC,
    a.waits DESC;
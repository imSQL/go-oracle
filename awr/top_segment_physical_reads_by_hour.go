SELECT
    b.*
FROM
    (
        SELECT
            TO_CHAR(a.end_interval_time,'MM/DD/YYYY HH24:MI:SS') AS
            end_interval_time,
            a.object_name,
            a.object_type,
            a.logical_reads  AS LOGICAL_READS,
            a.physical_reads AS PHYSICAL_READS,
            RANK() OVER (PARTITION BY a.end_interval_time ORDER BY
            a.logical_reads DESC) AS LOGICAL_READ_RANK,
            RANK() OVER (PARTITION BY a.end_interval_time ORDER BY
            a.physical_reads DESC) AS PHYSICAL_READ_RANK
        FROM
            (
                SELECT
                    dhs.end_interval_time,
                    dhsso.object_name,
                    dhsso.object_type,
                    dhss.logical_reads_delta  AS logical_reads,
                    dhss.physical_reads_delta AS physical_reads
                FROM
                    dba_hist_seg_stat dhss,
                    dba_hist_seg_stat_obj dhsso,
                    dba_hist_snapshot dhs
                WHERE
                    dhss.snap_id = dhs.snap_id
                AND dhss.obj#    = dhsso.obj#
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
            a
        ORDER BY
            end_interval_time,
            physical_read_rank
    )
    b
WHERE
    b.logical_read_rank  <= :rank
 OR b.physical_read_rank <= :rank;
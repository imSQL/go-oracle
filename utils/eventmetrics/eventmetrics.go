package eventmetrics

const ViewEventMetrics = `
SELECT
    n.wait_class wait_class,
    n.name wait_name,
    m.wait_count cnt,
    ROUND(10*m.time_waited/NULLIF(m.wait_count,0),3) avgms
FROM
    v$eventmetric m,
    v$event_name n
WHERE
    m.event_id    =n.event_id
AND n.wait_class <> 'Idle'
AND m.wait_count  > 0
ORDER BY 1
`

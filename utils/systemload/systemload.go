package systemload

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"pdefcon-for-oracle/utils/sqlutils"
)

type SystemLoad struct {
	sr        sqlutils.Result
	DbHandler *sql.DB
}

const ViewSystemLoad = `
SELECT
    n.wait_class,
    ROUND(m.time_waited/m.intsize_csec,3) aas
FROM
    v$waitclassmetric m,
    v$system_wait_class n
WHERE
    m.wait_class_id =n.wait_class_id
AND n.wait_class   != 'Idle'
UNION
SELECT
    'CPU',
    ROUND(VALUE/100,3) aas
FROM
    v$sysmetric
WHERE
    metric_name='CPU Usage Per Sec'
AND group_id   =2
UNION
SELECT
    'CPU_OS',
    ROUND((prcnt.busy*parameter.cpu_count)/100,3) - aas.cpu
FROM
    (
        SELECT
            VALUE busy
        FROM
            v$sysmetric
        WHERE
            metric_name='Host CPU Utilization (%)'
        AND group_id   =2
    )
    prcnt,
    (
        SELECT
            VALUE cpu_count
        FROM
            v$parameter
        WHERE
            NAME='cpu_count'
    )
    parameter,
    (
        SELECT
            'CPU',
            ROUND(VALUE/100,3) cpu
        FROM
            v$sysmetric
        WHERE
            metric_name='CPU Usage Per Sec'
        AND group_id   =2
    )
    aas
`

func (sl *SystemLoad) GetSystemLoad() {
	sl.sr.GetMetric(sl.DbHandler, ViewSystemLoad)
}

func (sl *SystemLoad) PrintMetrics() {
	fmt.Println(sl.sr)
}

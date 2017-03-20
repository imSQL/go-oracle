package systemload

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"os"
	"pdefcon-for-oracle/utils/sqlutils"
	"sort"
	"strings"
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
	map_length := len(sl.sr)
	current_hostname, _ := os.Hostname()
	sorted_keys := make([]int, 0)
	for k, _ := range sl.sr {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Ints(sorted_keys)

	fmt.Fprintf(os.Stdout, "Oracle,host=%s,region=SystemLoad ", current_hostname)
	for _, k := range sorted_keys {
		fmt.Fprintf(os.Stdout, "%s=%s", strings.Replace(strings.Replace(sl.sr[k]["WAIT_CLASS"], " ", "_", -1), "/", "", -1), strings.Replace(sl.sr[k]["AAS"], ".", "0.", -1))
		if k < map_length-1 {
			fmt.Fprintf(os.Stdout, ",")
		}
	}
	fmt.Println()
}

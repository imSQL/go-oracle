package instance

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"os"
	"pdefcon-for-oracle/utils/sqlutils"
	"strconv"
	"strings"
)

type Instances struct {
	instances sqlutils.Result
	DbHandler *sql.DB
}

const ViewInstance = `
SELECT
	INSTANCE_NUMBER,
	INSTANCE_NAME,
	HOST_NAME,
	VERSION,
	STARTUP_TIME,
	STATUS,
	PARALLEL,
	THREAD#,
	ARCHIVER,
	LOG_SWITCH_WAIT,
	LOGINS,
	SHUTDOWN_PENDING,
	DATABASE_STATUS,
	INSTANCE_ROLE,
	ACTIVE_STATE,
	BLOCKED
FROM v$instance
`

func (ins *Instances) GetMetrics() {
	ins.instances.GetMetric(ins.DbHandler, ViewInstance)
}

func (ins *Instances) PrintMetrics() {
	current_hostname, _ := os.Hostname()

	for _, v := range ins.instances {
		fmt.Fprintf(os.Stdout, "OracleInstance,host=%s,instancename=%s ", current_hostname, v["INSTANCE_NAME"])
		length_v := len(v)
		counter := 0
		for ak, av := range v {
			if counter == length_v-1 {
				if _, ok := strconv.ParseInt(av, 10, 64); ok != nil {
					fmt.Fprintf(os.Stdout, "%s=%q", strings.Replace(strings.ToLower(ak), "#", "", -1), strings.ToLower(av))
				} else {
					fmt.Fprintf(os.Stdout, "%s=%s", strings.Replace(strings.ToLower(ak), "#", "", -1), strings.ToLower(av))

				}
			} else {
				if _, ok := strconv.ParseInt(av, 10, 64); ok != nil {
					fmt.Fprintf(os.Stdout, "%s=%q,", strings.Replace(strings.ToLower(ak), "#", "", -1), strings.ToLower(av))
				} else {
					fmt.Fprintf(os.Stdout, "%s=%s,", strings.Replace(strings.ToLower(ak), "#", "", -1), strings.ToLower(av))

				}
			}
			counter++
		}
		fmt.Fprintf(os.Stdout, "\n")
	}
}

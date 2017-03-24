package datafile

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"os"
	"pdefcon-for-oracle/utils/sqlutils"
	"strconv"
	"strings"
)

type Datafile struct {
	datafiles sqlutils.Result
	DbHandler *sql.DB
}

const ViewDataFiles = `
SELECT 
	FILE_NAME,
	FILE_ID,
	TABLESPACE_NAME,
	BYTES,
	BLOCKS,
	STATUS,
	RELATIVE_FNO,
	AUTOEXTENSIBLE,
	MAXBYTES,
	MAXBLOCKS,
	INCREMENT_BY,
	USER_BYTES,
	USER_BLOCKS,
	ONLINE_STATUS
FROM
    DBA_DATA_FILES
`

func (dfs *Datafile) GetMetrics() {
	dfs.datafiles.GetMetric(dfs.DbHandler, ViewDataFiles)
}

func (dfs *Datafile) PrintMetrics() {
	current_hostname, _ := os.Hostname()

	for _, v := range dfs.datafiles {
		datafilename := strings.Split(v["FILE_NAME"], "/")
		fmt.Fprintf(os.Stdout, "OracleDatafiles,host=%s,datafilename=%s ", current_hostname, datafilename[len(datafilename)-1])
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

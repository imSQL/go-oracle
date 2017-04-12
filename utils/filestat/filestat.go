package filestat

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"os"
	"pdefcon-for-oracle/utils/sqlutils"
	"strconv"
	"strings"
)

type Filestat struct {
	filestat  sqlutils.Result
	DbHandler *sql.DB
}

const (
	FileStat = `
SELECT
    d.name as filename,
    t.name as tbsname,
    f.PHYRDS,
    f.PHYBLKRD,
    f.OPTIMIZED_PHYBLKRD,
    f.SINGLEBLKRDS,
    f.PHYWRTS,
    f.PHYBLKWRT,
    f.READTIM,
    f.WRITETIM,
    f.SINGLEBLKRDTIM,
    f.AVGIOTIM,
    f.LSTIOTIM,
    f.MINIOTIM,
    f.MAXIORTM,
    f.MAXIOWTM
FROM
    v$datafile d,
    v$filestat f,
    v$tablespace t
WHERE
    f.file# = d.file#
AND d.ts#   = t.ts#
`
)

func (fds *Filestat) GetMetrics() {
	fds.filestat.GetMetric(fds.DbHandler, FileStat)
}

func (fds *Filestat) PrintMetrics() {
	current_hostname, _ := os.Hostname()

	for _, v := range fds.filestat {
		datafilename := strings.Split(v["FILENAME"], "/")
		fmt.Fprintf(os.Stdout, "OracleFileStat,host=%s,tablespace=%s,datafile=%s ", current_hostname, strings.ToLower(v["TBSNAME"]), datafilename[len(datafilename)-1])
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

package tablespace

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"os"
	"pdefcon-for-oracle/utils/sqlutils"
	"strconv"
	"strings"
)

type Tablespace struct {
	tablespaces sqlutils.Result
	DbHandler   *sql.DB
}

const ViewTableSpace = `
SELECT
	TABLESPACE_NAME,
	BLOCK_SIZE,
	INITIAL_EXTENT,
	NEXT_EXTENT,
	MIN_EXTENTS,
	MAX_EXTENTS,
	MAX_SIZE,
	PCT_INCREASE,
	MIN_EXTLEN,
	STATUS,
	CONTENTS,
	LOGGING,
	FORCE_LOGGING,
	EXTENT_MANAGEMENT,
	ALLOCATION_TYPE,
	PLUGGED_IN,
	SEGMENT_SPACE_MANAGEMENT,
	DEF_TAB_COMPRESSION,
	RETENTION,
	BIGFILE,
	PREDICATE_EVALUATION,
	ENCRYPTED,
	COMPRESS_FOR
FROM
DBA_TABLESPACES
`

func (tbs *Tablespace) GetMetrics() {
	tbs.tablespaces.GetMetric(tbs.DbHandler, ViewTableSpace)
}

func (tbs *Tablespace) PrintMetrics() {
	current_hostname, _ := os.Hostname()

	for _, v := range tbs.tablespaces {
		fmt.Fprintf(os.Stdout, "OracleTablespace,host=%s,tablespacename=%s ", current_hostname, v["TABLESPACE_NAME"])
		length_v := len(v)
		counter := 0
		for ak, av := range v {
			if len(av) == 0 {
				av = "0"
			}
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

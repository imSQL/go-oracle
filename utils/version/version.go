package version

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"os"
	"pdefcon-for-oracle/utils/sqlutils"
	//	"sort"
	"strings"
)

type Version struct {
	vers      sqlutils.Result
	DbHandler *sql.DB
}

const ViewVersion = `SELECT * FROM V$VERSION`

func (vrs *Version) GetMetrics() {
	vrs.vers.GetMetric(vrs.DbHandler, ViewVersion)
}

func (vrs *Version) PrintMetrics() {
	current_hostname, _ := os.Hostname()
	oracle := strings.Split(vrs.vers[0]["BANNER"], " ")
	fmt.Fprintf(os.Stdout, "OracleProductionInfo,host=%s,region=Version OracleName=%q,OracleVersion=%q,OracleArch=%q\n", current_hostname, strings.Join(oracle[0:5], " "), oracle[6], oracle[8])
}

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

const (
	ViewVersion = `SELECT * FROM V$VERSION`
)

type OracleProductionInfo struct {
	OracleName    string
	OracleVersion string
	OracleArch    string
}

type Version struct {
	OracleProductionInfo
	vers      sqlutils.Result
	DbHandler *sql.DB
}

func (vrs *Version) GetMetrics() {
	vrs.vers.GetMetric(vrs.DbHandler, ViewVersion)
	OracleVersionInfo := strings.Split(vrs.vers[0]["BANNER"], " ")
	vrs.OracleProductionInfo.OracleName = strings.Join(OracleVersionInfo[0:5], " ")
	vrs.OracleProductionInfo.OracleVersion = OracleVersionInfo[6]
	vrs.OracleProductionInfo.OracleArch = OracleVersionInfo[8]
}

func (vrs *Version) PrintMetrics() {
	current_hostname, _ := os.Hostname()
	fmt.Fprintf(os.Stdout,
		"OracleProductionInfo,host=%s,region=Version OracleName=%q,OracleVersion=%q,OracleArch=%q\n",
		current_hostname,
		vrs.OracleProductionInfo.OracleName,
		vrs.OracleProductionInfo.OracleVersion,
		vrs.OracleProductionInfo.OracleArch)
}

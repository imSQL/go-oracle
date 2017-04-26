package pga

import (
	"database/sql"
	"fmt"
	"os"
	"pdefcon-for-oracle/utils/sqlutils"
	"strconv"
	"strings"
)

var (
	escaper = strings.NewReplacer(` `, `_`, `(`, `_`, `)`, `_`, `/`, `_`, `\`, `_`)
)

type OraclePGA struct {
	pganame  string
	pgavalue float64
	pgaunit  string
}

type Cursor struct {
	pga       []OraclePGA
	cursor    sqlutils.Result
	DbHandler *sql.DB
}

const (
	PGA       = `SELECT NAME,VALUE,UNIT FROM V$PGASTAT`
	PGA_ADVIS = `
SELECT
    'LATEST ADVISORY'         AS "PGA ADVISORY",
    b.pga_target_for_estimate AS "PGA TARGET FOR EST",
    b.pga_target_factor       AS "PGA FACTOR",
    b.estd_extra_bytes_rw     AS "EST EXTRA BYTES RW",
    b.estd_overalloc_count    AS "EST OVERALLOC"
FROM
    (
        SELECT
            a.snap_id,
            a.pga_target_for_estimate,
            a.pga_target_factor,
            a.estd_extra_bytes_rw,
            a.estd_overalloc_count,
            RANK () OVER (PARTITION BY a.snap_id ORDER BY pga_target_factor) AS
            RANK
        FROM
            (
                SELECT
                    snap_id,
                    pga_target_for_estimate,
                    TO_NUMBER (pga_target_factor) AS pga_target_factor,
                    estd_extra_bytes_rw,
                    estd_overalloc_count
                FROM
                    dba_hist_pga_target_advice
                WHERE
                    snap_id =
                    (
                        SELECT
                            MAX (snap_id)
                        FROM
                            dba_hist_snapshot
                    )
                AND
                    (
                        TO_NUMBER (pga_target_factor) = 1
                     OR
                        (
                            TO_NUMBER (pga_target_factor) != 1
                        AND estd_overalloc_count           = 0
                        )
                    )
            )
            a
    )
    b
WHERE
    RANK <= 2`
)

func (cs *Cursor) GetMetrics() {
	cs.cursor.GetMetric(cs.DbHandler, PGA)
	for _, val := range cs.cursor {
		tmp := new(OraclePGA)
		for ak, av := range val {
			switch ak {
			case "NAME":
				tmp.pganame = escaper.Replace(av)
			case "VALUE":
				tmp.pgavalue, _ = strconv.ParseFloat(av, 64)
			case "UNIT":
				tmp.pgaunit = av
			default:
				fmt.Println("Nothing")
			}
		}
		cs.pga = append(cs.pga, *tmp)

	}
}

func (cs *Cursor) PrintMetrics() {
	current_hostname, _ := os.Hostname()
	for _, av := range cs.pga {
		fmt.Fprintf(os.Stdout, "OraclePGA,host=%s,pganame=%s pganame=%q,pgavalue=%.2f\n", current_hostname, av.pganame, av.pganame, av.pgavalue)
	}
}

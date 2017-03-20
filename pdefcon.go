package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-oci8"
)

type ID string

const (
	ViewSystemLoad = `
	select n.wait_class, round(m.time_waited/m.INTSIZE_CSEC,3) AAS
        from   v$waitclassmetric  m, v$system_wait_class n
        where m.wait_class_id=n.wait_class_id and n.wait_class != 'Idle'
        union
        select  'CPU', round(value/100,3) AAS
        from v$sysmetric where metric_name='CPU Usage Per Sec' and group_id=2
        union select 'CPU_OS', round((prcnt.busy*parameter.cpu_count)/100,3) - aas.cpu
        from
            ( select value busy
                from v$sysmetric
                where metric_name='Host CPU Utilization (%)'
                and group_id=2 ) prcnt,
                ( select value cpu_count from v$parameter where name='cpu_count' )  parameter,
                ( select  'CPU', round(value/100,3) cpu from v$sysmetric where metric_name='CPU Usage Per Sec' and group_id=2) aas`

	ViewEventMetrics = `select
	n.wait_class wait_class,
       	n.name wait_name,
       	m.wait_count cnt,
       	round(10*m.time_waited/nullif(m.wait_count,0),3) avgms
	from v$eventmetric m,
     	v$event_name n
	where m.event_id=n.event_id
  	and n.wait_class <> 'Idle' and m.wait_count > 0 order by 1 `

	ViewDatabase                 = `SELECT NAME,CREATED,LOG_MODE,OPEN_MODE FROM V$DATABASE`
	ViewInstance                 = `SELECT HOST_NAME,INSTANCE_NAME,VERSION FROM V$INSTANCE`
	ViewVersion                  = `SELECT * FROM V$VERSION`
	ViewControlfile              = `SELECT * FROM V$CONTROLFILE`
	ViewOnlineLogs               = `SELECT GROUP#,MEMBERS,BYTES,STATUS,ARCHIVED FROM V$LOG`
	ViewOnlineLogsFiles          = `SELECT * FROM V$LOGFILE`
	ViewTableSpace               = `SELECT TABLESPACE_NAME,BLOCK_SIZE,STATUS,CONTENTS,LOGGING FROM DBA_TABLESPACES`
	ViewDataFiles                = `SELECT FILE_ID,FILE_NAME,TABLESPACE_NAME,STATUS,BYTES FROM DBA_DATA_FILES`
	ViewUsers                    = `SELECT USERNAME,CREATED FROM DBA_USERS`
	ViewControlFileRecordSection = `SELECT TYPE,RECORD_SIZE,RECORDS_TOTAL,RECORDS_USED FROM V$CONTROLFILE_RECORD_SECTION`
	ViewSGA                      = `SELECT * FROM V$SGAINFO`
	ViewSpparameter              = `SELECT NAME,VALUE FROM V$SPPARAMETER`
)

func (id ID) Scan(src interface{}) error {
	fmt.Println(src)
	return nil
}

func getDSN() string {
	var dsn string
	if len(os.Args) > 1 {
		dsn = os.Args[1]
		if dsn != "" {
			return dsn
		}
	}
	dsn = os.Getenv("GO_OCI8_CONNECT_STRING")
	if dsn != "" {
		return dsn
	}
	fmt.Fprintln(os.Stderr, `Please specifiy connection parameter in GO_OCI8_CONNECT_STRING environment variable,
or as the first argument! (The format is user/name@host:port/sid)`)
	return "scott/tiger@XE"
}

func main() {
	os.Setenv("NLS_LANG", "")

	db, err := sql.Open("oci8", getDSN())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	var wait_class string
	rows, err := db.Query(query_text1)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer rows.Close()

	cols, _ := rows.Columns()

	values := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))

	for i := range values {
		scans[i] = &values[i]
	}

	result := make(map[int]map[string]string)
	i := 0

	for rows.Next() {
		if err := rows.Scan(scans...); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		row := make(map[string]string)

		for k, v := range values {
			key := cols[k]
			row[key] = string(v)
		}

		result[i] = row
		i++
	}

	fmt.Println(result)

	fmt.Println(wait_class)
}

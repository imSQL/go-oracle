package users

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-oci8"
	"os"
	"pdefcon-for-oracle/utils/sqlutils"
	//"sort"
	"strconv"
	"strings"
)

type Users struct {
	usr       sqlutils.Result
	DbHandler *sql.DB
}

const ViewUsers = `
SELECT
    username,
    sid,
    serial#,
    status,
    schemaname,
    osuser,
    machine,
    port,
    program,
    type,
    logon_time,
    event,
    wait_class,
    seconds_in_wait,
    state,
    service_name
FROM
    v$session
WHERE type = 'USER'
`

func (urs *Users) GetMetrics() {
	urs.usr.GetMetric(urs.DbHandler, ViewUsers)
}

func (urs *Users) GetTotalCounter() int {
	return (len(urs.usr))
}

func (urs *Users) GetTotalCounterByUsername() map[string]int {
	result := make(map[string]int)
	for _, v := range urs.usr {
		if len(v["USERNAME"]) == 0 {
			v["USERNAME"] = "UNKNOWN"
		}
		if _, ok := result[v["USERNAME"]]; ok {
			result[v["USERNAME"]]++
		} else {
			result[v["USERNAME"]] = 1
		}
	}
	return result
}

func (urs *Users) GetTotalCounterByStatus() map[string]int {
	result := make(map[string]int)
	for _, v := range urs.usr {
		if len(v["STATUS"]) == 0 {
			v["STATUS"] = "UNKNOWN"
		}
		if _, ok := result[v["STATUS"]]; ok {
			result[v["STATUS"]]++
		} else {
			result[v["STATUS"]] = 1
		}
	}
	return result
}

func (urs *Users) GetTotalCounterByOsuser() map[string]int {
	result := make(map[string]int)
	for _, v := range urs.usr {
		if len(v["OSUSER"]) == 0 {
			v["OSUSER"] = "UNKNOWN"
		}
		if _, ok := result[v["OSUSER"]]; ok {
			result[v["OSUSER"]]++
		} else {
			result[v["OSUSER"]] = 1
		}
	}
	return result
}

func (urs *Users) GetTotalCounterByMachine() map[string]int {
	result := make(map[string]int)
	for _, v := range urs.usr {
		if len(v["MACHINE"]) == 0 {
			v["MACHINE"] = "UNKNOWN"
		}
		if _, ok := result[v["MACHINE"]]; ok {
			result[v["MACHINE"]]++
		} else {
			result[v["MACHINE"]] = 1
		}
	}
	return result
}

func (urs *Users) GetTotalCounterByProgram() map[string]int {
	result := make(map[string]int)
	for _, v := range urs.usr {
		if len(v["PROGRAM"]) == 0 {
			v["PROGRAM"] = "UNKNOWN"
		}
		if _, ok := result[v["PROGRAM"]]; ok {
			result[v["PROGRAM"]]++
		} else {
			result[v["PROGRAM"]] = 1
		}
	}
	return result
}

func (urs *Users) GetTotalCounterBySchemaname() map[string]int {
	result := make(map[string]int)
	for _, v := range urs.usr {
		if len(v["SCHEMANAME"]) == 0 {
			v["SCHEMANAME"] = "UNKNOWN"
		}
		if _, ok := result[v["SCHEMANAME"]]; ok {
			result[v["SCHEMANAME"]]++
		} else {
			result[v["SCHEMANAME"]] = 1
		}
	}
	return result
}

func (urs *Users) GetTotalCounterByService() map[string]int {
	result := make(map[string]int)
	for _, v := range urs.usr {
		if len(v["SERVICE_NAME"]) == 0 {
			v["SERVICE_NAME"] = "UNKNOWN"
		}
		if _, ok := result[v["SERVICE_NAME"]]; ok {
			result[v["SERVICE_NAME"]]++
		} else {
			result[v["SERVICE_NAME"]] = 1
		}
	}
	return result
}

func (urs *Users) PrintMetrics() {
	current_hostname, _ := os.Hostname()
	fmt.Fprintf(os.Stdout, "OracleConnections,host=%s,region=Connections Total=%s\n", current_hostname, strconv.Itoa(urs.GetTotalCounter()))
	for k, v := range urs.GetTotalCounterByUsername() {
		fmt.Fprintf(os.Stdout, "OracleConnections,host=%s,region=ConnectionsByUser,username=%s Total=%s\n", current_hostname, k, strconv.Itoa(v))
	}
	for k, v := range urs.GetTotalCounterByStatus() {
		fmt.Fprintf(os.Stdout, "OracleConnections,host=%s,region=ConnectionsByStatus,status=%s Total=%s\n", current_hostname, k, strconv.Itoa(v))
	}
	for k, v := range urs.GetTotalCounterByOsuser() {
		fmt.Fprintf(os.Stdout, "OracleConnections,host=%s,region=ConnectionsByOsuser,osuser=%s Total=%s\n", current_hostname, k, strconv.Itoa(v))
	}
	for k, v := range urs.GetTotalCounterByMachine() {
		fmt.Fprintf(os.Stdout, "OracleConnections,host=%s,region=ConnectionsByMachine,machine=%s Total=%s\n", current_hostname, k, strconv.Itoa(v))
	}
	for k, v := range urs.GetTotalCounterByProgram() {
		fmt.Fprintf(os.Stdout, "OracleConnections,host=%s,region=ConnectionsByProgram,program=%s Total=%s\n", current_hostname, strings.Replace(k, " ", "_", -1), strconv.Itoa(v))
	}
	for k, v := range urs.GetTotalCounterBySchemaname() {
		fmt.Fprintf(os.Stdout, "OracleConnections,host=%s,region=ConnectionsBySchema,schema=%s Total=%s\n", current_hostname, strings.Replace(k, " ", "_", -1), strconv.Itoa(v))
	}
	for k, v := range urs.GetTotalCounterByService() {
		fmt.Fprintf(os.Stdout, "OracleConnections,host=%s,region=ConnectionsByService,service=%s Total=%s\n", current_hostname, strings.Replace(k, " ", "_", -1), strconv.Itoa(v))
	}
}

package controlfiles

const (
	ViewControlfile              = `SELECT * FROM V$CONTROLFILE`
	ViewControlFileRecordSection = `SELECT TYPE,RECORD_SIZE,RECORDS_TOTAL,RECORDS_USED FROM V$CONTROLFILE_RECORD_SECTION`
)

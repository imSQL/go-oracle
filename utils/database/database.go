package database

import (
	"fmt"
	"os"
)

const ViewDatabase = `SELECT NAME,CREATED,LOG_MODE,OPEN_MODE FROM V$DATABASE`

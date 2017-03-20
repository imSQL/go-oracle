package instance

import (
	"fmt"
	"os"
)

const ViewInstance = `SELECT HOST_NAME,INSTANCE_NAME,VERSION FROM V$INSTANCE`

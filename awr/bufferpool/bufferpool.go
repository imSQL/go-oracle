package bufferpool

import (
	"fmt"
)

const (
	BufferPool = `
SELECT
    NAME,
    PHYSICAL_READS,
    DB_BLOCK_GETS,
    CONSISTENT_GETS,
    SUBSTR(1-(PHYSICAL_READS/(DB_BLOCK_GETS+CONSISTENT_GETS)),1,5) HIT
FROM
    V$BUFFER_POOL_STATISTICS
    `
)

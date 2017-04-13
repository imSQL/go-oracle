package librarycache

import (
	"fmt"
)

const (
	LibraryCacheHitRatio = `
SELECT SUBSTR(SUM(pinhits)/SUM(pins),1,5) library_cache_hit_ratio
FROM V$LIBRARYCACHE
`
	LibraryCacheReloads = `
SELECT namespace,pins,pinhits,reloads
FROM V$LIBRARYCACHE
`
)

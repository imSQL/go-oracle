package dictcache

import (
	"fmt"
)

const (
	DictCache = `
SELECT 
    SUBSTR(SUM(gets - getmisses -fixed)) / SUM(gets),1,5) dict_cache_hit
FROM v$rowcache
`
)

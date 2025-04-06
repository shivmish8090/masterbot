package database

import (
	"fmt"
	"sync"
)

var cache sync.Map

func GetCacheKey(prefix string, id int64) string {
	return fmt.Sprintf("%s:%d", prefix, id)
}
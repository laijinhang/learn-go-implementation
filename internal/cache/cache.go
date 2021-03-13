package cache

import "time"

type Cache struct {
	dir string
	now func() time.Time
}

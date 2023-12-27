package adapters

import (
	"time"
)

type Entry struct {
	Kind      string
	Timestamp time.Time
}

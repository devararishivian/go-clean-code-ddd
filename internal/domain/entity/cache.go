package entity

import "time"

type Cache struct {
	Key   string
	Value any
	TTL   time.Duration
}

package registry

import "time"

type Backend struct {
	Id       string
	Address  string
	LastSeen time.Time
}

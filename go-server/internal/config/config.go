package config

import "time"

const (
	DBPath             = "test.db"
	ServerAddr         = "0.0.0.0:8080"
	MaxActiveUsers     = 100
	ActivateDelay      = 100 * time.Millisecond
	SlowTaskIterations = 200_000
)

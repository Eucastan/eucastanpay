package startup

import "time"

type Result struct {
	Name     string
	Success  bool
	Duration time.Duration
	Error    error
}

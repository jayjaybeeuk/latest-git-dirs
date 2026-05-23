package model

import "time"

type Result struct {
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
	Unix      int64     `json:"unix"`
	Source    string    `json:"source"`
}

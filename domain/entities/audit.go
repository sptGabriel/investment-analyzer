package entities

import "time"

type AuditLog struct {
	IP           string
	Port         string
	Params       string
	Result       string
	ErrorMessage string
	Timestamp    time.Time
}

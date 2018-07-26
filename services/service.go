package services

import (
	"time"
)

// Service represents one service that is to be checked by a checker
type Service struct {
	Name        string                 `json:"name,omitempty"`
	CheckerName string                 `json:"checker_name,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Interval    int                    `json:"interval,omitempty"`
	LastStarted time.Time              `json:"-"`
}

// Services is the list of all active services
var Services []*Service

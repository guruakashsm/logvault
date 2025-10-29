package model

import (
	"time"
)

type Segment struct {
	ID        int             `json:"id"`
	FileName  string          `json:"filename"`
	FilePath  string          `json:"filepath"`
	Logs      []*LogEntry     `json:"logs"`
	Index     *SegmentIndex   `json:"index"`
	StartTime time.Time       `json:"start_time"`
	EndTime   time.Time       `json:"end_time"`
	Metrics   *SegmentMetrics `json:"metrics"`
}

type SegmentIndex struct {
	ByLevel     map[LogLevel][]int `json:"by_level"`
	ByComponent map[string][]int   `json:"by_component"`
	ByHost      map[string][]int   `json:"by_host"`
	ByRequestID map[string][]int   `json:"by_request_id"`
}



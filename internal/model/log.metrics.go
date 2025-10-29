package model

type SegmentMetrics struct {
	ByLevel     map[LogLevel]int
	ByComponent map[string]int
	ByHost      map[string]int
	ByRequestID map[string]int
	TotalLogs   int
}

type LogStoreMetrics struct {
	ByLevel       map[LogLevel]int
	ByComponent   map[string]int
	ByHost        map[string]int
	ByRequestID   map[string]int
	TotalSegments int
	TotalLogs     int
}

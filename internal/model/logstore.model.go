package model

type LogStore struct {
	Segments []*Segment       `json:"segments"`
	Metrics  *LogStoreMetrics `json:"metrics"`
}

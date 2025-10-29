package model

type LogStoreIndex struct {
	ByLevel     map[LogLevel][]int64 `json:"by_level"`
	ByComponent map[string][]int64   `json:"by_component"`
	ByHost      map[string][]int64   `json:"by_host"`
	ByRequestID map[string][]int64   `json:"by_request_id"`
}

package filter

import "github.com/guruakashsm/logvault/internal/model"

type Filter interface {
	Filter(filter LogFilter) []model.LogEntry
	filterBySegment(filter LogFilter, segment *model.Segment) []model.LogEntry
}

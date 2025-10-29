package metrics

import "github.com/guruakashsm/logvault/model"

type MetricsUpdater interface {
	UpdateSegment(metrics *model.SegmentMetrics, entry *model.LogEntry)
	UpdateStore(metrics *model.LogStoreMetrics, segment *model.Segment)
}

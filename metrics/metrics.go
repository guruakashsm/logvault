package metrics

import "github.com/guruakashsm/logvault/model"

type DefaultMetricsUpdater struct{}

func NewMetricsUpdater() MetricsUpdater {
	return &DefaultMetricsUpdater{}
}

func (u *DefaultMetricsUpdater) UpdateSegment(metrics *model.SegmentMetrics, entry *model.LogEntry) {
	metrics.ByLevel[entry.Level]++
	metrics.ByComponent[entry.Component]++
	metrics.ByHost[entry.Host]++
	metrics.ByRequestID[entry.RequestID]++
	metrics.TotalLogs++
}

func (u *DefaultMetricsUpdater) UpdateStore(metrics *model.LogStoreMetrics, segment *model.Segment) {
	for logLevel := range segment.Metrics.ByLevel {
		metrics.ByLevel[logLevel] += segment.Metrics.ByLevel[logLevel]
	}

	for component := range segment.Metrics.ByComponent {
		metrics.ByComponent[component] += segment.Metrics.ByComponent[component]
	}

	for host := range segment.Metrics.ByHost {
		metrics.ByHost[host] += segment.Metrics.ByHost[host]
	}

	for requestID := range segment.Metrics.ByRequestID {
		metrics.ByRequestID[requestID] += segment.Metrics.ByRequestID[requestID]
	}

	metrics.TotalSegments++
	metrics.TotalLogs += segment.Metrics.TotalLogs
}

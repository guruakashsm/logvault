package segmenter

import (
	"github.com/guruakashsm/logvault/internal/metrics"
	"github.com/guruakashsm/logvault/internal/model"
	"github.com/guruakashsm/logvault/internal/parser"
	"github.com/guruakashsm/logvault/internal/segmenter/indexer"
)

type DefaultSegmentHandler struct {
	parser         parser.LogParser
	metricsUpdater metrics.MetricsUpdater
	segment        *model.Segment
	indexer        indexer.Indexer
}

func NewSegmentHandler(parser parser.LogParser, metricsUpdater metrics.MetricsUpdater, segment *model.Segment) SegmentHandler {
	return &DefaultSegmentHandler{
		parser:         parser,
		metricsUpdater: metricsUpdater,
		segment:        segment,
		indexer:        indexer.NewIndexer(segment),
	}
}

func (h *DefaultSegmentHandler) AddLine(line string) error {
	entry, err := h.parser.ParseLine(line)
	if err != nil {
		return err
	}

	h.metricsUpdater.UpdateSegment(h.segment.Metrics, entry)
	h.segment.Logs = append(h.segment.Logs, entry)

	err = h.indexer.AddIndex(entry, h.segment.Metrics.TotalLogs-1)
	if err != nil {
		return err
	}
	

	return nil
}

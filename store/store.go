package store

import (
	"github.com/guruakashsm/logvault/metrics"
	"github.com/guruakashsm/logvault/model"
)

type DefaultStoreHandler struct {
	store          *model.LogStore
	metricsUpdater metrics.MetricsUpdater
}

func NewStoreHandler(store *model.LogStore, metricsUpdater metrics.MetricsUpdater) *DefaultStoreHandler {
	return &DefaultStoreHandler{
		store:          store,
		metricsUpdater: metricsUpdater,
	}
}

func (s *DefaultStoreHandler) AddSegment(segment *model.Segment) error {
	s.store.Segments = append(s.store.Segments, segment)
	s.metricsUpdater.UpdateStore(s.store.Metrics, segment)
	return nil
}

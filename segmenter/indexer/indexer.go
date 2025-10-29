package indexer

import "github.com/guruakashsm/logvault/model"

type DefaultIndexer struct {
	Segment *model.Segment
}

func NewIndexer(segment *model.Segment) Indexer {
	return &DefaultIndexer{Segment: segment}
}

func (i *DefaultIndexer) AddIndex(entry *model.LogEntry, index int) error {
	i.Segment.Index.ByLevel[entry.Level] = append(i.Segment.Index.ByLevel[entry.Level], index)
	i.Segment.Index.ByComponent[entry.Component] = append(i.Segment.Index.ByComponent[entry.Component], index)
	i.Segment.Index.ByHost[entry.Host] = append(i.Segment.Index.ByHost[entry.Host], index)
	i.Segment.Index.ByRequestID[entry.RequestID] = append(i.Segment.Index.ByRequestID[entry.RequestID], index)

	return nil
}

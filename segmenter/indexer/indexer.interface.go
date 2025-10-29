package indexer

import "github.com/guruakashsm/logvault/model"

type Indexer interface {
	AddIndex(segment *model.LogEntry, index int) error
}

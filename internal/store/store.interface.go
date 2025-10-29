package store

import "github.com/guruakashsm/logvault/internal/model"

type Store interface {
	AddSegment(segment *model.Segment) error
}



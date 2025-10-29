package store

import "github.com/guruakashsm/logvault/model"

type Store interface {
	AddSegment(segment *model.Segment) error
}

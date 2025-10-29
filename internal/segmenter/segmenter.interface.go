package segmenter

type SegmentHandler interface {
	AddLine(line string) error
}

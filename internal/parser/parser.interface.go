package parser

import "github.com/guruakashsm/logvault/internal/model"

type LogParser interface {
	ParseLine(line string) (*model.LogEntry, error)
}

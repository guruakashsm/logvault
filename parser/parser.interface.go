package parser

import "github.com/guruakashsm/logvault/model"

type LogParser interface {
	ParseLine(line string) (*model.LogEntry, error)
}

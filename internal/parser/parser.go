package parser

import (
	"fmt"
	"log/slog"
	"regexp"

	"github.com/guruakashsm/logvault/internal/constant"
	"github.com/guruakashsm/logvault/internal/model"
	"github.com/guruakashsm/logvault/utils"
)

type RegexParser struct {
	regex *regexp.Regexp
}

func NewRegexParser(pattern string) *RegexParser {
	return &RegexParser{
		regex: regexp.MustCompile(pattern),
	}
}

func (p *RegexParser) ParseLine(line string) (*model.LogEntry, error) {
	matches := p.regex.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("log line does not match pattern")
	}

	timestamp, err := utils.ParseTime(matches[p.regex.SubexpIndex(constant.Timestamp)], constant.TimeFormats...)
	if err != nil {
		slog.Error("Failed to parse timestamp", "error", err)
		return nil, err
	}

	return &model.LogEntry{
		Level:     model.LogLevel(matches[p.regex.SubexpIndex(constant.Level)]),
		Component: matches[p.regex.SubexpIndex(constant.Component)],
		Host:      matches[p.regex.SubexpIndex(constant.Host)],
		RequestID: matches[p.regex.SubexpIndex(constant.RequestID)],
		Message:   matches[p.regex.SubexpIndex(constant.Message)],
		Metadata:  make(map[string]string),
		Log:       line,
		Timestamp: *timestamp,
	}, nil
}

package filter

import (
	"time"

	"github.com/guruakashsm/logvault/internal/model"
)

type LogFilter struct {
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`

	Levels     []model.LogLevel `json:"levels,omitempty"`
	Components []string         `json:"components,omitempty"`
	Hosts      []string         `json:"hosts,omitempty"`
	RequestIDs []string         `json:"request_ids,omitempty"`

	MessageContains *string `json:"message_contains,omitempty"`
	LogContains     *string `json:"log_contains,omitempty"`

	// Metadata map[string]string `json:"metadata,omitempty"`
}

type DefaultFilterHandler struct {
	LogStore *model.LogStore `json:"log_store"`
}

func NewFilterHandler(logStore *model.LogStore) *DefaultFilterHandler {
	return &DefaultFilterHandler{
		LogStore: logStore,
	}
}

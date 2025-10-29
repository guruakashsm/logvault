package filter

import "github.com/guruakashsm/logvault/internal/model"

func (f *DefaultFilterHandler) Remove(filter LogFilter) []model.LogEntry {
	return []model.LogEntry{}
}

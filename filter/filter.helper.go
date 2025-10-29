package filter

import "github.com/guruakashsm/logvault/model"

func (f *DefaultFilterHandler) Remove(filter LogFilter) []model.LogEntry {
	return []model.LogEntry{}
}

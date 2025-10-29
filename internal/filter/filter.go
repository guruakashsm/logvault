package filter

import (
	"strings"

	"github.com/guruakashsm/logvault/internal/model"
)

func (f *DefaultFilterHandler) Filter(filter LogFilter) []model.LogEntry {
	filteredSegments := make([]*model.Segment, 0, len(f.LogStore.Segments))

	levelSet := toSet(filter.Levels)
	componentSet := toSet(filter.Components)
	hostSet := toSet(filter.Hosts)
	reqIDSet := toSet(filter.RequestIDs)

	for _, seg := range f.LogStore.Segments {
		
		if filter.StartTime != nil && seg.EndTime.Before(*filter.StartTime) {
			continue
		}

		if filter.EndTime != nil && seg.StartTime.After(*filter.EndTime) {
			continue
		}

		if len(levelSet) > 0 && !hasAny(seg.Metrics.ByLevel, levelSet) {
			continue
		}

		if len(componentSet) > 0 && !hasAny(seg.Metrics.ByComponent, componentSet) {
			continue
		}

		if len(hostSet) > 0 && !hasAny(seg.Metrics.ByHost, hostSet) {
			continue
		}
		
		if len(reqIDSet) > 0 && !hasAny(seg.Metrics.ByRequestID, reqIDSet) {
			continue
		}

		filteredSegments = append(filteredSegments, seg)
	}

	var filteredEntries []model.LogEntry
	for _, seg := range filteredSegments {
		filteredEntries = append(filteredEntries, f.FilterBySegment(filter, seg)...)
	}

	return filteredEntries
}

func (f *DefaultFilterHandler) FilterBySegment(filter LogFilter, segment *model.Segment) []model.LogEntry {
	var filteredEntries []model.LogEntry

	levelSet := toSet(filter.Levels)
	componentSet := toSet(filter.Components)
	hostSet := toSet(filter.Hosts)
	reqIDSet := toSet(filter.RequestIDs)

	for _, entry := range segment.Logs {
		if entry == nil {
			continue
		}

		if filter.StartTime != nil && entry.Timestamp.Before(*filter.StartTime) {
			continue
		}
		if filter.EndTime != nil && entry.Timestamp.After(*filter.EndTime) {
			continue
		}
		if len(levelSet) > 0 && !contains(levelSet, entry.Level) {
			continue
		}
		if len(componentSet) > 0 && !contains(componentSet, entry.Component) {
			continue
		}
		if len(hostSet) > 0 && !contains(hostSet, entry.Host) {
			continue
		}
		if len(reqIDSet) > 0 && !contains(reqIDSet, entry.RequestID) {
			continue
		}
		if filter.MessageContains != nil && !strings.Contains(entry.Message, *filter.MessageContains) {
			continue
		}
		if filter.LogContains != nil && !strings.Contains(entry.Log, *filter.LogContains) {
			continue
		}

		filteredEntries = append(filteredEntries, *entry)
	}

	return filteredEntries
}

func toSet[T comparable](arr []T) map[T]struct{} {
	if len(arr) == 0 {
		return nil
	}
	set := make(map[T]struct{}, len(arr))
	for _, v := range arr {
		set[v] = struct{}{}
	}
	return set
}

func contains[T comparable](set map[T]struct{}, v T) bool {
	if len(set) == 0 {
		return true
	}
	_, ok := set[v]
	return ok
}

func hasAny[K comparable](m map[K]int, set map[K]struct{}) bool {
	for key := range set {
		if m[key] > 0 {
			return true
		}
	}
	return false
}

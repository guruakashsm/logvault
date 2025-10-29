package utils

import (
	"time"

	iError "github.com/guruakashsm/logvault/error"
)

func ParseTime(s string, formats ...string) (*time.Time, error) {
	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			return &t, nil
		}
	}
	return nil, iError.ErrFailedtoParseTime
}

package utils

import (
	"os"
	"time"
)

var timezone *time.Location

func ShortDatetimeWithTZ(datetimeStr string) string {
	parsedTime, err := time.Parse(time.RFC3339, datetimeStr)
	if err != nil {
		return "(invalid)"
	}

	localTime := parsedTime.In(timezone)
	return localTime.Format("2006-01-02 15:04")
}

func init() {
	tz := os.Getenv("TZ")
	if tz == "" {
		timezone = time.Local
	} else {
		var err error
		timezone, err = time.LoadLocation(tz)
		if err != nil {
			timezone = time.Local
		}
	}
}

package utils

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func JsonFromStringArray(source []string) string {
	if source == nil || len(source) == 0 {
		return "[]"
	} else {
		jsonBytes, err := json.Marshal(source)
		if err != nil {
			logrus.Fatalf("Failed to marshal string array: %v", err)
		}
		return string(jsonBytes)
	}
}

func StringArrayFromJson(source string) []string {
	if source == "" {
		return []string{}
	}
	var result []string
	if err := json.Unmarshal([]byte(source), &result); err != nil {
		logrus.Fatalf("Failed to unmarshal string array: %v", err)
	}
	return result
}

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

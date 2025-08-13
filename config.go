package datepicker

import (
	"time"
)

type Config struct {
	FirstWeekdayIsMo bool
	OutputFormat     string
	StartAt          time.Time
	HideHelp         bool
}

func DefaultConfig() Config {
	today, _ := time.Parse("2006/01/02", time.Now().Format("2006/01/02"))
	return Config{
		FirstWeekdayIsMo: true,
		OutputFormat:     "2006/01/02",
		StartAt:          today,
		HideHelp:         false,
	}
}

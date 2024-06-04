package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Hour int

const (
	MinHour    = 0
	MaxHour    = 23
	minutes    = 1
	MinsInHour = 60 * minutes
)

func GetHour(hours, minutes int) Hour {
	return (Hour)(hours*60 + minutes)
}

func ParseHour(s string) (hour Hour, err error) {

	intHour := 0

	s = strings.ToLower(s)

	if strings.Contains(s, "pm") {
		intHour += 12 * MinsInHour
	}

	s = strings.Replace(s, "am", "", 1)
	s = strings.Replace(s, "pm", "", 1)

	s = strings.ReplaceAll(s, "-", ":")
	s = strings.ReplaceAll(s, ".", ":")

	nums := strings.Split(s, ":")

	if len(nums) != 2 {
		return hour, errors.New("invalid hour, should be \"12:00am\"")
	}

	hours, hoursErr := strconv.Atoi(nums[0])
	minutes, minutesErr := strconv.Atoi(nums[1])

	if hoursErr != nil {
		return hour, errors.New("error parsing hours")
	}
	if minutesErr != nil {
		return hour, errors.New("error parsing minutes")
	}

	intHour += MinsInHour * (hours % 12)
	intHour += minutes

	hour = (Hour)(intHour)
	return
}

func (hour Hour) AsDuration() time.Duration {
	return time.Duration((time.Duration)(hour) * time.Minute)
}

package timespec

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func GetTimeStamp(spec string) (timestamp time.Time, err error) {
	// is it a unix timestamp?
	i64, err := strconv.ParseInt(spec, 10, 32)
	if err != nil {
		return time.Unix(i64, 0), nil
	}
	now := time.Now()
	re := regexp.MustCompile("(\\+|-)?([0-9]+)(second|minute|min|m|hour|h|day|d|D|week|w|month|mo)s?")
	matches := re.FindStringSubmatch(spec)
	if len(matches) == 0 {
		err = errors.New(fmt.Sprintf("could not parse '%s'", spec))
		return
	}
	duration_i64, _ := strconv.ParseUint(matches[1]+matches[2], 10, 32)
	var duration time.Duration
	// not always technically correct, but it doesn't need to be
	// because we just want timestamps that are more or less correct
	// to get a given timerange
	switch matches[3] {
	case "second":
		duration = time.Second
	case "minute", "min", "m":
		duration = time.Minute
	case "hour", "h":
		duration = time.Hour
	case "day", "d", "D":
		duration = time.Hour * 24
	case "week", "w":
		duration = time.Hour * 24 * 7
	case "month", "mo":
		duration = time.Hour * 24 * 30
	}
	timestamp = now.Add(duration * time.Duration(duration_i64))
	return timestamp, nil
}

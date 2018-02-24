package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// This files is comes from moby's code
// https://github.com/moby/moby/blob/3a633a712c8bbb863fe7e57ec132dd87a9c4eff7/api/types/time/timestamp.go

// These are additional predefined layouts for use in Time.Format and Time.Parse
// with --since and --until parameters for `docker logs` and `docker events`
const (
	rFC3339Local     = "2006-01-02T15:04:05"           // RFC3339 with local timezone
	rFC3339NanoLocal = "2006-01-02T15:04:05.999999999" // RFC3339Nano with local timezone
	dateWithZone     = "2006-01-02Z07:00"              // RFC3339 with time at 00:00:00
	dateLocal        = "2006-01-02"                    // RFC3339 with local timezone and time at 00:00:00
)

var relativeRe = regexp.MustCompile(`(\d+)\s?(m|minute|minutes|h|hour|hours|d|day|days|w|week|weeks)(?: ago)?`)

func parseTime(value string) (time.Time, error) {
	m := relativeRe.FindAllStringSubmatch(strings.TrimSpace(value), 1)
	if m != nil {
		t, err := parseRelativeTime(m[0], time.Now())
		if err == nil {
			return t, err
		}
		// if err, try other parse
	}

	switch value {
	case "all":
		return time.Unix(0, 0), nil
	case "", "now":
		return time.Now(), nil
	}

	return getTimestamp(value, time.Now())
}

// GetTimestamp tries to parse given string as golang duration,
// then RFC3339 time and finally as a Unix timestamp. If
// any of these were successful, it returns a Unix timestamp
// as string otherwise returns the given value back.
// In case of duration input, the returned timestamp is computed
// as the given reference time minus the amount of the duration.
func getTimestamp(value string, reference time.Time) (time.Time, error) {
	if d, err := time.ParseDuration(value); value != "0" && err == nil {
		return reference.Add(-d), nil
	}

	var format string
	// if the string has a Z or a + or three dashes use parse otherwise use parseinlocation
	parseInLocation := !(strings.ContainsAny(value, "zZ+") || strings.Count(value, "-") == 3)

	if strings.Contains(value, ".") {
		if parseInLocation {
			format = rFC3339NanoLocal
		} else {
			format = time.RFC3339Nano
		}
	} else if strings.Contains(value, "T") {
		// we want the number of colons in the T portion of the timestamp
		tcolons := strings.Count(value, ":")
		// if parseInLocation is off and we have a +/- zone offset (not Z) then
		// there will be an extra colon in the input for the tz offset subtract that
		// colon from the tcolons count
		if !parseInLocation && !strings.ContainsAny(value, "zZ") && tcolons > 0 {
			tcolons--
		}
		if parseInLocation {
			switch tcolons {
			case 0:
				format = "2006-01-02T15"
			case 1:
				format = "2006-01-02T15:04"
			default:
				format = rFC3339Local
			}
		} else {
			switch tcolons {
			case 0:
				format = "2006-01-02T15Z07:00"
			case 1:
				format = "2006-01-02T15:04Z07:00"
			default:
				format = time.RFC3339
			}
		}
	} else if parseInLocation {
		format = dateLocal
	} else {
		format = dateWithZone
	}

	var t time.Time
	var err error

	if parseInLocation {
		t, err = time.ParseInLocation(format, value, time.FixedZone(reference.Zone()))
	} else {
		t, err = time.Parse(format, value)
	}

	if err != nil {
		if strings.Contains(value, "-") {
			return time.Unix(0, 0), err // was probably an RFC3339 like timestamp but the parser failed with an error
		}
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return time.Unix(0, 0), err
		}

		t = time.Unix(int64(intVal), 0)
	}

	return t, nil
}

func parseRelativeTime(m []string, t time.Time) (time.Time, error) {
	if len(m) != 3 {
		return t, fmt.Errorf("invalid Relative Time")
	}
	p, err := strconv.Atoi(m[1])
	if err != nil {
		return t, err
	}
	var d time.Duration
	switch m[2] {
	case "m", "minute", "minutes":
		d = -time.Duration(int(p)) * time.Minute
	case "h", "hour", "hours":
		d = -time.Duration(int(p)) * time.Hour
	case "d", "day", "days":
		d = -time.Duration(int(p)) * 24 * time.Hour
	case "w", "week", "weeks":
		d = -time.Duration(int(p)) * 24 * time.Hour * 7
	}

	return t.Add(d), nil
}

package main

import (
	"fmt"
	"testing"
	"time"
)

func TestParseTimeSpecial(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		for _, all := range []string{"all", "ALL"} {
			p, err := parseTime(all)
			if err != nil {
				t.Error(err)
			}
			if p.UTC().Format(time.RFC3339) != "1970-01-01T00:00:00Z" {
				t.Error(fmt.Errorf("not expected: %s", p.UTC().Format(time.RFC3339)))
			}
		}
	})
	t.Run("empty", func(t *testing.T) {
		for _, now := range []string{"", "now"} {
			p, err := parseTime(now)
			if err != nil {
				t.Error(err)
			}
			if p.UTC().Format(time.RFC3339) == "1970-01-01T00:00:00Z" {
				t.Error(fmt.Errorf("not expected: %s", p.UTC().Format(time.RFC3339)))
			}
		}
	})
	t.Run("error", func(t *testing.T) {
		p, err := parseTime("HOGEHOGE")
		if err == nil {
			t.Error(fmt.Errorf("err should happen: %s", p))
		}
	})
}
func TestParseRelativeTime(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2018-10-30T10:10:00Z")

	var tests = []struct {
		in       []string
		expected string
	}{
		{[]string{"1m", "1", "m"}, "2018-10-30T10:09:00Z"},
		{[]string{"1min", "1", "min"}, "2018-10-30T10:09:00Z"},
		{[]string{"1h", "1", "h"}, "2018-10-30T09:10:00Z"},
		{[]string{"1hour", "1", "hour"}, "2018-10-30T09:10:00Z"},
		{[]string{"1hours", "1", "hours"}, "2018-10-30T09:10:00Z"},
		{[]string{"1d", "1", "d"}, "2018-10-29T10:10:00Z"},
		{[]string{"1day", "1", "day"}, "2018-10-29T10:10:00Z"},
		{[]string{"1days", "1", "days"}, "2018-10-29T10:10:00Z"},
		{[]string{"1w", "1", "w"}, "2018-10-23T10:10:00Z"},
		{[]string{"1week", "1", "week"}, "2018-10-23T10:10:00Z"},
		{[]string{"1weeks", "1", "weeks"}, "2018-10-23T10:10:00Z"},
	}

	for _, test := range tests {
		p, err := parseRelativeTime(test.in, now)
		if err != nil {
			t.Error(err)
		}
		if p.Format(time.RFC3339) != test.expected {
			t.Error(fmt.Errorf("%s is not expected: %s", p.Format(time.RFC3339), test.expected))
		}
	}
}

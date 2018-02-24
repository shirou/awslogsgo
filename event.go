package main

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

// Event represents a log event
type Event struct {
	Event         cloudwatchlogs.FilteredLogEvent
	Id            string
	Group         string
	Stream        string
	Message       string
	IngestionTime time.Time
	Timestamp     time.Time
}

func fromFilteredLogEvent(group string, src cloudwatchlogs.FilteredLogEvent) Event {
	return Event{
		Event:         src,
		Id:            *src.EventId,
		Group:         group,
		Stream:        *src.LogStreamName,
		Message:       *src.Message,
		IngestionTime: time.Unix(0, *src.IngestionTime*1000000),
		Timestamp:     time.Unix(0, *src.Timestamp*1000000),
	}
}

package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/hashicorp/golang-lru"
)

const (
	maxEventsBuffer = 10000
	maxEventsCache  = 100000
	watchSleepTime  = 1 // SleepTime when watch is enabled
)

// When sleep
var watchStart = -time.Duration(3) * time.Minute

// CloudwatchLogsReader is responsible for fetching logs for a particular log
// group
type CloudwatchLogsReader struct {
	svc          *cloudwatchlogs.CloudWatchLogs
	config       aws.Config
	logGroupName string
	start        time.Time
	end          time.Time
	filter       string
	streamPrefix string
	eventCache   *lru.Cache
}

// NewCloudwatchLogsReader takes a group and optionally a stream prefix, start and
// end time, and returns a reader for any logs that match those parameters.
func NewCloudwatchLogsReader(config aws.Config,
	group, streamPrefix, filter string,
	start time.Time, end time.Time) (*CloudwatchLogsReader, error) {
	svc := cloudwatchlogs.New(config)

	// check that group is exists
	if groupExists(config, group) == false {
		return nil, fmt.Errorf("group %s does not exists", group)
	}
	cache, err := lru.New(maxEventsCache)
	if err != nil {
		return nil, err
	}

	reader := &CloudwatchLogsReader{
		config:       config,
		svc:          svc,
		logGroupName: group,
		start:        start,
		end:          end,
		filter:       filter,
		streamPrefix: streamPrefix,
		eventCache:   cache,
	}

	return reader, nil
}

func (reader *CloudwatchLogsReader) Stream(watch bool) (chan Event, error) {
	stream := make(chan Event, maxEventsBuffer)

	go reader.startStream(stream, watch)

	return stream, nil
}

func (reader *CloudwatchLogsReader) startStream(stream chan Event, watch bool) {
	ss, err := ListStreams(reader.config,
		reader.logGroupName, reader.streamPrefix,
		reader.start, reader.end)
	if err != nil {
		fmt.Println(err)
		close(stream)
		return
	}

	// FilterLogEventsInput can not use more than 100 streams.
	if len(ss) > 100 {
		ss = ss[0:99]
	}

	params := &cloudwatchlogs.FilterLogEventsInput{
		StartTime:      aws.Int64(aws.TimeUnixMilli(reader.start)),
		EndTime:        aws.Int64(aws.TimeUnixMilli(reader.end)),
		LogGroupName:   aws.String(reader.logGroupName),
		LogStreamNames: ss,
		Interleaved:    aws.Bool(true),
	}
	if reader.filter != "" {
		params.FilterPattern = aws.String(reader.filter)
	}

LOOP:
	req := reader.svc.FilterLogEventsRequest(params)
	p := req.Paginate()
	for p.Next() {
		page := p.CurrentPage()
		for _, e := range page.Events {
			if _, ok := reader.eventCache.Peek(*e.EventId); !ok {
				stream <- fromFilteredLogEvent(reader.logGroupName, e)
				reader.eventCache.Add(*e.EventId, nil)
			}
		}
	}
	if err := p.Err(); err != nil {
		fmt.Println(err)
	}

	if watch {
		time.Sleep(watchSleepTime)
		params.StartTime = aws.Int64(aws.TimeUnixMilli(time.Now().Add(watchStart)))
		params.EndTime = aws.Int64(aws.TimeUnixMilli(time.Now()))
		goto LOOP
	}

	close(stream)
}

func groupExists(config aws.Config, group string) bool {
	svc := cloudwatchlogs.New(config)
	describeLogGroupsInput := &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(group),
	}
	req := svc.DescribeLogGroupsRequest(describeLogGroupsInput)
	p, err := req.Send()
	if err != nil {
		return false
	}
	if len(p.LogGroups) == 0 {
		return false
	}
	return true
}

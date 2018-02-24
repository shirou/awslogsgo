package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

const (
	MaxEventsBuffer = 10000
)

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
}

// NewCloudwatchLogsReader takes a group and optionally a stream prefix, start and
// end time, and returns a reader for any logs that match those parameters.
func NewCloudwatchLogsReader(config aws.Config,
	group, streamPrefix, filter string,
	start time.Time, end time.Time) (*CloudwatchLogsReader, error) {
	svc := cloudwatchlogs.New(config)

	// check that group is exists
	if groupExists(config, group) == false {
		return nil, fmt.Errorf("group %s is not exists", group)
	}
	reader := &CloudwatchLogsReader{
		config:       config,
		svc:          svc,
		logGroupName: group,
		start:        start,
		end:          end,
		filter:       filter,
		streamPrefix: streamPrefix,
	}

	return reader, nil
}

func (reader *CloudwatchLogsReader) Stream(follow bool) (chan Event, error) {
	stream := make(chan Event, MaxEventsBuffer)

	go reader.startStream(stream, follow)

	return stream, nil
}

func (reader *CloudwatchLogsReader) startStream(stream chan Event, follow bool) {
	ss, err := ListStreams(reader.config,
		reader.logGroupName, reader.streamPrefix,
		reader.start, reader.end)
	if err != nil {
		fmt.Println(err)
		close(stream)
		return
	}
	params := &cloudwatchlogs.FilterLogEventsInput{
		StartTime:      aws.Int64(aws.TimeUnixMilli(reader.start)),
		EndTime:        aws.Int64(aws.TimeUnixMilli(reader.end)),
		LogGroupName:   aws.String(reader.logGroupName),
		LogStreamNames: ss,
	}
	if reader.filter != "" {
		params.FilterPattern = aws.String(reader.filter)
	}

	req := reader.svc.FilterLogEventsRequest(params)
	p := req.Paginate()
	for p.Next() {
		page := p.CurrentPage()
		if len(page.Events) == 0 && follow == false {
			close(stream)
			return
		}
		for _, e := range page.Events {
			stream <- fromFilteredLogEvent(reader.logGroupName, e)
		}
	}
	if err := p.Err(); err != nil {
		fmt.Println(err)
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

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func ListGroup(config aws.Config, prefix string) error {
	svc := cloudwatchlogs.New(config)

	describeLogGroupsInput := &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(prefix),
	}

	ctx, cancel := context.WithTimeout(context.Background(), callDeadLine)
	defer cancel()

	req := svc.DescribeLogGroupsRequest(describeLogGroupsInput)
	p := req.Paginate()
	for p.Next(ctx) {
		page := p.CurrentPage()
		if len(page.LogGroups) == 0 {
			return fmt.Errorf("Could not find log group '%s'", prefix)
		}
		for _, group := range page.LogGroups {
			fmt.Println(*group.LogGroupName)
		}

	}
	if err := p.Err(); err != nil {
		return err
	}
	return nil
}

func ListStreams(config aws.Config, group, prefix string, start, end time.Time) ([]string, error) {
	svc := cloudwatchlogs.New(config)

	describeLogStreamsInput := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(group),
	}
	if prefix != "" {
		describeLogStreamsInput.LogStreamNamePrefix = aws.String(prefix)
	}

	ctx, cancel := context.WithTimeout(context.Background(), callDeadLine)
	defer cancel()

	ret := make([]string, 0)
	req := svc.DescribeLogStreamsRequest(describeLogStreamsInput)
	p := req.Paginate()
	for p.Next(ctx) {
		page := p.CurrentPage()
		if len(page.LogStreams) == 0 {
			return nil, fmt.Errorf("Could not find log streams for %s, '%s'", group, prefix)
		}
		for _, group := range page.LogStreams {
			ret = append(ret, *group.LogStreamName)
		}
	}
	if err := p.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
)

func awsConfig() (aws.Config, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return aws.Config{}, err
	}

	//	return &aws.Config{MaxRetries: aws.Int(10)}, nil
	return cfg, err
}

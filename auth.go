package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/urfave/cli"
)

func awsConfigArgs(c *cli.Context) (string, string, string, string, string) {
	return c.String("aws-access-key-id"),
		c.String("aws-secret-access-key"),
		c.String("aws-session-token"),
		c.String("profile"),
		c.String("aws-region")
}

func awsConfig(keyid, secretid, token, profile, region string) (aws.Config, error) {
	configs := make([]external.Config, 0)

	if keyid != "" && secretid != "" {
		credentials := aws.Credentials{
			AccessKeyID:     keyid,
			SecretAccessKey: secretid,
		}
		configs = append(configs, external.WithCredentialsValue(credentials))
	} else if token != "" {
		credentials := aws.Credentials{
			SessionToken: token,
		}
		configs = append(configs, external.WithCredentialsValue(credentials))
	}
	if profile != "" {
		configs = append(configs, external.WithSharedConfigProfile(profile))
	}

	if region != "" {
		configs = append(configs, external.WithRegion(region))
	}

	cfg, err := external.LoadDefaultAWSConfig(configs...)
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, err
}

package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	bucket     = "ginput"
	cloudfront = "https://d17enza3bfujl8.cloudfront.net"
)

func Images() ([]string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	}))
	svc := s3.New(sess)
	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		return nil, err
	}

	var images []string
	for _, i := range result.Contents {
		images = append(images, fmt.Sprintf("%s/%s", cloudfront, *i.Key))
	}

	return images, nil
}

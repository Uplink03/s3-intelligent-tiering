package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func listObjects(svc *s3.S3) {
	marker := (*string)(nil)

	for {
		fmt.Printf("Getting listing with marker %q\n", aws.StringValue(marker))

		result, errListObjects := svc.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(settings.Bucket),
			Marker: marker,
		})
		if errListObjects != error(nil) {
			exitErrorf("Unable to list items in bucket, %v", errListObjects)
		}

		for _, entry := range result.Contents {
			if aws.StringValue(entry.StorageClass) != "INTELLIGENT_TIERING" {
				data.addPath(aws.StringValue(entry.Key))
			}
		}

		if !aws.BoolValue(result.IsTruncated) {
			break
		}

		if m := result.NextMarker; aws.StringValue(m) != "" {
			marker = m
		} else {
			marker = result.Contents[len(result.Contents)-1].Key
		}
	}
}

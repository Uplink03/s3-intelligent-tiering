package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, errSession := session.NewSession(&aws.Config{
		Region: aws.String(settings.Region)},
	)
	if errSession != error(nil) {
		exitErrorf("Could not create new session, %v", errSession)
	}

	// Create S3 service client
	svc := s3.New(sess)

	handler := handlers[settings.RunMode]
	if handler != nil {
		handler(svc)
	}
}

func exitErrorf(msg string, args ...interface{}) {
	panic(fmt.Sprintf(msg+"\n", args...))
}

type runModeHandler func(svc *s3.S3)

var handlers = map[RunMode]runModeHandler{
	RunModeInvalid: func(*s3.S3) {
		exitErrorf("Invalid run mode")
	},
	RunModeListObjects:   listObjects,
	RunModeChangeObjects: changeObjects,
}

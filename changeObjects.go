package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func changeObjects(svc *s3.S3) {
	ch := make(chan PathRecord)
	wg := &sync.WaitGroup{}

	for i := 0; i < settings.WorkerCount; i++ {
		go processFiles(svc, ch, wg)
	}

	data.getUnprocessedPaths(ch)

	wg.Wait()
}

func processFiles(svc *s3.S3, ch <-chan PathRecord, wg *sync.WaitGroup) {
	wg.Add(1)

	for record := range ch {
		op := &s3.CopyObjectInput{
			CopySource:   aws.String(url.QueryEscape(settings.Bucket + "/" + record.Path)),
			Bucket:       aws.String(settings.Bucket),
			Key:          aws.String(record.Path),
			StorageClass: aws.String(s3.StorageClassIntelligentTiering),
		}
		_, errCopyObject := svc.CopyObject(op)
		if errCopyObject != error(nil) {
			fmt.Fprintf(os.Stderr, "Unable to copy item %q, %v\n", aws.StringValue(op.Key), errCopyObject)
		}
		data.setPathAsProcessed(record.Id, errCopyObject)
	}

	wg.Done()
}

package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"strings"
	"sync"
)

const (
	s3Endpoint = "http://localhost:4566"
	bucketName = "test-2"
	keyPrefix  = "test-2/"
)

func createBucket(c *s3.S3, bucketName string) {
	out, err := c.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
		ACL:    aws.String(s3.BucketCannedACLPublicRead),
	})
	if err != nil {
		log.Printf("failed to create bucket: %v", err)
		return
	}

	log.Printf("bucket created: %v", out)
}

func listBuckets(c *s3.S3) {
	buckets, err := c.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Printf("failed to list buckets: %v", err)
		return
	}

	log.Printf("found %d buckets\n", len(buckets.Buckets))
	for _, b := range buckets.Buckets {
		log.Printf(
			"bucket: `%s` with creation date `%s`\n",
			aws.StringValue(b.Name),
			aws.TimeValue(b.CreationDate),
		)
	}
}

func viewBucket(c *s3.S3, bucketName string) {
	objects := make([]*s3.Object, 0)
	err := c.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(bucketName),
	}, func(p *s3.ListObjectsOutput, lastPage bool) bool {
		for _, o := range p.Contents {
			objects = append(objects, o)
		}

		return true
	})

	if err != nil {
		log.Printf("failed to list objects: %v", err)
		return
	}

	log.Printf("found %d objects in bucket `%s`\n", len(objects), bucketName)
	for _, o := range objects {
		log.Printf(
			"object: `%s` with size `%d`\n",
			aws.StringValue(o.Key),
			aws.Int64Value(o.Size),
		)
	}
}

func putObject(c *s3.S3, bucketName, keyName string, body io.ReadSeeker) {
	out, err := c.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
		Body:   body,
	})
	if err != nil {
		log.Printf("failed to put object: %v", err)
		return
	}

	log.Printf("object created: %v", out)
}

func main() {
	s := session.Must(session.NewSession(&aws.Config{
		Endpoint:         aws.String(s3Endpoint),
		Region:           aws.String(endpoints.EuCentral1RegionID),
		S3ForcePathStyle: aws.Bool(true),
	}))

	c := s3.New(s)

	createBucket(c, bucketName)
	listBuckets(c)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("%s%d", keyPrefix, i)
			content := strings.NewReader(fmt.Sprintf("Object content %d", i))
			putObject(c, bucketName, key, content)
		}(i)
	}
	wg.Wait()

	viewBucket(c, bucketName)
}

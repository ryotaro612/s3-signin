package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go-v2/config"
)

var profile string
var bucket string
var key string

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s generates a signed URL of a s3 object.\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVar(&profile, "profile", "default", "AWS profile to use")
	flag.StringVar(&bucket, "bucket", "", "A S3 bucket name")
	flag.StringVar(&key, "key", "", "A S3 object key")
}

// Presigner encapsulates the Amazon Simple Storage Service (Amazon S3) presign actions
// used in the examples.
// It contains PresignClient, a client that is used to presign requests to Amazon S3.
// Presigned requests contain temporary credentials and can be made from any HTTP client.
type presigner struct {
	client *s3.PresignClient
}

// GetObject makes a presigned request that can be used to get an object from a bucket.
// The presigned request is valid for the specified number of seconds.
func (presigner presigner) GetObject(
	ctx context.Context, bucketName string, objectKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.client.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		log.Printf("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}

func makePresinder(ctx context.Context, awsProfile string) (presigner, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile))
	if err != nil {
		return presigner{}, err
	}
	s3Client := s3.NewFromConfig(cfg)
	return presigner{s3.NewPresignClient(s3Client)}, nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	var err error
	flag.Parse()
	if profile == "" {
		err = errors.New("-profile <aws profile> is required")
		return
	}
	if bucket == "" {
		err = errors.New("-bucket <bucket> is required")
		return
	}
	if key == "" {
		err = errors.New("-key <key> is required")
		return
	}

	defer func() {
		if err != nil {
			logger.Error("An error occurred", "error", err)
			os.Exit(1)
		}
	}()
	ctx := context.Background()
	logger.Info("profile", "profile", profile)
	presigner, err := makePresinder(ctx, profile)
	if err != nil {
		return
	}
	req, err := presigner.GetObject(ctx, bucket, key, 60*60*24)
	if err != nil {
		return
	}
	logger.Info("generate a presigned url", "url", req.URL)

}

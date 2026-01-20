package r2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func Session(cfg Config) (*s3.S3, *s3manager.Uploader, error) {
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("auto"),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
		Endpoint:    aws.String(endpoint),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("session failed: %w", err)
	}

	client := s3.New(sess)
	uploader := s3manager.NewUploaderWithClient(client)

	return client, uploader, nil
}

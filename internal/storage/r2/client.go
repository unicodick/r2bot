package r2

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/unicodick/r2bot/internal/storage"
)

type Client struct {
	uploader  *manager.Uploader
	bucket    string
	publicURL string
	keyGen    *KeyGenerator
}

func NewClient(cfg Config) (*Client, error) {
	_, uploader, err := Session(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		uploader:  uploader,
		bucket:    cfg.Bucket,
		publicURL: cfg.PublicURL,
		keyGen:    NewKeyGenerator(),
	}, nil
}

func (c *Client) UploadFile(filename string, reader io.Reader, size int64) (*storage.FileInfo, error) {
	key := c.keyGen.Generate(filename)

	input := &s3.PutObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
		Body:   reader,
	}

	_, err := c.uploader.Upload(context.Background(), input)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	return &storage.FileInfo{
		Name: filename,
		Size: size,
		URL:  fmt.Sprintf("%s/%s", c.publicURL, key),
	}, nil
}

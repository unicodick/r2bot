package storage

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    90 * time.Second,
		DisableCompression: true,
	},
}

type R2Client struct {
	client    *s3.S3
	uploader  *s3manager.Uploader
	bucket    string
	publicURL string
}

type FileInfo struct {
	Name string
	Size int64
	URL  string
}

func NewR2Client(accountID, accessKey, secretKey, bucket, publicURL string) (*R2Client, error) {
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("auto"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:    aws.String(endpoint),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	client := s3.New(sess)
	uploader := s3manager.NewUploaderWithClient(client)

	return &R2Client{
		client:    client,
		uploader:  uploader,
		bucket:    bucket,
		publicURL: strings.TrimSuffix(publicURL, "/"),
	}, nil
}

func (r *R2Client) UploadFile(filename string, reader io.Reader, size int64) (*FileInfo, error) {
	key := r.generateKey(filename)

	input := &s3manager.UploadInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
		Body:   reader,
	}

	_, err := r.uploader.Upload(input)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	fileInfo := &FileInfo{
		Name: filename,
		Size: size,
		URL:  fmt.Sprintf("%s/%s", r.publicURL, key),
	}

	return fileInfo, nil
}

func (r *R2Client) DownloadFile(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file data: %w", err)
	}

	return data, nil
}

func (r *R2Client) generateKey(filename string) string {
	timestamp := time.Now().Unix()
	ext := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, ext)

	// sanitize filename
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")

	return fmt.Sprintf("%d_%s%s", timestamp, name, ext)
}

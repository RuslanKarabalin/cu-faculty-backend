package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	Endpoint string
	// PublicEndpoint is the internet-reachable base URL used when signing
	// download URLs handed to clients. Falls back to Endpoint when empty.
	PublicEndpoint string
	Region         string
	AccessKey      string
	SecretKey      string
	Bucket         string
	UsePathStyle   bool
	PresignTTL     time.Duration
}

type Client struct {
	s3         *s3.Client
	presign    *s3.PresignClient
	bucket     string
	presignTTL time.Duration
}

func New(ctx context.Context, cfg Config) (*Client, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	newS3 := func(endpoint string) *s3.Client {
		return s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			o.UsePathStyle = cfg.UsePathStyle
		})
	}

	// Uploads go server-to-store over the internal endpoint; download URLs are
	// signed against the public endpoint so browsers can reach (and validate) them.
	publicEndpoint := cfg.PublicEndpoint
	if publicEndpoint == "" {
		publicEndpoint = cfg.Endpoint
	}

	return &Client{
		s3:         newS3(cfg.Endpoint),
		presign:    s3.NewPresignClient(newS3(publicEndpoint)),
		bucket:     cfg.Bucket,
		presignTTL: cfg.PresignTTL,
	}, nil
}

func (c *Client) Upload(ctx context.Context, key, contentType string, body io.Reader, size int64) error {
	input := &s3.PutObjectInput{
		Bucket:        aws.String(c.bucket),
		Key:           aws.String(key),
		Body:          body,
		ContentLength: aws.Int64(size),
	}
	if contentType != "" {
		input.ContentType = aws.String(contentType)
	}
	if _, err := c.s3.PutObject(ctx, input); err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

func (c *Client) PresignDownload(ctx context.Context, key string) (string, error) {
	req, err := c.presign.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(c.presignTTL))
	if err != nil {
		return "", fmt.Errorf("presign download: %w", err)
	}
	return req.URL, nil
}

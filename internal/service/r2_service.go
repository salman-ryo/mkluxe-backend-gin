package service

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"mkluxe-backend/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Service struct {
	s3Client      *s3.Client
	presignClient *s3.PresignClient
	bucketName    string
	publicBaseURL string
}

func NewR2Service(cfg *config.Config) (*R2Service, error) {
	r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2AccountID)

	sdkConfig, err := s3config.LoadDefaultConfig(
		context.Background(),
		s3config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.R2AccessKey,
				cfg.R2SecretKey,
				"",
			),
		),
		s3config.WithRegion("auto"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load SDK config: %w", err)
	}

	s3Client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(r2Endpoint)
	})

	return &R2Service{
		s3Client:      s3Client,
		presignClient: s3.NewPresignClient(s3Client),
		bucketName:    cfg.R2BucketName,
		publicBaseURL: strings.TrimRight(cfg.R2PublicBaseURL, "/"),
	}, nil
}

// GetPresignedUploadURL generates a presigned upload URL.
// If objectKey is empty, a unique key is generated.
// If objectKey is provided, it will be used directly.
func (s *R2Service) GetPresignedUploadURL(
	ctx context.Context,
	fileName string,
	contentType string,
	objectKey string,
	expires time.Duration,
) (uploadURL, publicURL, finalKey string, err error) {

	// Generate unique key if none supplied
	if objectKey == "" {
		ext := filepath.Ext(fileName)
		base := strings.TrimSuffix(fileName, ext)

		objectKey = fmt.Sprintf(
			"uploads/%d_%s%s",
			time.Now().UnixNano(),
			base,
			ext,
		)
	}

	presignedReq, err := s.presignClient.PresignPutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:      aws.String(s.bucketName),
			Key:         aws.String(objectKey),
			ContentType: aws.String(contentType),
		},
		func(opts *s3.PresignOptions) {
			opts.Expires = expires
		},
	)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate presigned PUT URL: %w", err)
	}

	publicURL = fmt.Sprintf("%s/%s", s.publicBaseURL, objectKey)

	return presignedReq.URL, publicURL, objectKey, nil
}

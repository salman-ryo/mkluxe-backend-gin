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
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
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

// ExtractObjectKey parses a media URL and extracts the R2 object key if it belongs to this R2 storage.
// Returns an empty string if the URL is external or does not match our R2 uploads pattern.
func (s *R2Service) ExtractObjectKey(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}

	var candidateKey string

	// 1. Direct match with configured publicBaseURL
	if s.publicBaseURL != "" && strings.HasPrefix(rawURL, s.publicBaseURL+"/") {
		candidateKey = strings.TrimPrefix(rawURL, s.publicBaseURL+"/")
	} else if strings.Contains(rawURL, "r2.cloudflarestorage.com") || strings.Contains(rawURL, "r2.dev") {
		// 2. Generic R2 endpoint format, e.g. https://<subdomain>.r2.dev/uploads/...
		idx := strings.Index(rawURL, "/uploads/")
		if idx != -1 {
			candidateKey = rawURL[idx+1:]
		}
	} else if strings.HasPrefix(rawURL, "/uploads/") {
		// 3. Relative path starting with /uploads/
		candidateKey = strings.TrimPrefix(rawURL, "/")
	} else if strings.HasPrefix(rawURL, "uploads/") {
		candidateKey = rawURL
	}

	// Remove any query parameters or fragments if present
	if qIdx := strings.Index(candidateKey, "?"); qIdx != -1 {
		candidateKey = candidateKey[:qIdx]
	}
	if fIdx := strings.Index(candidateKey, "#"); fIdx != -1 {
		candidateKey = candidateKey[:fIdx]
	}

	// Safety check: candidate key MUST start with "uploads/" and not contain path traversal ".."
	if strings.HasPrefix(candidateKey, "uploads/") && !strings.Contains(candidateKey, "..") && len(candidateKey) > len("uploads/") {
		return candidateKey
	}

	return ""
}

// DeleteObjects deletes multiple objects by their keys from the R2 bucket in batches of up to 1000.
func (s *R2Service) DeleteObjects(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	// Deduplicate keys
	uniqueKeysMap := make(map[string]struct{}, len(keys))
	var uniqueKeys []types.ObjectIdentifier
	for _, k := range keys {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		if _, exists := uniqueKeysMap[k]; !exists {
			uniqueKeysMap[k] = struct{}{}
			uniqueKeys = append(uniqueKeys, types.ObjectIdentifier{
				Key: aws.String(k),
			})
		}
	}

	if len(uniqueKeys) == 0 {
		return nil
	}

	const batchSize = 1000
	for i := 0; i < len(uniqueKeys); i += batchSize {
		end := i + batchSize
		if end > len(uniqueKeys) {
			end = len(uniqueKeys)
		}

		batch := uniqueKeys[i:end]
		input := &s3.DeleteObjectsInput{
			Bucket: aws.String(s.bucketName),
			Delete: &types.Delete{
				Objects: batch,
				Quiet:   aws.Bool(true),
			},
		}

		output, err := s.s3Client.DeleteObjects(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to delete objects batch from R2: %w", err)
		}
		if len(output.Errors) > 0 {
			var errMsgs []string
			for _, objErr := range output.Errors {
				keyStr := ""
				if objErr.Key != nil {
					keyStr = *objErr.Key
				}
				codeStr := ""
				if objErr.Code != nil {
					codeStr = *objErr.Code
				}
				errMsgs = append(errMsgs, fmt.Sprintf("%s: %s", keyStr, codeStr))
			}
			return fmt.Errorf("R2 deletion errors: %s", strings.Join(errMsgs, "; "))
		}
	}

	return nil
}

// DeleteMediaByURLs extracts valid R2 object keys from URLs and deletes them from R2.
func (s *R2Service) DeleteMediaByURLs(ctx context.Context, urls []string) error {
	if len(urls) == 0 {
		return nil
	}

	var keys []string
	for _, u := range urls {
		key := s.ExtractObjectKey(u)
		if key != "" {
			keys = append(keys, key)
		}
	}

	if len(keys) == 0 {
		return nil
	}

	return s.DeleteObjects(ctx, keys)
}


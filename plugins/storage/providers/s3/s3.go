package s3

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-zepto/zepto/plugins/storage"
)

var ErrAwsAcessKeyIdRequired = errors.New("[s3 storage] access key id required")
var ErrAwsSecretAccessKeyRequired = errors.New("[s3 storage] secret access key required")
var ErrAwsRegionRequired = errors.New("[s3 storage] region required")
var ErrAwsBucketRequired = errors.New("[s3 storage] bucket required")

// Make sure S3Storage implements storage.Storage
var _ storage.Storage = &S3Storage{}

type S3Storage struct {
	bucket   string
	uploader *s3manager.Uploader
}

type Options struct {
	Bucket          string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

func getEnv(key string, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func newAwsSession(opts Options) *session.Session {
	accessKeyID := getEnv("AWS_ACCESS_KEY_ID", opts.AccessKeyID)
	if accessKeyID == "" {
		panic(ErrAwsAcessKeyIdRequired)
	}
	secretAccessKey := getEnv("AWS_SECRET_ACCESS_KEY", opts.SecretAccessKey)
	if secretAccessKey == "" {
		panic(ErrAwsSecretAccessKeyRequired)
	}
	region := getEnv("AWS_REGION", opts.Region)
	if region == "" {
		panic(ErrAwsAcessKeyIdRequired)
	}
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
			Credentials: credentials.NewStaticCredentials(
				accessKeyID,
				secretAccessKey,
				"",
			),
		})
	if err != nil {
		panic(err)
	}
	return sess
}

func NewS3Storage(opts Options) *S3Storage {
	bucket := getEnv("AWS_BUCKET", opts.Bucket)
	if bucket == "" {
		panic(ErrAwsBucketRequired)
	}
	sess := newAwsSession(opts)
	uploader := s3manager.NewUploader(sess)
	return &S3Storage{
		bucket:   bucket,
		uploader: uploader,
	}
}

func s3AclFromAccessType(at storage.AccessType) *string {
	if at == storage.Public {
		return aws.String("public-read")
	}
	return aws.String("private")
}

func (s *S3Storage) Create(ctx context.Context, opts storage.CreateOptions) (*storage.FileInfo, error) {
	res, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(s.bucket),
		ACL:         s3AclFromAccessType(opts.AccessType),
		Key:         aws.String(opts.Key),
		ContentType: &opts.ContentType,
		Body:        opts.Body,
	})
	if err != nil {
		return nil, err
	}
	return &storage.FileInfo{
		Key:         opts.Key,
		Url:         res.Location,
		AccessType:  opts.AccessType,
		ContentType: opts.ContentType,
		Size:        opts.Size,
	}, nil
}

func (s *S3Storage) Destroy(ctx context.Context, key string) error {
	_, err := s.uploader.S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	return err
}

func (s *S3Storage) GenerateSignedURL(ctx context.Context, key string, expirationTime time.Duration) (string, error) {
	req, _ := s.uploader.S3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &key,
	})
	exp := 15 * time.Minute
	if expirationTime > 0 {
		exp = expirationTime
	}
	urlStr, err := req.Presign(exp)
	return urlStr, err
}

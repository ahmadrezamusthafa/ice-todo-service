package storage

import (
	"github.com/ahmadrezamusthafa/ice-todo-service/config"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Connector struct {
	config     *config.Config
	logger     logger.Logger
	session    *session.Session
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func NewS3Connector(config *config.Config, logger logger.Logger) *S3Connector {
	return &S3Connector{
		config: config,
		logger: logger,
	}
}

func (c *S3Connector) Connect() (*s3manager.Uploader, *s3manager.Downloader, error) {
	if c.uploader != nil && c.downloader != nil {
		return c.uploader, c.downloader, nil
	}

	awsConfig := &aws.Config{
		Region:           aws.String(c.config.S3.Region),
		Endpoint:         aws.String(c.config.S3.Endpoint),
		Credentials:      credentials.NewStaticCredentials(c.config.S3.AccessKey, c.config.S3.SecretKey, ""),
		DisableSSL:       aws.Bool(!c.config.S3.UseSSL),
		S3ForcePathStyle: aws.Bool(true),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		c.logger.Error("Failed to create AWS session: %v", err)
		return nil, nil, err
	}

	c.session = sess
	c.uploader = s3manager.NewUploader(sess)
	c.downloader = s3manager.NewDownloader(sess)

	s3Client := s3.New(sess)

	_, err = s3Client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(c.config.S3.Bucket),
	})

	if err != nil {

		c.logger.Info("Creating bucket: %s", c.config.S3.Bucket)

		createInput := &s3.CreateBucketInput{
			Bucket: aws.String(c.config.S3.Bucket),
		}

		if c.config.S3.Region != "us-east-1" {
			createInput.CreateBucketConfiguration = &s3.CreateBucketConfiguration{
				LocationConstraint: aws.String(c.config.S3.Region),
			}
		}

		_, err = s3Client.CreateBucket(createInput)
		if err != nil {

			if aErr, ok := err.(awserr.Error); ok {
				switch aErr.Code() {
				case s3.ErrCodeBucketAlreadyExists, s3.ErrCodeBucketAlreadyOwnedByYou:
					c.logger.Info("Bucket already exists: %s", c.config.S3.Bucket)
				default:
					c.logger.Error("Failed to create bucket: %v", err)
					return nil, nil, err
				}
			} else {
				c.logger.Error("Failed to create bucket: %v", err)
				return nil, nil, err
			}
		} else {
			c.logger.Info("Successfully created bucket: %s", c.config.S3.Bucket)
		}
	} else {
		c.logger.Info("Bucket already exists: %s", c.config.S3.Bucket)
	}

	c.logger.Info("Connected to S3")

	return c.uploader, c.downloader, nil
}

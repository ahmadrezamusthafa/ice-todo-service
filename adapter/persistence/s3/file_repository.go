package s3

import (
	"bytes"
	"fmt"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/apperrors"
	"io"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

type UploaderInterface interface {
	Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error)
}

type DownloaderInterface interface {
	Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error)
}

type FileRepository struct {
	uploader   UploaderInterface
	downloader DownloaderInterface
	bucketName string
}

func NewFileRepository(uploader UploaderInterface, downloader DownloaderInterface, bucketName string) *FileRepository {
	return &FileRepository{
		uploader:   uploader,
		downloader: downloader,
		bucketName: bucketName,
	}
}

func getContentType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}

func (r *FileRepository) Upload(fileName string, fileSize int64, fileContent io.Reader) (string, error) {
	fileID := uuid.New().String()

	content, err := io.ReadAll(fileContent)
	if err != nil {
		return "", apperrors.NewStorageError("failed to read file content", err)
	}

	_, err = r.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(fileID),
		Body:        bytes.NewReader(content),
		ContentType: aws.String(getContentType(fileName)),
		Metadata: map[string]*string{
			"original-filename": aws.String(fileName),
		},
	})
	if err != nil {
		if aErr, ok := err.(awserr.Error); ok {
			return "", apperrors.NewStorageError(
				fmt.Sprintf("failed to upload file to S3: %s - %s", aErr.Code(), aErr.Message()),
				aErr,
			)
		}
		return "", apperrors.NewStorageError("failed to upload file to S3", err)
	}

	return fileID, nil
}

func (r *FileRepository) Get(fileID string) ([]byte, error) {
	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := r.downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(fileID),
	})
	if err != nil {
		if aErr, ok := err.(awserr.Error); ok {
			switch aErr.Code() {
			case s3.ErrCodeNoSuchKey:
				return nil, apperrors.NewNotFoundError(fmt.Sprintf("file with id %s not found", fileID), aErr)
			case s3.ErrCodeNoSuchBucket:
				return nil, apperrors.NewStorageError(fmt.Sprintf("bucket %s not found", r.bucketName), aErr)
			default:
				return nil, apperrors.NewStorageError(
					fmt.Sprintf("failed to download file from S3: %s - %s", aErr.Code(), aErr.Message()),
					aErr,
				)
			}
		}
		return nil, apperrors.NewStorageError("failed to download file from S3", err)
	}

	return buf.Bytes(), nil
}

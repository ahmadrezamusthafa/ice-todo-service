package s3_test

import (
	"bytes"
	"fmt"
	"testing"

	s3adapter "github.com/ahmadrezamusthafa/ice-todo-service/adapter/persistence/s3"
	"github.com/ahmadrezamusthafa/ice-todo-service/config"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/logger"
	"github.com/ahmadrezamusthafa/ice-todo-service/infrastructure/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

func setupTestS3() (*storage.S3Connector, error) {
	cfg := config.NewConfig()
	log := logger.NewStandardLogger()

	s3Connector := storage.NewS3Connector(cfg, log)
	_, _, err := s3Connector.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to S3: %w", err)
	}

	return s3Connector, nil
}

func cleanupTestData(s3Connector *storage.S3Connector, fileIDs []string) {
	cfg := config.NewConfig()
	sess := s3Connector.GetSession()
	s3Client := s3.New(sess)

	for _, fileID := range fileIDs {
		_, _ = s3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(cfg.S3.Bucket),
			Key:    aws.String(fileID),
		})
	}
}

func BenchmarkFileRepositoryUpload(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	s3Connector, err := setupTestS3()
	if err != nil {
		b.Fatalf("Failed to setup test S3: %v", err)
	}

	cfg := config.NewConfig()
	uploader, downloader, err := s3Connector.Connect()
	if err != nil {
		b.Fatalf("Failed to get S3 uploader and downloader: %v", err)
	}

	repo := s3adapter.NewFileRepository(uploader, downloader, cfg.S3.Bucket)

	var uploadedFileIDs []string
	defer func() {
		cleanupTestData(s3Connector, uploadedFileIDs)
	}()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		fileName := fmt.Sprintf("benchmark-test-%d.txt", i)
		fileContent := fmt.Sprintf("This is benchmark test content for iteration %d", i)
		fileSize := int64(len(fileContent))
		fileReader := bytes.NewReader([]byte(fileContent))

		b.StartTimer()

		fileID, err := repo.Upload(fileName, fileSize, fileReader)
		if err != nil {
			b.Fatalf("Failed to upload file: %v", err)
		}

		b.StopTimer()
		uploadedFileIDs = append(uploadedFileIDs, fileID)
	}
}

func BenchmarkFileRepositoryUploadParallel(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	s3Connector, err := setupTestS3()
	if err != nil {
		b.Fatalf("Failed to setup test S3: %v", err)
	}

	cfg := config.NewConfig()
	uploader, downloader, err := s3Connector.Connect()
	if err != nil {
		b.Fatalf("Failed to get S3 uploader and downloader: %v", err)
	}

	repo := s3adapter.NewFileRepository(uploader, downloader, cfg.S3.Bucket)

	idChan := make(chan string, b.N)
	var uploadedFileIDs []string

	defer func() {
		close(idChan)

		for id := range idChan {
			uploadedFileIDs = append(uploadedFileIDs, id)
		}

		cleanupTestData(s3Connector, uploadedFileIDs)
	}()

	go func() {
		for id := range idChan {
			uploadedFileIDs = append(uploadedFileIDs, id)
		}
	}()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			uniqID := uuid.New().String()
			fileName := fmt.Sprintf("benchmark-test-parallel-%s.txt", uniqID)
			fileContent := fmt.Sprintf("This is parallel benchmark test content for %s", uniqID)
			fileSize := int64(len(fileContent))
			fileReader := bytes.NewReader([]byte(fileContent))

			fileID, err := repo.Upload(fileName, fileSize, fileReader)
			if err != nil {
				b.Fatalf("Failed to upload file in parallel benchmark: %v", err)
			}

			idChan <- fileID
		}
	})
}

package s3_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/ahmadrezamusthafa/ice-todo-service/adapter/persistence/s3"
	"github.com/ahmadrezamusthafa/ice-todo-service/domain/apperrors"
	mock_s3_adapter "github.com/ahmadrezamusthafa/ice-todo-service/mock/adapter/persistence/s3"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUploader := mock_s3_adapter.NewMockUploaderInterface(ctrl)
	mockDownloader := mock_s3_adapter.NewMockDownloaderInterface(ctrl)
	bucketName := "test-bucket"

	repo := s3.NewFileRepository(mockUploader, mockDownloader, bucketName)

	testCases := []struct {
		name           string
		fileName       string
		fileSize       int64
		fileContent    io.Reader
		setupMock      func()
		expectedFileID string
		expectedError  error
	}{
		{
			name:        "Success",
			fileName:    "test.jpg",
			fileSize:    100,
			fileContent: bytes.NewReader([]byte("test content")),
			setupMock: func() {
				mockUploader.EXPECT().
					Upload(gomock.Any(), gomock.Any()).
					Return(&s3manager.UploadOutput{}, nil)
			},
			expectedFileID: "",
			expectedError:  nil,
		},
		{
			name:        "Error reading file content",
			fileName:    "test.txt",
			fileSize:    100,
			fileContent: &errorReader{},
			setupMock: func() {

			},
			expectedFileID: "",
			expectedError:  apperrors.NewStorageError("failed to read file content", errors.New("read error")),
		},
		{
			name:        "S3 Upload Error",
			fileName:    "test.png",
			fileSize:    100,
			fileContent: bytes.NewReader([]byte("test content")),
			setupMock: func() {
				s3Err := awserr.New("InternalError", "Internal S3 Error", nil)
				mockUploader.EXPECT().
					Upload(gomock.Any(), gomock.Any()).
					Return(nil, s3Err)
			},
			expectedFileID: "",
			expectedError:  apperrors.NewStorageError("failed to upload file to S3: InternalError - Internal S3 Error", nil),
		},
		{
			name:        "Generic Upload Error",
			fileName:    "test.pdf",
			fileSize:    100,
			fileContent: bytes.NewReader([]byte("test content")),
			setupMock: func() {
				mockUploader.EXPECT().
					Upload(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("network error"))
			},
			expectedFileID: "",
			expectedError:  apperrors.NewStorageError("failed to upload file to S3", errors.New("network error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			fileID, err := repo.Upload(tc.fileName, tc.fileSize, tc.fileContent)

			if tc.expectedError != nil {
				assert.Error(t, err)

				expectedAppErr, expectedIsAppErr := tc.expectedError.(*apperrors.AppError)
				actualAppErr, actualIsAppErr := err.(*apperrors.AppError)

				if expectedIsAppErr && actualIsAppErr {

					assert.Equal(t, expectedAppErr.Type, actualAppErr.Type)
					assert.Equal(t, expectedAppErr.Message, actualAppErr.Message)
				} else {

					assert.Contains(t, err.Error(), tc.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, fileID)
				_, err := uuid.Parse(fileID)
				assert.NoError(t, err, "FileID should be a valid UUID")
			}
			assert.Equal(t, tc.expectedFileID == "", fileID == "" || fileID != "")
		})
	}
}

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUploader := mock_s3_adapter.NewMockUploaderInterface(ctrl)
	mockDownloader := mock_s3_adapter.NewMockDownloaderInterface(ctrl)
	bucketName := "test-bucket"

	repo := s3.NewFileRepository(mockUploader, mockDownloader, bucketName)

	testCases := []struct {
		name          string
		fileID        string
		setupMock     func()
		expectedData  []byte
		expectedError error
	}{
		{
			name:   "Success",
			fileID: "file123",
			setupMock: func() {
				mockDownloader.EXPECT().
					Download(gomock.Any(), gomock.Any(), gomock.Any()).
					Do(func(w io.WriterAt, input *awss3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error) {

						assert.Equal(t, bucketName, *input.Bucket)
						assert.Equal(t, "file123", *input.Key)

						data := []byte("test file content")
						n, err := w.WriteAt(data, 0)
						assert.NoError(t, err)
						return int64(n), nil
					})
			},
			expectedData:  []byte("test file content"),
			expectedError: nil,
		},
		{
			name:   "File Not Found Error",
			fileID: "nonexistent",
			setupMock: func() {
				s3Err := awserr.New(awss3.ErrCodeNoSuchKey, "The specified key does not exist.", nil)
				mockDownloader.EXPECT().
					Download(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int64(0), s3Err)
			},
			expectedData:  nil,
			expectedError: apperrors.NewNotFoundError(fmt.Sprintf("file with id %s not found", "nonexistent"), nil),
		},
		{
			name:   "Bucket Not Found Error",
			fileID: "file123",
			setupMock: func() {
				s3Err := awserr.New(awss3.ErrCodeNoSuchBucket, "The specified bucket does not exist.", nil)
				mockDownloader.EXPECT().
					Download(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int64(0), s3Err)
			},
			expectedData:  nil,
			expectedError: apperrors.NewStorageError(fmt.Sprintf("bucket %s not found", bucketName), nil),
		},
		{
			name:   "Other AWS Error",
			fileID: "file123",
			setupMock: func() {
				s3Err := awserr.New("InternalError", "Internal S3 Error", nil)
				mockDownloader.EXPECT().
					Download(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int64(0), s3Err)
			},
			expectedData:  nil,
			expectedError: apperrors.NewStorageError("failed to download file from S3: InternalError - Internal S3 Error", nil),
		},
		{
			name:   "Generic Error",
			fileID: "file123",
			setupMock: func() {
				mockDownloader.EXPECT().
					Download(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(int64(0), errors.New("network error"))
			},
			expectedData:  nil,
			expectedError: apperrors.NewStorageError("failed to download file from S3", errors.New("network error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			data, err := repo.Get(tc.fileID)

			if tc.expectedError != nil {
				assert.Error(t, err)

				expectedAppErr, expectedIsAppErr := tc.expectedError.(*apperrors.AppError)
				actualAppErr, actualIsAppErr := err.(*apperrors.AppError)

				if expectedIsAppErr && actualIsAppErr {

					assert.Equal(t, expectedAppErr.Type, actualAppErr.Type)
					assert.Equal(t, expectedAppErr.Message, actualAppErr.Message)
				} else {

					assert.Contains(t, err.Error(), tc.expectedError.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedData, data)
			}
		})
	}
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

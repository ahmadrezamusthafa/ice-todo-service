package repository

import (
	"io"
)

type FileRepository interface {
	Upload(fileName string, fileSize int64, fileContent io.Reader) (string, error)
	Get(fileID string) ([]byte, error)
}

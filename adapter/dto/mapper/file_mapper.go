package mapper

import (
	"fmt"
	"mime/multipart"
)

type FileUploadResponse struct {
	FileID   string `json:"fileId"`
	FileName string `json:"fileName"`
	Size     string `json:"size"`
}

type FileDownloadResponse struct {
	FileID string `json:"fileId"`
	Size   int    `json:"size"`
}

func FileToUploadResponse(fileID string, file *multipart.FileHeader) *FileUploadResponse {
	return &FileUploadResponse{
		FileID:   fileID,
		FileName: file.Filename,
		Size:     fmt.Sprintf("%d", file.Size),
	}
}

package validator

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"
)

type TodoValidator struct {
	MaxDescriptionLength int
}

func NewTodoValidator() *TodoValidator {
	return &TodoValidator{
		MaxDescriptionLength: 1000,
	}
}

func (v *TodoValidator) ValidateDescription(description string) error {
	if description == "" {
		return fmt.Errorf("description is required")
	}

	if len(description) > v.MaxDescriptionLength {
		return fmt.Errorf("description is too long (maximum %d characters)", v.MaxDescriptionLength)
	}

	return nil
}

func (v *TodoValidator) ValidateDueDate(dueDate time.Time) error {
	if dueDate.IsZero() {
		return fmt.Errorf("due date is required")
	}

	if dueDate.Before(time.Now()) {
		return fmt.Errorf("due date must be in the future")
	}

	return nil
}

type FileValidator struct {
	MaxFileSize  int64
	AllowedTypes map[string]bool
}

func NewFileValidator() *FileValidator {
	return &FileValidator{
		MaxFileSize: 10 * 1024 * 1024,
		AllowedTypes: map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".pdf":  true,
			".txt":  true,
		},
	}
}

func (v *FileValidator) ValidateFile(header *multipart.FileHeader) error {

	if err := v.ValidateFileSize(header.Size); err != nil {
		return err
	}

	ext := filepath.Ext(header.Filename)
	return v.ValidateFileType(ext)
}

func (v *FileValidator) ValidateFileSize(size int64) error {
	if size > v.MaxFileSize {
		return fmt.Errorf("file too large (maximum %d bytes)", v.MaxFileSize)
	}
	return nil
}

func (v *FileValidator) ValidateFileType(ext string) error {
	if !v.AllowedTypes[ext] {
		return fmt.Errorf("invalid file type: %s", ext)
	}
	return nil
}

package filestore

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrFileTooLarge    = errors.New("file size exceeds maximum limit")
	ErrInvalidFileType = errors.New("invalid file type, only images are allowed")
)

// FileService mendefinisikan service untuk upload file
type FileService interface {
	UploadFile(file *multipart.FileHeader) (string, error)
}

// LocalFileService implementasi file service dengan penyimpanan lokal
type LocalFileService struct {
	uploadDir string
	maxSize   int64
}

// NewLocalFileService membuat instance baru LocalFileService
func NewLocalFileService(uploadDir string, maxSize int64) FileService {
	// Ensure upload directory exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	return &LocalFileService{
		uploadDir: uploadDir,
		maxSize:   maxSize,
	}
}

// UploadFile mengupload file ke penyimpanan lokal
func (s *LocalFileService) UploadFile(file *multipart.FileHeader) (string, error) {
	if file.Size > s.maxSize {
		return "", ErrFileTooLarge
	}

	// Validasi jenis file (hanya gambar)
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isImageExtension(ext) {
		return "", ErrInvalidFileType
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filepath := filepath.Join(s.uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Open source file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Copy source to destination
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return filename, nil
}

// isImageExtension memeriksa apakah ekstensi file adalah ekstensi gambar yang valid
func isImageExtension(ext string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

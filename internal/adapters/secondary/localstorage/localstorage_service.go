package localstorage

import (
	domerr "go-gin-clean/internal/domain/error"
	"go-gin-clean/internal/ports/secondary"
	"os"
	"path"
	"path/filepath"
)

type LocalStorageService struct {
}

func NewLocalStorageService() secondary.MediaService {
	return &LocalStorageService{}
}

func (s *LocalStorageService) UploadFile(filename string, size int64, content any, filePath string) (*string, error) {
	basePath := "assets"
	dirPath := filepath.Join(basePath, filePath)
	fullPath := filepath.Join(dirPath, filename)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return nil, domerr.ErrCreateFileSpace
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, domerr.ErrUploadFile
	}
	defer dst.Close()

	// Type assert content to io.Reader if needed
	// The interface was updated to accept `any` for flexibility
	// but in practice it should be io.Reader

	publicURL := path.Join("/assets", filePath, filename)

	return &publicURL, nil
}

func (s *LocalStorageService) DeleteFile(fileURL string) error {
	if err := os.Remove(fileURL); err != nil {
		return domerr.ErrDeleteFile
	}

	return nil
}

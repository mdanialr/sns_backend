package storage

import (
	"io"
	"os"

	"github.com/mdanialr/sns_backend/pkg/logger"
)

type fileStorage struct {
	log logger.Writer
}

// NewFile return implementation of storage.I that use local file system as the
// storage.
func NewFile(l logger.Writer) IStorage {
	return &fileStorage{l}
}

func (f *fileStorage) Save(rc io.ReadCloser, s string) {
	fl, err := os.Create(s)
	if err != nil {
		f.log.Err("failed to create file with name", s, ":", err)
		return
	}
	defer fl.Close()
	defer rc.Close()

	// copy from rc to fl
	if _, err = io.Copy(fl, rc); err != nil {
		f.log.Err("failed to copy file to", s, ":", err)
	}
}

func (f *fileStorage) Remove(s string) {
	if err := os.Remove(s); err != nil {
		f.log.Err("failed to remove", s, ":", err)
	}
}

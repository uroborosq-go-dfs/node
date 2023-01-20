package fs

import "io"

type Storage interface {
	Save(path string, reader io.Reader) error
	Remove(path string) error
	GetFile(path string) (io.Reader, error)
	GetFileSize(path string) (int64, error)
	GetPathList() []string
	GetStorageSize() int64
}

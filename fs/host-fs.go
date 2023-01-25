package fs

import (
	"errors"
	"io"
	"os"
)

type HostFS struct {
	path  string
	paths []string
	size  int64
}

var _ Storage = (*HostFS)(nil)

func New(path string) Storage {
	return &HostFS{path: path, paths: make([]string, 0), size: 0}
}

func (h *HostFS) Save(path string, reader io.Reader) error {
	err := os.MkdirAll(h.path, os.ModePerm)
	file, err := os.OpenFile(h.path+string(os.PathSeparator)+path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	bufferSize := uint64(1024)
	buffer := make([]byte, bufferSize)
	fileSize := int64(0)
	for {
		read, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		} else if read == 0 {
			break
		} else if err == io.EOF {
			break
		}
		write, err := file.Write(buffer[:read])
		if err != nil {
			return err
		} else if write != read {
			return errors.New("can't write to file")
		}
		fileSize += int64(write)
	}
	h.size += fileSize
	h.paths = append(h.paths, path)
	return nil
}

func (h *HostFS) GetFile(path string) (io.Reader, error) {
	file, err := os.Open(h.path + string(os.PathSeparator) + path)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (h *HostFS) GetFileSize(path string) (int64, error) {
	info, err := os.Stat(h.path + string(os.PathSeparator) + path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

func (h *HostFS) GetPathList() []string {
	cPaths := make([]string, len(h.paths))
	copy(cPaths, h.paths)
	return cPaths
}

func (h *HostFS) GetStorageSize() int64 {
	return h.size
}

func (h *HostFS) Remove(path string) error {
	return os.RemoveAll(h.path + string(os.PathSeparator) + path)
}

package service

import (
	"io"
	"node/fs"
)

type NodeService struct {
	storage fs.Storage
}

func New(storage fs.Storage) *NodeService {
	return &NodeService{storage: storage}
}

func (n *NodeService) AddFile(partialPath string, reader io.Reader) error {
	return n.storage.Save(partialPath, reader)
}

func (n *NodeService) GetFile(partialPath string) (io.Reader, error) {
	return n.storage.GetFile(partialPath)
}

func (n *NodeService) GetPathList() []string {
	return n.storage.GetPathList()
}

func (n *NodeService) GetNodeSize() int64 {
	return n.storage.GetStorageSize()
}

func (n *NodeService) GetFileSize(partialPath string) (int64, error) {
	return n.storage.GetFileSize(partialPath)
}

func (n *NodeService) RemoveFile(partialPath string) error {
	return n.storage.Remove(partialPath)
}

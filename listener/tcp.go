package listener

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/uroborosq-go-dfs/node/service"
	"net"

	code "github.com/uroborosq-go-dfs/models/tcp-operation-code"
)

func CreateTcpListener(port string, service *service.NodeService) NodeListener {
	return &TcpListener{port: port, service: service}
}

type TcpListener struct {
	port    string
	service *service.NodeService
}

var _ NodeListener = (*TcpListener)(nil)

func (t *TcpListener) Listen() error {
	listener, err := net.Listen("tcp", "localhost:"+t.port)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		operationCode := make([]byte, 1)
		reader := bufio.NewReader(conn)
		read, err := reader.Read(operationCode)

		if err != nil {
			fmt.Println(err.Error())
		} else if read != 1 {
			fmt.Println("can't read from the stream")
		}
		switch operationCode[0] {
		case code.SendFile:
			{
				lenBytes := make([]byte, 4)
				read, err = reader.Read(lenBytes)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				pathLen := binary.BigEndian.Uint32(lenBytes)
				path := make([]byte, pathLen)
				read, err = reader.Read(path)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				lenBytes = make([]byte, 8)
				read, err = reader.Read(lenBytes)
				if err != nil {
					fmt.Println(err.Error())
					break
				}

				err = t.service.AddFile(string(path), reader)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		case code.RequestFile:
			{
				lenBytes := make([]byte, 4)
				read, err = reader.Read(lenBytes)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				pathLen := binary.BigEndian.Uint32(lenBytes)
				path := make([]byte, pathLen)
				read, err = reader.Read(path)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				fileReader, err := t.service.GetFile(string(path))
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				fileSize, err := t.service.GetFileSize(string(path))
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				bufferSize := int64(1024)
				buffer := make([]byte, bufferSize)

				fileSizeBytes := make([]byte, 8)
				binary.BigEndian.PutUint64(fileSizeBytes, uint64(fileSize))
				write, err := conn.Write(fileSizeBytes)
				if err != nil {
					fmt.Println(err.Error())
					break
				} else if write != 8 {
					fmt.Println("can't read to the stream")
					break
				}
				for i := int64(0); i < fileSize; i += bufferSize {
					read, err = fileReader.Read(buffer)
					if err != nil {
						fmt.Println(err.Error())
						break
					}
					write, err = conn.Write(buffer[:read])
					if err != nil {
						fmt.Println(err.Error())
						break
					} else if write != read {
						fmt.Println("can't read to the stream")
						break
					}
				}
			}
		case code.RequestList:
			{
				paths := t.service.GetPathList()
				size := uint64(0)
				for i := 0; i < len(paths); i++ {
					size += uint64(len(paths[i])) + 1
				}
				sizeBytes := make([]byte, 8)
				binary.BigEndian.PutUint64(sizeBytes, size)
				write, err := conn.Write(sizeBytes)
				if err != nil {
					fmt.Println(err.Error())
					break
				} else if write != 8 {
					fmt.Println("can't write to the stream")
					break
				}
				for i := 0; i < len(paths); i++ {
					_, err := conn.Write([]byte(paths[i] + "\000"))
					if err != nil {
						fmt.Println(err.Error())
						break
					}
				}
			}
		case code.RequestSize:
			{
				size := t.service.GetNodeSize()
				sizeBytes := make([]byte, 8)
				binary.BigEndian.PutUint64(sizeBytes, uint64(size))
				write, err := conn.Write(sizeBytes)
				if err != nil {
					fmt.Println(err.Error())
				} else if write != 8 {
					fmt.Println("can't write to the stream")
				}
			}
		case code.RemoveFile:
			{
				lenBytes := make([]byte, 4)
				read, err = reader.Read(lenBytes)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				pathLen := binary.BigEndian.Uint32(lenBytes)
				path := make([]byte, pathLen)
				read, err = reader.Read(path)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				err = t.service.RemoveFile(string(path))
				if err != nil {
					fmt.Println(err.Error())
				}
			}
		}
	}
}

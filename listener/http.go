package listener

import (
	"encoding/json"
	"github.com/uroborosq-go-dfs/node/service"
	"io"
	"log"
	"net/http"
)

type HttpListener struct {
	port        string
	nodeService *service.NodeService
}

func (h *HttpListener) Listen() error {
	http.HandleFunc("/file/send", h.handleSendingFile)
	http.HandleFunc("/file/request", h.handleRequestingFile)
	http.HandleFunc("/list", h.handleRequestingListFiles)
	http.HandleFunc("/size", h.handleRequestingFileSize)
	http.HandleFunc("/file/remove", h.handleRemovingFileSize)
	return http.ListenAndServe(h.port, nil)
}

func (h *HttpListener) handleSendingFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()
	err = h.nodeService.AddFile(handler.Filename, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *HttpListener) handleRequestingFile(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("query argument file is empty"))
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}
	reader, err := h.nodeService.GetFile(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	buffer := make([]byte, 1024)
	for {
		read, err := reader.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		_, err = w.Write(buffer[:read])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (h *HttpListener) handleRequestingListFiles(w http.ResponseWriter, r *http.Request) {
	paths := h.nodeService.GetPathList()
	jsonStr, err := json.Marshal(paths)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *HttpListener) handleRequestingFileSize(w http.ResponseWriter, r *http.Request) {
	size := h.nodeService.GetNodeSize()
	jsonStr, err := json.Marshal(size)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonStr)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *HttpListener) handleRemovingFileSize(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	if filePath == "" {
		_, err := w.Write([]byte("query argument file is empty"))
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	}
	err := h.nodeService.RemoveFile(filePath)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Package upload -> uploads files to the server/container
package upload

import (
	"io"
	"net/http"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

const ChunkSize = int64(10 * 1024 * 1024) // 10MB

type fileUpload struct {
	uploadID         int64
	fileName         int
	filePath         string
	size             int64
	transferredBytes int64
	contentType      string
}

type progressBar struct {
	size             int64
	transferredBytes int64
	reader           *io.Reader
}

func (p *progressBar) Read(bytes []byte) (n int, err error) {
	n, err = (*p.reader).Read(bytes)
	p.transferredBytes += int64(n)
	return n, err
}

// constructor
func newFileUpload(uploadID int64, fileName int, filePath string) (*fileUpload, error) {
	size, err := getSize(filePath)
	if err != nil {
		return nil, err
	}
	contentType, err := mimetype.DetectFile(filePath)
	if err != nil {
		return nil, err
	}
	res := new(fileUpload{
		uploadID:    uploadID,
		fileName:    fileName,
		filePath:    filePath,
		size:        size,
		contentType: contentType.String(),
	})
	return res, nil
}

func (f *fileUpload) largeFileUpload() {
	for i := int64(0); i < f.size; i += ChunkSize {
	}
}

func (f *fileUpload) uploadChunk() {}

func (f *fileUpload) smallFileUpload(url string) (*http.Response, error) {
	body, err := os.Open(f.filePath)
	if err != nil {
		return nil, err
	}
	return http.Post(url, f.contentType, body)
}

func getSize(path string) (int64, error) {
	f, err := os.Stat(path)
	if err != nil || f.IsDir() || f.Mode() == os.ModeSymlink {
		return 0, err
	}

	return f.Size(), nil
}

package lib

import (
	"fmt"
	"os"
	"strings"
)

// FileHandler handles list file.
type FileHandler struct {
	filePath string
	fp       *os.File
}

// NewFileHandler returns initialized *FileHandler
func NewFileHandler(filePath string) (*FileHandler, error) {
	info, err := os.Stat(filePath)
	if err == nil && info.IsDir() {
		return nil, fmt.Errorf("'%s' is dir, please set file path", filePath)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	return &FileHandler{
		filePath: filePath,
		fp:       f,
	}, nil
}

func (f *FileHandler) Close() error {
	return f.fp.Close()
}

// AppendRows appends rows into the file.
func (f *FileHandler) AppendRows(rows []string) error {
	fp := f.fp
	if _, err := fp.WriteString(strings.Join(rows, "\n") + "\n"); err != nil {
		return err
	}
	return fp.Sync()
}

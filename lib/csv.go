package lib

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CSVHandler handles CSV file.
type CSVHandler struct {
	header    []string
	headerMap map[string]int
	reader    *csv.Reader
	fp        *os.File
}

// NewCSVHandler returns initialized *CSVHandler
func NewCSVHandler(file string) (*CSVHandler, error) {
	info, err := os.Stat(file)
	if err == nil && info.IsDir() {
		return nil, fmt.Errorf("'%s' is dir, please set file path", file)
	}

	// load file
	fp, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	// load csv
	reader := csv.NewReader(fp)
	reader.LazyQuotes = true
	switch filepath.Ext(file) {
	case ".tsv":
		reader.Comma = '\t'
	}

	header, err := reader.Read()
	if err != nil {
		fp.Close()
		return nil, err
	}

	headerMap := make(map[string]int)
	for i, col := range header {
		headerMap[col] = i
	}

	return &CSVHandler{
		header:    header,
		headerMap: headerMap,
		reader:    reader,
		fp:        fp,
	}, nil
}

func (f *CSVHandler) Close() error {
	return f.fp.Close()
}

// Read reads a line from CSV file.
func (f *CSVHandler) Read() (map[string]string, error) {
	if f.reader == nil {
		return nil, fmt.Errorf("f.reader is nil")
	}

	line, err := f.reader.Read()
	switch {
	case err == io.EOF:
		return nil, nil
	case err != nil:
		return nil, err
	}

	header := f.header
	result := make(map[string]string)
	for i, col := range line {
		result[header[i]] = col
	}
	return result, nil
}

// CheckHeaders checks header columns.
func (f *CSVHandler) CheckHeaders(cols ...string) error {
	for _, col := range cols {
		if _, ok := f.headerMap[col]; !ok {
			return fmt.Errorf("cannot find header in csv file: [%s]", col)
		}
	}
	return nil
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func makeDir(path string) error {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return nil
	}

	return os.MkdirAll(path, os.ModePerm)
}

package types

import (
	"encoding/csv"
	"os"
	"sync"
)

type InternalWriter struct {
	writer *csv.Writer
	file   *os.File
	lock   sync.Mutex
}

func (w *InternalWriter) WriteHeader() error {
	w.writer.Write(Headers())
	w.writer.Flush()
	err := w.writer.Error()
	if err != nil {
		w.CloseFile()
		return err
	}
	return nil
}

func NewWriter() (*InternalWriter, error) {
	file, err := os.OpenFile("database.csv", os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	cw := csv.NewWriter(file)

	return &InternalWriter{
		writer: cw,
		file:   file,
	}, nil
}

func (w *InternalWriter) CloseFile() {
	w.file.Close()
}

func (w *InternalWriter) WriteMovie(movie Movie) error {
	s := movie.String()
	w.lock.Lock()
	w.writer.Write(s)
	w.writer.Flush()
	w.lock.Unlock()
	err := w.writer.Error()
	if err != nil {
		w.CloseFile()
		return err
	}

	return nil
}

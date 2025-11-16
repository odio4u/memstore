package wal

import (
	"bufio"
	"os"
	"sync"
)

type WALer struct {
	mu     sync.Mutex
	f      *os.File
	path   string
	writer *bufio.Writer
}

func OpenWAl(path string) (*WALer, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &WALer{
		f:      f,
		path:   path,
		writer: bufio.NewWriter(f),
	}, nil
}

func (w *WALer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.writer.Flush(); err != nil {
		return err
	}
	return w.f.Close()
}

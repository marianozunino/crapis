package aof

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"

	"github.com/marianozunino/crapis/internal/resp"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// Start a goroutine to sync AOF to disk every 1 second
	go aof.sync()

	return aof, nil
}

func (aof *Aof) sync() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()
		}
	}
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	return aof.file.Close()
}

func (aof *Aof) Write(value resp.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (aof *Aof) Read(f func(io.Reader) error) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	return f(aof.rd)
}

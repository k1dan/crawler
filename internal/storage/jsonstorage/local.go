package jsonstorage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/k1dan/crawler/internal/storage"
)

type JSONStorage struct {
	storeChan chan storage.File
	savePath  string
}

// New func initialize new JSONStorage Instance
func New(savePath string) JSONStorage {
	return JSONStorage{
		storeChan: make(chan storage.File, 20),
		savePath:  savePath,
	}
}

// Listen func starts asynchronously listen for new files to save
func (s JSONStorage) Listen(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go s.listen(ctx, wg)
}

// listen func monitor storeChan for new files to save
func (s JSONStorage) listen(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case file, ok := <-s.storeChan:
			if !ok {
				wg.Done()
				return
			}
			s.save(ctx, file)
		}
	}
}

// Close func closes file channel
func (s JSONStorage) Close() {
	close(s.storeChan)
}

// Save func adds new file to channel
func (s JSONStorage) Save(file storage.File) {
	s.storeChan <- file
}

// save func implements Storage interface, saves file as json
func (s JSONStorage) save(ctx context.Context, f storage.File) error {
	savePath := fmt.Sprintf("%s/%s.json", s.savePath, f.FileName)
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.Encode(f.Item)
	return nil
}

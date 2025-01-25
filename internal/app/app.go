package app

import (
	"sync"
	"sql-injection-server/internal/httpServer"
	"sql-injection-server/internal/filestorage"
)

func Run() error {
	var wg sync.WaitGroup

	filestorage.Start()
	wg.Add(1)
	go httpServer.Start(&wg)
	wg.Wait()
	return nil
}
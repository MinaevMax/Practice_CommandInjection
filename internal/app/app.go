package app

import (
	"sync"
	"command-injection-server/internal/httpServer"
	"command-injection-server/internal/filestorage"
	"time"
)

func Run() error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	ticker := time.NewTicker(time.Minute * 1)

	filestorage.Start()
	mu.Lock()
	filestorage.TmpRefresh()
	mu.Unlock()
	go func(){
		for{
			select{
			case <-ticker.C:
					mu.Lock()
					filestorage.TmpRefresh()
					mu.Unlock()
			}
		}
	} ()
	wg.Add(1)
	go httpServer.Start(&wg, &mu)
	wg.Wait()
	return nil
}
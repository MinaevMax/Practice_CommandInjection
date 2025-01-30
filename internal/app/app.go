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
	passnum := filestorage.TmpRefresh()
	mu.Unlock()
	server := httpServer.Server{Mu: mu, Passnum: passnum}
	go func(){
		for{
			select{
			case <-ticker.C:
					mu.Lock()
					passnum := filestorage.TmpRefresh()
					go server.UpdateAdminPassword(passnum)
					mu.Unlock()
			}
		}
	} ()
	wg.Add(1)
	go server.Start(&wg)
	wg.Wait()
	return nil
}
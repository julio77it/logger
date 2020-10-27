package main

import (
	log "github.com/julio77it/logger"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func test() {
	wg.Add(1)
	defer wg.Done()
	log.Trace("This a int %d", 0)
	log.Debug("This a int %d", 1)
	log.Info("This a int %d", 2)
	log.Warning("This a int %d", 3)
	log.Error("This a int %d", 4)
}
func final() {
	wg.Add(1)
	defer wg.Done()
	log.Fatal("This a int %d", 5)
}

func main() {
	log.NewLogger(os.Stdout, true, 100)
	log.SetLevel(log.TraceLvl)
	for i := 0; i < 10; i++ {
		go test()
		go test()
		go test()
		go test()
		go test()
	}
	wg.Wait()

	go final()

	time.Sleep(time.Second)
}

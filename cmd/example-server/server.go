package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	example_server "github.com/wangweihong/eazycloud/internal/exampleserver"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	example_server.NewApp("exampleserver").Run()
}

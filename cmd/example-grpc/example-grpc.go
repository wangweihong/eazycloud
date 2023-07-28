package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/wangweihong/eazycloud/internal/examplegrpc"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	examplegrpc.NewApp("example-grpc").Run()
}

package main

import (
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/wangweihong/eazycloud/internal/kubeagent"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	kubeagent.NewApp("kubeagent").Run()
}

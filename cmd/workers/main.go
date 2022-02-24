package main

import (
	"os"

	"github.com/jerryan999/CryptoAlert/workers"
)

func main() {
	// run workers
	worker := os.Args[1]
	switch worker {
	case "email":
		workers.SendEmailWorker()
	case "watch":
		workers.WatchCryptoWorker()
	}
}

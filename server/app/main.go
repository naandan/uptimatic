package main

import (
	"fmt"
	"os"
	"uptimatic/cmd/server"
	"uptimatic/cmd/worker"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: nama_app [server|worker|scheduler]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "server":
		server.Start()
	case "worker":
		worker.Start()
	// case "scheduler":
	// 	scheduler.Start()
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}

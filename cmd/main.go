package main

import (
	"fmt"

	"github.com/sergivb01/auditor-worker/worker"
)

func main() {
	srv, err := worker.NewWorker()
	if err != nil {
		fmt.Println(err)
		return
	}

	go sendJobs(srv)

	if err := srv.Start(); err != nil {
		fmt.Println(err)
	}
}

func sendJobs(srv *worker.Worker) {
	fmt.Println("starting to send jobs...")
}

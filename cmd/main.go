package main

import (
	"fmt"

	"github.com/sergivb01/auditor-worker/internal"
)

func main() {
	srv, err := worker.NewWorker()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := srv.Listen(); err != nil {
		fmt.Println(err)
	}
}

package worker

import (
	"fmt"
	"time"

	"github.com/thoas/bokchoy"
)

func handleQueueJob(r *bokchoy.Request) error {
	fmt.Println("Receive request:", r)
	fmt.Println("Request context:", r.Context())
	fmt.Printf("Payload type: %T\n", r.Task.Payload)
	fmt.Printf("Payload: %+v\n", r.Task.Payload)
	fmt.Println("ID:", r.Task.ID)

	r.Task.Result = map[string]interface{}{
		"Time": time.Now(),
		"Name": "Sergi",
	}
	return nil
}
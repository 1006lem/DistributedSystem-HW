package remote_write

import (
	"distributed-system/server/remote-write/api"
	"distributed-system/server/storage"
	"fmt"
	"log"
	"sync"
)

// Main is the main entry point for the remote-write functionality
func Main() {
	// Initialize a wait group
	var wg sync.WaitGroup

	// Initialize JSON file handler
	handler, err := storage.New()
	if err != nil {
		msg := fmt.Sprintf("could not initialize storage handler: %v", err)
		log.Fatal(msg)
	}

	// Assign the storage manager
	storage.Manager = handler

	// Increment the wait group before starting the API handler
	wg.Add(1)

	// Initialize rest API handler
	apiHandler := api.New()
	apiHandler.InitRoutes()
	// Start Server
	apiHandler.Run(&wg)

	// log.Printf("Initialized storage successfully (%s)", storage.Manager.GetHandlerName())

	// Wait for both API to complete
	wg.Wait()
}

package main

import (
	"distributed-system/server/env"
	localWrite "distributed-system/server/local-write"
	remoteWrite "distributed-system/server/remote-write"
	"log"
	"strings"
	"sync"
)

func printASCIIArtLine(line string) {
	log.Printf("| %s |", line)
}

func main() {
	// Initialize a wait group
	var wg sync.WaitGroup

	log.Println("+" + strings.Repeat("-", 30) + "+")

	printASCIIArtLine(" ___   ___  _ __ __   __  ___  _ __ ")
	printASCIIArtLine("/ __| / _ \\| '__|\\ \\ / / / _ \\| '__|")
	printASCIIArtLine("\\__ \\|  __/| |    \\ V / |  __/| |   ")
	printASCIIArtLine("|___/ \\___||_|     \\_/   \\___||_|   ")

	log.Println("+" + strings.Repeat("-", 30) + "+")

	wg.Add(1)

	// Load environment variables
	env.LoadEnvVariables()

	// Start server (local-write or remote-write) through env
	go func() {
		defer wg.Done()
		switch env.Config.Sync {
		case "local-write":
			localWrite.Main()
		case "remote-write":
			remoteWrite.Main()
		default:
			log.Fatal("Invalid argument. Please check json file ..")
		}
	}()

	// Wait for the server to finish
	wg.Wait()

}

package main

import (
	"distributed-system/client/example"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func printASCIIArtLine(line string) {
	log.Printf("| %s |", line)
}

func main() {
	log.Println("+" + strings.Repeat("-", 30) + "+")

	printASCIIArtLine("       _  _               _   ")
	printASCIIArtLine("      | |(_)             | |  ")
	printASCIIArtLine("  ___ | | _   ___  _ __  | |_ ")
	printASCIIArtLine(" / __|| || | / _ \\| '_ \\ | __|")
	printASCIIArtLine("|| (__ | || ||  __/| | | || |_ ")
	printASCIIArtLine("| \\___||_||_| \\___||_| |_| \\__|")

	log.Println("+" + strings.Repeat("-", 30) + "+")

	log.Println("Scenario with 'class info. Distributed System' .. ")

	method, err := LoadEnvVariables()
	if err != nil {
		log.Fatalf("Wrong ENV ! Please set the correct ENV value (get, post, put, patch, delete)")
	}
	args := os.Args[2:] // Get command-line arguments

	switch method {
	case "GET":
		if len(args) == 0 {
			example.GetAllNote()
		} else if len(args) == 1 {
			// Assuming the argument is an integer ID
			id, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatalf("Wrong ID ! Please set the correct ID")
			}
			example.GetNote(id)
		} else {
			log.Fatalf("Invalid number of arguments for GET method")
		}
	case "POST":
		log.Println("POST...")
		example.PostNote()
		return
	case "PATCH":
		example.PatchNote()
		return
	case "PUT":
		example.PutNote()
		return
	case "DELETE":
		example.DeleteNote()
		return
	default:
		log.Println("method: %s", method)
		log.Fatalf("Wrong ID ! Please set the correct ID")
		return
	}

}

func LoadEnvVariables() (string, error) {
	envVar := os.Args[1]
	if envVar == "" {
		return "", fmt.Errorf("environment variable not set")
	}

	// Convert envVar to uppercase
	envVar = strings.ToUpper(envVar)

	return envVar, nil
}

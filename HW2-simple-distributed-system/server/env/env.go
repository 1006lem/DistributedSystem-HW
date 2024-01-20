package env

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Address struct {
	IP   string
	Port int
}
type replicasAddress struct {
	Addresses []Address
}

type Value struct {
	ServicePort int      `json:"servicePort"`
	Sync        string   `json:"sync"`
	Replicas    []string `json:"replicas"`
}

var Config Value
var IsPM bool
var PrevPM int // variable for track prev PM (index for Replicas)
var ReplicaAddresses replicasAddress
var ThisServer int // server's index

func LoadEnvVariables() {
	log.Println("Loading environment variable")

	jsonFilePath := os.Args[1]

	// read json file
	file, err := os.Open(jsonFilePath)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	for _, replica := range Config.Replicas {
		parts := strings.Split(replica, ":")
		if len(parts) != 2 {
			log.Fatalf("Invalid Replica format: %s", parts)
			return
		}
		var newAddress Address

		// Extract IP and port
		newAddress.IP = parts[0]

		port, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatalf("Invalid port format: %s", parts[1])
			return
		}
		newAddress.Port = port
		ReplicaAddresses.Addresses = append(ReplicaAddresses.Addresses, newAddress)
	}

	// Check if the first replica is the current host
	if len(Config.Replicas) > 0 && ReplicaAddresses.Addresses[0].IP == "127.0.0.1" && ReplicaAddresses.Addresses[0].Port == Config.ServicePort {
		IsPM = true
	} else {
		IsPM = false
	}

	// Determine the index of this server within the replica set.
	for index, address := range ReplicaAddresses.Addresses {
		if address.IP == "127.0.0.1" && address.Port == Config.ServicePort {
			ThisServer = index
		}
	}

	// Setting previous PM index (for primary based local-write)
	PrevPM = 0

	log.Println("Service Port:", Config.ServicePort)
	log.Println("Sync:", Config.Sync)
	log.Println("Replicas:")
	for _, address := range ReplicaAddresses.Addresses {
		log.Printf(" ip: %s, port: %d\n", address.IP, address.Port)
	}
	log.Println("isPM: ", IsPM)
	log.Println("indexPM: ", PrevPM)
}

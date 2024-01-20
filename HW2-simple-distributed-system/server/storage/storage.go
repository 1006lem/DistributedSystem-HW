package storage

import (
	"distributed-system/server/common"
	"distributed-system/server/env"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"sync"
)

var Manager Handler

// Handler interface for json file
type Handler interface {
	// read
	GetAllNotes() ([]common.Note, error)
	GetNote(id int) (common.Note, error)

	// write
	OpenAndWriteData(Notes []common.Note) error
	PostNote(newNote common.Note) (common.Note, error)
	GenerateID() (id int, error error)
	PostWithID(id int, newNote common.Note) (common.Note, error)

	PutNote(id int, overrideNote common.Note) (common.Note, error)
	PatchNote(id int, editNote common.Note) (common.Note, error)
	DeleteNote(id int) (common.Note, error)
}

// fileHandler implements Handler (interface)
type fileHandler struct {
	mu       sync.Mutex // Mutex for synchronization, if needed
	filePath string     // File path for data storage
}

// New creates a json file handler according to environment variables
func New() (Handler, error) {
	port := env.Config.ServicePort
	dataPath := "value" + strconv.Itoa(port) + ".json"

	switch env.IsPM {
	case true:
		var initData []common.Note
		initData, err := readInitData("./storage/init-value.json")
		if err != nil {
			return nil, fmt.Errorf("failed to read init data: %v", err)
		}

		err = saveValueData(initData, dataPath)
		if err != nil {
			return nil, fmt.Errorf("failed to save value data: %v", err)
		}

		return newFileHandler(dataPath), nil
	default:
		return newFileHandler(dataPath), nil
	}
}

func readInitData(initDataPath string) ([]common.Note, error) {
	data, err := ioutil.ReadFile(initDataPath)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %v", initDataPath, err)
	}

	var initData []common.Note
	err = json.Unmarshal(data, &initData)
	if err != nil {
		return nil, fmt.Errorf("error decoding %s: %v", initDataPath, err)
	}

	return initData, nil
}

func saveValueData(initData []common.Note, dataPath string) error {
	// Convert the []common.Note to JSON format
	jsonData, err := json.Marshal(initData)
	if err != nil {
		return fmt.Errorf("error encoding init data to JSON: %v", err)
	}

	// Write the JSON data to the file
	err = ioutil.WriteFile(dataPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("error writing data to value.json: %v", err)
	}

	return nil
}

func newFileHandler(filePath string) Handler {
	// Create a new handler
	return &fileHandler{filePath: filePath}
}

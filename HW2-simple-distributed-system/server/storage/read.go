package storage

import (
	"distributed-system/server/common"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// GetAllNotes get all Notes from json file and return all Notes
func (fh *fileHandler) GetAllNotes() ([]common.Note, error) {
	return readNotesFromFile(fh.filePath)
}

// GetNote get Note with id from json file and return the Note
func (fh *fileHandler) GetNote(id int) (common.Note, error) {
	Notes, err := readNotesFromFile(fh.filePath)
	if err != nil {
		return common.Note{}, err
	}

	// Find the Note with the specified ID
	for _, Note := range Notes {
		if Note.ID == id {
			return Note, nil
		}
	}

	// If Note with the specified ID is not found
	return common.Note{}, fmt.Errorf("there is no note with id %d", id)
}

// Helper function to get Note from json file
func readNotesFromFile(filePath string) ([]common.Note, error) {
	// Read data from the specified file path
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Decode JSON data to []common.Note
	var Notes []common.Note
	err = json.Unmarshal(data, &Notes)
	if err != nil {
		return nil, err
	}

	return Notes, nil
}

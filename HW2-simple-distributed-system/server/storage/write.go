package storage

import (
	"distributed-system/server/common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

// PostNote generates a new Note and store it to json file and return posted Note
func (fh *fileHandler) PostNote(newNote common.Note) (common.Note, error) {
	// Lock to ensure exclusive access to the file during write operation
	fh.mu.Lock()
	//defer fh.mu.Unlock()

	// Read data from the specified file path
	data, err := ioutil.ReadFile(fh.filePath)
	if err != nil {
		log.Printf("error reading file: %v", err)
		return common.Note{}, fmt.Errorf("error reading file: %v", err)
	}

	// Decode JSON data to []common.Note
	var Notes []common.Note
	if err := json.Unmarshal(data, &Notes); err != nil {
		log.Printf("error decoding JSON data: %v\n", err)
		return common.Note{}, fmt.Errorf("error decoding JSON data: %v", err)
	}

	// Assign a new ID to the new Note
	newID := generateNewID(Notes)
	newNote.ID = newID

	// Append the new Note to the existing list
	Notes = append(Notes, newNote)

	fh.mu.Unlock()
	// Write the updated data back to the file
	if err := fh.OpenAndWriteData(Notes); err != nil {
		return common.Note{}, fmt.Errorf("error writing data to file: %v", err)
	}

	return newNote, nil
}

// GenerateID generates an ID for new Note
func (fh *fileHandler) GenerateID() (int, error) {
	// Lock to ensure exclusive access to the file during read operation
	fh.mu.Lock()
	defer fh.mu.Unlock()

	// Read data from the specified file path
	data, err := ioutil.ReadFile(fh.filePath)
	if err != nil {
		return 0, err
	}

	// Decode JSON data to []common.Note
	var Notes []common.Note
	err = json.Unmarshal(data, &Notes)
	if err != nil {
		return 0, err
	}

	// Return the next index as the ID
	return len(Notes) + 1, nil
}

// PostWithID generates a new Note with specific id and store it to json file and return posted Note
func (fh *fileHandler) PostWithID(id int, newNote common.Note) (common.Note, error) {
	// Lock to ensure exclusive access to the file during write operation
	fh.mu.Lock()
	defer fh.mu.Unlock()

	// Read data from the specified file path
	data, err := ioutil.ReadFile(fh.filePath)
	if err != nil {
		return common.Note{}, err
	}

	// Decode JSON data to []common.Note
	var Notes []common.Note
	err = json.Unmarshal(data, &Notes)
	if err != nil {
		return common.Note{}, err
	}

	// Check if the specified ID is already in use
	for _, Note := range Notes {
		if Note.ID == id {
			//log.Printf("ID %d is already in use", id)
			return common.Note{}, fmt.Errorf("ID %d is already in use", id)
		}
	}

	// Create a new Note with the specified ID and combine it with newNote
	newNoteWithID := common.Note{
		ID:    id,
		Title: newNote.Title,
		Body:  newNote.Body,
	}

	// Append the new Note to the existing list
	Notes = append(Notes, newNoteWithID)

	// Write the updated data back to the file
	err = fh.OpenAndWriteData(Notes)
	if err != nil {
		return common.Note{}, err
	}

	return newNoteWithID, nil
}

// PutNote overwrites existing Note from json file with given Note and return put Note
func (fh *fileHandler) PutNote(id int, overwriteNote common.Note) (common.Note, error) {
	// Read data from the specified file path
	data, err := ioutil.ReadFile(fh.filePath)
	if err != nil {
		return common.Note{}, err
	}

	// Decode JSON data to []common.Note
	var Notes []common.Note
	err = json.Unmarshal(data, &Notes)
	if err != nil {
		return common.Note{}, err
	}

	// Find the Note with the specified ID
	var updated bool
	for i, Note := range Notes {
		if Note.ID == id {
			Notes[i] = overwriteNote
			overwriteNote.ID = id // Set the ID in the overwriteNote
			updated = true
			break
		}
	}

	// If Note with the specified ID is not found, return an error
	if !updated {
		return common.Note{}, fmt.Errorf("Note with ID %d not found", id)
	}

	// Write the updated data back to the file
	err = fh.OpenAndWriteData(Notes)
	if err != nil {
		return common.Note{}, err
	}

	return overwriteNote, nil
}

// PatchNote update existing Note from json file with given Note and return patched Note
func (fh *fileHandler) PatchNote(id int, editNote common.Note) (common.Note, error) {
	// Read data from the specified file path
	data, err := ioutil.ReadFile(fh.filePath)
	if err != nil {
		return common.Note{}, err
	}

	// Decode JSON data to []common.Note
	var Notes []common.Note
	err = json.Unmarshal(data, &Notes)
	if err != nil {
		return common.Note{}, err
	}

	// Find the Note with the specified ID
	var updated bool
	var updatedNote common.Note
	for i, Note := range Notes {
		if Note.ID == id {
			// Apply the updates to the Note using reflection
			updateFields(&Notes[i], editNote)
			updated = true
			updatedNote = Notes[i]
			break
		}
	}

	// If Note with the specified ID is not found, return an error
	if !updated {
		return common.Note{}, fmt.Errorf("Note with ID %d not found", id)
	}
	// Set the ID in the updatedNote
	updatedNote.ID = id

	// Write the updated data back to the file
	err = fh.OpenAndWriteData(Notes)
	if err != nil {
		return common.Note{}, err
	}

	return updatedNote, nil
}

// DeleteNote delete a Note from json file and return deleted Note
func (fh *fileHandler) DeleteNote(id int) (common.Note, error) {
	// Read data from the specified file path
	data, err := ioutil.ReadFile(fh.filePath)
	if err != nil {
		return common.Note{}, err
	}

	// Decode JSON data to []common.Note
	var Notes []common.Note
	err = json.Unmarshal(data, &Notes)
	if err != nil {
		return common.Note{}, err
	}

	// Find the Note with the specified ID
	var deletedNote common.Note
	for i, Note := range Notes {
		if Note.ID == id {
			// Remove the Note from the list
			deletedNote = Notes[i]
			Notes = append(Notes[:i], Notes[i+1:]...)
			break
		}
	}

	// If Note with the specified ID is not found, return an error
	if deletedNote.ID == 0 {
		return common.Note{}, fmt.Errorf("Note with ID %d not found", id)
	}

	// Write the updated data back to the file
	err = fh.OpenAndWriteData(Notes)
	if err != nil {
		return common.Note{}, err
	}

	return deletedNote, nil
}

// Helper function to generate a new unique ID (index)
func generateNewID(Notes []common.Note) int {
	highestID := 0
	for _, Note := range Notes {
		if Note.ID > highestID {
			highestID = Note.ID
		}
	}
	return highestID + 1
}

// Helper function to update fields of a struct using reflection
func updateFields(target *common.Note, updates common.Note) {
	targetValue := reflect.ValueOf(target).Elem()
	updatesValue := reflect.ValueOf(updates)

	for i := 0; i < targetValue.NumField(); i++ {
		field := targetValue.Type().Field(i).Name
		if value := updatesValue.FieldByName(field); value.IsValid() {
			targetValue.Field(i).Set(value)
		}
	}
}

// OpenAndWriteData opens json file and write Note into json
func (fh *fileHandler) OpenAndWriteData(Notes []common.Note) error {
	// Lock to ensure exclusive access to the file during write operation
	fh.mu.Lock()
	defer fh.mu.Unlock()

	// Open the file in append mode, or create it if it doesn't exist
	file, err := os.OpenFile(fh.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode []common.Note to JSON format
	jsonData, err := json.Marshal(Notes)
	if err != nil {
		return err
	}

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

package api

import (
	"bytes"
	"distributed-system/server/common"
	"distributed-system/server/env"
	"distributed-system/server/storage"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type BroadcastData struct {
	Note  *common.Note
	Error error
}

var BroadcastChannel = make(chan BroadcastData, 1)

const (
	wrongURIMsg             = "Wrong URI (wrong id)"
	wrongBodyMsg            = "Wrong body (JSON format)"
	retrieveNoteErrorMsg    = "Failed to retrieve post_notes.json.go"
	forwardErrorMsg         = "Failed to forward request to PM"
	broadcastingErrorMsg    = "Failed to broadcast request to Replica"
	readFromPMErrorMsg      = "Failed to read response from PM"
	readFromReplicaErrorMsg = "Failed to read response from Replica"
	generateNewNoteErrorMsg
)

func backupRequestToReplica(method string, uri string, url string, noteRequest *common.NoteRequest) (err error) {
	var payloadBytes []byte = nil
	// Marshal the JSON request body
	if method == "POST" || method == "PUT" || method == "PATCH" { // with request body
		var err error
		payloadBytes, err = json.Marshal(noteRequest)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// Create a new PUT request with the specified URL, request body, and content type.
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))

	if err != nil {
		log.Println(err)
		return err
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	log.Printf("REPLICA [REQUEST] Tell backups to update")
	broadcastingResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer broadcastingResponse.Body.Close()
	log.Printf("REPLICA [REPLY] Acknowledge update")
	return nil
}

func backupRequestToAllReplica(method string, uri string, noteRequest *common.NoteRequest) (err error) {
	var wg sync.WaitGroup
	errCh := make(chan error, len(env.ReplicaAddresses.Addresses)-1)

	for _, address := range env.ReplicaAddresses.Addresses[1:] {
		// backup request for all replica
		wg.Add(1)
		go func(addrInfo env.Address) {
			defer wg.Done()

			ReplicaIp := addrInfo.IP
			ReplicaPort := addrInfo.Port
			addr := ReplicaIp + ":" + strconv.Itoa(ReplicaPort)
			url := "http://" + addr + uri

			errFromReplica := backupRequestToReplica(method, uri, url, noteRequest)
			if errFromReplica != nil {
				errCh <- errFromReplica
			}
		}(address)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Check if there were any errors

	for err := range errCh {
		if err != nil {
			log.Println("Error from replica:", err)
			return err
		}
	}

	return nil
}

// postNoteRequest generates a new note based on the client's request
func postNoteRequest(noteRequest common.NoteRequest) (returnBody *common.Note, error error) {
	// typecasting common.NoteRequest to common.Note
	newNote := common.Note{
		Title: noteRequest.Title,
		Body:  noteRequest.Body,
	}

	// requestNote common.Note
	note, err := storage.Manager.PostNote(newNote)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}

	return &note, nil
}

// putNoteRequest overwrites existing Note from json file as a client's note
func putNoteRequest(noteRequest common.NoteRequest, id int) (returnBody *common.Note, error error) {
	// typecasting common.NoteRequest to common.Note
	newNote := common.Note{
		ID:    id,
		Title: noteRequest.Title,
		Body:  noteRequest.Body,
	}

	// requestNote common.Note
	note, err := storage.Manager.PutNote(id, newNote)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	return &note, nil
}

// patchNoteRequest updates a note from existing json file as a client's note
func patchNoteRequest(noteRequest common.NoteRequest, id int) (returnBody *common.Note, error error) {
	// typecasting common.NoteRequest to common.Note
	newNote := common.Note{
		ID:    id,
		Title: noteRequest.Title,
		Body:  noteRequest.Body,
	}

	// requestNote common.Note
	note, err := storage.Manager.PatchNote(id, newNote)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	return &note, nil

}

// deleteNoteRequest deletes a Note from json file with specific id and return ok msg
func deleteNoteRequest(id int) (returnBody *common.Note, error error) {
	note, err := storage.Manager.DeleteNote(id)

	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	return &note, nil
}

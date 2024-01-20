package api

import (
	"bytes"
	"distributed-system/server/common"
	"distributed-system/server/env"
	"distributed-system/server/storage"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// getAllNotesHandler returns all post_notes.json.go in json file
func (h *Handler) getAllNotesHandler(c *gin.Context) {
	notes, err := storage.Manager.GetAllNotes()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, common.Error{
			Msg:    retrieveNoteErrorMsg,
			Method: "GET",
			URI:    "/note",
			Body:   "",
		})
		return
	}

	// Return the post_notes.json.go in the response

	c.JSON(http.StatusOK, notes)
}

// getNoteHandler returns a note with specific id
func (h *Handler) getNoteHandler(c *gin.Context) {
	// Extract the 'id' parameter from the URL
	idString := c.Param("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		uri := "/note/" + idString
		// Handle the error (invalid or non-integer id)
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongURIMsg,
			Method: "GET",
			URI:    uri,
			Body:   "",
		})
		return
	}

	note, err := storage.Manager.GetNote(id)
	if err != nil {
		log.Println(err)
		uri := "/note/" + idString
		c.JSON(http.StatusInternalServerError, common.Error{
			Msg:    retrieveNoteErrorMsg,
			Method: "GET",
			URI:    uri,
			Body:   "",
		})
		return
	}

	c.JSON(http.StatusOK, note)
}

// postNoteHandler get client's note and store it to json file
func (h *Handler) postNoteHandler(c *gin.Context) {
	var requestBody common.NoteRequest

	// Bind the JSON request body to the NoteRequest struct
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongBodyMsg,
			Method: "POST",
			URI:    "/note",
			Body:   "",
		})
		return
	}

	guessID, _ := storage.Manager.GenerateID()

	// if this Server is not PM, then Forward a request to PM
	if !env.IsPM {
		// Extract IP and port
		PMip := env.ReplicaAddresses.Addresses[0].IP
		PMport := env.ReplicaAddresses.Addresses[0].Port
		addr := PMip + ":" + strconv.Itoa(PMport)
		url := "http://" + addr + "/primary"

		// Marshal the JSON request body
		payloadBytes, err := json.Marshal(requestBody)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    wrongBodyMsg,
				Method: "POST",
				URI:    "/note",
				Body:   requestBody.ToString(),
			})
			return
		}

		// Send POST request to PM
		log.Printf("REPLICA [REQUEST] forward request to Primary")
		response, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    forwardErrorMsg,
				Method: "POST",
				URI:    "/primary",
				Body:   requestBody.ToString(),
			})
			return
		}

		log.Printf("REPLICA [REPLY] forward request to Primary")

		// Check if the status code is not OK (200)
		if response.StatusCode != http.StatusOK {
			log.Println("Unexpected status code:", response.StatusCode)
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    forwardErrorMsg,
				Method: "POST",
				URI:    "/primary",
				Body:   requestBody.ToString(),
			})
			return
		}

		// Forward the PM's response to the client
		// Call postBackupHandler with the captured response body
		// BroadcastChannel로 오는 데이터가 있다면 그 데이터를 c.JSON(http.StatusOK, 채널데이터 ) 이렇게 보내고 싶다

		for {
			select {
			case data := <-BroadcastChannel:
				// Handle the broadcasting data
				if data.Note == nil && data.Error == nil {
					continue
				}
				if data.Note.ID != guessID {
					continue
				}
				if data.Note != nil {
					// Send the captured response body through the channel to client
					c.JSON(http.StatusOK, data.Note)
					return
				} else if data.Error != nil {
					log.Println(err)
					c.JSON(http.StatusInternalServerError, common.Error{
						// TODO: channel error
						Msg:    readFromPMErrorMsg,
						Method: "POST",
						URI:    "/note",
						Body:   requestBody.ToString(),
					})
					return
				}
				continue

			default:
				// If no data is received from BroadcastChannel, continue checking
				continue
			}
		}

	} else {
		returnBody, err := postNoteRequest(requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    err.Error(),
				Method: "POST",
				URI:    "/note",
				Body:   requestBody.ToString(),
			})
			return
		}
		// backup request to all Replica
		err = backupRequestToAllReplica("POST", "/backup", &requestBody)
		if err != nil {
			return
		}
		// POST reply to Client

		c.JSON(http.StatusOK, &returnBody)
	}

}

// putPrimaryHandler overwrites existing Note from json file as a client's note
func (h *Handler) putNoteHandler(c *gin.Context) {
	// Extract the 'id' parameter from the URL
	idString := c.Param("id")

	// Convert the 'id' parameter to an integer, assuming it's an integer
	id, err := strconv.Atoi(idString)

	uri := "/note/" + idString
	if err != nil {
		log.Println(err)
		// Handle the error (invalid or non-integer id)
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongURIMsg,
			Method: "PUT",
			URI:    uri,
			Body:   "",
		})
		return
	}

	var requestBody common.NoteRequest

	// Bind the JSON request body to the NoteRequest struct
	if err := c.BindJSON(&requestBody); err != nil {
		log.Println(err)

		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongBodyMsg,
			Method: "PUT",
			URI:    uri,
			Body:   "",
		})
		return
	}

	// if this Server is not PM, then Forward a request to PM
	if !env.IsPM {
		// Extract IP and port
		PMip := env.ReplicaAddresses.Addresses[0].IP
		PMport := env.ReplicaAddresses.Addresses[0].Port
		addr := PMip + ":" + strconv.Itoa(PMport)
		uri := "/primary/" + idString
		url := "http://" + addr + uri

		// Marshal the JSON request body
		payloadBytes, err := json.Marshal(requestBody)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    wrongBodyMsg,
				Method: "PUT",
				URI:    uri,
				Body:   requestBody.ToString(),
			})
			return
		}

		// Create a new PUT request with the specified URL, request body, and content type.
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    forwardErrorMsg,
				Method: "PUT",
				URI:    uri,
				Body:   requestBody.ToString(),
			})
			return
		}

		// Set the content type header
		req.Header.Set("Content-Type", "application/json")

		log.Printf("REPLICA [REQUEST] forward request to Primary")
		putResponse, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    forwardErrorMsg,
				Method: "PUT",
				URI:    uri,
				Body:   requestBody.ToString(),
			})
			return
		}
		defer putResponse.Body.Close()
		log.Printf("REPLICA [REPLY] forward request to Primary")

		// Read the response body
		_, err = ioutil.ReadAll(putResponse.Body)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    readFromPMErrorMsg,
				Method: "PUT",
				URI:    uri,
				Body:   requestBody.ToString(),
			})
			return
		}

		for {
			select {
			case data := <-BroadcastChannel:
				// Handle the broadcasting data
				if data.Note == nil && data.Error == nil {
					continue
				}
				if data.Note.ID != id {
					continue
				}
				if data.Note != nil {
					// Send the captured response body through the channel to client
					c.JSON(http.StatusOK, data.Note)
					return
				} else if data.Error != nil {
					log.Println(err)
					c.JSON(http.StatusInternalServerError, common.Error{
						// TODO: channel error
						Msg:    readFromPMErrorMsg,
						Method: "PUT",
						URI:    uri,
						Body:   requestBody.ToString(),
					})
					return
				}
				continue
			default:
				// If no data is received from BroadcastChannel, continue checking
				continue
			}
		}
	} else {
		returnBody, err := putNoteRequest(requestBody, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    err.Error(),
				Method: "PUT",
				URI:    uri,
				Body:   requestBody.ToString(),
			})
			return
		}
		// backup request to all Replica
		backupURI := "/backup/" + idString
		err = backupRequestToAllReplica("PUT", backupURI, &requestBody)
		if err != nil {
			return
		}
		// PUT reply to Client

		c.JSON(http.StatusOK, returnBody)
	}

}

// patchNoteHandler updates existing Note from json file with given Note and return patched Note
func (h *Handler) patchNoteHandler(c *gin.Context) {
	// Extract the 'id' parameter from the URL
	idString := c.Param("id")

	// Convert the 'id' parameter to an integer, assuming it's an integer
	id, err := strconv.Atoi(idString)

	uri := "/note/" + idString

	if err != nil {
		log.Println(err)
		// Handle the error (invalid or non-integer id)
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongURIMsg,
			Method: "PATCH",
			URI:    uri,
			Body:   "",
		})
		return
	}

	var requestBody common.NoteRequest

	// Bind the JSON request body to the NoteRequest struct
	if err := c.BindJSON(&requestBody); err != nil {
		log.Println(err)

		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongBodyMsg,
			Method: "PATCH",
			URI:    uri,
			Body:   "",
		})
		return
	}

	// if this Server is not PM, then Forward a request to PM
	if !env.IsPM {
		// Extract IP and port
		PMip := env.ReplicaAddresses.Addresses[0].IP
		PMport := env.ReplicaAddresses.Addresses[0].Port
		addr := PMip + ":" + strconv.Itoa(PMport)
		uri := "/primary/" + idString
		url := "http://" + addr + uri

		// Marshal the JSON request body
		payloadBytes, err := json.Marshal(requestBody)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    wrongBodyMsg,
				Method: "PATCH",
				URI:    uri,
				Body:   requestBody.ToString(),
			})
			return
		}

		// Create a new PUT request with the specified URL, request body, and content type.
		req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    forwardErrorMsg,
				Method: "PATCH",
				URI:    uri,
				Body:   requestBody.ToString(),
			})
			return
		}

		// Set the content type header
		req.Header.Set("Content-Type", "application/json")

		log.Printf("REPLICA [REQUEST] forward request to Primary")
		patchResponse, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    forwardErrorMsg,
				Method: "PATCH",
				URI:    uri,
				Body:   requestBody.ToString(),
			})
			return
		}
		defer patchResponse.Body.Close()

		log.Printf("REPLICA [REPLY] forward request to Primary")

		// Read the response body
		_, err = ioutil.ReadAll(patchResponse.Body)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    readFromPMErrorMsg,
				Method: "PATCH",
				URI:    uri,
				Body:   requestBody.ToString(),
			})
			return
		}

		for {
			select {
			case data := <-BroadcastChannel:
				// Handle the broadcasting data
				if data.Note == nil && data.Error == nil {
					continue
				}
				if data.Note.ID != id {
					continue
				}
				if data.Note != nil {
					// Send the captured response body through the channel to client

					c.JSON(http.StatusOK, data.Note)
					return
				} else if data.Error != nil {
					log.Println(err)
					c.JSON(http.StatusInternalServerError, common.Error{
						// TODO: channel error
						Msg:    readFromPMErrorMsg,
						Method: "PATCH",
						URI:    uri,
						Body:   requestBody.ToString(),
					})
					return
				}
				continue
			default:
				// If no data is received from BroadcastChannel, continue checking
				continue
			}
		}
	} else {
		returnBody, err := patchNoteRequest(requestBody, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    err.Error(),
				Method: "PATCH",
				URI:    uri,
				Body:   requestBody.ToString(),
			})
			return
		}
		// backup request to all Replica
		backupURI := "/backup/" + idString
		err = backupRequestToAllReplica("PATCH", backupURI, &requestBody)
		if err != nil {
			return
		}

		// PATCH reply to Client
		c.JSON(http.StatusOK, returnBody)

	}
}

// deleteNoteHandler deletes a Note from json file with specific id and return ok msg
func (h *Handler) deleteNoteHandler(c *gin.Context) {
	// Extract the 'id' parameter from the URL
	idString := c.Param("id")

	// Convert the 'id' parameter to an integer, assuming it's an integer
	id, err := strconv.Atoi(idString)

	uri := "/note/" + idString

	if err != nil {
		log.Println(err)
		// Handle the error (invalid or non-integer id)
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongURIMsg,
			Method: "DELETE",
			URI:    uri,
			Body:   "",
		})
		return
	}

	// if this Server is not PM, then Forward a request to PM
	if !env.IsPM {
		// Extract IP and port
		PMip := env.ReplicaAddresses.Addresses[0].IP
		PMport := env.ReplicaAddresses.Addresses[0].Port
		addr := PMip + ":" + strconv.Itoa(PMport)
		uri := "/primary/" + idString
		url := "http://" + addr + uri

		// Create a new PUT request with the specified URL, request body, and content type.
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    forwardErrorMsg,
				Method: "DELETE",
				URI:    uri,
				Body:   "",
			})
			return
		}

		log.Printf("REPLICA [REQUEST] forward request to Primary")
		deleteResponse, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    forwardErrorMsg,
				Method: "DELETE",
				URI:    uri,
				Body:   "",
			})
			return
		}
		defer deleteResponse.Body.Close()
		log.Printf("REPLICA [REPLY] forward request to Primary")

		// Read the response body
		_, err = ioutil.ReadAll(deleteResponse.Body)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    readFromPMErrorMsg,
				Method: "DELETE",
				URI:    uri,
				Body:   "",
			})
			return
		}

		for {
			select {
			case data := <-BroadcastChannel:
				// Handle the broadcasting data
				if data.Note == nil && data.Error == nil {
					continue
				}

				if data.Note.ID != id {
					continue
				}
				if data.Note != nil {
					// Send the captured response body through the channel to client
					returnOKMsg := common.DeleteResponse{
						Msg: "OK",
					}
					// DELETE reply to Client
					c.JSON(http.StatusOK, returnOKMsg)
					return
				} else if data.Error != nil {
					log.Println(err)
					c.JSON(http.StatusInternalServerError, common.Error{
						// TODO: channel error
						Msg:    readFromPMErrorMsg,
						Method: "DELETE",
						URI:    uri,
					})
					return
				}
				continue
			default:
				// If no data is received from BroadcastChannel, continue checking
				continue
			}
		}
	} else {
		_, err := deleteNoteRequest(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{
				Msg:    err.Error(),
				Method: "DELETE",
				URI:    uri,
				Body:   "",
			})
			return
		}
		// backup request to all Replica
		backupURI := "/backup/" + idString
		err = backupRequestToAllReplica("DELETE", backupURI, nil)
		if err != nil {
			return
		}

		// Return ok msg
		returnOKMsg := common.DeleteResponse{
			Msg: "OK",
		}
		// DELETE reply to Client
		c.JSON(http.StatusOK, returnOKMsg)
		return
	}
}

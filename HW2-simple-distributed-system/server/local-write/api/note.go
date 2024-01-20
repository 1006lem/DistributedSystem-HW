package api

import (
	"distributed-system/server/common"
	"distributed-system/server/env"
	"distributed-system/server/storage"
	"github.com/gin-gonic/gin"
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

	// Get Data from prev pm with given id
	guessID, _ := storage.Manager.GenerateID()
	_, err := GETDataFromPrevPM(guessID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{
			Msg:    err.Error(),
			Method: "POST",
			URI:    "/note",
			Body:   requestBody.ToString(),
		})
		return
	}

	// Add new note into local storage
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

	// POST reply to Client
	c.JSON(http.StatusOK, &returnBody)

	// Update prev_pm (as this server)
	env.PrevPM = env.ThisServer

	//log.Println("PM is me")
	// backup request to all Replica
	err = backupRequestToAllReplica("POST", "/backup", &requestBody)
	if err != nil {
		return
	}
	return
}

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

	// PUT reply to Client
	c.JSON(http.StatusOK, returnBody)

	// Update prev_pm (as this server)
	env.PrevPM = env.ThisServer

	// backup request to all Replica
	backupURI := "/backup/" + idString
	err = backupRequestToAllReplica("PUT", backupURI, &requestBody)
	if err != nil {
		return
	}
	return
}

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
	// PATCH reply to Client
	c.JSON(http.StatusOK, returnBody)

	// Update prev_pm (as this server)
	env.PrevPM = env.ThisServer

	// backup request to all Replica
	backupURI := "/backup/" + idString
	err = backupRequestToAllReplica("PATCH", backupURI, &requestBody)
	if err != nil {
		return
	}
	return
}

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

	_, err = deleteNoteRequest(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{
			Msg:    err.Error(),
			Method: "DELETE",
			URI:    uri,
			Body:   "",
		})
		return
	}
	// Return ok msg
	returnOKMsg := common.DeleteResponse{
		Msg: "OK",
	}
	// DELETE reply to Client
	c.JSON(http.StatusOK, returnOKMsg)

	// Update prev_pm (as this server)
	env.PrevPM = env.ThisServer

	// backup request to all Replica
	backupURI := "/backup/" + idString
	err = backupRequestToAllReplica("DELETE", backupURI, nil)
	if err != nil {
		return
	}

	return
}

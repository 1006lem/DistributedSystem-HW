package api

import (
	"distributed-system/server/common"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// TODO: line 91 return with c(internal server err) -> err
// postPrimaryHandler generates a new note based on the client's request and stores it in a JSON file.
func (h *Handler) postPrimaryHandler(c *gin.Context) {
	var requestBody common.NoteRequest

	// Bind the JSON request body to the NoteRequest struct
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    "Wrong body (JSON format)",
			Method: "POST",
			URI:    "/note",
			Body:   "",
		})
		return
	}

	_, err := postNoteRequest(requestBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{
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
	c.JSON(http.StatusOK, nil)

}

// putPrimaryHandler overwrites a note from a json file as a client's note
func (h *Handler) putPrimaryHandler(c *gin.Context) {
	// Extract the 'id' parameter from the URL
	idString := c.Param("id")
	uri := "/note/" + idString

	// Convert the 'id' parameter to an integer, assuming it's an integer
	id, err := strconv.Atoi(idString)
	if err != nil {
		// Handle the error (invalid or non-integer id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or non-integer ID"})
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

	_, err = putNoteRequest(requestBody, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    err.Error(),
			Method: "PUT",
			URI:    uri,
			Body:   requestBody.ToString(),
		})
		return
	}

	// backup request to all Replica
	backupUri := "/backup/" + idString
	err = backupRequestToAllReplica("PUT", backupUri, &requestBody)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, nil)
}

// patchNoteHandler updates existing Note from json file with given Note and return patched Note
func (h *Handler) patchPrimaryHandler(c *gin.Context) {
	// Extract the 'id' parameter from the URL
	idString := c.Param("id")
	uri := "/note/" + idString

	// Convert the 'id' parameter to an integer, assuming it's an integer
	id, err := strconv.Atoi(idString)
	if err != nil {
		// Handle the error (invalid or non-integer id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or non-integer ID"})
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

	_, err = patchNoteRequest(requestBody, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    err.Error(),
			Method: "PATCH",
			URI:    uri,
			Body:   requestBody.ToString(),
		})
		return
	}

	// backup request to all Replica
	backupUri := "/backup/" + idString
	err = backupRequestToAllReplica("PATCH", backupUri, &requestBody)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, nil)
}

// deletePrimaryHandler deletes a Note from json file with specific id and return ok msg
func (h *Handler) deletePrimaryHandler(c *gin.Context) {
	// Extract the 'id' parameter from the URL
	idString := c.Param("id")
	uri := "/note/" + idString
	// Convert the 'id' parameter to an integer, assuming it's an integer
	id, err := strconv.Atoi(idString)
	if err != nil {
		// Handle the error (invalid or non-integer id)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or non-integer ID"})
		return
	}

	_, err = deleteNoteRequest(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    err.Error(),
			Method: "DELETE",
			URI:    uri,
			Body:   "",
		})
		return
	}

	// backup request to all Replica
	backupUri := "/backup/" + idString
	err = backupRequestToAllReplica("DELETE", backupUri, nil)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, nil)
}

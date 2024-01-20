package api

import (
	"distributed-system/server/common"
	"distributed-system/server/env"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) postBackupHandler(c *gin.Context) {
	var requestBody common.NoteRequestWithNewPM
	// Bind the JSON request body to the NoteRequest struct
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongBodyMsg,
			Method: "POST",
			URI:    "/backup",
			Body:   "",
		})
		return
	}

	requestNoteBody := common.NoteRequest{
		Title: requestBody.Title,
		Body:  requestBody.Body,
	}

	// change local storage data
	returnBody, err := postNoteRequest(requestNoteBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{
			Msg:    err.Error(),
			Method: "POST",
			URI:    "/backup",
			Body:   requestNoteBody.ToString(),
		})
		return
	}

	// Update
	c.JSON(http.StatusOK, returnBody)
	//log.Println("PM is", requestBody.Pm)
	env.PrevPM = requestBody.Pm

	return
}

// getNoteHandler return a memo with specific id
func (h *Handler) putBackupHandler(c *gin.Context) {
	// Extract the 'id' parameter from the URL
	idString := c.Param("id")

	// Convert the 'id' parameter to an integer
	id, err := strconv.Atoi(idString)
	if err != nil {
		// Handle the error (invalid or non-integer id)
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongURIMsg,
			Method: "PUT",
			URI:    "/backup" + strconv.Itoa(id),
			Body:   "",
		})
		return
	}

	var requestBody common.NoteRequestWithNewPM

	// Bind the JSON request body to the NoteRequest struct
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongBodyMsg,
			Method: "PUT",
			URI:    "/backup" + strconv.Itoa(id),
			Body:   "",
		})
		return
	}

	requestNoteBody := common.NoteRequest{
		Title: requestBody.Title,
		Body:  requestBody.Body,
	}

	// change local storage data
	returnBody, err := putNoteRequest(requestNoteBody, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{
			Msg:    err.Error(),
			Method: "PUT",
			URI:    "/backup" + strconv.Itoa(id),
			Body:   requestNoteBody.ToString(),
		})
		return
	}

	// Update
	c.JSON(http.StatusOK, returnBody)

	env.PrevPM = requestBody.Pm

	return
}

// getNoteHandler return a memo with specific id
func (h *Handler) patchBackupHandler(c *gin.Context) {
	// Extract the 'id' parameter from the URL
	idString := c.Param("id")

	// Convert the 'id' parameter to an integer
	id, err := strconv.Atoi(idString)
	if err != nil {
		// Handle the error (invalid or non-integer id)
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongURIMsg,
			Method: "PATCH",
			URI:    "/backup" + strconv.Itoa(id),
			Body:   "",
		})
		return
	}

	var requestBody common.NoteRequestWithNewPM

	// Bind the JSON request body to the NoteRequest struct
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongBodyMsg,
			Method: "PATCH",
			URI:    "/backup" + strconv.Itoa(id),
			Body:   "",
		})
		return
	}

	requestNoteBody := common.NoteRequest{
		Title: requestBody.Title,
		Body:  requestBody.Body,
	}

	// change local storage data
	returnBody, err := patchNoteRequest(requestNoteBody, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{
			Msg:    err.Error(),
			Method: "PATCH",
			URI:    "/backup" + strconv.Itoa(id),
			Body:   requestNoteBody.ToString(),
		})
		return
	}

	c.JSON(http.StatusOK, returnBody)

	env.PrevPM = requestBody.Pm

	return
}

// getNoteHandler return a memo with specific id
func (h *Handler) deleteBackupHandler(c *gin.Context) {
	// Extract the 'id' parameter from the URL
	idString := c.Param("id")

	// Convert the 'id' parameter to an integer
	id, err := strconv.Atoi(idString)
	if err != nil {
		// Handle the error (invalid or non-integer id)
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongURIMsg,
			Method: "DELETE",
			URI:    "/backup" + strconv.Itoa(id),
			Body:   "",
		})
		return
	}

	var requestBody common.NoteRequestWithNewPM
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongBodyMsg,
			Method: "DELETE",
			URI:    "/backup",
			Body:   "",
		})
		return
	}

	// change local storage data
	returnBody, err := deleteNoteRequest(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{
			Msg:    err.Error(),
			Method: "DELETE",
			URI:    "/backup" + strconv.Itoa(id),
			Body:   "",
		})
		return
	}

	c.JSON(http.StatusOK, returnBody)

	env.PrevPM = requestBody.Pm

	return
}

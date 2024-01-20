package api

import (
	"distributed-system/server/common"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) postBackupHandler(c *gin.Context) {
	var requestBody common.NoteRequest
	// 초기화
	initialBroadcastData := BroadcastData{
		Note:  nil,
		Error: nil,
	}
	go func() {
		BroadcastChannel <- initialBroadcastData
	}()

	// Bind the JSON request body to the NoteRequest struct
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongBodyMsg,
			Method: "POST",
			URI:    "/backup",
			Body:   "",
		})
		// Send err through the channel
		broadcastData := BroadcastData{
			Note:  nil,
			Error: err,
		}

		BroadcastChannel <- broadcastData
		return
	}

	// change local storage data
	returnBody, err := postNoteRequest(requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{
			Msg:    err.Error(),
			Method: "POST",
			URI:    "/backup",
			Body:   requestBody.ToString(),
		})
		// Send err through the channel
		broadcastData := BroadcastData{
			Note:  nil,
			Error: err,
		}

		BroadcastChannel <- broadcastData
		return
	}

	c.JSON(http.StatusOK, returnBody)

	// Send the captured response body through the channel
	broadcastData := BroadcastData{
		Note:  returnBody,
		Error: nil,
	}

	go func() {
		BroadcastChannel <- broadcastData
	}()
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

	var requestBody common.NoteRequest
	// 초기화
	initialBroadcastData := BroadcastData{
		Note:  nil,
		Error: nil,
	}
	go func() {
		BroadcastChannel <- initialBroadcastData
	}()

	// Bind the JSON request body to the NoteRequest struct
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongBodyMsg,
			Method: "PUT",
			URI:    "/backup" + strconv.Itoa(id),
			Body:   "",
		})
		// Send err through the channel
		broadcastData := BroadcastData{
			Note:  nil,
			Error: err,
		}

		BroadcastChannel <- broadcastData
		return
	}

	// change local storage data
	returnBody, err := putNoteRequest(requestBody, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{
			Msg:    err.Error(),
			Method: "PUT",
			URI:    "/backup" + strconv.Itoa(id),
			Body:   requestBody.ToString(),
		})
		// Send err through the channel
		broadcastData := BroadcastData{
			Note:  nil,
			Error: err,
		}

		BroadcastChannel <- broadcastData
		return
	}

	c.JSON(http.StatusOK, returnBody)

	// Send the captured response body through the channel
	broadcastData := BroadcastData{
		Note:  returnBody,
		Error: nil,
	}

	go func() {
		BroadcastChannel <- broadcastData

	}()
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

	var requestBody common.NoteRequest
	// 초기화
	initialBroadcastData := BroadcastData{
		Note:  nil,
		Error: nil,
	}
	go func() {
		BroadcastChannel <- initialBroadcastData
	}()

	// Bind the JSON request body to the NoteRequest struct
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{
			Msg:    wrongBodyMsg,
			Method: "PATCH",
			URI:    "/backup" + strconv.Itoa(id),
			Body:   "",
		})
		// Send err through the channel
		broadcastData := BroadcastData{
			Note:  nil,
			Error: err,
		}

		BroadcastChannel <- broadcastData
		return
	}

	// change local storage data
	returnBody, err := patchNoteRequest(requestBody, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{
			Msg:    err.Error(),
			Method: "PATCH",
			URI:    "/backup" + strconv.Itoa(id),
			Body:   requestBody.ToString(),
		})
		// Send err through the channel
		broadcastData := BroadcastData{
			Note:  nil,
			Error: err,
		}

		BroadcastChannel <- broadcastData
		return
	}

	c.JSON(http.StatusOK, returnBody)

	// Send the captured response body through the channel
	broadcastData := BroadcastData{
		Note:  returnBody,
		Error: nil,
	}

	go func() {
		BroadcastChannel <- broadcastData
	}()
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

	// 초기화
	initialBroadcastData := BroadcastData{
		Note:  nil,
		Error: nil,
	}
	go func() {
		BroadcastChannel <- initialBroadcastData
	}()

	// change local storage data
	returnBody, err := deleteNoteRequest(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{
			Msg:    err.Error(),
			Method: "DELETE",
			URI:    "/backup" + strconv.Itoa(id),
			Body:   "",
		})
		// Send err through the channel
		broadcastData := BroadcastData{
			Note:  nil,
			Error: err,
		}

		BroadcastChannel <- broadcastData
		return
	}

	c.JSON(http.StatusOK, returnBody)

	// Send the captured response body through the channel
	broadcastData := BroadcastData{
		Note:  returnBody,
		Error: nil,
	}

	go func() {
		BroadcastChannel <- broadcastData
	}()
	return

}

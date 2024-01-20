package api

import (
	"distributed-system/server/common"
	"distributed-system/server/storage"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// getNoteHandler return a memo with specific id
func (h *Handler) getPrimaryHandler(c *gin.Context) {
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

package common

import "fmt"

// NoteRequest represents the data structure used for creating/editing a note
type NoteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// ToString converts NoteRequest to String
func (n NoteRequest) ToString() string {
	return fmt.Sprintf("Title: %s, Body: %s", n.Title, n.Body)
}

// NoteRequestWithNewPM represents the data structure used for creating/editing a note and update PM
type NoteRequestWithNewPM struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Pm    int    `json:"pm"`
}

package common

// Notes is a collection of NoteResponse objects
type Notes []Note

// Note represents contents of a note
type Note struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Msg   string `json:"msg"`
}

type Error struct {
	Msg    string `json:"msg"`
	Method string `json:"method"`
	URI    string `json:"uri"`
	Body   string `json:"body"`
}

type DeleteResponse struct {
	Msg string `json:"Msg"`
}

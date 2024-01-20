package example

import (
	"distributed-system/server/common"
	"distributed-system/server/env"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func DeleteNote() {
	// id 1
	Delete(1, addresses[0])

	// id 2
	Delete(2, addresses[1])

	// id 2
	Delete(3, addresses[2])
}

func Delete(id int, address env.Address) {
	uri := "/note/" + strconv.Itoa(id)
	url := "http://" + address.IP + ":" + strconv.Itoa(address.Port) + uri
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Error creating DELETE request:", err)
		return
	}

	// Send the DELETE request
	log.Printf("CLIENT [REQUEST] [DELETE] %s\n", uri)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending DELETE request:", err)
		return
	}
	defer resp.Body.Close()

	type DeleteResponse struct {
		Msg string `json:"msg"`
	}
	var deleteResponse DeleteResponse

	// 응답 확인
	if resp.StatusCode == http.StatusOK {
		err := json.NewDecoder(resp.Body).Decode(&deleteResponse)
		if err != nil {
			log.Println("Error decoding JSON response:", err)
			return
		}
		log.Printf("CLIENT [REPLY] [DELETE] %s %+v\n", uri, deleteResponse)
		log.Printf("  Response from %s:%d\n", address.IP, address.Port)
	} else if resp.StatusCode == http.StatusInternalServerError {
		var errorBody common.Error
		err := json.NewDecoder(resp.Body).Decode(&errorBody)
		if err != nil {
			log.Println("Error decoding JSON response:", err)
			return
		}
		log.Printf("Error (Internal Server Err): %+v\n", errorBody)
		log.Printf("  Response from %s:%d\n", address.IP, address.Port)
	}
}

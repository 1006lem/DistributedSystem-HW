package example

import (
	"distributed-system/server/common"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func GetAllNote() {
	// GET 요청 (8081, 8082, 8083)
	for _, address := range addresses {
		url := "http://" + address.IP + ":" + strconv.Itoa(address.Port) + "/note"

		log.Printf("CLIENT [REQUEST] [GET] /note\n")
		resp, err := http.Get(url)
		if err != nil {
			log.Println("Error sending GET request:", err)
			return
		}
		defer resp.Body.Close()

		// 응답 확인
		if resp.StatusCode == http.StatusOK {
			var notes common.Notes
			err := json.NewDecoder(resp.Body).Decode(&notes)
			if err != nil {
				log.Println("Error decoding JSON response:", err)
				return
			}
			log.Printf("CLIENT [REPLY] [GET] /note %+v\n", notes)
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
}

func GetNote(id int) {
	// GET 요청 (8081, 8082, 8083)
	for _, address := range addresses {
		uri := "/note/" + strconv.Itoa(id)
		url := "http://" + address.IP + ":" + strconv.Itoa(address.Port) + uri

		log.Printf("CLIENT [REQUEST] [GET] %s\n", uri)
		resp, err := http.Get(url)
		if err != nil {
			log.Println("Error sending GET request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			var note common.Note
			err := json.NewDecoder(resp.Body).Decode(&note)
			if err != nil {
				log.Println("Error decoding JSON response:", err)
				return
			}
			log.Printf("CLIENT [REPLY] [GET] %s %+v\n", uri, note)
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
}

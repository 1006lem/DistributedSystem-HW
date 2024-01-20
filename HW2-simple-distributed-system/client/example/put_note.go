package example

import (
	"bytes"
	"distributed-system/server/common"
	"distributed-system/server/env"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func Put(buf []byte, address env.Address, id int, note common.Note) {
	uri := "/note/" + strconv.Itoa(id)
	url := "http://" + address.IP + ":" + strconv.Itoa(address.Port) + uri
	// Create a new PATCH request with the specified URL, request body, and content type.
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println("Error creating PUT request:", err)
		return
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Send the PATCH request
	log.Printf("CLIENT [REQUEST] [PUT] %s %+v\n", uri, note)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending PUT request:", err)
		return
	}
	defer resp.Body.Close()

	// 응답 확인
	if resp.StatusCode == http.StatusOK {
		var note common.Note
		err := json.NewDecoder(resp.Body).Decode(&note)
		if err != nil {
			fmt.Println("Error decoding JSON response:", err)
			return
		}
		log.Printf("CLIENT [REPLY] [PUT] %s %+v\n", uri, note)
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

func PutNote() {
	// JSON 파일 읽기
	jsonData, err := ioutil.ReadFile("./example/json/put_notes.json")
	if err != nil {
		fmt.Println("JSON 파일 읽기 오류:", err)
		return
	}

	// JSON 데이터를 슬라이스로 언마샬링
	var notes []common.Note
	err = json.Unmarshal(jsonData, &notes)
	if err != nil {
		fmt.Println("JSON 언마샬링 오류:", err)
		return
	}

	if len(notes) > 0 {
		Note1 := notes[0]
		Note2 := notes[1]
		Note3 := notes[2]

		// JSON 데이터를 바이트 슬라이스로 변환
		payloadBytes0, err := json.Marshal(Note1)
		if err != nil {
			fmt.Println("JSON 마샬링 오류:", err)
			return
		}
		// JSON 데이터를 바이트 슬라이스로 변환
		payloadBytes1, err := json.Marshal(Note2)
		if err != nil {
			fmt.Println("JSON 마샬링 오류:", err)
			return
		}
		// JSON 데이터를 바이트 슬라이스로 변환
		payloadBytes2, err := json.Marshal(Note3)
		if err != nil {
			fmt.Println("JSON 마샬링 오류:", err)
			return
		}

		// 초기 데이터 (id 1, 2, 3) 수정
		Put(payloadBytes0, addresses[0], 1, Note1)
		Put(payloadBytes1, addresses[1], 2, Note2)
		Put(payloadBytes2, addresses[2], 3, Note3)
	}
}

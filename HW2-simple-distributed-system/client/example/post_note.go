package example

import (
	"bytes"
	"distributed-system/server/common"
	"distributed-system/server/env"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var addresses = []env.Address{
	{IP: "127.0.0.1", Port: 8081},
	{IP: "127.0.0.1", Port: 8082},
	{IP: "127.0.0.1", Port: 8083},
}

func PostNote() {
	// JSON 파일 읽기
	jsonData, err := ioutil.ReadFile("./example/json/post_notes.json")
	if err != nil {
		log.Println("JSON 파일 읽기 오류:", err)
		return
	}

	// JSON 데이터를 슬라이스로 언마샬링
	var notes []common.Note
	err = json.Unmarshal(jsonData, &notes)
	if err != nil {
		log.Println("JSON 언마샬링 오류:", err)
		return
	}

	if len(notes) > 0 {
		for i, note := range notes {
			// JSON 데이터를 바이트 슬라이스로 변환
			payloadBytes, err := json.Marshal(note)
			if err != nil {
				log.Println("JSON 마샬링 오류:", err)
				return
			}

			// 이전 POST가 완료될 때까지 기다림
			Post(payloadBytes, addresses[i], note)
		}
	}
}

func Post(buf []byte, address env.Address, note common.Note) {
	url := "http://" + address.IP + ":" + strconv.Itoa(address.Port) + "/note"

	log.Printf("CLIENT [REQUEST] [POST] /note %+v\n", note)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()

	// 응답 확인
	if resp.StatusCode == http.StatusOK {
		var note common.Note
		err := json.NewDecoder(resp.Body).Decode(&note)
		if err != nil {
			log.Println("Error decoding JSON response:", err)
			return
		}

		log.Printf("CLIENT [REPLY] [POST] /note %+v\n", note)
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

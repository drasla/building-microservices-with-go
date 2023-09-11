package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type helloWorldResponse struct {
	// Go에서는 소문자로 시작되는 프로퍼티는 export가 불가능하기 때문에 필요할 경우 태그를 이용함
	// json 출력 필드를 "message"르 변경함
	Message string `json:"message"`
	// 구조체에는 포함되나, json에서 제외하기 위해서는 태그에 "-" 를 쓴다
	Author string `json:"-"`
	// 값이 비어 있으면 필드를 출력하지 않도록 하려면 omitempty 를 추가함
	Date string `json:",omitempty"`
	//출력을 문자열로 변환하고 이름을 id로 변경함
	Id int `json:"id, string"`
}

type helloWorldRequest struct {
	Name string `json:"name"`
}

func hellowWorldHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var request helloWorldRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	response := helloWorldResponse{Message: "Hello " + request.Name}

	encoder := json.NewEncoder(w)
	encoder.Encode(response)
}

func main() {
	port := 8080

	http.HandleFunc("/helloworld", hellowWorldHandler)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

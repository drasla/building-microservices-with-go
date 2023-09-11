package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"
)

type validationContextKey string

type helloWorldRequest struct {
	Name string `json:"name"`
}

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

type validationHandler struct {
	next http.Handler
}

const port = 8080

func newValidationHandler(next http.Handler) http.Handler {
	return validationHandler{next: next}
}

func (h validationHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request helloWorldRequest
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&request)
	if err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}

	c := context.WithValue(r.Context(), validationContextKey("name"), request.Name)
	r = r.WithContext(c)

	h.next.ServeHTTP(rw, r)
}

type helloWorldHandler struct{}

func newHelloWorldHandler() http.Handler {
	return helloWorldHandler{}
}

func (h helloWorldHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	name := r.Context().Value(validationContextKey("name")).(string)
	response := helloWorldResponse{Message: "Hello " + name}
	encoder := json.NewEncoder(rw)
	encoder.Encode(response)
}

func fetchGoogle(t *testing.T) {
	r, _ := http.NewRequest("GET", "https://google.com", nil)

	timeoutRequest, cancelFunc := context.WithTimeout(r.Context(), 1*time.Millisecond)
	defer cancelFunc()

	r = r.WithContext(timeoutRequest)

	_, err := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func server() {
	handler := newValidationHandler(newHelloWorldHandler())

	http.Handle("/helloworld", handler)

	log.Printf("Server starting on port %v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func main() {
	server()
}

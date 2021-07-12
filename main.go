package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const port string = ":4000"

type URLDescription struct {
	URL string `json:"url"`
	Method string `json:"-"`
	Description string `json:"description"`
	Payload string `json:"payload,omitempty"`
}

func documentation(rw http.ResponseWriter, r *http.Request)  {
	data:= []URLDescription{
		{
			URL: "/",
			Method : "GET",
			Description: "See Documentation",
		},
		{
			URL: "/blocks",
			Method: "POST",
			Description:"Add A Block",
			Payload:"data:string",
		},
	}
	rw.Header().Add("Content-Type", "application/json")

	// Marshal은 Value를 JSON으로 변환해준다. 근데 브라우저는 Header에 JSON이라고 명시하지 않으면 JSON인걸 모른다.
	// b, err := json.Marshal(data)
	// utils.HandleErr(err)
	// fmt.Fprintf(rw,"%s",b)

	//위 3줄과 아래 1줄이 같음
	json.NewEncoder(rw).Encode(data)
}

func main(){
	http.HandleFunc("/", documentation)
	fmt.Printf("Listening on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
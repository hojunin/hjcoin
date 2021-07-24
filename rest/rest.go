package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hojunin/hjcoin/blockchain"
	"github.com/hojunin/hjcoin/utils"
)


var port string

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"` 
}

type urlDescription struct {
	URL url `json:"url"`
	Method string `json:"method"`
	Description string `json:"description"`
	Payload string `json:"payload,omitempty"`
}

type addBlockBody struct{
	Message string
}

type url string

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s",port, u)
	return []byte(url), nil
}

func documentation(rw http.ResponseWriter, r *http.Request)  {
	data:= []urlDescription{
		{
			URL: url("/"),
			Method : "GET",
			Description: "See Documentation",
		},
		{
			URL: url("/blocks"),
			Method : "GET",
			Description: "See All Blocks",
		},
		{
			URL: url("/blocks"),
			Method: "POST",
			Description:"Add A Block",
			Payload:"data:string",
		},
		{
			URL: url("/blocks/{hash}"),
			Method : "GET",
			Description: "See A Block",
		},
	}
	rw.Header().Add("Content-Type", "application/json")

	b, err := json.Marshal(data)
	utils.HandleErr(err)
	fmt.Fprintf(rw,"%s",b)
}

func block(rw http.ResponseWriter, r*http.Request)  {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block ,err:= blockchain.FindBlock(hash)

	encoder := json.NewEncoder(rw)
	if err==blockchain.ErrNotFound{
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	}else{
		encoder.Encode(block)
	}
}



func blocks(rw http.ResponseWriter, r *http.Request)  {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.Blockchain().Blocks())
	case "POST":
		var addBlockBody addBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.Blockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	default:
		break
	}
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler  {
	return http.HandlerFunc(func (rw http.ResponseWriter, r *http.Request)  {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
	
}

func Start(aport int)  {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aport)
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
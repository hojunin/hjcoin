package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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
			URL: url("/blocks/{id}"),
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
	id, err := strconv.Atoi(vars["height"])
	utils.HandleErr(err)
	block ,err:= blockchain.GetBlockchain().GetBlock(id)

	encoder := json.NewEncoder(rw)
	if err==blockchain.ErrNotFound{
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	}else{
		encoder.Encode(block)
	}
}



func blocks (rw http.ResponseWriter, r *http.Request)  {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlock())
		break
	case "POST":
		var addBlockBody addBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	default:
		break
	}
}

func Start(aport int)  {
	handler := mux.NewRouter()
	port = fmt.Sprintf(":%d", aport)
	handler.HandleFunc("/", documentation).Methods("GET")
	handler.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	handler.HandleFunc("/blocks/{height:[0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, handler))
}
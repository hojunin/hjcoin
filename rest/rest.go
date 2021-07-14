package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hojunin/hjcoin/blockchain"
	"github.com/hojunin/hjcoin/utils"
)


const port string = ":4000"

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

func Start()  {
	http.HandleFunc("/", documentation)
	http.HandleFunc("/blocks", blocks)
	fmt.Printf("Listening on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
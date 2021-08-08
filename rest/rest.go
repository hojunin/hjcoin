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

type balanceResponse struct {
	Address string `json:"address"`
	Balance int 	`json:"balance"`
}

type addTxPayload struct {
	To string
	Amount int
}

type urlDescription struct {
	URL url `json:"url"`
	Method string `json:"method"`
	Description string `json:"description"`
	Payload string `json:"payload,omitempty"`
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
			URL: url("/status"),
			Method : "GET",
			Description: "See Status",
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
		{
			URL : url("/balance/{address}"),
			Method: "GET",
			Description: "Get TxOuts for an Address",
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
		// var addBlockBody addBlockBody
		// utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.Blockchain().AddBlock()
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

func status(rw http.ResponseWriter, r *http.Request)  {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blockchain()))
}
func balacne(rw http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")
	switch total {
	case "true":
		amount := blockchain.Blockchain().BalanceByAddress(address)
		utils.HandleErr(json.NewEncoder(rw).Encode(balanceResponse{address, amount})) 
	default:
		utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Blockchain().TxOutsByAddress(address)))
	}
}
func mempool(rw http.ResponseWriter, r *http.Request)  {
	utils.HandleErr(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}
func transactions (rw http.ResponseWriter, r *http.Request)  {
	var payload addTxPayload
	utils.HandleErr(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err!=nil{
		json.NewEncoder(rw).Encode(errorResponse{"Not Enough Funds"})
	}
	rw.WriteHeader(http.StatusCreated)
}

func Start(aport int)  {
	router := mux.NewRouter()
	port = fmt.Sprintf(":%d", aport)
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/status" ,status)
	router.HandleFunc("/balance/{address}", balacne)
	router.HandleFunc("/mempool" ,mempool)
	router.HandleFunc("/transactions" ,transactions)
	fmt.Printf("Listening on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
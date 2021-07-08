package main

import (
	"fmt"

	"github.com/hojunin/hjcoin/blockchain"
)

func main(){
	chain := blockchain.GetBlockchain()
	chain.AddBlock("First")
	chain.AddBlock("Second")
	chain.AddBlock("Third")

	for _, block := range chain.AllBlock(){
		fmt.Printf("Data is %s\n", block.Data)
		fmt.Printf("Hash is %s\n", block.Hash)
		fmt.Printf("PrevHash is %s\n", block.PrevHash)
	}
}
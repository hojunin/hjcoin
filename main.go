package main

import (
	"github.com/hojunin/hjcoin/cli"
	"github.com/hojunin/hjcoin/db"
)



func main(){
	defer db.Close()
	cli.Start()
	
}
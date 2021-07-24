package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/hojunin/hjcoin/explorer"
	"github.com/hojunin/hjcoin/rest"
)


func usage()  {
	fmt.Println("Welcome To hjcoin\n")	
	fmt.Println("-port    ---->  Set the PORT of the server")
	fmt.Println("-mode    ---->  Choose between 'html' and 'rest'")
	runtime.Goexit()
}

func Start()  {
	if len(os.Args) ==1 {
		usage()
	}
	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")

	flag.Parse()


	switch *mode{
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}
}
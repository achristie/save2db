package main

import (
	"flag"
	"log"

	platts "github.com/achristie/save2db/pkg/platts"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// read cmd line arguments
	APIKey := flag.String("apikey", "NULL", "API Key to call API with")
	Username := flag.String("username", "NULL", "Username to get a token")
	Password := flag.String("password", "NULL", "Password associated with Username")
	flag.Parse()

	client := platts.NewClient(APIKey, Username, Password)

	MDCs, err := client.GetSubscribedMDC()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range MDCs {
		log.Printf("MDC: %s, Symbol Count: %d", v.MDC, v.SymbolCount)
	}

}

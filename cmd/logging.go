package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func initLogging(cmd *cobra.Command, args []string) {
	f, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// defer f.Close()
	log.SetOutput(f)
}

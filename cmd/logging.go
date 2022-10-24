package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func initLogging(cmd *cobra.Command, args []string) {
	_, err := os.OpenFile("log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	// config.infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// config.errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// defer f.Close()
	// config.infoLog.SetOutput(f)
	// config.errorLog.SetOutput(f)
}

package main

import (
	"log"
	"os"

	_ "github.com/garry-sharp/enclave-assessment/pkg/logger"

	"github.com/garry-sharp/enclave-assessment/pkg/logger"

	"github.com/garry-sharp/enclave-assessment/pkg/cli"
)

func main() {

	// Set multiwriter logger to file and stdout
	fn := os.Getenv("LOG_FILE")
	if fn == "" {
		fn = "app.log"
	}
	l, err := logger.NewFileAndStdOutLogger(fn)
	if err != nil {
		log.Fatalln(err)
	}
	logger.SetLogger(l)

	cmd, err := cli.LoadCLI()
	if err != nil {
		logger.Error("Failed to load CLI")
		logger.Error(err)
		os.Exit(1)
	}

	cmd.Execute()
}

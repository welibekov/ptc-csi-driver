package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func main() {
	// Define the flags
	endpoint := flag.String("endpoint", "", "The endpoint to connect to (must start with 'unix://')")
	logLevel := flag.String("log-level", "info", "The logging level (debug, info, warning, error, fatal, panic)")

	// Parse the flags
	flag.Parse()

	// Validate the endpoint
	if !strings.HasPrefix(*endpoint, "unix://") {
		fmt.Fprintf(os.Stderr, "Error: endpoint must start with 'unix://'\n")
		flag.Usage()
		os.Exit(1)
	}

	// Set up logrus logging level
	level, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid log level '%s'\n", *logLevel)
		flag.Usage()
		os.Exit(1)
	}
	logrus.SetLevel(level)

	// Log the endpoint and log level
	logrus.Infof("Connecting to endpoint: %s", *endpoint)
	logrus.Infof("Log level set to: %s", *logLevel)

	// Placeholder for the main functionality
	logrus.Warn("Main functionality not yet implemented")
}

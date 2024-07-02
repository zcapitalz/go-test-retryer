package main

import (
	"log"
	"os"

	rt "github.com/zcapitalz/go-test-retryer"
)

func main() {
	errorLogger := log.New(os.Stderr, "", 0)

	cfg, err := rt.NewConfigFromArgs(os.Args)
	if err != nil {
		errorLogger.Println(err)
		os.Exit(2)
	}

	retryer := rt.NewRetryer(cfg, os.Stdout, os.Stderr)

	err = retryer.Run()
	if _, isTestError := err.(rt.TestError); err != nil && !isTestError {
		errorLogger.Println(err)
	}
	switch err := err.(type) {
	case nil:
		os.Exit(0)
	case rt.TestError:
		os.Exit(err.TestExitCode())
	default:
		os.Exit(1)
	}
}

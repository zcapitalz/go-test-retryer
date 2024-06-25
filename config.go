package retryer

import "flag"

type Config struct {
	testOutputTypeJSON bool
	maxRetriesPerTest  int
	maxTotalRetries    int
	testCommandName    string
	testArgs           string
	verbose            bool
	shellPath          string
}

func NewConfigFromArgs(args []string) (Config, error) {
	cfg := Config{}

	flag.BoolVar(&cfg.testOutputTypeJSON, "json", false, "parse go test output as json")
	flag.IntVar(&cfg.maxRetriesPerTest, "retries-per-test", 0, "maximum retries per test")
	flag.IntVar(&cfg.maxTotalRetries, "total-retries", 0, "maximum retries for all tests")
	flag.BoolVar(&cfg.verbose, "verbose", false, "verbose mode")
	flag.StringVar(&cfg.testCommandName, "test-command-name", "go test", `test command name`)
	flag.StringVar(&cfg.testArgs, "test-args", "", "test arguments")
	flag.StringVar(&cfg.shellPath, "shell", "/bin/bash", "path to shell")
	flag.Parse()

	if cfg.maxTotalRetries < 0 || cfg.maxRetriesPerTest < 0 {
		return Config{}, InvalidParameterError{"Retries amount should be non-negative"}
	}

	return cfg, nil
}

func (cfg *Config) isTotalRetriesLimitEnabled() bool {
	return cfg.maxTotalRetries != 0
}

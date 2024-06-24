package test

import (
	"flag"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

var (
	configPath string
	config     *Config
)

const (
	logMessage = "working..."
)

type Config struct {
	FlakyTestFailuresLeft int `yaml:"flaky_test_failures_left"`
}

func init() {
	flag.StringVar(&configPath, "config-path", "", "path to config")
}

func TestMain(m *testing.M) {
	readConfig()
	os.Exit(m.Run())
}

func TestSuccess(t *testing.T) {
	t.Log(logMessage)
}

func TestFail(t *testing.T) {
	t.Log(logMessage)
	t.FailNow()
}

func TestFlaky(t *testing.T) {
	require.NotNil(t, config, "no config")
	t.Log(logMessage)
	if config.FlakyTestFailuresLeft != 0 {
		config.FlakyTestFailuresLeft--
		err := writeConfig()
		if err != nil {
			t.Logf("Log file write error: %v", err)
		}
		t.FailNow()
	}
}

func readConfig() {
	flag.Parse()

	if configPath == "" {
		return
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Read config file error: %v", err)
	}

	config = new(Config)
	err = yaml.Unmarshal(configBytes, config)
	if err != nil {
		log.Fatalf("Parse config file error: %v", err)
	}
}

func writeConfig() error {
	configBytes, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, configBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

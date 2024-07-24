package retryer

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	"github.com/jstemmer/go-junit-report/v2/gtr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var debug bool

func init() {
	flag.BoolVar(&debug, "debug", false, "print debug info")
}

func Test(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, testFromTestCase(tc))

		tc.name += "JSON"
		tc.cfg.testOutputTypeJSON = true
		tc.cfg.testArgs += " -json"
		for i, command := range tc.expectedCommands {
			tc.expectedCommands[i] = command + " -json"
		}
		t.Run(tc.name, testFromTestCase(tc))
	}
}

func testFromTestCase(tc testCase) func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		expectedStdoutStr, expectedStderrStr := getTestCaseExpectedOutputStr(t, tc)

		if tc.cfg.testOutputTypeJSON {
			debugLogf(t, "Expected output:\n%v\n%v\n", expectedStdoutStr, expectedStderrStr)
		}

		if tc.testConfig != "" {
			testConfigPath := createTestConfigFile(t, tc.testConfig)
			defer os.Remove(testConfigPath)
			tc.cfg.testArgs += " -config-path=" + testConfigPath
		}

		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		output := new(buffer)
		retryerArgs := fmt.Sprintf(`-json=%v -total-retries=%v -retries-per-test=%v -test-command-name="%v" -verbose=%v`,
			tc.cfg.testOutputTypeJSON, tc.cfg.maxTotalRetries, tc.cfg.maxRetriesPerTest, tc.cfg.testCommandName, tc.cfg.verbose)
		command := fmt.Sprintf(`go run ./cmd/go-test-retryer/main.go %v -test-args="%v"`, retryerArgs, escapeQuotes(tc.cfg.testArgs))
		debugLogf(t, "Command:\n%v\n", command)
		exitCode, err := runCommand(
			fmt.Sprintf(`go run ./cmd/go-test-retryer/main.go %v -test-args="%v"`, retryerArgs, escapeQuotes(tc.cfg.testArgs)),
			io.MultiWriter(stdout, output),
			io.MultiWriter(stderr, output))
		if tc.cfg.testOutputTypeJSON {
			debugLogf(t, "Actual output:\n%v\n", output.String())
		}
		if err != nil {
			_, ok := err.(*exec.ExitError)
			require.True(t, ok, fmt.Sprintf("command execution error: %v", err))
		}
		assert.Equal(t, tc.expectedExitCode, exitCode)

		checkStdout(t, strings.NewReader(expectedStdoutStr), stdout, tc.cfg.testOutputTypeJSON)
		checkStderr(t, strings.NewReader(expectedStderrStr), stderr)
	}
}

func getTestCaseExpectedOutputStr(t *testing.T, tc testCase) (string, string) {
	testConfigPath := ""
	if tc.testConfig != "" {
		testConfigPath = createTestConfigFile(t, tc.testConfig)
		defer os.Remove(testConfigPath)
	}

	expectedStdout := new(bytes.Buffer)
	expectedStderr := new(bytes.Buffer)
	for _, command := range tc.expectedCommands {
		if tc.testConfig != "" {
			command += " -config-path=" + testConfigPath
		}

		_, err := runCommand(command, expectedStdout, expectedStderr)
		if err != nil {
			_, ok := err.(*exec.ExitError)
			require.True(t, ok, fmt.Sprintf("unexpected exec error: %v", err))
		}
	}

	return readerToString(t, expectedStdout), readerToString(t, expectedStderr)
}

func runCommand(command string, stdout, stderr io.Writer) (int, error) {
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err != nil {
		return cmd.ProcessState.ExitCode(), err
	}

	return cmd.ProcessState.ExitCode(), nil
}

func checkStdout(t *testing.T, expectedStdout, actualStdout io.Reader, outputTypeJSON bool) {
	if outputTypeJSON {
		actualReport, err := parseTestReport(actualStdout, outputTypeJSON)
		require.NoError(t, err)
		expectedReport, err := parseTestReport(expectedStdout, outputTypeJSON)
		require.NoError(t, err)

		requireEqualReports(t, expectedReport, actualReport)
	}
}

func checkStderr(t *testing.T, expectedStderr, actualStderr io.Reader) {
	actualStderrStr := readerToString(t, actualStderr)
	expectedStderrStr := readerToString(t, expectedStderr)
	r := regexp.MustCompile(`exit status 1\n?\z`)
	actualStderrStr = r.ReplaceAllString(actualStderrStr, "")
	assert.Equal(t, expectedStderrStr, actualStderrStr)
}

func requireEqualReports(t *testing.T, expectedReport, actualReport gtr.Report) {
	require.Equal(t, len(expectedReport.Packages), len(actualReport.Packages))
	for i := 0; i < len(expectedReport.Packages); i++ {
		expectedPackage := expectedReport.Packages[i]
		actualPackage := actualReport.Packages[i]
		requireEqualReportPackages(t, expectedPackage, actualPackage)
	}
}

func requireEqualReportPackages(t *testing.T, expectedPackage, actualPackage gtr.Package) {
	require.Equal(t, len(expectedPackage.Tests), len(actualPackage.Tests))
	for i := 0; i < len(expectedPackage.Tests); i++ {
		expectedTest := expectedPackage.Tests[i]
		actualTest := actualPackage.Tests[i]
		require.Equal(t, expectedTest.Result, actualTest.Result)
		for i, expectedOutputStr := range expectedTest.Output {
			require.Greater(t, len(actualTest.Output), i)
			require.Equal(t, expectedOutputStr, actualTest.Output[i])
		}
	}
}

func createTestConfigFile(t *testing.T, config string) (configPath string) {
	f, err := os.CreateTemp("", "test_config_*.yaml")
	require.NoError(t, err)
	_, err = f.WriteString(config)
	require.NoError(t, err)
	f.Close()
	return f.Name()
}

func readerToString(t *testing.T, reader io.Reader) string {
	builder := new(strings.Builder)
	_, err := io.Copy(builder, reader)
	require.NoError(t, err)
	return builder.String()
}

func escapeQuotes(str string) string {
	return regexp.MustCompile(`"`).ReplaceAllString(str, `\"`)
}

func debugLogf(t *testing.T, format string, args ...any) {
	if debug {
		t.Logf(format, args...)
	}
}

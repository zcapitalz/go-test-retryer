package retryer

//nolint:unused
type testCase struct {
	name             string
	retryerCfg       Config
	testCfg          string
	expectedExitCode int
	expectedCommands []string
}

// Add only test cases for plain mode (-json=false).
// Tests cases for json mode are generated automatically.
//
//nolint:unused
var testCases = []testCase{
	{
		name: "SuccessfulTestRetriesNotAllowed",
		retryerCfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  0,
			maxTotalRetries:    0,
			testCommandName:    "go test",
			testArgs:           "-v -run=^TestSuccess$ -count=1 github.com/zcapitalz/go-test-retryer/test",
			verbose:            false,
			shellPath:          "/bin/bash",
		},
		expectedExitCode: 0,
		expectedCommands: []string{
			"go test -v -run=^TestSuccess$ -count=1 github.com/zcapitalz/go-test-retryer/test",
		},
	},
	{
		name: "SuccessfulTestRetriesAllowed",
		retryerCfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  1,
			maxTotalRetries:    1,
			testCommandName:    "go test",
			testArgs:           "-v -run=^TestSuccess$ -count=1 github.com/zcapitalz/go-test-retryer/test",
			verbose:            false,
			shellPath:          "/bin/bash",
		},
		expectedExitCode: 0,
		expectedCommands: []string{
			"go test -v -run=^TestSuccess$ -count=1 github.com/zcapitalz/go-test-retryer/test",
		},
	},
	{
		name: "FailedTestRetriesNotAllowed1",
		retryerCfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  0,
			maxTotalRetries:    0,
			testCommandName:    "go test",
			testArgs:           "-v -run=^TestFail$ -count=1 github.com/zcapitalz/go-test-retryer/test",
			verbose:            false,
			shellPath:          "/bin/bash",
		},
		expectedExitCode: 1,
		expectedCommands: []string{
			"go test -v -run=^TestFail$ -count=1 github.com/zcapitalz/go-test-retryer/test",
		},
	},
	{
		name: "FailedTestRetriesNotAllowed2",
		retryerCfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  0,
			maxTotalRetries:    1,
			testCommandName:    "go test",
			testArgs:           "-v -run=^TestFail$ -count=1 github.com/zcapitalz/go-test-retryer/test",
			verbose:            false,
			shellPath:          "/bin/bash",
		},
		expectedExitCode: 1,
		expectedCommands: []string{
			"go test -v -run=^TestFail$ -count=1 github.com/zcapitalz/go-test-retryer/test",
		},
	},
	{
		name: "FailedTestRetriesAllowed1",
		retryerCfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  1,
			maxTotalRetries:    1,
			testCommandName:    "go test",
			testArgs:           "-v -run=^TestFail$ -count=1 github.com/zcapitalz/go-test-retryer/test",
			verbose:            false,
			shellPath:          "/bin/bash",
		},
		expectedExitCode: 1,
		expectedCommands: []string{
			"go test -v -run=^TestFail$ -count=1 github.com/zcapitalz/go-test-retryer/test",
			"go test -v -run=^TestFail$ -count=1 github.com/zcapitalz/go-test-retryer/test",
		},
	},
	{
		name: "FailedTestRetriesAllowed2",
		retryerCfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  1,
			maxTotalRetries:    0,
			testCommandName:    "go test",
			testArgs:           "-v -run=^TestFail$ -count=1 github.com/zcapitalz/go-test-retryer/test",
			verbose:            false,
			shellPath:          "/bin/bash",
		},
		expectedExitCode: 1,
		expectedCommands: []string{
			"go test -v -run=^TestFail$ -count=1 github.com/zcapitalz/go-test-retryer/test",
			"go test -v -run=^TestFail$ -count=1 github.com/zcapitalz/go-test-retryer/test",
		},
	},
	{
		name: "NotCompilableRetriesNotAllowed",
		retryerCfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  0,
			maxTotalRetries:    0,
			testCommandName:    "go test",
			testArgs:           "-v -count=1 github.com/zcapitalz/go-test-retryer/test/notcompilable",
			verbose:            false,
			shellPath:          "/bin/bash",
		},
		expectedExitCode: 1,
		expectedCommands: []string{
			"go test -v -count=1 github.com/zcapitalz/go-test-retryer/test/notcompilable",
		},
	},
	{
		name: "NotCompilableRetriesAllowed",
		retryerCfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  1,
			maxTotalRetries:    1,
			testCommandName:    "go test",
			testArgs:           "-v -count=1 github.com/zcapitalz/go-test-retryer/test/notcompilable",
			verbose:            false,
			shellPath:          "/bin/bash",
		},
		expectedExitCode: 1,
		expectedCommands: []string{
			"go test -v -count=1 github.com/zcapitalz/go-test-retryer/test/notcompilable",
		},
	},
	{
		name: "FlakyTestRetriesAllowed",
		retryerCfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  2,
			maxTotalRetries:    2,
			testCommandName:    "go test",
			testArgs:           "-v -count=1 -run=^TestFlaky$ github.com/zcapitalz/go-test-retryer/test",
			verbose:            false,
			shellPath:          "/bin/bash",
		},
		testCfg:          "flaky_test_failures_left: 2",
		expectedExitCode: 0,
		expectedCommands: []string{
			"go test -v -count=1 -run=^TestFlaky$ github.com/zcapitalz/go-test-retryer/test",
			"go test -v -count=1 -run=^TestFlaky$ github.com/zcapitalz/go-test-retryer/test",
			"go test -v -count=1 -run=^TestFlaky$ github.com/zcapitalz/go-test-retryer/test",
		},
	},
	{
		name: "FlakyTestNotEnoughRetries",
		retryerCfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  1,
			maxTotalRetries:    1,
			testCommandName:    "go test",
			testArgs:           "-v -count=1 -run=^TestFlaky$ github.com/zcapitalz/go-test-retryer/test",
			verbose:            false,
			shellPath:          "/bin/bash",
		},
		testCfg:          "flaky_test_failures_left: 2",
		expectedExitCode: 1,
		expectedCommands: []string{
			"go test -v -count=1 -run=^TestFlaky$ github.com/zcapitalz/go-test-retryer/test",
			"go test -v -count=1 -run=^TestFlaky$ github.com/zcapitalz/go-test-retryer/test",
		},
	},
	{
		name: "FlakyAndFailedTest",
		retryerCfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  1,
			maxTotalRetries:    2,
			testCommandName:    "go test",
			testArgs:           `-v -count=1 -run="^(TestFlaky|TestFail)$" github.com/zcapitalz/go-test-retryer/test`,
			verbose:            false,
			shellPath:          "/bin/bash",
		},
		testCfg:          "flaky_test_failures_left: 1",
		expectedExitCode: 1,
		expectedCommands: []string{
			`go test -v -count=1 -run="^(TestFlaky|TestFail)" github.com/zcapitalz/go-test-retryer/test`,
			`go test -v -count=1 -run="^(TestFlaky|TestFail)" github.com/zcapitalz/go-test-retryer/test`,
		},
	},
}

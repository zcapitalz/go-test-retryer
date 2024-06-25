package retryer

//nolint:unused
type testCase struct {
	name             string
	cfg              Config
	testConfig       string
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
		cfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  0,
			maxTotalRetries:    0,
			testCommandName:    "go test",
			testArgs:           "-v -run=^TestSuccess$ -count=1 go-test-retryer/test",
			verbose:            false,
		},
		expectedExitCode: 0,
		expectedCommands: []string{
			"go test -v -run=^TestSuccess$ -count=1 go-test-retryer/test",
		},
	},
	{
		name: "SuccessfulTestRetriesAllowed",
		cfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  1,
			maxTotalRetries:    1,
			testCommandName:    "go test",
			testArgs:           "-v -run=^TestSuccess$ -count=1 go-test-retryer/test",
			verbose:            false,
		},
		expectedExitCode: 0,
		expectedCommands: []string{
			"go test -v -run=^TestSuccess$ -count=1 go-test-retryer/test",
		},
	},
	{
		name: "FailedTestRetriesNotAllowed1",
		cfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  0,
			maxTotalRetries:    0,
			testCommandName:    "go test",
			testArgs:           "-v -run=^TestFail$ -count=1 go-test-retryer/test",
			verbose:            false,
		},
		expectedExitCode: 1,
		expectedCommands: []string{
			"go test -v -run=^TestFail$ -count=1 go-test-retryer/test",
		},
	},
	{
		name: "FailedTestRetriesNotAllowed2",
		cfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  0,
			maxTotalRetries:    1,
			testCommandName:    "go test",
			testArgs:           "-v -run=^TestFail$ -count=1 go-test-retryer/test",
			verbose:            false,
		},
		expectedExitCode: 1,
		expectedCommands: []string{
			"go test -v -run=^TestFail$ -count=1 go-test-retryer/test",
		},
	},
	{
		name: "FailedTestRetriesAllowed1",
		cfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  1,
			maxTotalRetries:    1,
			testCommandName:    "go test",
			testArgs:           "-v -run=^TestFail$ -count=1 go-test-retryer/test",
			verbose:            false,
		},
		expectedExitCode: 1,
		expectedCommands: []string{
			"go test -v -run=^TestFail$ -count=1 go-test-retryer/test",
			"go test -v -run=^TestFail$ -count=1 go-test-retryer/test",
		},
	},
	{
		name: "FailedTestRetriesAllowed2",
		cfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  1,
			maxTotalRetries:    0,
			testCommandName:    "go test",
			testArgs:           "-v -run=^TestFail$ -count=1 go-test-retryer/test",
			verbose:            false,
		},
		expectedExitCode: 1,
		expectedCommands: []string{
			"go test -v -run=^TestFail$ -count=1 go-test-retryer/test",
			"go test -v -run=^TestFail$ -count=1 go-test-retryer/test",
		},
	},
	{
		name: "NotCompilableRetriesNotAllowed",
		cfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  0,
			maxTotalRetries:    0,
			testCommandName:    "go test",
			testArgs:           "-v -count=1 go-test-retryer/test/notcompilable",
			verbose:            false,
		},
		expectedExitCode: 1,
		expectedCommands: []string{
			"go test -v -count=1 go-test-retryer/test/notcompilable",
		},
	},
	{
		name: "NotCompilableRetriesAllowed",
		cfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  1,
			maxTotalRetries:    1,
			testCommandName:    "go test",
			testArgs:           "-v -count=1 go-test-retryer/test/notcompilable",
			verbose:            false,
		},
		expectedExitCode: 1,
		expectedCommands: []string{
			"go test -v -count=1 go-test-retryer/test/notcompilable",
		},
	},
	{
		name: "FlakyTestRetriesAllowed",
		cfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  2,
			maxTotalRetries:    2,
			testCommandName:    "go test",
			testArgs:           "-v -count=1 -run=^TestFlaky$ go-test-retryer/test",
			verbose:            false,
		},
		testConfig:       "flaky_test_failures_left: 2",
		expectedExitCode: 0,
		expectedCommands: []string{
			"go test -v -count=1 -run=^TestFlaky$ go-test-retryer/test",
			"go test -v -count=1 -run=^TestFlaky$ go-test-retryer/test",
			"go test -v -count=1 -run=^TestFlaky$ go-test-retryer/test",
		},
	},
	{
		name: "FlakyTestNotEnoughRetries",
		cfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  1,
			maxTotalRetries:    1,
			testCommandName:    "go test",
			testArgs:           "-v -count=1 -run=^TestFlaky$ go-test-retryer/test",
			verbose:            false,
		},
		testConfig:       "flaky_test_failures_left: 2",
		expectedExitCode: 1,
		expectedCommands: []string{
			"go test -v -count=1 -run=^TestFlaky$ go-test-retryer/test",
			"go test -v -count=1 -run=^TestFlaky$ go-test-retryer/test",
		},
	},
	{
		name: "FlakyAndFailedTest",
		cfg: Config{
			testOutputTypeJSON: false,
			maxRetriesPerTest:  1,
			maxTotalRetries:    2,
			testCommandName:    "go test",
			testArgs:           `-v -count=1 -run="^(TestFlaky|TestFail)$" go-test-retryer/test`,
			verbose:            false,
		},
		testConfig:       "flaky_test_failures_left: 1",
		expectedExitCode: 1,
		expectedCommands: []string{
			`go test -v -count=1 -run="^(TestFlaky|TestFail)" go-test-retryer/test`,
			`go test -v -count=1 -run="^(TestFlaky|TestFail)" go-test-retryer/test`,
		},
	},
}

package retryer

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/jstemmer/go-junit-report/v2/gtr"
	"github.com/jstemmer/go-junit-report/v2/parser/gotest"
	"github.com/pkg/errors"
)

type Retryer struct {
	cfg                    Config
	stdout                 io.Writer
	stderr                 io.Writer
	totalRetriesLeft       int
	totalSuccessfulRetries int
	totalRetriesPerTest    map[string]int
	everFailedTests        map[string]struct{}
	lastFailedTests        []string
	lastTestExitCode       int
	firstRun               bool
	failedAnyPackageBuild  bool
}

func NewRetryer(cfg Config, stdout, stderr io.Writer) *Retryer {
	return &Retryer{
		cfg:                    cfg,
		stdout:                 stdout,
		stderr:                 stderr,
		totalRetriesLeft:       cfg.maxTotalRetries,
		totalRetriesPerTest:    make(map[string]int),
		everFailedTests:        make(map[string]struct{}),
		totalSuccessfulRetries: 0,
		lastTestExitCode:       -1,
		firstRun:               true,
	}
}

func (r *Retryer) Run() (err error) {
	if r.cfg.maxRetriesPerTest == 0 {
		r.log("No retries allowed, going to run tests and exit")
		err := r.test(r.cfg.testArgs, r.stdout, r.stderr)
		if exitError, ok := err.(*exec.ExitError); ok {
			return TestError{exitCode: exitError.ExitCode()}
		}
		return err
	}

	r.log("Initial run of tests")
	err = r.testAndUpdateState(r.cfg.testArgs)
	if err != nil {
		return err
	}

	for len(r.lastFailedTests) > 0 {
		testsToRetry := r.selectTestsForRetry()
		if len(testsToRetry) == 0 {
			break
		}

		testRunArgParts := make([]string, 0, len(testsToRetry))
		for _, t := range testsToRetry {
			testRunArgParts = append(testRunArgParts, "("+t+")")
		}

		r.log("Retrying tests")
		testArgs := r.cfg.testArgs + ` "--test.run=^(` + strings.Join(testRunArgParts, "|") + `)$"`
		err := r.testAndUpdateState(testArgs)
		if err != nil {
			return err
		}
	}

	totalRetries := r.cfg.maxTotalRetries - r.totalRetriesLeft
	r.log("Total retries:", totalRetries)
	if totalRetries > 0 {
		successfultRetriesPercentage := float64(r.totalSuccessfulRetries) / float64(totalRetries) * 100
		r.logf("Successful retries: %v (%.0f%%)\n", r.totalSuccessfulRetries, successfultRetriesPercentage)
		r.log("Retries per test:", r.totalRetriesPerTest)
	}

	if len(r.everFailedTests) != r.totalSuccessfulRetries {
		return TestError{exitCode: r.lastTestExitCode}
	}
	if r.failedAnyPackageBuild {
		return TestError{exitCode: 1}
	}

	return nil
}

func (r *Retryer) testAndUpdateState(testArgs string) error {
	outputBuffer := new(buffer)

	err := r.test(
		testArgs,
		io.MultiWriter(r.stdout, outputBuffer),
		io.MultiWriter(r.stderr, outputBuffer))

	if exitError, ok := err.(*exec.ExitError); err != nil && ok {
		r.lastTestExitCode = exitError.ExitCode()
	} else if err != nil && !ok {
		r.lastFailedTests = nil
		r.lastTestExitCode = -1
		return errors.Wrap(err, "run tests")
	}

	report, err := parseTestReport(outputBuffer, r.cfg.testOutputTypeJSON)
	if err != nil {
		return errors.Wrap(err, "parse test output")
	}
	r.updateStateWithTestReport(report)

	return nil
}

func (r *Retryer) test(testArgs string, stdout, stderr io.Writer) error {
	command := exec.Command("/bin/bash", "-c", r.cfg.testCommandName+" "+testArgs)
	r.log("Running command:", strings.Join(command.Args, " "))
	command.Stdout = stdout
	command.Stderr = stderr
	return command.Run()
}

func (r *Retryer) selectTestsForRetry() (testsToRetry []string) {
	for i := 0; i < len(r.lastFailedTests); i++ {
		if r.cfg.isTotalRetriesLimitEnabled() && r.totalRetriesLeft == 0 {
			break
		}

		failedTest := r.lastFailedTests[i]
		if retries := r.totalRetriesPerTest[failedTest]; retries < r.cfg.maxRetriesPerTest {
			testsToRetry = append(testsToRetry, failedTest)
			r.totalRetriesPerTest[failedTest] = retries + 1
			r.totalRetriesLeft--
		}
	}
	return testsToRetry
}

func (r *Retryer) updateStateWithTestReport(report gtr.Report) {
	tests := testsFromReport(report)
	tests = filter(tests, isRootTest)

	r.lastFailedTests = testNamesFromTests(filter(tests, isFailedTest))
	r.log("Failed tests:", r.lastFailedTests)

	for _, failedTest := range r.lastFailedTests {
		r.everFailedTests[failedTest] = struct{}{}
	}

	if !r.firstRun {
		passedTests := filter(tests, isPassedTest)
		r.totalSuccessfulRetries += len(passedTests)
		r.log("Passed retries:", testNamesFromTests(passedTests))
	} else {
		r.firstRun = false
	}

	r.failedAnyPackageBuild = r.failedAnyPackageBuild || anyBuildErrorsInReport(report)
}

func (r *Retryer) log(args ...any) {
	if r.cfg.verbose {
		fmt.Fprintln(r.stdout, args...)
	}
}

func (r *Retryer) logf(format string, args ...any) {
	if r.cfg.verbose {
		fmt.Fprintf(r.stdout, format, args...)
	}
}

func parseTestReport(r io.Reader, testOutputTypeJSON bool) (gtr.Report, error) {
	if testOutputTypeJSON {
		return gotest.NewJSONParser().Parse(r)
	}
	return gotest.NewParser().Parse(r)
}

func testsFromReport(report gtr.Report) []gtr.Test {
	tests := make([]gtr.Test, 0)
	for _, pkg := range report.Packages {
		tests = append(tests, pkg.Tests...)
	}
	return tests
}

func anyBuildErrorsInReport(report gtr.Report) bool {
	for _, pkg := range report.Packages {
		if len(pkg.BuildError.Output) > 0 {
			return true
		}
	}

	return false
}

func testNamesFromTests(tests []gtr.Test) []string {
	testNames := make([]string, 0)
	for _, test := range tests {
		testNames = append(testNames, test.Name)
	}
	return testNames
}

func isFailedTest(test gtr.Test) bool {
	return test.Result == gtr.Fail
}

func isPassedTest(test gtr.Test) bool {
	return test.Result == gtr.Pass
}

func isRootTest(test gtr.Test) bool {
	return test.Level == 0
}

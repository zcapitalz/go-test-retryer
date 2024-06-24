package retryer

import "fmt"

type TestError struct {
	exitCode int
}

func (err TestError) Error() string     { return fmt.Sprintf("exit code: %v", int(err.exitCode)) }
func (err TestError) TestExitCode() int { return int(err.exitCode) }

type InvalidParameterError struct {
	message string
}

func (err InvalidParameterError) Error() string { return err.message }

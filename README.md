# Go Test Retryer

CLI tool that runs `go test`(or another command from `--test-command-name`) with arguments from `--test-args`, parses test results from output and retries failed tests according to `--retries-per-test` and `--total-reries` limits. If `--total-retries` is 0, then no global limit is applied.

Testing commands are run using provided `--shell`(default "/bin/bash") with -c option.
<br><br>

**Installation**:
```
go install github.com/zcapitalz/go-test-retryer/cmd/go-test-retryer@latest
```
<br>

**Usage**:
- --test-args string  
&emsp;&emsp;test arguments  
- --test-command-name string  
&emsp;&emsp;test command name (default "go test")  
- --retries-per-test int  
&emsp;&emsp;maximum retries per test  
- --total-retries int  
&emsp;&emsp;maximum retries for all tests  
- --json bool  
&emsp;&emsp;parse go test output as json  
- --verbose bool  
&emsp;&emsp;verbose mode
- --shell string  
&emsp;&emsp;path to shell (default "/bin/bash")  
<br>

**Exit codes**:
- `0`: if all tests that failed during any run are successfuly retried
- `1`: for unexpected error
- `*`: whatever was returned by last run otherwise
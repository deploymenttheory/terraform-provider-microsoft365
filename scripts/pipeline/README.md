# Pipeline Scripts

This directory contains scripts used in the GitHub Actions CI/CD pipeline for the terraform-provider-microsoft365 project.

## 📜 Scripts Overview

### ✅ Active Scripts

#### `run-tests.py`
**Purpose:** Runs acceptance tests and captures test failures  
**Usage:** `./run-tests.py <type> [service] [coverage-file] [test-output-file]`

**Parameters:**
- `type`: Type of tests to run (`provider-core`, `resources`, `datasources`)
- `service`: Service name (required for resources/datasources tests)
- `coverage-file`: Output file for coverage data (default: `coverage.txt`)
- `test-output-file`: Output file for test logs (default: `test-output.log`)

**Features:**
- Runs Go tests with race detection
- Captures test output for failure analysis
- Generates JSON report of failing tests (`test-failures.json`)
- Allows tests to fail gracefully (uses `continue-on-error`)
- Includes test context (error messages, stack traces)

**Output:**
- `test-failures.json`: Structured JSON with failing test details
- `test-output-*.log`: Raw test output logs
- `coverage-*.txt`: Coverage profile

#### `manage-test-issues.py`
**Purpose:** Manages GitHub issues for test failures (create, update, close)  
**Usage:** `./manage-test-issues.py <owner> <repo> <run-id> <failures-json> [successes-json]`

**Parameters:**
- `owner`: GitHub repository owner
- `repo`: Repository name
- `run-id`: Workflow run ID
- `failures-json`: Path to test failures JSON file
- `successes-json`: (Optional) Path to test successes JSON file

**Features:**
- **Full Issue Lifecycle**: Creates, updates, and closes issues automatically
- **De-duplication**: Detects and updates existing issues
- **Factual Reporting**: Contains only actionable information (test name, error, service, date, workflow link)
- **Recurring Detection**: Adds `recurring` label for repeated failures
- **Auto-close Resolved**: Closes issues when tests pass (if successes provided)
- **Automatic Labels**: Tags with `test-failure`, `automated`, `recurring`

**Issue Contents:**
- Test name (as title)
- Service area
- Failure date and timestamp
- Workflow run ID and URL
- Error output
- Links to test source and logs

**Behavior:**
- **New Failure**: Creates issue with test name as title
- **Existing Failure**: Adds comment with timestamp and latest error
- **Resolved Test**: Automatically closes issue when test passes

#### `map-credentials.py`
**Purpose:** Maps service-specific credentials to environment variables  
**Usage:** `./map-credentials.py <service>`

**Parameters:**
- `service`: Service name (e.g., `device_and_app_management`, `groups`)

**Features:**
- Maps service-specific `M365_CLIENT_ID_*` and `M365_CLIENT_SECRET_*` to generic `M365_CLIENT_ID` and `M365_CLIENT_SECRET`
- Sets `SKIP_TESTS=true` if credentials are not configured
- Enables per-service credential management

#### `merge-test-results.py`
**Purpose:** Merges multiple test result JSON files into a single file  
**Usage:** `./merge-test-results.py <artifacts-dir> <output-file> <filename-to-merge>`

**Parameters:**
- `artifacts-dir`: Directory containing downloaded artifacts
- `output-file`: Output merged JSON file
- `filename-to-merge`: Name of files to merge (e.g., `test-failures.json`, `test-successes.json`)

**Features:**
- Recursively finds all matching JSON files in artifacts directory
- Merges JSON arrays into single consolidated file
- Handles empty files gracefully
- Used for both failures and successes

#### `detect-job-failures.py`
**Purpose:** Detects job-level failures using GitHub API  
**Usage:** `./detect-job-failures.py <owner> <repo> <run-id> [output-file]`

**Parameters:**
- `owner`: GitHub repository owner
- `repo`: Repository name
- `run-id`: Workflow run ID
- `output-file`: Output JSON file (default: `job-failures.json`)

**Features:**
- **Comprehensive Detection**: Identifies timeouts, OOM errors, runner failures, infrastructure issues
- **API-Based**: Uses GitHub API to query job statuses and steps
- **Intelligent Filtering**: Distinguishes job-level failures from expected test failures
- **Detailed Context**: Captures failed step, job duration, runner info

**Detection Types:**
- `timeout`: Job exceeded maximum execution time
- `out_of_memory`: OOM kill detected
- `runner_failure`: GitHub Actions runner/infrastructure issue
- `infrastructure`: Setup or dependency failure
- `step_failure`: Non-test step failure
- `cancelled`: Job was cancelled

#### `manage-job-failure-issues.py`
**Purpose:** Manages GitHub issues for job-level failures  
**Usage:** `./manage-job-failure-issues.py <owner> <repo> <run-id> [job-failures-json]`

**Parameters:**
- `owner`: GitHub repository owner
- `repo`: Repository name
- `run-id`: Workflow run ID
- `job-failures-json`: Path to job failures JSON file (default: `job-failures.json`)

**Features:**
- **Automatic Issue Creation**: Creates issues for infrastructure failures
- **Severity Indicators**: Uses emojis to indicate failure type (⏱️ timeout, 💥 OOM, 🚨 infrastructure)
- **Troubleshooting Hints**: Includes specific guidance based on failure type
- **De-duplication**: Updates existing issues for recurring failures
- **Labels**: `job-failure`, `infrastructure`, `automated`, `recurring`

**Issue Contents:**
- Job name and failure type
- Failed step details
- Job ID, runner info, timestamps
- Direct link to job logs
- Troubleshooting guidance specific to failure type

## 🔄 Workflow Integration

### Nightly Test Workflow Flow

```mermaid
graph TD
    A[Run Tests] -->|Success or Failure| B[Parse Test Results]
    B --> C[Upload test-failures.json]
    B --> D[Upload test-successes.json]
    C --> E[Summary Job: Download Artifacts]
    D --> E
    E --> F[Merge Test Results]
    F --> G[Manage Test Issues]
    G --> H[Detect Job-Level Failures via API]
    H --> I{Job Failures?}
    I -->|Yes| J[Create/Update Job Failure Issues]
    I -->|No| K[Complete]
    
    G --> L{Test Failures?}
    L -->|Yes| M[Create/Update Test Issues]
    L -->|No| N[Close Resolved Issues]
    M --> H
    N --> H
```

### Test Failure Report Lifecycle

1. **Detection**: Test fails in nightly run
2. **Artifact Upload**: `test-failures.json` and `test-successes.json` uploaded as artifacts
3. **Aggregation**: Summary job downloads all artifacts
4. **De-duplication**: Script checks for existing issues by test name
5. **Issue Management**:
   - **If test fails**: Create new issue or update existing with `recurring` label
   - **If test passes**: Automatically close previously opened issue
6. **Resolution**: Issues auto-close when tests pass

### Job Failure Detection Lifecycle

1. **API Query**: Summary job queries GitHub API for all job statuses
2. **Analysis**: Each job analyzed for:
   - Job-level failures (timeout, OOM, cancelled)
   - Infrastructure failures (runner issues, setup failures)
   - Non-test step failures (dependency installation, checkout)
3. **Filtering**: Distinguishes job failures from expected test failures
4. **Issue Creation**:
   - **If new job failure**: Create issue with failure type emoji and details
   - **If recurring**: Update existing issue and add `recurring` label
5. **Context**: Include failed step, job logs, troubleshooting guidance

### Example Test Failure Issue

**Title:** `TestAccAndroidPolicyResource_Lifecycle`

**Body:**
```markdown
## Test Failure

**Test:** `TestAccAndroidPolicyResource_Lifecycle`  
**Service:** `resources/device_and_app_management`  
**Date:** 2025-11-15  
**Workflow:** [19383092062](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/19383092062)

### Error Output

```
--- FAIL: TestAccAndroidPolicyResource_Lifecycle (5.23s)
    resource_test.go:45: Error applying: 
    Error: Provider produced inconsistent result after apply
    
    When applying changes to microsoft365_graph_beta_device_and_app_management_android_managed_device_app_configuration_policy.test,
    provider "registry.terraform.io/deploymenttheory/microsoft365" produced an unexpected new value:
    .assignments: was null, but now cty.SetValEmpty(cty.Object(map[string]cty.Type{"target":cty.Object...
```

### Resources

- [Workflow Logs](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/19383092062)
- [Test Source](../../internal/services/resources/device_and_app_management)

---
*Automated report from nightly tests*
```

### Example Job Failure Issue

**Title:** `💥 Job Failure: Test Resources - device_and_app_management`

**Body:**
```markdown
## Job-Level Failure

**Job:** `Test Resources - device_and_app_management`  
**Failure Type:** Out Of Memory  
**Failed Step:** `Run tests for resources/device_and_app_management`  
**Date:** 2025-11-15  
**Workflow:** [19383092062](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/19383092062)

### Details

- **Job ID:** 123456789
- **Conclusion:** failure
- **Runner:** ubuntu-24.04-arm-16core
- **Started:** 2025-11-15T02:30:45Z
- **Completed:** 2025-11-15T02:45:12Z

### Job Logs

[View Job Logs](https://github.com/deploymenttheory/terraform-provider-microsoft365/actions/runs/19383092062/job/123456789)

### Possible Causes

- Job ran out of memory (OOM)
- Consider reducing parallel test execution
- May need larger runner or memory optimization
- Check if `-race` flag is causing excessive memory usage

---
*Automated report from nightly test pipeline*
```

## 🏷️ Issue Labels

### Test Failure Labels

Issues managed by `manage-test-issues.py` use these labels:

- `test-failure`: Identifies failing test issues
- `automated`: Automatically generated
- `recurring`: Added when test fails multiple times

### Job Failure Labels

Issues managed by `manage-job-failure-issues.py` use these labels:

- `job-failure`: Identifies job-level infrastructure failures
- `infrastructure`: Runner or setup issues
- `automated`: Automatically generated
- `recurring`: Added when job fails multiple times

### Label Lifecycle

**Initial:** `test-failure`, `automated`  
**Recurring:** Adds `recurring` label automatically  
**Manual Labels:** Add as needed:
- `bug`: Code defect
- `flaky-test`: Intermittent failures  
- `wontfix`: Test issue, but not fixing

## 🔍 Troubleshooting

### No Reports Created

**Problem:** Tests fail but no reports are created  
**Causes:**
- `test-failures.json` is empty
- No failing tests matched the `--- FAIL:` pattern
- `GITHUB_TOKEN` permissions insufficient

**Solution:**
1. Check test output logs in artifacts
2. Verify `test-failures.json` contains data
3. Ensure workflow has `issues: write` permission

### Duplicate Reports

**Problem:** Multiple reports created for same test  
**Causes:**
- Test name variation (e.g., TestFoo vs TestFoo/subtest)
- Report title search not matching

**Solution:**
- Issue titles are the exact test name
- Script searches for `in:title "TestName"`
- Ensure test names are stable

### Missing Context

**Problem:** Report created but no failure details  
**Causes:**
- Test output parsing failed
- Context extraction logic didn't match output format

**Solution:**
- Check `test-output-*.log` artifacts
- Verify Go test output follows standard format
- Update parse logic in `run-tests.py` if needed

## 🚀 Future Improvements

Potential enhancements:

1. **Auto-close Resolved Issues**: Close issues when test passes in subsequent run
2. **Flaky Test Detection**: Track failure frequency, auto-label intermittent failures
3. **Trend Analysis**: Add metrics on failure rates and patterns over time
4. **Test Retry Logic**: Automatically retry failed tests once before reporting
5. **Slack/Teams Integration**: Notify team channels of new failures

## 📚 Related Documentation

- [GitHub Actions Workflow](.github/workflows/nightly-tests.yml)
- [Testing Guide](../../docs/TESTING.md)
- [Contributing Guidelines](../../CONTRIBUTING.md)


package utilityDeploymentScheduler_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	utilityDeploymentScheduler "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/datasources/utility/deployment_scheduler"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	// DataSource type name from the datasource package
	dataSourceType = utilityDeploymentScheduler.DataSourceName
)

// ===========================================================================
// Time Condition Tests
// ===========================================================================

// WHY: Gate should open immediately because delay_start_time_by = 0 and deployment_start_time
// is in the past (2024-01-01), so condition is met and scope_id should be released
func TestUnitDeploymentSchedulerDataSource_TimeCondition_ImmediateRelease(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_time_condition_immediate.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").HasValue("12345678-1234-1234-1234-123456789abc"),
					check.That("data."+dataSourceType+".test").Key("status_message").Exists(),
					check.That("data."+dataSourceType+".test").Key("id").IsSet(),
					resource.TestCheckOutput("condition_met", "true"),
					resource.TestCheckOutput("released_scope_id", "12345678-1234-1234-1234-123456789abc"),
				),
			},
		},
	})
}

// PURPOSE: Verify gate remains closed when deployment_start_time is in the future
// EXPECTED: condition_met=false, released_scope_id=null, status shows delay not elapsed
// WHY: deployment_start_time is 2099 (73 years in future), so delay_start_time_by of 48h cannot be met
func TestUnitDeploymentSchedulerDataSource_TimeCondition_NotMetYet(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_time_condition_not_met.tf"),
				Check: resource.ComposeTestCheckFunc(
					// Datasource attributes
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("false"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").DoesNotExist(),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile("(?i)\\[GATE CLOSED\\].*FAIL.*Time.*Delay not elapsed")),

					// Output assertions - condition_met is the only output when gate is closed
					resource.TestCheckOutput("condition_met", "false"),
					// released_scope_id output is null/absent when gate closed - verified via datasource assertion above
				),
			},
		},
	})
}

// PURPOSE: Verify validation rejects negative delay_start_time_by
// EXPECTED: Terraform plan fails with validation error
// WHY: delay_start_time_by must be non-negative (can't delay by negative hours)
func TestUnitDeploymentSchedulerDataSource_TimeCondition_InvalidNegativeOffset(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("03_time_condition_negative_offset.tf"),
				ExpectError: regexp.MustCompile("delay_start_time_by must be >= 0"),
			},
		},
	})
}

// PURPOSE: Test absolute_earliest time constraint - gate won't open before specified time
// EXPECTED: condition_met=true, released_scope_id=GUID, gate open
// WHY: absolute_earliest is 2025, current time is 2026, so constraint is met
func TestUnitDeploymentSchedulerDataSource_TimeCondition_AbsoluteEarliest(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("27_time_condition_absolute_earliest.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").HasValue("12345678-1234-1234-1234-123456789abc"),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile(`(?i)\[GATE OPEN\].*Time`)),

					resource.TestCheckOutput("condition_met", "true"),
					resource.TestCheckOutput("released_scope_id", "12345678-1234-1234-1234-123456789abc"),
				),
			},
		},
	})
}

// PURPOSE: Test absolute_latest time constraint - gate permanently closes after this time
// EXPECTED: condition_met=true, released_scope_id=GUID, gate open
// WHY: absolute_latest is 2027, current time is 2026, so still before deadline
func TestUnitDeploymentSchedulerDataSource_TimeCondition_AbsoluteLatest(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("28_time_condition_absolute_latest.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").HasValue("12345678-1234-1234-1234-123456789abc"),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile(`(?i)\[GATE OPEN\].*Time`)),

					resource.TestCheckOutput("condition_met", "true"),
					resource.TestCheckOutput("released_scope_id", "12345678-1234-1234-1234-123456789abc"),
				),
			},
		},
	})
}

// PURPOSE: Test max_open_duration_hours - gate auto-closes after being open for specified time
// EXPECTED: condition_met=false, released_scope_id=null, gate closed
// WHY: Gate opened at deployment_start_time (2024-01-01), been open ~2 years (17544h), max is 2 years (17520h), exceeded
func TestUnitDeploymentSchedulerDataSource_TimeCondition_MaxDuration(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("29_time_condition_max_duration.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("false"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").DoesNotExist(),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile(`(?i)\[GATE CLOSED\].*FAIL.*Time`)),

					resource.TestCheckOutput("condition_met", "false"),
					// released_scope_id output is null/absent when gate closed
				),
			},
		},
	})
}

// PURPOSE: Test combination of delay_start_time_by, absolute_earliest, and absolute_latest
// EXPECTED: condition_met=true, released_scope_id=GUID, gate open
// WHY: delay=0 (met), absolute_earliest=2025 (2026 > 2025, met), absolute_latest=2027 (2026 < 2027, met)
func TestUnitDeploymentSchedulerDataSource_TimeCondition_CombinedAdvanced(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("30_time_condition_combined_advanced.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").HasValue("12345678-1234-1234-1234-123456789abc"),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile(`(?i)\[GATE OPEN\].*Time`)),

					resource.TestCheckOutput("condition_met", "true"),
					resource.TestCheckOutput("released_scope_id", "12345678-1234-1234-1234-123456789abc"),
				),
			},
		},
	})
}

// ===========================================================================
// Inclusion Time Windows Tests
// ===========================================================================

// PURPOSE: Test inclusion window with day_of_week constraint (all weekdays)
// EXPECTED: Varies - condition_met=true on weekdays (Mon-Fri), false on weekends
// WHY: Current day is Thursday (2026-01-16), which is in the inclusion window, so gate opens
func TestUnitDeploymentSchedulerDataSource_InclusionWindow_WithinDayOfWeek(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_inclusion_window_day_of_week.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").Exists(),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile(`(?i)\[(GATE OPEN|GATE CLOSED)\].*(PASS|FAIL).*Inclusion Window`)),

					// Output values vary by day - verified via datasource key assertions above
				),
			},
		},
	})
}

// PURPOSE: Test inclusion window with time_of_day constraint covering entire day (00:00-23:59)
// EXPECTED: condition_met=true, released_scope_id=GUID, always passes
// WHY: Window covers 24 hours, so current time always falls within it
func TestUnitDeploymentSchedulerDataSource_InclusionWindow_WithinTimeOfDay(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_inclusion_window_time_of_day.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").HasValue("12345678-1234-1234-1234-123456789abc"),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile(`(?i)\[GATE OPEN\].*Inclusion Window`)),

					resource.TestCheckOutput("condition_met", "true"),
					resource.TestCheckOutput("released_scope_id", "12345678-1234-1234-1234-123456789abc"),
				),
			},
		},
	})
}

// PURPOSE: Test inclusion window with date_start/date_end range constraint
// EXPECTED: condition_met=true, released_scope_id=GUID, gate open
// WHY: Current date (2026-01-16) falls within date range (2026-01-01 to 2026-12-31)
func TestUnitDeploymentSchedulerDataSource_InclusionWindow_WithinDateRange(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("06_inclusion_window_date_range.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").HasValue("12345678-1234-1234-1234-123456789abc"),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile(`(?i)\[GATE OPEN\].*Inclusion Window`)),

					resource.TestCheckOutput("condition_met", "true"),
					resource.TestCheckOutput("released_scope_id", "12345678-1234-1234-1234-123456789abc"),
				),
			},
		},
	})
}

// PURPOSE: Test multiple inclusion windows - passes if ANY window matches
// EXPECTED: Varies - depends on current day and time
// WHY: Windows are (Mon/Wed/Fri 09-12) OR (Tue/Thu 14-17). Result varies by when test runs.
func TestUnitDeploymentSchedulerDataSource_InclusionWindow_MultipleWindows(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("07_inclusion_window_multiple.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").Exists(),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile(`(?i)\[(GATE OPEN|GATE CLOSED)\].*(PASS|FAIL).*Inclusion Window`)),

					// Output values vary by day and time - verified via datasource key assertions above
				),
			},
		},
	})
}

// ===========================================================================
// Exclusion Time Windows Tests
// ===========================================================================

// PURPOSE: Test exclusion window blocking deployment on specific days (Sat/Sun)
// EXPECTED: Varies - condition_met=false on weekends, true on weekdays
// WHY: Current day is Thursday (2026-01-16), not in exclusion window, so gate opens
func TestUnitDeploymentSchedulerDataSource_ExclusionWindow_BlockedByDayOfWeek(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("08_exclusion_window_day_of_week.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").Exists(),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile(`(?i)\[(GATE OPEN|GATE CLOSED)\].*(PASS|FAIL).*Exclusion Window`)),

					// Output values vary by day - verified via datasource key assertions above
				),
			},
		},
	})
}

// PURPOSE: Test exclusion window blocking deployment during date range (holiday freeze)
// EXPECTED: condition_met=true, released_scope_id=GUID, gate open
// WHY: Current date (2026-01-16) is NOT in exclusion range (2026-12-20 to 2027-01-05), so not blocked
func TestUnitDeploymentSchedulerDataSource_ExclusionWindow_BlockedByDateRange(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("09_exclusion_window_date_range.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").HasValue("12345678-1234-1234-1234-123456789abc"),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile(`(?i)\[GATE OPEN\].*PASS.*Exclusion Window`)),

					resource.TestCheckOutput("condition_met", "true"),
					resource.TestCheckOutput("released_scope_id", "12345678-1234-1234-1234-123456789abc"),
				),
			},
		},
	})
}

// ===========================================================================
// Manual Override Tests
// ===========================================================================

// PURPOSE: Test manual_override flag forces gate open regardless of time conditions
// EXPECTED: condition_met=true, released_scope_id=GUID, gate forced open
// WHY: manual_override=true bypasses all conditions (deployment_start_time is 2099, time not met, but overridden)
func TestUnitDeploymentSchedulerDataSource_ManualOverride_Enabled(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("10_manual_override_enabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").HasValue("12345678-1234-1234-1234-123456789abc"),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile("(?i)Manual override enabled")),

					resource.TestCheckOutput("condition_met", "true"),
					resource.TestCheckOutput("released_scope_id", "12345678-1234-1234-1234-123456789abc"),
				),
			},
		},
	})
}

// PURPOSE: Test manual_override bypasses ALL conditions including time and exclusion windows
// EXPECTED: condition_met=true, released_scope_id=GUID, gate forced open
// WHY: manual_override=true bypasses everything (deployment_start_time=2099, delay=168h, exclusion active on all days)
func TestUnitDeploymentSchedulerDataSource_ManualOverride_BypassesAllConditions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("11_manual_override_bypasses_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").HasValue("12345678-1234-1234-1234-123456789abc"),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile("(?i)Manual override enabled")),

					resource.TestCheckOutput("condition_met", "true"),
					resource.TestCheckOutput("released_scope_id", "12345678-1234-1234-1234-123456789abc"),
				),
			},
		},
	})
}

// ===========================================================================
// Dependency Gate Tests
// ===========================================================================

// PURPOSE: Test dependency gate when prerequisite hasn't opened yet
// EXPECTED: condition_met=false, released_scope_id=null, gate closed
// WHY: deployment_start_time=2099 (future), so time condition fails first (before checking dependency)
func TestUnitDeploymentSchedulerDataSource_Dependency_NotMetPrerequisiteNotOpen(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("12_dependency_prerequisite_not_open.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("false"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").DoesNotExist(),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile("(?i)\\[GATE CLOSED\\].*FAIL.*Time")),

					resource.TestCheckOutput("condition_met", "false"),
					// released_scope_id output is null/absent when gate closed
				),
			},
		},
	})
}

// PURPOSE: Test dependency gate when prerequisite opened but hasn't been open long enough
// EXPECTED: condition_met=false, released_scope_id=null, gate closed
// WHY: deployment_start_time=2099 (future), so time condition fails first (before checking dependency)
func TestUnitDeploymentSchedulerDataSource_Dependency_NotMetInsufficientTime(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("13_dependency_insufficient_time.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("false"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").DoesNotExist(),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile("(?i)\\[GATE CLOSED\\].*FAIL.*Time")),

					resource.TestCheckOutput("condition_met", "false"),
					// released_scope_id output is null/absent when gate closed
				),
			},
		},
	})
}

// PURPOSE: Test dependency gate when prerequisite requirements are satisfied
// EXPECTED: condition_met=true, released_scope_id=GUID, gate open
// WHY: Prerequisite delay=0, minimum_open_hours=0, deployment_start_time=2024 (past 2 years ago), dependency satisfied
func TestUnitDeploymentSchedulerDataSource_Dependency_Met(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("14_dependency_met.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").HasValue("12345678-1234-1234-1234-123456789abc"),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile("(?i)(PASS.*Dependency|Prerequisite open)")),

					resource.TestCheckOutput("condition_met", "true"),
					resource.TestCheckOutput("released_scope_id", "12345678-1234-1234-1234-123456789abc"),
				),
			},
		},
	})
}

// PURPOSE: Verify validation rejects negative minimum_open_hours in dependency
// EXPECTED: Terraform plan fails with validation error
// WHY: minimum_open_hours must be non-negative (can't require negative hours open)
func TestUnitDeploymentSchedulerDataSource_Dependency_InvalidNegativeMinimumHours(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("15_dependency_negative_minimum.tf"),
				ExpectError: regexp.MustCompile("minimum_open_hours must be >= 0"),
			},
		},
	})
}

// ===========================================================================
// Scope ID Validation Tests
// ===========================================================================

// PURPOSE: Test using singular scope_id (single GUID to release)
// EXPECTED: condition_met=true, released_scope_id=GUID, released_scope_ids=null
// WHY: Using scope_id releases single GUID value, conditions met (deployment_start_time=2024, delay=0)
func TestUnitDeploymentSchedulerDataSource_ScopeId_Singular(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("16_scope_id_singular.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("scope_id").HasValue("12345678-1234-1234-1234-123456789abc"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").HasValue("12345678-1234-1234-1234-123456789abc"),
					check.That("data."+dataSourceType+".test").Key("released_scope_ids.#").DoesNotExist(),

					resource.TestCheckOutput("condition_met", "true"),
					resource.TestCheckOutput("released_scope_id", "12345678-1234-1234-1234-123456789abc"),
				),
			},
		},
	})
}

// PURPOSE: Test using plural scope_ids (multiple GUIDs to release)
// EXPECTED: condition_met=true, released_scope_ids=[2 GUIDs], released_scope_id=null
// WHY: Using scope_ids releases list of GUIDs, conditions met (deployment_start_time=2024, delay=0)
func TestUnitDeploymentSchedulerDataSource_ScopeIds_Plural(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("17_scope_ids_plural.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("scope_ids.#").HasValue("2"),
					check.That("data."+dataSourceType+".test").Key("released_scope_ids.#").HasValue("2"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").DoesNotExist(),

					resource.TestCheckOutput("condition_met", "true"),
					// released_scope_ids is a list - verified via datasource key assertions above
				),
			},
		},
	})
}

// PURPOSE: Verify validation rejects when both scope_id and scope_ids are provided
// EXPECTED: Terraform plan fails with validation error
// WHY: Must use either scope_id OR scope_ids, not both (mutually exclusive)
func TestUnitDeploymentSchedulerDataSource_ScopeId_BothProvided(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("18_scope_both_provided.tf"),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
		},
	})
}

// PURPOSE: Verify validation rejects when neither scope_id nor scope_ids are provided
// EXPECTED: Terraform plan fails with validation error
// WHY: Must provide either scope_id OR scope_ids (at least one required)
func TestUnitDeploymentSchedulerDataSource_ScopeId_NeitherProvided(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("19_scope_neither_provided.tf"),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
		},
	})
}

// ===========================================================================
// Combination Tests
// ===========================================================================

// PURPOSE: Test combination of time condition AND inclusion window - both must pass
// EXPECTED: Varies - depends on current day (weekdays only)
// WHY: Time met (deployment_start_time=2024, delay=0), inclusion requires Mon-Fri. Today is Thu, so passes.
func TestUnitDeploymentSchedulerDataSource_Combination_TimeAndInclusionWindow(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("20_combination_time_and_inclusion.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").Exists(),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile(`(?i)\[(GATE OPEN|GATE CLOSED)\].*Time.*Inclusion Window`)),

					// Output values vary by day - verified via datasource key assertions above
				),
			},
		},
	})
}

// PURPOSE: Test combination of time condition with exclusion window - exclusion blocks if active
// EXPECTED: Varies - condition_met=false on weekends (blocked), true on weekdays
// WHY: Time met (deployment_start_time=2024, delay=0), exclusion blocks Sat/Sun. Today is Thu, so passes.
func TestUnitDeploymentSchedulerDataSource_Combination_TimeAndExclusionWindow(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("21_combination_time_and_exclusion.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").Exists(),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile(`(?i)\[(GATE OPEN|GATE CLOSED)\].*Time.*Exclusion Window`)),

					// Output values vary by day - verified via datasource key assertions above
				),
			},
		},
	})
}

// PURPOSE: Test ALL condition types together - time, inclusion, exclusion, dependency
// EXPECTED: Varies - depends on current day, time, and exclusion date range
// WHY: Complex: time=met, inclusion=Mon-Fri 09-17, exclusion=Dec20-Jan5, dependency=met. Varies by when test runs.
func TestUnitDeploymentSchedulerDataSource_Combination_AllConditions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("22_combination_all_conditions.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").Exists(),
					check.That("data."+dataSourceType+".test").Key("status_message").MatchesRegex(regexp.MustCompile(`(?i)\[(GATE OPEN|GATE CLOSED)\].*Time.*Inclusion.*Exclusion.*Dependency`)),

					// Output values vary by day, time, and date - verified via datasource key assertions above
				),
			},
		},
	})
}

// ===========================================================================
// Released Values Tests
// ===========================================================================

// PURPOSE: Test released values when gate is open - scope_id should be released
// EXPECTED: condition_met=true, released_scope_id=GUID, released_scope_ids=null
// WHY: All conditions met (deployment_start_time=2024, delay=0), gate open releases scope_id
func TestUnitDeploymentSchedulerDataSource_ReleasedValues_GateOpen(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("23_released_values_gate_open.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("true"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").HasValue("12345678-1234-1234-1234-123456789abc"),
					check.That("data."+dataSourceType+".test").Key("released_scope_ids.#").DoesNotExist(),

					resource.TestCheckOutput("condition_met", "true"),
					resource.TestCheckOutput("released_scope_id", "12345678-1234-1234-1234-123456789abc"),
				),
			},
		},
	})
}

// PURPOSE: Test released values when gate is closed - scope_id should be null
// EXPECTED: condition_met=false, released_scope_id=null, released_scope_ids=null
// WHY: Conditions not met (deployment_start_time=2099, delay=48h, in future), gate closed withholds values
func TestUnitDeploymentSchedulerDataSource_ReleasedValues_GateClosed(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("24_released_values_gate_closed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That("data."+dataSourceType+".test").Key("condition_met").HasValue("false"),
					check.That("data."+dataSourceType+".test").Key("released_scope_id").DoesNotExist(),
					check.That("data."+dataSourceType+".test").Key("released_scope_ids.#").DoesNotExist(),

					resource.TestCheckOutput("condition_met", "false"),
					// released_scope_id output is null/absent when gate closed
				),
			},
		},
	})
}

// ===========================================================================
// Validation Tests
// ===========================================================================

// PURPOSE: Verify GUID validation rejects invalid scope_id format
// EXPECTED: Terraform plan fails with validation error
// WHY: scope_id must be valid GUID (xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx), "not-a-valid-guid" is invalid
func TestUnitDeploymentSchedulerDataSource_Validation_InvalidScopeIdFormat(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("25_invalid_scope_id_format.tf"),
				ExpectError: regexp.MustCompile("Must be a valid GUID format"),
			},
		},
	})
}

// PURPOSE: Verify GUID validation rejects invalid scope_ids format
// EXPECTED: Terraform plan fails with validation error
// WHY: Each scope_ids element must be valid GUID, list contains invalid entries like "not-a-valid-guid"
func TestUnitDeploymentSchedulerDataSource_Validation_InvalidScopeIdsFormat(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("26_invalid_scope_ids_format.tf"),
				ExpectError: regexp.MustCompile("Each scope ID must be a valid GUID format"),
			},
		},
	})
}

// ===========================================================================
// Configuration functions
// ===========================================================================

// loadUnitTestTerraform loads a Terraform configuration file from the unit test directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

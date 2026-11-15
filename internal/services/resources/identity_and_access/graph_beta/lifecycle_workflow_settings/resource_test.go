package graphBetaIdentityAndAccessLifecycleWorkflowSettings_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	lifecycleWorkflowSettingsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/lifecycle_workflow_settings/mocks"
)

func setupMockEnvironment() (*mocks.Mocks, *lifecycleWorkflowSettingsMocks.LifecycleWorkflowSettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	lifecycleWorkflowSettingsMock := &lifecycleWorkflowSettingsMocks.LifecycleWorkflowSettingsMock{}
	lifecycleWorkflowSettingsMock.RegisterMocks()
	return mockClient, lifecycleWorkflowSettingsMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *lifecycleWorkflowSettingsMocks.LifecycleWorkflowSettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	lifecycleWorkflowSettingsMock := &lifecycleWorkflowSettingsMocks.LifecycleWorkflowSettingsMock{}
	lifecycleWorkflowSettingsMock.RegisterErrorMocks()
	return mockClient, lifecycleWorkflowSettingsMock
}

func setupNotFoundMockEnvironment() (*mocks.Mocks, *lifecycleWorkflowSettingsMocks.LifecycleWorkflowSettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	lifecycleWorkflowSettingsMock := &lifecycleWorkflowSettingsMocks.LifecycleWorkflowSettingsMock{}
	lifecycleWorkflowSettingsMock.RegisterNotFoundMocks()
	return mockClient, lifecycleWorkflowSettingsMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// importStep verifies import by the fixed singleton ID.
func importStep() resource.TestStep {
	return resource.TestStep{
		ResourceName:            resourceType + ".test",
		ImportState:             true,
		ImportStateId:           "lifecycleWorkflowSettings",
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"timeouts"},
	}
}

// checkLastPatchBody asserts a value in the most recent PATCH request body sent to the mock.
func checkLastPatchBody(mock *lifecycleWorkflowSettingsMocks.LifecycleWorkflowSettingsMock, check func(body map[string]any) error) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		patches := mock.PatchRequests()
		if len(patches) == 0 {
			return fmt.Errorf("expected at least one PATCH request")
		}
		return check(patches[len(patches)-1])
	}
}

// checkSettingsRemainAfterDestroy asserts that destroy left the settings in the mock tenant unchanged.
func checkSettingsRemainAfterDestroy(mock *lifecycleWorkflowSettingsMocks.LifecycleWorkflowSettingsMock, expectedInterval float64) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		settings := mock.Settings()
		if settings["workflowScheduleIntervalInHours"] != expectedInterval {
			return fmt.Errorf("expected destroy to leave workflowScheduleIntervalInHours %v, got %v", expectedInterval, settings["workflowScheduleIntervalInHours"])
		}
		return nil
	}
}

// Minimal: only the schedule interval is configured; the email settings come from the tenant.
func TestUnitResourceLifecycleWorkflowSettings_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy:             checkSettingsRemainAfterDestroy(lifecycleWorkflowSettingsMock, 6),
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("lifecycleWorkflowSettings"),
					check.That(resourceType+".test").Key("workflow_schedule_interval_in_hours").HasValue("6"),
					check.That(resourceType+".test").Key("email_settings.sender_domain").HasValue("microsoft.com"),
					check.That(resourceType+".test").Key("email_settings.use_company_branding").HasValue("false"),
					checkLastPatchBody(lifecycleWorkflowSettingsMock, func(body map[string]any) error {
						if body["workflowScheduleIntervalInHours"] != float64(6) {
							return fmt.Errorf("expected workflowScheduleIntervalInHours 6, got %v", body["workflowScheduleIntervalInHours"])
						}
						if _, ok := body["emailSettings"]; ok {
							return fmt.Errorf("expected unconfigured emailSettings to be omitted, got %v", body)
						}
						return nil
					}),
				),
			},
			importStep(),
		},
	})
}

// Maximal: every configurable value is set.
func TestUnitResourceLifecycleWorkflowSettings_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("workflow_schedule_interval_in_hours").HasValue("12"),
					check.That(resourceType+".test").Key("email_settings.sender_domain").HasValue("contoso.com"),
					check.That(resourceType+".test").Key("email_settings.use_company_branding").HasValue("true"),
				),
			},
			importStep(),
		},
	})
}

// Update: the interval and email settings change in place.
func TestUnitResourceLifecycleWorkflowSettings_03_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy:             checkSettingsRemainAfterDestroy(lifecycleWorkflowSettingsMock, 24),
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("workflow_schedule_interval_in_hours").HasValue("12"),
					check.That(resourceType+".test").Key("email_settings.use_company_branding").HasValue("true"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("workflow_schedule_interval_in_hours").HasValue("24"),
					check.That(resourceType+".test").Key("email_settings.sender_domain").HasValue("microsoft.com"),
					check.That(resourceType+".test").Key("email_settings.use_company_branding").HasValue("false"),
				),
			},
			importStep(),
		},
	})
}

// PreservesUnconfiguredValues: values that are not configured keep their tenant value, and a value that a
// previous apply set stays in place when the configuration stops managing it.
func TestUnitResourceLifecycleWorkflowSettings_04_PreservesUnconfiguredValues(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_branding_only.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("workflow_schedule_interval_in_hours").HasValue("3"),
					check.That(resourceType+".test").Key("email_settings.sender_domain").HasValue("microsoft.com"),
					check.That(resourceType+".test").Key("email_settings.use_company_branding").HasValue("true"),
					checkLastPatchBody(lifecycleWorkflowSettingsMock, func(body map[string]any) error {
						if _, ok := body["workflowScheduleIntervalInHours"]; ok {
							return fmt.Errorf("expected unconfigured workflowScheduleIntervalInHours to be omitted, got %v", body)
						}
						emailSettings, _ := body["emailSettings"].(map[string]any)
						if _, ok := emailSettings["senderDomain"]; ok {
							return fmt.Errorf("expected unconfigured senderDomain to be omitted, got %v", emailSettings)
						}
						if emailSettings["useCompanyBranding"] != true {
							return fmt.Errorf("expected useCompanyBranding true, got %v", emailSettings)
						}
						return nil
					}),
				),
			},
			{
				// Stop managing email_settings and manage the interval instead: the branding keeps the value from the previous apply.
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("workflow_schedule_interval_in_hours").HasValue("6"),
					check.That(resourceType+".test").Key("email_settings.sender_domain").HasValue("microsoft.com"),
					check.That(resourceType+".test").Key("email_settings.use_company_branding").HasValue("true"),
				),
			},
			importStep(),
		},
	})
}

// ErrorHandling: a 403 from Graph surfaces as a diagnostic.
func TestUnitResourceLifecycleWorkflowSettings_05_ErrorHandlingForbidden(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowSettingsMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`(?i)forbidden|insufficient privileges|403`),
			},
		},
	})
}

// ErrorHandling: a 404 from Graph surfaces as a diagnostic.
func TestUnitResourceLifecycleWorkflowSettings_06_ErrorHandlingNotFound(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowSettingsMock := setupNotFoundMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`(?i)not\s*found|404`),
			},
		},
	})
}

// ErrorHandling: a 400 from Graph for an unverified sender domain surfaces as a diagnostic.
func TestUnitResourceLifecycleWorkflowSettings_07_ErrorHandlingBadRequest(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_unverified_domain.tf"),
				ExpectError: regexp.MustCompile(`(?i)bad\s*request|400|not verified`),
			},
		},
	})
}

// ValidationErrors: invalid configuration and import IDs are rejected before any request is sent.
func TestUnitResourceLifecycleWorkflowSettings_08_ValidationErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, lifecycleWorkflowSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer lifecycleWorkflowSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_invalid_interval_low.tf"),
				ExpectError: regexp.MustCompile(`(?s)workflow_schedule_interval_in_hours.*between\s+1\s+and\s+24`),
			},
			{
				Config:      loadUnitTestTerraform("resource_invalid_interval_high.tf"),
				ExpectError: regexp.MustCompile(`(?s)workflow_schedule_interval_in_hours.*between\s+1\s+and\s+24`),
			},
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
			},
			{
				ResourceName:  resourceType + ".test",
				ImportState:   true,
				ImportStateId: "settings",
				ExpectError:   regexp.MustCompile(`(?s)Invalid\s+import\s+ID.*lifecycleWorkflowSettings`),
			},
		},
	})
}

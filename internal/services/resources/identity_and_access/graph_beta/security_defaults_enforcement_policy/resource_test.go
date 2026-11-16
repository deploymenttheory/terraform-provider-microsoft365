package graphBetaIdentityAndAccessSecurityDefaultsEnforcementPolicy_test

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
	securityDefaultsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/security_defaults_enforcement_policy/mocks"
)

const (
	expectedDisplayName = "Security Defaults"
	expectedDescription = "Security defaults is a set of basic identity security mechanisms recommended by Microsoft. When enabled, these recommendations will be automatically enforced in your organization. Administrators and users will be better protected from common identity related attacks."
)

func setupMockEnvironment() (*mocks.Mocks, *securityDefaultsMocks.SecurityDefaultsEnforcementPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	securityDefaultsMock := &securityDefaultsMocks.SecurityDefaultsEnforcementPolicyMock{}
	securityDefaultsMock.RegisterMocks()
	return mockClient, securityDefaultsMock
}

func setupConditionalAccessConflictMockEnvironment() (*mocks.Mocks, *securityDefaultsMocks.SecurityDefaultsEnforcementPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	securityDefaultsMock := &securityDefaultsMocks.SecurityDefaultsEnforcementPolicyMock{}
	securityDefaultsMock.RegisterConditionalAccessConflictMocks()
	return mockClient, securityDefaultsMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *securityDefaultsMocks.SecurityDefaultsEnforcementPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	securityDefaultsMock := &securityDefaultsMocks.SecurityDefaultsEnforcementPolicyMock{}
	securityDefaultsMock.RegisterErrorMocks()
	return mockClient, securityDefaultsMock
}

func setupNotFoundMockEnvironment() (*mocks.Mocks, *securityDefaultsMocks.SecurityDefaultsEnforcementPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	securityDefaultsMock := &securityDefaultsMocks.SecurityDefaultsEnforcementPolicyMock{}
	securityDefaultsMock.RegisterNotFoundMocks()
	return mockClient, securityDefaultsMock
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
		ImportStateId:           "identitySecurityDefaultsEnforcementPolicy",
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"timeouts"},
	}
}

// checkLastPatchBody asserts the most recent PATCH request body sent to the mock: only isEnabled is sent.
func checkLastPatchBody(mock *securityDefaultsMocks.SecurityDefaultsEnforcementPolicyMock, expectedIsEnabled bool) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		patches := mock.PatchRequests()
		if len(patches) == 0 {
			return fmt.Errorf("expected at least one PATCH request")
		}
		body := patches[len(patches)-1]
		if body["isEnabled"] != expectedIsEnabled {
			return fmt.Errorf("expected isEnabled %v, got %v", expectedIsEnabled, body)
		}
		for _, readOnly := range []string{"displayName", "description", "id"} {
			if _, ok := body[readOnly]; ok {
				return fmt.Errorf("expected read-only %s to be omitted, got %v", readOnly, body)
			}
		}
		return nil
	}
}

// checkPolicyRemainsAfterDestroy asserts that destroy left security defaults in the mock tenant unchanged.
func checkPolicyRemainsAfterDestroy(mock *securityDefaultsMocks.SecurityDefaultsEnforcementPolicyMock, expectedIsEnabled bool) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		policy := mock.Policy()
		if policy["isEnabled"] != expectedIsEnabled {
			return fmt.Errorf("expected destroy to leave isEnabled %v, got %v", expectedIsEnabled, policy["isEnabled"])
		}
		return nil
	}
}

// Disabled: the policy is adopted with security defaults disabled.
func TestUnitResourceSecurityDefaultsEnforcementPolicy_01_Disabled(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, securityDefaultsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer securityDefaultsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyRemainsAfterDestroy(securityDefaultsMock, false),
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_disabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("identitySecurityDefaultsEnforcementPolicy"),
					check.That(resourceType+".test").Key("is_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("display_name").HasValue(expectedDisplayName),
					check.That(resourceType+".test").Key("description").HasValue(expectedDescription),
					checkLastPatchBody(securityDefaultsMock, false),
				),
			},
			importStep(),
		},
	})
}

// Enabled: security defaults are enabled, and destroy leaves them enabled.
func TestUnitResourceSecurityDefaultsEnforcementPolicy_02_Enabled(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, securityDefaultsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer securityDefaultsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyRemainsAfterDestroy(securityDefaultsMock, true),
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_enabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("identitySecurityDefaultsEnforcementPolicy"),
					check.That(resourceType+".test").Key("is_enabled").HasValue("true"),
					check.That(resourceType+".test").Key("display_name").HasValue(expectedDisplayName),
					checkLastPatchBody(securityDefaultsMock, true),
				),
			},
			importStep(),
		},
	})
}

// Update: security defaults are enabled and then disabled in place.
func TestUnitResourceSecurityDefaultsEnforcementPolicy_03_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, securityDefaultsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer securityDefaultsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyRemainsAfterDestroy(securityDefaultsMock, false),
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_disabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("is_enabled").HasValue("false"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_enabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("is_enabled").HasValue("true"),
					checkLastPatchBody(securityDefaultsMock, true),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_disabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("is_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("display_name").HasValue(expectedDisplayName),
					checkLastPatchBody(securityDefaultsMock, false),
				),
			},
			importStep(),
		},
	})
}

// ConditionalAccessConflict: Graph rejects enabling security defaults while Conditional Access policies are
// enabled, and the resource surfaces Graph's error together with an explanation.
func TestUnitResourceSecurityDefaultsEnforcementPolicy_04_ConditionalAccessConflict(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, securityDefaultsMock := setupConditionalAccessConflictMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer securityDefaultsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_enabled.tf"),
				ExpectError: regexp.MustCompile(`(?s)Security defaults could not be enabled.*Conditional\s+Access`),
			},
			{
				// Disabling is still allowed while Conditional Access policies are enabled.
				Config: loadUnitTestTerraform("resource_disabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("is_enabled").HasValue("false"),
				),
			},
			{
				// Enabling through an update is rejected the same way.
				Config:      loadUnitTestTerraform("resource_enabled.tf"),
				ExpectError: regexp.MustCompile(`(?s)Security defaults could not be enabled.*Conditional\s+Access`),
			},
		},
	})
}

// ErrorHandling: a 403 from Graph surfaces as a diagnostic.
func TestUnitResourceSecurityDefaultsEnforcementPolicy_05_ErrorHandlingForbidden(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, securityDefaultsMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer securityDefaultsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_disabled.tf"),
				ExpectError: regexp.MustCompile(`(?i)forbidden|insufficient privileges|403`),
			},
		},
	})
}

// ErrorHandling: a 404 from Graph surfaces as a diagnostic.
func TestUnitResourceSecurityDefaultsEnforcementPolicy_06_ErrorHandlingNotFound(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, securityDefaultsMock := setupNotFoundMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer securityDefaultsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_disabled.tf"),
				ExpectError: regexp.MustCompile(`(?i)not\s*found|404`),
			},
		},
	})
}

// ValidationErrors: invalid configuration and import IDs are rejected before any request is sent.
func TestUnitResourceSecurityDefaultsEnforcementPolicy_07_ValidationErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, securityDefaultsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer securityDefaultsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_missing_is_enabled.tf"),
				ExpectError: regexp.MustCompile(`(?s)is_enabled.*required`),
			},
			{
				Config: loadUnitTestTerraform("resource_disabled.tf"),
			},
			{
				ResourceName:  resourceType + ".test",
				ImportState:   true,
				ImportStateId: "00000000-0000-0000-0000-000000000005",
				ExpectError:   regexp.MustCompile(`(?s)Invalid\s+import\s+ID.*identitySecurityDefaultsEnforcementPolicy`),
			},
		},
	})
}

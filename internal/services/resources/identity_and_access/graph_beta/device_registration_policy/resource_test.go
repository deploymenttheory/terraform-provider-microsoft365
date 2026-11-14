package graphBetaIdentityAndAccessDeviceRegistrationPolicy_test

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
	deviceRegistrationPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/device_registration_policy/mocks"
)

func setupMockEnvironment() (*mocks.Mocks, *deviceRegistrationPolicyMocks.DeviceRegistrationPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	deviceRegistrationPolicyMock := &deviceRegistrationPolicyMocks.DeviceRegistrationPolicyMock{}
	deviceRegistrationPolicyMock.RegisterMocks()
	return mockClient, deviceRegistrationPolicyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *deviceRegistrationPolicyMocks.DeviceRegistrationPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	deviceRegistrationPolicyMock := &deviceRegistrationPolicyMocks.DeviceRegistrationPolicyMock{}
	deviceRegistrationPolicyMock.RegisterErrorMocks()
	return mockClient, deviceRegistrationPolicyMock
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
		ImportStateId:           "deviceRegistrationPolicy",
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"timeouts"},
	}
}

// checkLastPutBody asserts a value in the most recent PUT request body sent to the mock.
func checkLastPutBody(mock *deviceRegistrationPolicyMocks.DeviceRegistrationPolicyMock, check func(body map[string]any) error) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		puts := mock.PutRequests()
		if len(puts) == 0 {
			return fmt.Errorf("expected at least one PUT request")
		}
		return check(puts[len(puts)-1])
	}
}

// Minimal: only user_device_quota is configured; every other value comes from the tenant policy.
func TestUnitResourceDeviceRegistrationPolicy_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deviceRegistrationPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceRegistrationPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("deviceRegistrationPolicy"),
					check.That(resourceType+".test").Key("display_name").HasValue("Device Registration Policy"),
					check.That(resourceType+".test").Key("user_device_quota").HasValue("25"),
					check.That(resourceType+".test").Key("multi_factor_auth_configuration").HasValue("notRequired"),
					check.That(resourceType+".test").Key("local_admin_password.is_enabled").HasValue("true"),
					check.That(resourceType+".test").Key("azure_ad_registration.is_admin_configurable").HasValue("false"),
					check.That(resourceType+".test").Key("azure_ad_registration.allowed_to_register.type").HasValue("all"),
					check.That(resourceType+".test").Key("azure_ad_join.is_admin_configurable").HasValue("true"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.type").HasValue("enumerated"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.#").HasValue("2"),
					check.That(resourceType+".test").Key("azure_ad_join.local_admins.enable_global_admins").HasValue("true"),
					check.That(resourceType+".test").Key("azure_ad_join.local_admins.registering_users.type").HasValue("all"),
				),
			},
			importStep(),
		},
	})
}

// Maximal: every configurable value is set.
func TestUnitResourceDeviceRegistrationPolicy_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deviceRegistrationPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceRegistrationPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("user_device_quota").HasValue("20"),
					check.That(resourceType+".test").Key("multi_factor_auth_configuration").HasValue("required"),
					check.That(resourceType+".test").Key("local_admin_password.is_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("azure_ad_registration.allowed_to_register.type").HasValue("all"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.type").HasValue("enumerated"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.users.#").HasValue("1"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.users.*").ContainsTypeSetElement("33333333-3333-3333-3333-333333333333"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.#").HasValue("1"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.*").ContainsTypeSetElement("44444444-4444-4444-4444-444444444444"),
					check.That(resourceType+".test").Key("azure_ad_join.local_admins.enable_global_admins").HasValue("false"),
					check.That(resourceType+".test").Key("azure_ad_join.local_admins.registering_users.type").HasValue("none"),
				),
			},
			importStep(),
		},
	})
}

// QuotaUpdate: user_device_quota changes in place.
func TestUnitResourceDeviceRegistrationPolicy_03_QuotaUpdate(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deviceRegistrationPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceRegistrationPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("user_device_quota").HasValue("25"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_quota_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("user_device_quota").HasValue("30"),
					check.That(resourceType+".test").Key("local_admin_password.is_enabled").HasValue("true"),
				),
			},
			importStep(),
		},
	})
}

// MembershipTypeChanges: allowed_to_join moves from all to enumerated to none.
func TestUnitResourceDeviceRegistrationPolicy_04_MembershipTypeChanges(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deviceRegistrationPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceRegistrationPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_join_all.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.type").HasValue("all"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.#").DoesNotExist(),
					check.That(resourceType+".test").Key("azure_ad_join.local_admins.registering_users.type").HasValue("all"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_join_enumerated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.type").HasValue("enumerated"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.#").HasValue("1"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.*").ContainsTypeSetElement("55555555-5555-5555-5555-555555555555"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.users.#").DoesNotExist(),
					checkLastPutBody(deviceRegistrationPolicyMock, func(body map[string]any) error {
						join, _ := body["azureADJoin"].(map[string]any)
						allowed, _ := join["allowedToJoin"].(map[string]any)
						if _, ok := allowed["users"]; ok {
							return fmt.Errorf("expected unset users to be omitted from the PUT body, got %v", allowed)
						}
						return nil
					}),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_join_none.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.type").HasValue("none"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.#").DoesNotExist(),
				),
			},
			importStep(),
		},
	})
}

// LAPSToggle: local_admin_password.is_enabled is switched off and on again.
func TestUnitResourceDeviceRegistrationPolicy_05_LAPSToggle(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deviceRegistrationPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceRegistrationPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_laps_disabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("local_admin_password.is_enabled").HasValue("false"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_laps_enabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("local_admin_password.is_enabled").HasValue("true"),
				),
			},
			importStep(),
		},
	})
}

// PreservesUnconfiguredAttributes: the full-replace PUT sends back every value that is not configured,
// including values a previous apply set and the configuration later stopped managing.
func TestUnitResourceDeviceRegistrationPolicy_06_PreservesUnconfiguredAttributes(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deviceRegistrationPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceRegistrationPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_laps_disabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("local_admin_password.is_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("user_device_quota").HasValue("50"),
					check.That(resourceType+".test").Key("multi_factor_auth_configuration").HasValue("notRequired"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.type").HasValue("enumerated"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.#").HasValue("2"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.*").ContainsTypeSetElement("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.*").ContainsTypeSetElement("22222222-2222-2222-2222-222222222222"),
					check.That(resourceType+".test").Key("azure_ad_join.local_admins.enable_global_admins").HasValue("true"),
					checkLastPutBody(deviceRegistrationPolicyMock, func(body map[string]any) error {
						if body["userDeviceQuota"] != float64(50) {
							return fmt.Errorf("expected userDeviceQuota 50 to be sent back, got %v", body["userDeviceQuota"])
						}
						if body["multiFactorAuthConfiguration"] != "notRequired" {
							return fmt.Errorf("expected multiFactorAuthConfiguration notRequired, got %v", body["multiFactorAuthConfiguration"])
						}
						registration, _ := body["azureADRegistration"].(map[string]any)
						if registration["isAdminConfigurable"] != false {
							return fmt.Errorf("expected azureADRegistration.isAdminConfigurable false to be sent back, got %v", registration)
						}
						join, _ := body["azureADJoin"].(map[string]any)
						allowed, _ := join["allowedToJoin"].(map[string]any)
						groups, _ := allowed["groups"].([]any)
						if len(groups) != 2 {
							return fmt.Errorf("expected the 2 existing allowedToJoin groups to be sent back, got %v", allowed)
						}
						return nil
					}),
				),
			},
			{
				// Stop managing local_admin_password and manage the quota instead: LAPS keeps the value
				// from the previous apply and the join groups are still preserved.
				Config: loadUnitTestTerraform("resource_quota_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("user_device_quota").HasValue("30"),
					check.That(resourceType+".test").Key("local_admin_password.is_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("azure_ad_join.allowed_to_join.groups.#").HasValue("2"),
				),
			},
		},
	})
}

// ErrorMocks: Graph errors surface as diagnostics.
func TestUnitResourceDeviceRegistrationPolicy_07_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deviceRegistrationPolicyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceRegistrationPolicyMock.CleanupMockState()

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

// ValidationErrors: invalid configuration is rejected before any request is sent.
func TestUnitResourceDeviceRegistrationPolicy_08_ValidationErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deviceRegistrationPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceRegistrationPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_invalid_membership.tf"),
				ExpectError: regexp.MustCompile(`(?s)groups\s+can\s+only\s+be\s+set\s+when\s+type\s+is\s+"enumerated"`),
			},
			{
				Config:      loadUnitTestTerraform("resource_invalid_quota.tf"),
				ExpectError: regexp.MustCompile(`(?s)user_device_quota.*at\s+least\s+0`),
			},
		},
	})
}

// RegistrationNotAdminConfigurable: changing azure_ad_registration fails when Intune manages registration.
func TestUnitResourceDeviceRegistrationPolicy_09_RegistrationNotAdminConfigurable(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deviceRegistrationPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deviceRegistrationPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create: no state yet, so the check runs at apply time against the policy read before the PUT.
				Config:      loadUnitTestTerraform("resource_registration_change.tf"),
				ExpectError: regexp.MustCompile(`(?s)allowed_to_register\s+cannot\s+be\s+changed.*isAdminConfigurable\s+=\s+false`),
			},
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
			},
			{
				// Update: the check runs at plan time against the state.
				Config:      loadUnitTestTerraform("resource_registration_change.tf"),
				ExpectError: regexp.MustCompile(`(?s)azure_ad_registration\s+cannot\s+be\s+changed.*isAdminConfigurable\s+=\s+false`),
			},
		},
	})
}

package graphBetaM365AdminSharepointSettings_test

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	sharepointSettingsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/m365_admin/graph_beta/sharepoint_settings/mocks"
)

func setupMockEnvironment() (*mocks.Mocks, *sharepointSettingsMocks.SharepointSettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	sharepointSettingsMock := &sharepointSettingsMocks.SharepointSettingsMock{}
	sharepointSettingsMock.RegisterMocks()
	return mockClient, sharepointSettingsMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *sharepointSettingsMocks.SharepointSettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	sharepointSettingsMock := &sharepointSettingsMocks.SharepointSettingsMock{}
	sharepointSettingsMock.RegisterErrorMocks()
	return mockClient, sharepointSettingsMock
}

func setupNotFoundMockEnvironment() (*mocks.Mocks, *sharepointSettingsMocks.SharepointSettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	sharepointSettingsMock := &sharepointSettingsMocks.SharepointSettingsMock{}
	sharepointSettingsMock.RegisterNotFoundMocks()
	return mockClient, sharepointSettingsMock
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
		ImportStateId:           "sharepointSettings",
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"timeouts"},
	}
}

// checkLastPatchBody asserts a value in the most recent PATCH request body sent to the mock.
func checkLastPatchBody(mock *sharepointSettingsMocks.SharepointSettingsMock, check func(body map[string]any) error) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		patches := mock.PatchRequests()
		if len(patches) == 0 {
			return fmt.Errorf("expected at least one PATCH request")
		}
		return check(patches[len(patches)-1])
	}
}

// checkSettingsRemainAfterDestroy asserts that destroy left the settings in the mock tenant unchanged.
func checkSettingsRemainAfterDestroy(mock *sharepointSettingsMocks.SharepointSettingsMock, expectedRetention float64) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		settings := mock.Settings()
		if settings["deletedUserPersonalSiteRetentionPeriodInDays"] != expectedRetention {
			return fmt.Errorf("expected destroy to leave deletedUserPersonalSiteRetentionPeriodInDays %v, got %v", expectedRetention, settings["deletedUserPersonalSiteRetentionPeriodInDays"])
		}
		return nil
	}
}

// Minimal: only the retention period is configured; every other setting comes from the tenant and is not sent.
func TestUnitResourceSharepointSettings_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, sharepointSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer sharepointSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy:             checkSettingsRemainAfterDestroy(sharepointSettingsMock, 60),
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue("sharepointSettings"),
					check.That(resourceType+".test").Key("deleted_user_personal_site_retention_period_in_days").HasValue("60"),
					check.That(resourceType+".test").Key("sharing_capability").HasValue("externalUserAndGuestSharing"),
					check.That(resourceType+".test").Key("is_mac_sync_app_enabled").HasValue("true"),
					check.That(resourceType+".test").Key("available_managed_paths_for_site_creation.#").HasValue("2"),
					check.That(resourceType+".test").Key("excluded_file_extensions_for_sync_app.#").HasValue("1"),
					check.That(resourceType+".test").Key("sharing_allowed_domain_list.#").HasValue("0"),
					check.That(resourceType+".test").Key("idle_session_sign_out.is_enabled").HasValue("true"),
					check.That(resourceType+".test").Key("idle_session_sign_out.warn_after_in_seconds").HasValue("3300"),
					check.That(resourceType+".test").Key("is_require_accepting_user_to_match_invited_user_enabled").DoesNotExist(),
					checkLastPatchBody(sharepointSettingsMock, func(body map[string]any) error {
						expected := map[string]any{"deletedUserPersonalSiteRetentionPeriodInDays": float64(60)}
						if !reflect.DeepEqual(body, expected) {
							return fmt.Errorf("expected only the configured retention period to be sent, got %v", body)
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
func TestUnitResourceSharepointSettings_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, sharepointSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer sharepointSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy:             checkSettingsRemainAfterDestroy(sharepointSettingsMock, 365),
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("allowed_domain_guids_for_sync_app.#").HasValue("1"),
					check.That(resourceType+".test").Key("allowed_domain_guids_for_sync_app.*").ContainsTypeSetElement("bdd1ab9b-3fd0-4def-a761-ec8d7471732c"),
					check.That(resourceType+".test").Key("deleted_user_personal_site_retention_period_in_days").HasValue("365"),
					check.That(resourceType+".test").Key("excluded_file_extensions_for_sync_app.#").HasValue("2"),
					check.That(resourceType+".test").Key("image_tagging_option").HasValue("enhanced"),
					check.That(resourceType+".test").Key("is_commenting_on_site_pages_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("is_loop_enabled").HasValue("true"),
					check.That(resourceType+".test").Key("is_require_accepting_user_to_match_invited_user_enabled").HasValue("true"),
					check.That(resourceType+".test").Key("is_sharepoint_newsfeed_enabled").HasValue("true"),
					check.That(resourceType+".test").Key("is_site_creation_ui_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("is_unmanaged_sync_app_for_tenant_restricted").HasValue("true"),
					check.That(resourceType+".test").Key("personal_site_default_storage_limit_in_mb").HasValue("512000"),
					check.That(resourceType+".test").Key("sharing_allowed_domain_list.#").HasValue("2"),
					check.That(resourceType+".test").Key("sharing_blocked_domain_list.#").HasValue("0"),
					check.That(resourceType+".test").Key("sharing_capability").HasValue("existingExternalUserSharingOnly"),
					check.That(resourceType+".test").Key("sharing_domain_restriction_mode").HasValue("allowList"),
					check.That(resourceType+".test").Key("site_creation_default_managed_path").HasValue("/teams/"),
					check.That(resourceType+".test").Key("site_creation_default_storage_limit_in_mb").HasValue("1048576"),
					check.That(resourceType+".test").Key("tenant_default_timezone").HasValue("(UTC) Coordinated Universal Time"),
					check.That(resourceType+".test").Key("idle_session_sign_out.is_enabled").HasValue("true"),
					check.That(resourceType+".test").Key("idle_session_sign_out.warn_after_in_seconds").HasValue("1800"),
					check.That(resourceType+".test").Key("idle_session_sign_out.sign_out_after_in_seconds").HasValue("3600"),
					checkLastPatchBody(sharepointSettingsMock, func(body map[string]any) error {
						if _, ok := body["availableManagedPathsForSiteCreation"]; ok {
							return fmt.Errorf("expected read-only availableManagedPathsForSiteCreation to be omitted, got %v", body)
						}
						if blocked, ok := body["sharingBlockedDomainList"].([]any); !ok || len(blocked) != 0 {
							return fmt.Errorf("expected an empty sharingBlockedDomainList to be sent as [], got %v", body["sharingBlockedDomainList"])
						}
						return nil
					}),
				),
			},
			importStep(),
		},
	})
}

// Update: every setting changes in place, and the sharing lists switch from allow to block.
func TestUnitResourceSharepointSettings_03_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, sharepointSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer sharepointSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy:             checkSettingsRemainAfterDestroy(sharepointSettingsMock, 90),
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("sharing_domain_restriction_mode").HasValue("allowList"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("allowed_domain_guids_for_sync_app.#").HasValue("0"),
					check.That(resourceType+".test").Key("deleted_user_personal_site_retention_period_in_days").HasValue("90"),
					check.That(resourceType+".test").Key("image_tagging_option").HasValue("disabled"),
					check.That(resourceType+".test").Key("is_require_accepting_user_to_match_invited_user_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("sharing_allowed_domain_list.#").HasValue("0"),
					check.That(resourceType+".test").Key("sharing_blocked_domain_list.#").HasValue("1"),
					check.That(resourceType+".test").Key("sharing_blocked_domain_list.*").ContainsTypeSetElement("fabrikam.com"),
					check.That(resourceType+".test").Key("sharing_capability").HasValue("externalUserSharingOnly"),
					check.That(resourceType+".test").Key("sharing_domain_restriction_mode").HasValue("blockList"),
					check.That(resourceType+".test").Key("idle_session_sign_out.is_enabled").HasValue("false"),
				),
			},
			importStep(),
		},
	})
}

// PreservesUnconfiguredValues: a partial idle_session_sign_out only sends the configured value, and a value that
// a previous apply set stays in place when the configuration stops managing it.
func TestUnitResourceSharepointSettings_04_PreservesUnconfiguredValues(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, sharepointSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer sharepointSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_idle_session_only.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("idle_session_sign_out.is_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("idle_session_sign_out.warn_after_in_seconds").HasValue("3300"),
					check.That(resourceType+".test").Key("idle_session_sign_out.sign_out_after_in_seconds").HasValue("3600"),
					check.That(resourceType+".test").Key("deleted_user_personal_site_retention_period_in_days").HasValue("30"),
					checkLastPatchBody(sharepointSettingsMock, func(body map[string]any) error {
						expected := map[string]any{"idleSessionSignOut": map[string]any{"isEnabled": false}}
						if !reflect.DeepEqual(body, expected) {
							return fmt.Errorf("expected only idleSessionSignOut.isEnabled to be sent, got %v", body)
						}
						return nil
					}),
				),
			},
			{
				// Stop managing idle_session_sign_out and manage the retention period instead: the policy keeps the value from the previous apply.
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("deleted_user_personal_site_retention_period_in_days").HasValue("60"),
					check.That(resourceType+".test").Key("idle_session_sign_out.is_enabled").HasValue("false"),
				),
			},
			importStep(),
		},
	})
}

// ClearList: a configured empty set is sent as [] and clears the list in the tenant.
func TestUnitResourceSharepointSettings_05_ClearList(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, sharepointSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer sharepointSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_clear_list.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("excluded_file_extensions_for_sync_app.#").HasValue("0"),
					checkLastPatchBody(sharepointSettingsMock, func(body map[string]any) error {
						expected := map[string]any{"excludedFileExtensionsForSyncApp": []any{}}
						if !reflect.DeepEqual(body, expected) {
							return fmt.Errorf("expected excludedFileExtensionsForSyncApp to be sent as [], got %v", body)
						}
						return nil
					}),
				),
			},
			importStep(),
		},
	})
}

// ErrorHandling: a 403 from Graph surfaces as a diagnostic.
func TestUnitResourceSharepointSettings_06_ErrorHandlingForbidden(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, sharepointSettingsMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer sharepointSettingsMock.CleanupMockState()

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
func TestUnitResourceSharepointSettings_07_ErrorHandlingNotFound(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, sharepointSettingsMock := setupNotFoundMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer sharepointSettingsMock.CleanupMockState()

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

// ErrorHandling: a 400 from Graph for an idle session warning after the sign-out surfaces as a diagnostic.
func TestUnitResourceSharepointSettings_08_ErrorHandlingBadRequest(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, sharepointSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer sharepointSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_invalid_idle_session.tf"),
				ExpectError: regexp.MustCompile(`(?i)bad\s*request|400|warnAfterInSeconds`),
			},
		},
	})
}

// ValidationErrors: invalid configuration and import IDs are rejected before any request is sent.
func TestUnitResourceSharepointSettings_09_ValidationErrors(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, sharepointSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer sharepointSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_invalid_retention_low.tf"),
				ExpectError: regexp.MustCompile(`(?s)deleted_user_personal_site_retention_period_in_days.*between\s+30\s+and\s+3650`),
			},
			{
				Config:      loadUnitTestTerraform("resource_invalid_retention_high.tf"),
				ExpectError: regexp.MustCompile(`(?s)deleted_user_personal_site_retention_period_in_days.*between\s+30\s+and\s+3650`),
			},
			{
				Config:      loadUnitTestTerraform("resource_invalid_sharing_capability.tf"),
				ExpectError: regexp.MustCompile(`(?s)sharing_capability.*value must be one of`),
			},
			{
				Config:      loadUnitTestTerraform("resource_invalid_domain_guid.tf"),
				ExpectError: regexp.MustCompile(`(?s)allowed_domain_guids_for_sync_app.*lower-case GUID`),
			},
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
			},
			{
				ResourceName:  resourceType + ".test",
				ImportState:   true,
				ImportStateId: "settings",
				ExpectError:   regexp.MustCompile(`(?s)Invalid\s+import\s+ID.*sharepointSettings`),
			},
		},
	})
}

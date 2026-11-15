package graphBetaRoleManagementPolicy_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	roleManagementPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/role_management_policy/mocks"
)

func setupMockEnvironment() (*mocks.Mocks, *roleManagementPolicyMocks.RoleManagementPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	roleManagementPolicyMock := &roleManagementPolicyMocks.RoleManagementPolicyMock{}
	roleManagementPolicyMock.RegisterMocks()
	return mockClient, roleManagementPolicyMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *roleManagementPolicyMocks.RoleManagementPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	roleManagementPolicyMock := &roleManagementPolicyMocks.RoleManagementPolicyMock{}
	roleManagementPolicyMock.RegisterErrorMocks()
	return mockClient, roleManagementPolicyMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// checkLastPatchRules asserts the rules sent in the most recent PATCH request, keyed by rule ID.
func checkLastPatchRules(mock *roleManagementPolicyMocks.RoleManagementPolicyMock, check func(rules map[string]map[string]any) error) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		patches := mock.PatchRequests()
		if len(patches) == 0 {
			return fmt.Errorf("expected at least one PATCH request")
		}
		items, _ := patches[len(patches)-1]["rules"].([]any)
		rules := map[string]map[string]any{}
		for _, item := range items {
			rule, _ := item.(map[string]any)
			id, _ := rule["id"].(string)
			rules[id] = rule
		}
		return check(rules)
	}
}

// checkImportedAllRules asserts that an import tracks every supported rule of the mock policy.
func checkImportedAllRules(states []*terraform.InstanceState) error {
	if len(states) != 1 {
		return fmt.Errorf("expected one imported instance, got %d", len(states))
	}
	attributes := states[0].Attributes
	expected := map[string]string{
		"id":                             roleManagementPolicyMocks.PolicyID,
		"policy_id":                      roleManagementPolicyMocks.PolicyID,
		"role_definition_id":             roleManagementPolicyMocks.RoleDefinitionID,
		"scope_id":                       "/",
		"scope_type":                     "DirectoryRole",
		"expiration_rules.%":             "3",
		"enablement_rules.%":             "3",
		"notification_rules.%":           "9",
		"authentication_context_rules.%": "1",
	}
	for key, want := range expected {
		if got := attributes[key]; got != want {
			return fmt.Errorf("expected imported %s to be %q, got %q", key, want, got)
		}
	}
	return nil
}

// Minimal: adopts the policy by role definition and manages one expiration rule.
func TestUnitResourceRoleManagementPolicy_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleManagementPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleManagementPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue(roleManagementPolicyMocks.PolicyID),
					check.That(resourceType+".test").Key("policy_id").HasValue(roleManagementPolicyMocks.PolicyID),
					check.That(resourceType+".test").Key("role_definition_id").HasValue(roleManagementPolicyMocks.RoleDefinitionID),
					check.That(resourceType+".test").Key("scope_id").HasValue("/"),
					check.That(resourceType+".test").Key("scope_type").HasValue("DirectoryRole"),
					check.That(resourceType+".test").Key("display_name").HasValue("DirectoryRole"),
					check.That(resourceType+".test").Key("is_organization_default").HasValue("false"),
					check.That(resourceType+".test").Key("last_modified_date_time").MatchesRegex(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T`)),
					check.That(resourceType+".test").Key("expiration_rules.%").HasValue("1"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Assignment.is_expiration_required").HasValue("false"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Assignment.maximum_duration").HasValue("P90D"),
					check.That(resourceType+".test").Key("enablement_rules.%").DoesNotExist(),
					check.That(resourceType+".test").Key("notification_rules.%").DoesNotExist(),
					check.That(resourceType+".test").Key("authentication_context_rules.%").DoesNotExist(),
					// Rules that are not configured are not sent, so Graph leaves them unchanged.
					checkLastPatchRules(roleManagementPolicyMock, func(rules map[string]map[string]any) error {
						if len(rules) != 1 {
							return fmt.Errorf("expected only the configured rule in the PATCH body, got %v", rules)
						}
						rule := rules["Expiration_Admin_Assignment"]
						if rule["@odata.type"] != "#microsoft.graph.unifiedRoleManagementPolicyExpirationRule" || rule["maximumDuration"] != "P90D" || rule["isExpirationRequired"] != false {
							return fmt.Errorf("unexpected expiration rule in the PATCH body: %v", rule)
						}
						if _, ok := rule["target"]; ok {
							return fmt.Errorf("expected the rule target to be omitted, got %v", rule)
						}
						return nil
					}),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     roleManagementPolicyMocks.PolicyID,
				ImportStateCheck:  checkImportedAllRules,
				ImportStateVerify: false,
			},
		},
	})
}

// Complete: manages every supported rule, so an import matches the applied state exactly.
func TestUnitResourceRoleManagementPolicy_02_Complete(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleManagementPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleManagementPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_complete.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("expiration_rules.%").HasValue("3"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Eligibility.is_expiration_required").HasValue("true"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Eligibility.maximum_duration").HasValue("P180D"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_EndUser_Assignment.maximum_duration").HasValue("PT8H"),
					check.That(resourceType+".test").Key("enablement_rules.%").HasValue("3"),
					check.That(resourceType+".test").Key("enablement_rules.Enablement_Admin_Eligibility.enabled_rules.#").HasValue("0"),
					check.That(resourceType+".test").Key("enablement_rules.Enablement_EndUser_Assignment.enabled_rules.#").HasValue("2"),
					check.That(resourceType+".test").Key("enablement_rules.Enablement_EndUser_Assignment.enabled_rules.*").ContainsTypeSetElement("MultiFactorAuthentication"),
					check.That(resourceType+".test").Key("notification_rules.%").HasValue("9"),
					check.That(resourceType+".test").Key("notification_rules.Notification_Admin_Admin_Assignment.notification_level").HasValue("Critical"),
					check.That(resourceType+".test").Key("notification_rules.Notification_Admin_Admin_Assignment.is_default_recipients_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("notification_rules.Notification_Admin_Admin_Assignment.notification_recipients.#").HasValue("2"),
					check.That(resourceType+".test").Key("notification_rules.Notification_Admin_Admin_Assignment.notification_recipients.*").ContainsTypeSetElement("security@example.com"),
					check.That(resourceType+".test").Key("notification_rules.Notification_Requestor_Admin_Assignment.notification_recipients.#").HasValue("0"),
					check.That(resourceType+".test").Key("authentication_context_rules.AuthenticationContext_EndUser_Assignment.is_enabled").HasValue("true"),
					check.That(resourceType+".test").Key("authentication_context_rules.AuthenticationContext_EndUser_Assignment.claim_value").HasValue("c1"),
					checkLastPatchRules(roleManagementPolicyMock, func(rules map[string]map[string]any) error {
						if len(rules) != 16 {
							return fmt.Errorf("expected the 16 configured rules in the PATCH body, got %d", len(rules))
						}
						if _, ok := rules["Approval_EndUser_Assignment"]; ok {
							return fmt.Errorf("expected the approval rule to be left out of the PATCH body")
						}
						notification := rules["Notification_Admin_Admin_Assignment"]
						if notification["recipientType"] != "Admin" || notification["notificationType"] != "Email" {
							return fmt.Errorf("expected the notification and recipient types to be copied from the policy, got %v", notification)
						}
						enablement := rules["Enablement_Admin_Eligibility"]
						if enabled, ok := enablement["enabledRules"].([]any); !ok || len(enabled) != 0 {
							return fmt.Errorf("expected an empty enabledRules array, got %v", enablement)
						}
						return nil
					}),
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateId:           roleManagementPolicyMocks.PolicyID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// ExpirationUpdate: changes a managed expiration rule in place and starts managing a second rule.
func TestUnitResourceRoleManagementPolicy_03_ExpirationUpdate(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleManagementPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleManagementPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("expiration_rules.Expiration_Admin_Assignment.maximum_duration").HasValue("P90D"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_expiration_update.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceType+".test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("expiration_rules.%").HasValue("2"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Assignment.is_expiration_required").HasValue("true"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Assignment.maximum_duration").HasValue("P60D"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_EndUser_Assignment.maximum_duration").HasValue("PT4H"),
					checkLastPatchRules(roleManagementPolicyMock, func(rules map[string]map[string]any) error {
						if len(rules) != 2 {
							return fmt.Errorf("expected the 2 configured rules in the PATCH body, got %v", rules)
						}
						return nil
					}),
				),
			},
			{
				// Stop managing a rule: it leaves state and is not changed in Graph.
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("expiration_rules.%").HasValue("1"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Assignment.maximum_duration").HasValue("P90D"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_EndUser_Assignment.maximum_duration").DoesNotExist(),
				),
			},
		},
	})
}

// ByPolicyID: adopts the policy by ID; the role definition is looked up from the policy assignments.
func TestUnitResourceRoleManagementPolicy_04_ByPolicyID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleManagementPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleManagementPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_by_policy_id.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").HasValue(roleManagementPolicyMocks.PolicyID),
					check.That(resourceType+".test").Key("role_definition_id").HasValue(roleManagementPolicyMocks.RoleDefinitionID),
					check.That(resourceType+".test").Key("scope_id").HasValue("/"),
					check.That(resourceType+".test").Key("scope_type").HasValue("DirectoryRole"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Assignment.maximum_duration").HasValue("P90D"),
				),
			},
			{
				// Switching to the equivalent role definition identification keeps the same policy without replacement.
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceType+".test", plancheck.ResourceActionNoop),
					},
				},
			},
		},
	})
}

// DurationNormalisation: Graph returns P1D for PT24H and PT8H for PT480M; state keeps the configured spelling.
func TestUnitResourceRoleManagementPolicy_05_DurationNormalisation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleManagementPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleManagementPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_duration_equivalent.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Assignment.maximum_duration").HasValue("PT24H"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_EndUser_Assignment.maximum_duration").HasValue("PT480M"),
				),
			},
			{
				Config:   loadUnitTestTerraform("resource_duration_equivalent.tf"),
				PlanOnly: true,
			},
			{
				// Respelling a duration to an equivalent value updates state to the new spelling.
				Config: loadUnitTestTerraform("resource_duration_canonical.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("expiration_rules.Expiration_Admin_Assignment.maximum_duration").HasValue("P1D"),
					check.That(resourceType+".test").Key("expiration_rules.Expiration_EndUser_Assignment.maximum_duration").HasValue("PT8H"),
				),
			},
		},
	})
}

// NotificationAndAuthenticationContext: manages notification, enablement and authentication context rules and returns
// them to their defaults.
func TestUnitResourceRoleManagementPolicy_06_NotificationAndAuthenticationContext(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleManagementPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleManagementPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_notification_and_auth_context.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("enablement_rules.Enablement_EndUser_Assignment.enabled_rules.#").HasValue("2"),
					check.That(resourceType+".test").Key("enablement_rules.Enablement_EndUser_Assignment.enabled_rules.*").ContainsTypeSetElement("Ticketing"),
					check.That(resourceType+".test").Key("notification_rules.Notification_Approver_EndUser_Assignment.notification_level").HasValue("Critical"),
					check.That(resourceType+".test").Key("notification_rules.Notification_Approver_EndUser_Assignment.is_default_recipients_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("notification_rules.Notification_Approver_EndUser_Assignment.notification_recipients.*").ContainsTypeSetElement("pim-approvers@example.com"),
					check.That(resourceType+".test").Key("authentication_context_rules.AuthenticationContext_EndUser_Assignment.is_enabled").HasValue("true"),
					check.That(resourceType+".test").Key("authentication_context_rules.AuthenticationContext_EndUser_Assignment.claim_value").HasValue("c4"),
					checkLastPatchRules(roleManagementPolicyMock, func(rules map[string]map[string]any) error {
						notification := rules["Notification_Approver_EndUser_Assignment"]
						if notification["recipientType"] != "Approver" {
							return fmt.Errorf("expected recipientType Approver to be copied from the policy, got %v", notification)
						}
						return nil
					}),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_notification_and_auth_context_reset.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("enablement_rules.Enablement_EndUser_Assignment.enabled_rules.#").HasValue("1"),
					check.That(resourceType+".test").Key("notification_rules.Notification_Approver_EndUser_Assignment.notification_level").HasValue("All"),
					check.That(resourceType+".test").Key("notification_rules.Notification_Approver_EndUser_Assignment.notification_recipients.#").HasValue("0"),
					check.That(resourceType+".test").Key("authentication_context_rules.AuthenticationContext_EndUser_Assignment.is_enabled").HasValue("false"),
					check.That(resourceType+".test").Key("authentication_context_rules.AuthenticationContext_EndUser_Assignment.claim_value").DoesNotExist(),
					checkLastPatchRules(roleManagementPolicyMock, func(rules map[string]map[string]any) error {
						notification := rules["Notification_Approver_EndUser_Assignment"]
						if recipients, ok := notification["notificationRecipients"].([]any); !ok || len(recipients) != 0 {
							return fmt.Errorf("expected an empty notificationRecipients array, got %v", notification)
						}
						if _, ok := rules["AuthenticationContext_EndUser_Assignment"]["claimValue"]; ok {
							return fmt.Errorf("expected claimValue to be omitted for a disabled rule, got %v", rules["AuthenticationContext_EndUser_Assignment"])
						}
						return nil
					}),
				),
			},
		},
	})
}

// Drift: a change made outside Terraform to a managed rule is detected and corrected.
func TestUnitResourceRoleManagementPolicy_07_Drift(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleManagementPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleManagementPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
			},
			{
				PreConfig: func() {
					roleManagementPolicyMock.SetRule("Expiration_Admin_Assignment", map[string]any{"maximumDuration": "P15D"})
				},
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceType+".test", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".test").Key("expiration_rules.Expiration_Admin_Assignment.maximum_duration").HasValue("P90D"),
				),
			},
		},
	})
}

// RuleValidation: rule IDs that do not exist on the policy, or that belong to another rule type, fail with a clear error.
func TestUnitResourceRoleManagementPolicy_08_RuleValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleManagementPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleManagementPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_unknown_rule.tf"),
				ExpectError: regexp.MustCompile(`(?s)expiration_rules:\s+rule\s+"Expiration_EndUser_Eligibility"\s+does\s+not\s+exist\s+on\s+the\s+policy`),
			},
			{
				Config:      loadUnitTestTerraform("resource_wrong_rule_type.tf"),
				ExpectError: regexp.MustCompile(`(?s)rule\s+"Enablement_Admin_Assignment"\s+is\s+a\s+unifiedRoleManagementPolicyEnablementRule,\s+not\s+a\s+unifiedRoleManagementPolicyExpirationRule`),
			},
			{
				Config:      loadUnitTestTerraform("resource_no_assignment.tf"),
				ExpectError: regexp.MustCompile(`(?s)expected\s+one\s+policy\s+assigned\s+to\s+role\s+definition\s+99999999-9999-9999-9999-999999999999.*found\s+0`),
			},
		},
	})
}

// ConfigValidation: invalid configuration is rejected before any request is sent.
func TestUnitResourceRoleManagementPolicy_09_ConfigValidation(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleManagementPolicyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleManagementPolicyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_invalid_identification.tf"),
				ExpectError: regexp.MustCompile(`(?s)Invalid\s+Attribute\s+Combination`),
			},
			{
				Config:      loadUnitTestTerraform("resource_invalid_duration.tf"),
				ExpectError: regexp.MustCompile(`(?s)invalid\s+ISO\s+8601\s+duration\s+"P2W"`),
			},
			{
				Config:      loadUnitTestTerraform("resource_invalid_claim_missing.tf"),
				ExpectError: regexp.MustCompile(`(?s)Missing\s+claim_value`),
			},
			{
				Config:      loadUnitTestTerraform("resource_invalid_claim_disabled.tf"),
				ExpectError: regexp.MustCompile(`(?s)Unexpected\s+claim_value`),
			},
			{
				Config:      loadUnitTestTerraform("resource_invalid_enablement.tf"),
				ExpectError: regexp.MustCompile(`(?s)value\s+must\s+be\s+one\s+of`),
			},
		},
	})
}

// ErrorHandling: Graph errors surface as diagnostics.
func TestUnitResourceRoleManagementPolicy_10_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, roleManagementPolicyMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer roleManagementPolicyMock.CleanupMockState()

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

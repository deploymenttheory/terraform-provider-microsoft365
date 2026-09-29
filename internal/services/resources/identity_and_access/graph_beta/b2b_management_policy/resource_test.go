package graphBetaIdentityAndAccessB2bManagementPolicy_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaB2bManagementPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/b2b_management_policy"
	b2bManagementPolicyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/b2b_management_policy/mocks"
)

var resourceType = "microsoft365_graph_beta_identity_and_access_b2b_management_policy"

func setupMockEnvironment() (*mocks.Mocks, *b2bManagementPolicyMocks.B2bManagementPolicyMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	policyMock := &b2bManagementPolicyMocks.B2bManagementPolicyMock{}
	policyMock.RegisterMocks()
	return mockClient, policyMock
}

// The mock rejects description with Graph's 404, so this also guards against sending it.
func TestUnitResourceB2bManagementPolicy_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceType+".minimal",
						"display_name",
						"unit-test-b2b-management-policy-min",
					),
					resource.TestCheckResourceAttr(
						resourceType+".minimal",
						"is_organization_default",
						"false",
					),
					resource.TestCheckResourceAttr(resourceType+".minimal", "definition.#", "1"),
					resource.TestCheckResourceAttr(
						resourceType+".minimal",
						"definition.0",
						`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"BlockedDomains":["example.com"]}}}`,
					),
					resource.TestMatchResourceAttr(
						resourceType+".minimal",
						"id",
						regexp.MustCompile(`^[0-9a-fA-F-]+$`),
					),
				),
			},
			{
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_02_Maximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceType+".maximal",
						"display_name",
						"unit-test-b2b-management-policy-max",
					),
					resource.TestCheckResourceAttr(
						resourceType+".maximal",
						"is_organization_default",
						"true",
					),
					resource.TestCheckResourceAttr(resourceType+".maximal", "definition.#", "1"),
					resource.TestMatchResourceAttr(
						resourceType+".maximal",
						"id",
						regexp.MustCompile(`^[0-9a-fA-F-]+$`),
					),
				),
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_03_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceType+".minimal",
						"display_name",
						"unit-test-b2b-management-policy-min",
					),
				),
			},
			{
				Config: testConfig("tests/terraform/unit/resource_minimal_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceType+".minimal",
						"display_name",
						"unit-test-b2b-management-policy-min-updated",
					),
					resource.TestCheckResourceAttr(
						resourceType+".minimal",
						"definition.0",
						`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["example.net"]}}}`,
					),
				),
			},
		},
	})
}

// Replicas return the pre-PATCH values for a few reads; the consistency predicate must wait them out.
func TestUnitResourceB2bManagementPolicy_04_EventualConsistency(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	policyMock.RegisterEventualConsistencyMocks(2)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_minimal.tf"),
			},
			{
				Config: testConfig("tests/terraform/unit/resource_minimal_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceType+".minimal",
						"display_name",
						"unit-test-b2b-management-policy-min-updated",
					),
				),
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_05_InvalidDefinition(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfig("tests/terraform/unit/resource_invalid_definition.tf"),
				ExpectError: regexp.MustCompile(`Invalid JSON String`),
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_06_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	policyMock.RegisterErrorMocks()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfig("tests/terraform/unit/resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`Property definition has an invalid value`),
			},
		},
	})
}

// A stale-replica 404 on refresh must not drop the policy from state (seen live: duplicate create).
func TestUnitResourceB2bManagementPolicy_07_StaleNotFoundOnRefresh(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_minimal.tf"),
			},
			{
				PreConfig: func() {
					policyMock.RegisterStaleNotFoundMocks(2)
				},
				Config:   testConfig("tests/terraform/unit/resource_minimal.tf"),
				PlanOnly: true,
			},
		},
	})
}

// A policy that is really gone is dropped once the confirmation window elapses.
func TestUnitResourceB2bManagementPolicy_08_DeletedOutOfBand(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()
	defer graphBetaB2bManagementPolicy.SetNotFoundConfirmationForTesting(
		2*time.Second,
		500*time.Millisecond,
	)()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_minimal.tf"),
			},
			{
				PreConfig: func() {
					policyMock.CleanupMockState()
				},
				Config:             testConfig("tests/terraform/unit/resource_minimal.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// Graph returns 500 for any change to isOrganizationDefault, so it must force replacement.
func TestUnitResourceB2bManagementPolicy_09_OrganizationDefaultRequiresReplace(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_minimal.tf"),
			},
			{
				Config: testConfig("tests/terraform/unit/resource_minimal_organization_default.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".minimal",
							plancheck.ResourceActionReplace,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceType+".minimal",
						"is_organization_default",
						"true",
					),
				),
			},
		},
	})
}

func testConfig(path string) string {
	unitTestConfig, err := helpers.ParseHCLFile(path)
	if err != nil {
		panic("failed to load B2B management policy config " + path + ": " + err.Error())
	}
	return unitTestConfig
}

func TestUnitResourceB2bManagementPolicy_10_AllowDomains(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mock.CleanupMockState()
	var id string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_20_allow_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_11_BlockDomains(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mock.CleanupMockState()
	var id string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_21_block_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"BlockedDomains":["blocked.example","untrusted.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_12_AnyDomain(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mock.CleanupMockState()
	var id string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_22_any_domain.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_13_AllowToBlock(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mock.CleanupMockState()
	var id string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_20_allow_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				Config: testConfig("tests/terraform/unit/resource_21_block_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"BlockedDomains":["blocked.example","untrusted.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_14_BlockToAllow(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mock.CleanupMockState()
	var id string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_21_block_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"BlockedDomains":["blocked.example","untrusted.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				Config: testConfig("tests/terraform/unit/resource_20_allow_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_15_RestrictedToAny(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mock.CleanupMockState()
	var id string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_20_allow_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				Config: testConfig("tests/terraform/unit/resource_22_any_domain.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_16_AnyToRestricted(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mock.CleanupMockState()
	var id string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_22_any_domain.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				Config: testConfig("tests/terraform/unit/resource_21_block_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"BlockedDomains":["blocked.example","untrusted.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_17_MinimalToMaximal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mock.CleanupMockState()
	var id string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_20_allow_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				Config: testConfig("tests/terraform/unit/resource_23_expanded_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example","partner.example","sub.partner.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_18_MaximalToMinimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mock.CleanupMockState()
	var id string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_23_expanded_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example","partner.example","sub.partner.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				Config: testConfig("tests/terraform/unit/resource_20_allow_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_19_OtherSettings(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mock.CleanupMockState()
	var id string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_24_other_settings.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"AutoRedeemPolicy":{"AdminConsentedForUsersIntoTenantIds":[],"NoAADConsentForUsersFromTenantsIds":[]},"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example"]},"PreviewPolicy":{"Features":[]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_20_Rename(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mock.CleanupMockState()
	var id string
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig("tests/terraform/unit/resource_20_allow_domains.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				Config: testConfig("tests/terraform/unit/resource_26_renamed.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").
						Key("display_name").
						HasValue("unit-test-b2b-renamed"),
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").
						Key("is_organization_default").
						HasValue("false"),
					check.That(resourceType+".test").
						Key("definition.0").
						HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
					func(s *terraform.State) error {
						current := s.RootModule().Resources[resourceType+".test"].Primary.ID
						if id != "" && current != id {
							return fmt.Errorf(
								"policy ID changed during update: %s to %s",
								id,
								current,
							)
						}
						id = current
						return nil
					},
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_21_FormattedJSON(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: testConfig("tests/terraform/unit/resource_25_formatted_json.tf")},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				Config:   testConfig("tests/terraform/unit/resource_25_formatted_json.tf"),
				PlanOnly: true,
			},
		},
	})
}

func TestUnitResourceB2bManagementPolicy_22_DomainDrift(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: testConfig("tests/terraform/unit/resource_20_allow_domains.tf")},
			{PreConfig: func() {
				policyMock.ChangeDefinition(
					`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["changed.example"]}}}`,
				)
			}, Config: testConfig("tests/terraform/unit/resource_20_allow_domains.tf"), PlanOnly: true, ExpectNonEmptyPlan: true},
			{
				Config: testConfig("tests/terraform/unit/resource_20_allow_domains.tf"),
				Check: check.That(resourceType + ".test").
					Key("definition.0").
					HasValue(`{"B2BManagementPolicy":{"InvitationsAllowedAndBlockedDomainsPolicy":{"AllowedDomains":["contoso.example","fabrikam.example"]}}}`),
			},
		},
	})
}

// Deactivating a default policy also requires replacement; cover the reverse transition.
func TestUnitResourceB2bManagementPolicy_23_DefaultToNonDefaultRequiresReplace(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: testConfig("tests/terraform/unit/resource_minimal_organization_default.tf")},
			{
				Config: testConfig("tests/terraform/unit/resource_minimal.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".minimal",
							plancheck.ResourceActionReplace,
						),
					},
				},
				Check: check.That(resourceType + ".minimal").
					Key("is_organization_default").
					HasValue("false"),
			},
		},
	})
}

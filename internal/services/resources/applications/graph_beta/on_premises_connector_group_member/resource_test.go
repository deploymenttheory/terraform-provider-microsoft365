package graphBetaApplicationsOnPremisesConnectorGroupMember_test

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
	connectorGroupMemberMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/on_premises_connector_group_member/mocks"
)

func setupMockEnvironment() (*mocks.Mocks, *connectorGroupMemberMocks.OnPremisesConnectorGroupMemberMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	connectorGroupMemberMock := &connectorGroupMemberMocks.OnPremisesConnectorGroupMemberMock{}
	connectorGroupMemberMock.RegisterMocks()
	return mockClient, connectorGroupMemberMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *connectorGroupMemberMocks.OnPremisesConnectorGroupMemberMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	connectorGroupMemberMock := &connectorGroupMemberMocks.OnPremisesConnectorGroupMemberMock{}
	connectorGroupMemberMock.RegisterErrorMocks()
	return mockClient, connectorGroupMemberMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// Assigns a connector to a group, imports it, and moves it back to the default group on destroy.
func TestUnitResourceOnPremisesConnectorGroupMember_01_Minimal(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorGroupMemberMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorGroupMemberMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy: func(_ *terraform.State) error {
			return expectConnectorGroup(connectorGroupMemberMock, connectorGroupMemberMocks.UnitTestDefaultConnectorGroupID)
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").HasValue(connectorGroupMemberMocks.UnitTestConnectorID),
					check.That(resourceType+".minimal").Key("connector_id").HasValue(connectorGroupMemberMocks.UnitTestConnectorID),
					check.That(resourceType+".minimal").Key("connector_group_id").HasValue(connectorGroupMemberMocks.UnitTestConnectorGroupAID),
					check.That(resourceType+".minimal").Key("connector_group_name").HasValue("Unit Test Connector Group A"),
					check.That(resourceType+".minimal").Key("machine_name").HasValue("unit-test-connector-01.contoso.local"),
					check.That(resourceType+".minimal").Key("external_ip").HasValue("203.0.113.10"),
				),
			},
			{
				ResourceName:            resourceType + ".minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// Moves the connector between groups in place.
func TestUnitResourceOnPremisesConnectorGroupMember_02_MoveBetweenGroups(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorGroupMemberMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorGroupMemberMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("connector_group_id").HasValue(connectorGroupMemberMocks.UnitTestConnectorGroupAID),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_updated.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceType+".minimal", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("id").HasValue(connectorGroupMemberMocks.UnitTestConnectorID),
					check.That(resourceType+".minimal").Key("connector_group_id").HasValue(connectorGroupMemberMocks.UnitTestConnectorGroupBID),
					check.That(resourceType+".minimal").Key("connector_group_name").HasValue("Unit Test Connector Group B"),
					func(_ *terraform.State) error {
						return expectConnectorGroup(connectorGroupMemberMock, connectorGroupMemberMocks.UnitTestConnectorGroupBID)
					},
				),
			},
		},
	})
}

// A move made outside Terraform is reported as drift and corrected on apply.
func TestUnitResourceOnPremisesConnectorGroupMember_03_ExternalMoveIsDrift(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorGroupMemberMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorGroupMemberMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
			},
			{
				PreConfig: func() {
					connectorGroupMemberMock.MoveConnector(connectorGroupMemberMocks.UnitTestConnectorID, connectorGroupMemberMocks.UnitTestConnectorGroupBID)
				},
				Config:             loadUnitTestTerraform("resource_minimal.tf"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("connector_group_id").HasValue(connectorGroupMemberMocks.UnitTestConnectorGroupAID),
					func(_ *terraform.State) error {
						return expectConnectorGroup(connectorGroupMemberMock, connectorGroupMemberMocks.UnitTestConnectorGroupAID)
					},
				),
			},
		},
	})
}

// Destroy refreshes first, so a move made outside Terraform is adopted and the
// connector is moved from its actual group back to the default group.
func TestUnitResourceOnPremisesConnectorGroupMember_04_DestroyAfterExternalMove(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorGroupMemberMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorGroupMemberMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy: func(_ *terraform.State) error {
			return expectConnectorGroup(connectorGroupMemberMock, connectorGroupMemberMocks.UnitTestDefaultConnectorGroupID)
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_minimal.tf"),
			},
			{
				PreConfig: func() {
					connectorGroupMemberMock.MoveConnector(connectorGroupMemberMocks.UnitTestConnectorID, connectorGroupMemberMocks.UnitTestConnectorGroupBID)
				},
				Config:  loadUnitTestTerraform("resource_minimal.tf"),
				Destroy: true,
			},
		},
	})
}

// Managing membership of the default group: destroy only removes state.
func TestUnitResourceOnPremisesConnectorGroupMember_05_DefaultGroup(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorGroupMemberMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorGroupMemberMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy: func(_ *terraform.State) error {
			return expectConnectorGroup(connectorGroupMemberMock, connectorGroupMemberMocks.UnitTestDefaultConnectorGroupID)
		},
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_default_group.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("connector_group_id").HasValue(connectorGroupMemberMocks.UnitTestDefaultConnectorGroupID),
					check.That(resourceType+".minimal").Key("connector_group_name").HasValue("Default"),
				),
			},
		},
	})
}

// A connector that does not exist surfaces the Graph 404 error.
func TestUnitResourceOnPremisesConnectorGroupMember_06_ConnectorNotFound(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorGroupMemberMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorGroupMemberMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_missing_connector.tf"),
				ExpectError: regexp.MustCompile(`(?i)not found`),
			},
		},
	})
}

// A 403 from Graph surfaces the permission error with the resource's write permissions.
func TestUnitResourceOnPremisesConnectorGroupMember_07_PermissionError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorGroupMemberMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorGroupMemberMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_minimal.tf"),
				ExpectError: regexp.MustCompile(`Directory\.ReadWrite\.All`),
			},
		},
	})
}

// Invalid GUIDs are rejected at plan time.
func TestUnitResourceOnPremisesConnectorGroupMember_08_InvalidGUID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, connectorGroupMemberMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer connectorGroupMemberMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "microsoft365_graph_beta_applications_on_premises_connector_group_member" "minimal" {
  connector_id       = "not-a-guid"
  connector_group_id = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
}`,
				ExpectError: regexp.MustCompile(`must be a valid GUID`),
			},
		},
	})
}

func expectConnectorGroup(connectorGroupMemberMock *connectorGroupMemberMocks.OnPremisesConnectorGroupMemberMock, expected string) error {
	if got := connectorGroupMemberMock.CurrentConnectorGroup(connectorGroupMemberMocks.UnitTestConnectorID); got != expected {
		return fmt.Errorf("connector group = %q, expected %q", got, expected)
	}
	return nil
}

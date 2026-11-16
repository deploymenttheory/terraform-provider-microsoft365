package graphBetaApplicationsOnPremisesConnectorGroupMember_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaApplicationsOnPremisesConnectorGroupMember "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/on_premises_connector_group_member"
)

// connectorIDEnvVar names an installed private network connector in the test
// tenant. Microsoft Graph cannot create connectors; they are registered by
// installing the connector agent on a server.
const connectorIDEnvVar = "M365_TEST_APP_PROXY_CONNECTOR_ID"

var (
	// Resource type name from the resource package
	resourceType = graphBetaApplicationsOnPremisesConnectorGroupMember.ResourceName

	// testResource is the test resource implementation for connector group memberships
	testResource = graphBetaApplicationsOnPremisesConnectorGroupMember.OnPremisesConnectorGroupMemberTestResource{}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// Moves a real connector into a test group, then into a second test group, and
// restores the connector's original group afterwards. Moving a connector can
// break published applications that rely on its current group, so the original
// group is snapshotted before the test and restored in CheckDestroy and cleanup.
func TestAccResourceOnPremisesConnectorGroupMember_01_Lifecycle(t *testing.T) {
	connectorID := os.Getenv(connectorIDEnvVar)
	if connectorID == "" {
		t.Skipf("%s not set, skipping acceptance test: the tenant needs an installed private network connector", connectorIDEnvVar)
	}
	t.Setenv("TF_VAR_connector_id", connectorID)

	var originalGroupID string

	t.Cleanup(func() {
		if originalGroupID != "" {
			if err := restoreConnectorGroup(connectorID, originalGroupID); err != nil {
				t.Errorf("failed to restore connector %s to connector group %s: %v", connectorID, originalGroupID, err)
			}
		}
	})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			groupID, err := currentConnectorGroupID(connectorID)
			if err != nil {
				t.Fatalf("failed to snapshot connector group for connector %s: %v", connectorID, err)
			}
			if groupID == "" {
				t.Fatalf("connector %s is not a member of any connector group", connectorID)
			}
			originalGroupID = groupID
			testlog.StepAction(resourceType, fmt.Sprintf("Snapshotted original connector group %s", originalGroupID))
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: func(state *terraform.State) error {
			if err := destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second)(state); err != nil {
				return err
			}
			testlog.StepAction(resourceType, fmt.Sprintf("Restoring original connector group %s", originalGroupID))
			return restoreConnectorGroup(connectorID, originalGroupID)
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Assigning connector to first test connector group")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").HasValue(connectorID),
					check.That(resourceType+".test").Key("connector_id").HasValue(connectorID),
					resource.TestCheckResourceAttrPair(resourceType+".test", "connector_group_id", "microsoft365_graph_beta_applications_on_premises_connector_group.first", "id"),
					resource.TestCheckResourceAttrPair(resourceType+".test", "connector_group_name", "microsoft365_graph_beta_applications_on_premises_connector_group.first", "name"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing connector group membership")
				},
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Moving connector to second test connector group")
				},
				Config: loadAcceptanceTestTerraform("resource_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").Key("id").HasValue(connectorID),
					resource.TestCheckResourceAttrPair(resourceType+".test", "connector_group_id", "microsoft365_graph_beta_applications_on_premises_connector_group.second", "id"),
					resource.TestCheckResourceAttrPair(resourceType+".test", "connector_group_name", "microsoft365_graph_beta_applications_on_premises_connector_group.second", "name"),
				),
			},
		},
	})
}

// currentConnectorGroupID returns the ID of the group the connector belongs to.
func currentConnectorGroupID(connectorID string) (string, error) {
	client, err := acceptance.TestGraphClient()
	if err != nil {
		return "", err
	}

	memberOf, err := client.OnPremisesPublishingProfiles().ByOnPremisesPublishingProfileId("applicationProxy").Connectors().ByConnectorId(connectorID).MemberOf().Get(context.Background(), nil)
	if err != nil {
		return "", err
	}

	if memberOf == nil || len(memberOf.GetValue()) == 0 || memberOf.GetValue()[0].GetId() == nil {
		return "", nil
	}

	return *memberOf.GetValue()[0].GetId(), nil
}

// restoreConnectorGroup moves the connector back to its original group when it
// is not already there.
func restoreConnectorGroup(connectorID, originalGroupID string) error {
	current, err := currentConnectorGroupID(connectorID)
	if err != nil {
		return err
	}
	if strings.EqualFold(current, originalGroupID) {
		return nil
	}

	client, err := acceptance.TestGraphClient()
	if err != nil {
		return err
	}

	odataID := fmt.Sprintf("%s/onPremisesPublishingProfiles/applicationProxy/connectorGroups/%s", strings.TrimRight(client.GetAdapter().GetBaseUrl(), "/"), originalGroupID)
	body := graphmodels.NewReferenceCreate()
	body.SetOdataId(&odataID)

	if err := client.OnPremisesPublishingProfiles().ByOnPremisesPublishingProfileId("applicationProxy").Connectors().ByConnectorId(connectorID).MemberOf().Ref().Post(context.Background(), body, nil); err != nil {
		return err
	}

	restored, err := currentConnectorGroupID(connectorID)
	if err != nil {
		return err
	}
	if !strings.EqualFold(restored, originalGroupID) {
		return fmt.Errorf("connector %s reports connector group %s after restore, expected %s", connectorID, restored, originalGroupID)
	}

	return nil
}

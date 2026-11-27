package graphBetaAgentCardManifest_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
)

// This API has no standalone POST or DELETE. Supply a dedicated existing manifest.
func TestAccResourceAgentCardManifest_Lifecycle(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("TF_ACC must be set to 1")
	}
	id := os.Getenv("M365_AGENT_CARD_MANIFEST_TEST_ID")
	if id == "" {
		t.Skip("Set M365_AGENT_CARD_MANIFEST_TEST_ID to a dedicated existing manifest; Graph provides no standalone manifest POST or DELETE")
	}
	mocks.TestAccPreCheck(t)
	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	original, err := client.AgentRegistry().AgentCardManifests().ByAgentCardManifestId(id).Get(ctx, nil)
	cancel()
	require.NoError(t, err)
	require.NotNil(t, original)
	require.NotNil(t, original.GetDisplayName(), "The fixture must have a display name to restore")
	savedName := *original.GetDisplayName()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		body := graphmodels.NewAgentCardManifest()
		body.SetDisplayName(&savedName)
		if _, err := client.AgentRegistry().AgentCardManifests().ByAgentCardManifestId(id).Patch(ctx, body, nil); err != nil {
			t.Errorf("restore manifest display name: %v", err)
		}
	})
	name := fmt.Sprintf("acctest-manifest-%d", time.Now().UnixNano())
	first := "provider \"microsoft365\" {}\n" + manifestConfig(id, name, "")
	updated := "provider \"microsoft365\" {}\n" + manifestConfig(id, name+"-updated", "")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: func(_ *terraform.State) error {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			remote, err := client.AgentRegistry().AgentCardManifests().ByAgentCardManifestId(id).Get(ctx, nil)
			if err != nil {
				return fmt.Errorf("manifest must survive Terraform destroy: %w", err)
			}
			if remote == nil || remote.GetId() == nil {
				return fmt.Errorf("manifest missing after Terraform destroy")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{Config: first, Check: resource.TestCheckResourceAttr(manifestType+".test", "display_name", name)},
			{Config: updated, Check: resource.TestCheckResourceAttr(manifestType+".test", "display_name", name+"-updated")},
			{Config: updated, PlanOnly: true},
			{ResourceName: manifestType + ".test", ImportState: true, ImportStateId: id, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"timeouts"}},
		},
	})
}

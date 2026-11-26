package graphBetaApplicationsTokenLifetimePolicy_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Graph exposes description through an inherited SDK model but rejects it on
// both POST and PATCH. Reject configuration before sending that invalid write.
func TestUnitResourceTokenLifetimePolicy_DescriptionIsReadOnly(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: `resource "microsoft365_graph_beta_applications_token_lifetime_policy" "unsupported" {
  display_name = "unsupported-description"
  description = "Graph rejects this field"
  definition = ["{\"TokenLifetimePolicy\":{\"Version\":1,\"AccessTokenLifetime\":\"01:00:00\"}}"]
}`,
			ExpectError: regexp.MustCompile(`(?s)(Invalid Configuration for Read-Only Attribute|Value for unconfigurable attribute)`),
		}},
	})
}

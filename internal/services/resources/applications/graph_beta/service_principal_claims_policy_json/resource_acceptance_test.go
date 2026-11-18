package graphBetaServicePrincipalClaimsPolicyJSON_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_claims_policy_json"
)

func TestAccResourceServicePrincipalClaimsPolicyJSON_Lifecycle(t *testing.T) {
	parent := claimsPolicyAcceptancePrincipal(t)
	minimal, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	require.NoError(t, err)
	complex, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_polymorphic.tf")
	require.NoError(t, err)
	config := fmt.Sprintf(minimal, parent)
	updated := strings.Replace(config, "fixture-value-1", "fixture-value-2", 1)
	address := resourceType + ".test"
	expected, err := helpers.ParseJSONFile("tests/responses/get_policy.json")
	require.NoError(t, err)
	resource.Test(t, resource.TestCase{PreCheck: func() { mocks.TestAccPreCheck(t) }, ProtoV6ProviderFactories: testAccProtoV6ProviderFactories, Steps: []resource.TestStep{
		{Config: config, Check: resource.TestCheckResourceAttr(address, "id", parent)},
		{ResourceName: address, ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"policy_json"}, ImportStateCheck: checkImportedPolicy(expected)},
		{
			ResourceName:     address,
			ImportState:      true,
			GenerateConfig:   true,
			ImportStateKind:  resource.ImportBlockWithID,
			ImportPlanChecks: resource.ImportPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionNoop)}},
		},
		{Config: updated},
		{Config: fmt.Sprintf(complex, parent)},
		{Config: config},
		{Config: config, PlanOnly: true},
	}})
	// Destroy must retain the custom policy; verify before deleting the owned parent.
	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	request, err := client.ServicePrincipals().ByServicePrincipalId(parent).ClaimsPolicy().ToGetRequestInformation(ctx, nil)
	require.NoError(t, err)
	remote, err := client.GetAdapter().SendPrimitive(ctx, request, "[]byte", abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue})
	require.NoError(t, err)
	content, ok := remote.([]byte)
	require.True(t, ok)
	equal, diags := policy.NewPolicyJSONValue(expected).StringSemanticEquals(ctx, policy.NewPolicyJSONValue(string(content)))
	require.False(t, diags.HasError())
	require.True(t, equal, "state-only destroy must retain the final configured policy")
}

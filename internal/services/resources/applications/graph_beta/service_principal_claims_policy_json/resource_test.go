package graphBetaServicePrincipalClaimsPolicyJSON_test

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_claims_policy_json"
	policymocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/service_principal_claims_policy_json/mocks"
)

var resourceType = policy.ResourceName

var errPolicyImport = errors.New("imported claims policy differs semantically from the configured policy")

func checkImportedPolicy(expected string) func([]*terraform.InstanceState) error {
	return func(states []*terraform.InstanceState) error {
		if len(states) != 1 {
			return errPolicyImport
		}
		equal, diags := policy.NewPolicyJSONValue(expected).StringSemanticEquals(context.Background(), policy.NewPolicyJSONValue(states[0].Attributes["policy_json"]))
		if diags.HasError() || !equal {
			return errPolicyImport
		}
		return nil
	}
}

func setupMockEnvironment(t *testing.T) *policymocks.ServicePrincipalClaimsPolicyJSONMock {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	mock := &policymocks.ServicePrincipalClaimsPolicyJSONMock{}
	mock.CleanupMockState()
	mock.RegisterMocks()
	t.Cleanup(func() { mock.CleanupMockState(); httpmock.DeactivateAndReset() })
	return mock
}

func loadUnitTestTerraform(t *testing.T, file string) string {
	t.Helper()
	value, err := helpers.ParseHCLFile("tests/terraform/unit/" + file)
	require.NoError(t, err)
	return value
}

func TestUnitResourceServicePrincipalClaimsPolicyJSON_Lifecycle(t *testing.T) {
	mock := setupMockEnvironment(t)
	config := loadUnitTestTerraform(t, "resource_minimal.tf")
	updated := strings.Replace(config, "fixture-value-1", "fixture-value-2", 1)
	address := resourceType + ".test"
	expected, err := helpers.ParseJSONFile("tests/responses/get_policy.json")
	require.NoError(t, err)
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{
		{
			Config: config,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(address, "id", policymocks.ServicePrincipalID),
				resource.TestCheckResourceAttr(address, "policy_id", policymocks.PolicyID),
			),
		},
		// ImportStateVerify compares JSON string bytes; verify the full custom
		// semantic value separately rather than ignoring meaningful policy data.
		{ResourceName: address, ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"policy_json"}, ImportStateCheck: checkImportedPolicy(expected)},
		{Config: updated},
		{Config: updated, PlanOnly: true},
	}})
	puts, deletes := mock.Requests()
	require.Len(t, puts, 2)
	require.Zero(t, deletes, "destroy must not call the unsupported DELETE endpoint")
	require.Equal(t, "#microsoft.graph.customClaim", puts[1]["claims"].([]any)[0].(map[string]any)["@odata.type"])
}

func TestUnitResourceServicePrincipalClaimsPolicyJSON_GeneratedImport(t *testing.T) {
	mock := setupMockEnvironment(t)
	require.NoError(t, mock.SeedFixture())
	address := resourceType + ".test"
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config:                 loadUnitTestTerraform(t, "import.tf"),
					ResourceName:           address,
					ImportState:            true,
					ImportStateId:          policymocks.ServicePrincipalID,
					GenerateConfig:         true,
					ImportStateKind:        resource.ImportBlockWithID,
					ImportStateConfigExact: true,
					ImportPlanChecks:       resource.ImportPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(address, plancheck.ResourceActionNoop)}},
				},
			},
		},
	)
}

func TestUnitResourceServicePrincipalClaimsPolicyJSON_RequiresExplicitFlags(t *testing.T) {
	setupMockEnvironment(t)
	config := regexp.MustCompile(`(?m)^\s*includeBasicClaimSet\s*=\s*true\s*$`).ReplaceAllString(loadUnitTestTerraform(t, "resource_minimal.tf"), "")
	resource.UnitTest(
		t,
		resource.TestCase{
			ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
			Steps:                    []resource.TestStep{{Config: config, ExpectError: regexp.MustCompile("Invalid claims policy JSON")}},
		},
	)
}

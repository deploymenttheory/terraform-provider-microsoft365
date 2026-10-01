package graphBetaAuthenticationFlowsPolicy_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/authentication_flows_policy"
)

var (
	resourceType = policy.ResourceName
	testResource = policy.AuthenticationFlowsPolicyTestResource{}
)

// loadAcceptanceTestTerraform loads a configuration with the shared provider setup.
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic(fmt.Sprintf("failed to load acceptance config %s: %s", filename, err))
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

// TestAccResourceAuthenticationFlowsPolicy_01_Lifecycle exercises the singleton and restores its original setting.
// CheckDestroy is nil because destroy must leave the policy in Microsoft Graph.
func TestAccResourceAuthenticationFlowsPolicy_01_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { authenticationFlowsPolicyPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Enabling self-service sign-up") },
				Config:    loadAcceptanceTestTerraform("resource_enabled.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").ExistsInGraph(testResource),
					check.That(resourceType+".test").
						Key("id").
						HasValue("authenticationFlowsPolicy"),
					check.That(resourceType+".test").
						Key("self_service_sign_up.is_enabled").
						HasValue("true"),
				),
			},
			{
				PreConfig: func() { testlog.StepAction(resourceType, "Disabling self-service sign-up") },
				Config:    loadAcceptanceTestTerraform("resource_disabled.tf"),
				Check: check.That(resourceType + ".test").
					Key("self_service_sign_up.is_enabled").
					HasValue("false"),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateId:           "authenticationFlowsPolicy",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{Config: loadAcceptanceTestTerraform("resource_disabled.tf"), PlanOnly: true},
		},
	})
}

// authenticationFlowsPolicyPreCheck captures the original writable setting before any mutations.
func authenticationFlowsPolicyPreCheck(t *testing.T) {
	t.Helper()
	mocks.TestAccPreCheck(t)
	client, err := acceptance.TestGraphClient()
	require.NoError(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	original, err := client.Policies().AuthenticationFlowsPolicy().Get(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, original)
	require.NotNil(t, original.GetSelfServiceSignUp())
	require.NotNil(
		t,
		original.GetSelfServiceSignUp().GetIsEnabled(),
		"a Boolean baseline is required for reversible testing",
	)
	expected := authenticationFlowsPolicySnapshot(original)
	enabled := *original.GetSelfServiceSignUp().GetIsEnabled()
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cleanupCancel()
		// Build a fresh SDK object so all restore fields are serialized as changes.
		body := graphmodels.NewAuthenticationFlowsPolicy()
		signUp := graphmodels.NewSelfServiceSignUpAuthenticationFlowConfiguration()
		signUp.SetIsEnabled(&enabled)
		body.SetSelfServiceSignUp(signUp)
		_, err := client.Policies().AuthenticationFlowsPolicy().Patch(cleanupCtx, body, nil)
		if err != nil {
			t.Errorf("restore original authentication flows policy: %v", err)
			return
		}
		consecutive := 0
		for {
			current, err := client.Policies().AuthenticationFlowsPolicy().Get(cleanupCtx, nil)
			if err == nil &&
				reflect.DeepEqual(expected, authenticationFlowsPolicySnapshot(current)) {
				consecutive++
				if consecutive >= 3 {
					return
				}
			} else {
				consecutive = 0
			}
			select {
			case <-cleanupCtx.Done():
				t.Error("original authentication flows policy did not verify after cleanup")
				return
			case <-time.After(2 * time.Second):
			}
		}
	})
}

func authenticationFlowsPolicySnapshot(
	value graphmodels.AuthenticationFlowsPolicyable,
) map[string]any {
	if value == nil || value.GetSelfServiceSignUp() == nil {
		return nil
	}
	return map[string]any{
		"id":          value.GetId(),
		"displayName": value.GetDisplayName(),
		"description": value.GetDescription(),
		"isEnabled":   value.GetSelfServiceSignUp().GetIsEnabled(),
	}
}

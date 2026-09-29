package graphBetaMacosDeviceCompliancePolicy

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

// TestUnitConstructScheduledActions verifies create and update send a usable rule
// name even when the optional/computed Terraform value is unknown during planning.
func TestUnitConstructScheduledActions(t *testing.T) {
	ctx := context.Background()
	action := graphmodels.NewDeviceComplianceActionItem()
	action.SetActionType(new(graphmodels.BLOCK_DEVICECOMPLIANCEACTIONTYPE))
	action.SetGracePeriodHours(new(int32(72)))
	configs, err := mapScheduledActionConfigurationsToState(
		ctx,
		[]graphmodels.DeviceComplianceActionItemable{action},
		nil,
	)
	require.NoError(t, err)
	for _, tc := range []struct {
		name  string
		value types.String
		want  string
	}{
		{"omitted", types.StringUnknown(), "PasswordRequired"},
		{"null", types.StringNull(), "PasswordRequired"},
		{"explicit", types.StringValue("PasswordRequired"), "PasswordRequired"},
		{"legacy", types.StringValue("unavailable"), "unavailable"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := ScheduledActionForRuleModel{
				RuleName:                      tc.value,
				ScheduledActionConfigurations: configs,
			}
			created, err := constructScheduledActionsForPolicyCreation(ctx, data)
			require.NoError(t, err)
			updated, err := constructDeviceComplianceScheduledActionForRulesWithPatchMethod(
				ctx,
				data,
			)
			require.NoError(t, err)
			for _, rules := range [][]graphmodels.DeviceComplianceScheduledActionForRuleable{created, updated.GetDeviceComplianceScheduledActionForRules()} {
				require.Len(t, rules, 1)
				require.Equal(t, tc.want, *rules[0].GetRuleName())
				require.Len(t, rules[0].GetScheduledActionConfigurations(), 1)
				require.Equal(
					t,
					int32(72),
					*rules[0].GetScheduledActionConfigurations()[0].GetGracePeriodHours(),
				)
			}
		})
	}
}

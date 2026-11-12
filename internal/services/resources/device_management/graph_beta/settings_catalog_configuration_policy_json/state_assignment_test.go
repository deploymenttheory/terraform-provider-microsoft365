package graphBetaSettingsCatalogConfigurationPolicyJson

import (
	"context"
	"testing"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func TestUnitSettingsCatalogAssignmentsPreserveEmptyState(t *testing.T) {
	assignmentType := SettingsCatalogConfigurationPolicyAssignmentType()
	assigned := types.SetValueMust(assignmentType, []attr.Value{
		types.ObjectValueMust(assignmentType.(types.ObjectType).AttrTypes, map[string]attr.Value{
			"type":     types.StringValue("allDevicesAssignmentTarget"),
			"group_id": types.StringNull(), "filter_id": types.StringNull(), "filter_type": types.StringValue("none"),
		}),
	})
	for _, testCase := range []struct {
		name  string
		prior types.Set
		null  bool
	}{
		{name: "omitted", prior: types.SetNull(assignmentType), null: true},
		{name: "explicit empty", prior: types.SetValueMust(assignmentType, []attr.Value{})},
		{name: "removed assignments", prior: assigned},
		{name: "unknown", prior: types.SetUnknown(assignmentType)},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			for _, remote := range [][]graphmodels.DeviceManagementConfigurationPolicyAssignmentable{nil, {}} {
				model := sharedmodels.SettingsCatalogJsonResourceModel{Assignments: testCase.prior}
				MapAssignmentsToTerraform(context.Background(), &model, remote)
				require.Equal(t, testCase.null, model.Assignments.IsNull())
				require.False(t, model.Assignments.IsUnknown())
				require.Empty(t, model.Assignments.Elements())
			}
		})
	}
}

package graphBetaNetworkForwardingProfile

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
	"github.com/stretchr/testify/require"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/identity_and_access"
)

// The v2 properties are absent from the SDK model, so they must survive the
// SDK's additional-data deserialization.
func TestMapRemoteStateToDataSourceForwardingProfileV2(t *testing.T) {
	raw := `{
	  "id": "5cbd52ab-66d6-4e65-a4af-f11efdffe4c8",
	  "name": "home",
	  "state": "enabled",
	  "trafficForwardingType": "private",
	  "priority": 150,
	  "isCustomProfile": true,
	  "localNetworkSettings": {"allowLocalSubnetBypass": true},
	  "forwardingSourceMatching": {"sessionType": "agent"},
	  "assignmentRules": {
	    "assignments": {
	      "@odata.type": "#microsoft.graph.networkaccess.forwardingProfileAssignmentsEnumeratedMembership",
	      "users": [],
	      "groups": [{"id": "a6ed41d3-45ac-4e12-8601-bdb56216e84e"}],
	      "devices": []
	    },
	    "devicePropertyFilters": {"platforms": "windows,macOS"}
	  }
	}`
	node, err := jsonserialization.NewJsonParseNode([]byte(raw))
	require.NoError(t, err)
	parsed, err := node.GetObjectValue(models.CreateForwardingProfileFromDiscriminatorValue)
	require.NoError(t, err)

	ctx := context.Background()
	mapped := MapRemoteStateToDataSource(ctx, parsed.(models.ForwardingProfileable))

	var local sharedmodels.LocalNetworkSettingsModel
	require.False(
		t,
		mapped.LocalNetworkSettings.As(ctx, &local, basetypes.ObjectAsOptions{}).HasError(),
	)
	require.True(t, local.AllowLocalSubnetBypass.ValueBool())

	var source sharedmodels.ForwardingSourceMatchingModel
	require.False(
		t,
		mapped.ForwardingSourceMatching.As(ctx, &source, basetypes.ObjectAsOptions{}).HasError(),
	)
	require.Equal(t, "agent", source.SessionType.ValueString())

	var rules sharedmodels.AssignmentRulesModel
	require.False(t, mapped.AssignmentRules.As(ctx, &rules, basetypes.ObjectAsOptions{}).HasError())
	require.Equal(t, "enumerated", rules.Membership.ValueString())
	require.True(t, rules.UserIDs.IsNull())
	require.Len(t, rules.GroupIDs.Elements(), 1)
	require.Len(t, rules.Platforms.Elements(), 2)
}

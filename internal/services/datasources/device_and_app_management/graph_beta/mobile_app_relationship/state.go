package graphBetaMobileAppRelationship

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a Graph relationship onto the data source model,
// keeping the lookup inputs the user configured.
func MapRemoteStateToDataSource(
	ctx context.Context,
	data graphmodels.MobileAppRelationshipable,
	config MobileAppRelationshipDataSourceModel,
) MobileAppRelationshipDataSourceModel {
	tflog.Debug(ctx, "Mapping mobile app relationship to data source state", map[string]any{
		"relationshipId": data.GetId(),
	})

	state := MobileAppRelationshipDataSourceModel{
		ID:                         convert.GraphToFrameworkString(data.GetId()),
		SourceID:                   convert.GraphToFrameworkString(data.GetSourceId()),
		RelationshipID:             convert.GraphToFrameworkString(data.GetId()),
		TargetID:                   convert.GraphToFrameworkString(data.GetTargetId()),
		TargetDisplayName:          convert.GraphToFrameworkString(data.GetTargetDisplayName()),
		ODataQuery:                 config.ODataQuery,
		TargetType:                 convert.GraphToFrameworkEnum(data.GetTargetType()),
		RelationshipType:           types.StringNull(),
		TargetDisplayVersion:       convert.GraphToFrameworkString(data.GetTargetDisplayVersion()),
		TargetPublisher:            convert.GraphToFrameworkString(data.GetTargetPublisher()),
		TargetPublisherDisplayName: convert.GraphToFrameworkString(data.GetTargetPublisherDisplayName()),
		SourceDisplayName:          convert.GraphToFrameworkString(data.GetSourceDisplayName()),
		SourceDisplayVersion:       convert.GraphToFrameworkString(data.GetSourceDisplayVersion()),
		SourcePublisherDisplayName: convert.GraphToFrameworkString(data.GetSourcePublisherDisplayName()),
		DependencyType:             types.StringNull(),
		DependentAppCount:          types.Int32Null(),
		DependsOnAppCount:          types.Int32Null(),
		SupersedenceType:           types.StringNull(),
		SupersededAppCount:         types.Int32Null(),
		SupersedingAppCount:        types.Int32Null(),
		Timeouts:                   config.Timeouts,
	}

	// Configured lookup values must survive verbatim; Graph may return different casing.
	if !config.SourceID.IsNull() && !config.SourceID.IsUnknown() {
		state.SourceID = config.SourceID
	}
	if !config.RelationshipID.IsNull() && !config.RelationshipID.IsUnknown() {
		state.RelationshipID = config.RelationshipID
	}
	if !config.TargetID.IsNull() && !config.TargetID.IsUnknown() {
		state.TargetID = config.TargetID
	}
	if !config.TargetDisplayName.IsNull() && !config.TargetDisplayName.IsUnknown() {
		state.TargetDisplayName = config.TargetDisplayName
	}

	switch relationship := data.(type) {
	case graphmodels.MobileAppDependencyable:
		state.RelationshipType = types.StringValue("dependency")
		state.DependencyType = convert.GraphToFrameworkEnum(relationship.GetDependencyType())
		state.DependentAppCount = convert.GraphToFrameworkInt32(relationship.GetDependentAppCount())
		state.DependsOnAppCount = convert.GraphToFrameworkInt32(relationship.GetDependsOnAppCount())
	case graphmodels.MobileAppSupersedenceable:
		state.RelationshipType = types.StringValue("supersedence")
		state.SupersedenceType = convert.GraphToFrameworkEnum(relationship.GetSupersedenceType())
		state.SupersededAppCount = convert.GraphToFrameworkInt32(relationship.GetSupersededAppCount())
		state.SupersedingAppCount = convert.GraphToFrameworkInt32(relationship.GetSupersedingAppCount())
	}

	return state
}

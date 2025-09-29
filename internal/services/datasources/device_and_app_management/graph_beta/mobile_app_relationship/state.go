package graphBetaMobileAppRelationship

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a mobile app relationship to a model
func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.MobileAppRelationshipable) MobileAppRelationshipModel {
	tflog.Debug(ctx, "Starting to map remote resource state to Terraform state", map[string]any{
		"resourceId": data.GetId(),
	})

	model := MobileAppRelationshipModel{
		ID:                         convert.GraphToFrameworkString(data.GetId()),
		TargetID:                   convert.GraphToFrameworkString(data.GetTargetId()),
		TargetDisplayName:          convert.GraphToFrameworkString(data.GetTargetDisplayName()),
		TargetDisplayVersion:       convert.GraphToFrameworkString(data.GetTargetDisplayVersion()),
		TargetPublisher:            convert.GraphToFrameworkString(data.GetTargetPublisher()),
		TargetPublisherDisplayName: convert.GraphToFrameworkString(data.GetTargetPublisherDisplayName()),
		SourceID:                   convert.GraphToFrameworkString(data.GetSourceId()),
		SourceDisplayName:          convert.GraphToFrameworkString(data.GetSourceDisplayName()),
		SourceDisplayVersion:       convert.GraphToFrameworkString(data.GetSourceDisplayVersion()),
		SourcePublisherDisplayName: convert.GraphToFrameworkString(data.GetSourcePublisherDisplayName()),
		TargetType:                 convert.GraphToFrameworkEnum(data.GetTargetType()),
	}

	return model
}

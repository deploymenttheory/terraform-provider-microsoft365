package graphBetaMobileAppSupersedence

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *MobileAppSupersedenceResourceModel) (graphmodels.MobileAppSupersedenceable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewMobileAppSupersedence()

	// Set required fields
	convert.FrameworkToGraphString(data.SourceID, requestBody.SetSourceId)
	convert.FrameworkToGraphString(data.TargetID, requestBody.SetTargetId)

	// Set optional fields
	convert.FrameworkToGraphString(data.TargetDisplayName, requestBody.SetTargetDisplayName)
	convert.FrameworkToGraphString(data.TargetDisplayVersion, requestBody.SetTargetDisplayVersion)
	convert.FrameworkToGraphString(data.TargetPublisher, requestBody.SetTargetPublisher)
	convert.FrameworkToGraphString(data.TargetPublisherDisplayName, requestBody.SetTargetPublisherDisplayName)
	convert.FrameworkToGraphString(data.SourceDisplayName, requestBody.SetSourceDisplayName)
	convert.FrameworkToGraphString(data.SourceDisplayVersion, requestBody.SetSourceDisplayVersion)
	convert.FrameworkToGraphString(data.SourcePublisherDisplayName, requestBody.SetSourcePublisherDisplayName)

	// Set supersedence type using the convert helper
	if err := convert.FrameworkToGraphEnum(data.SupersedenceType, graphmodels.ParseMobileAppSupersedenceType, func(v graphmodels.MobileAppSupersedenceType) {
		requestBody.SetSupersedenceType(&v)
	}); err != nil {
		return nil, fmt.Errorf("error setting supersedence type: %v", err)
	}

	if err := convert.FrameworkToGraphEnum(data.AppRelationshipType, graphmodels.ParseMobileAppRelationshipType, func(v graphmodels.MobileAppRelationshipType) {
		requestBody.SetTargetType(&v)
	}); err != nil {
		return nil, fmt.Errorf("error setting target type: %v", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

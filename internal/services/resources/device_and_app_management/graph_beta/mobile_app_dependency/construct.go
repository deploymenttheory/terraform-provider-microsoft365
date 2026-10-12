package graphBetaMobileAppDependency

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
)

var errInvalidDependencyType = errors.New("dependency_type must be detect or autoInstall")

// constructResource maps the Terraform schema to the SDK model.
func constructResource(
	ctx context.Context,
	data *MobileAppDependencyResourceModel,
) (graphmodels.MobileAppDependencyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	var dependencyType graphmodels.MobileAppDependencyType
	switch data.DependencyType.ValueString() {
	case "detect":
		dependencyType = graphmodels.DETECT_MOBILEAPPDEPENDENCYTYPE
	case "autoInstall":
		dependencyType = graphmodels.AUTOINSTALL_MOBILEAPPDEPENDENCYTYPE
	default:
		return nil, fmt.Errorf(
			"%w: %q",
			errInvalidDependencyType,
			data.DependencyType.ValueString(),
		)
	}
	requestBody := graphmodels.NewMobileAppDependency()

	// Set required fields; the source app is determined by the action endpoint.
	targetID := data.TargetID.ValueString()
	targetType := graphmodels.CHILD_MOBILEAPPRELATIONSHIPTYPE
	requestBody.SetTargetId(&targetID)
	requestBody.SetTargetType(&targetType)
	requestBody.SetDependencyType(&dependencyType)

	if err := constructors.DebugLogGraphObject(
		ctx,
		fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName),
		requestBody,
	); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

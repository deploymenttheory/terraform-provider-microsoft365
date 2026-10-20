package graphBetaMobileAppSupersedence

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
)

var errInvalidSupersedenceType = errors.New("supersedence_type must be update or replace")

// constructResource maps the Terraform schema to the SDK model
func constructResource(
	ctx context.Context,
	data *MobileAppSupersedenceResourceModel,
) (graphmodels.MobileAppSupersedenceable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	var supersedenceType graphmodels.MobileAppSupersedenceType
	switch data.SupersedenceType.ValueString() {
	case "update":
		supersedenceType = graphmodels.UPDATE_MOBILEAPPSUPERSEDENCETYPE
	case "replace":
		supersedenceType = graphmodels.REPLACE_MOBILEAPPSUPERSEDENCETYPE
	default:
		return nil, fmt.Errorf(
			"%w: %q",
			errInvalidSupersedenceType,
			data.SupersedenceType.ValueString(),
		)
	}
	requestBody := graphmodels.NewMobileAppSupersedence()
	targetID := data.TargetID.ValueString()
	targetType := graphmodels.CHILD_MOBILEAPPRELATIONSHIPTYPE
	requestBody.SetTargetId(&targetID)
	requestBody.SetTargetType(&targetType)
	requestBody.SetSupersedenceType(&supersedenceType)

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

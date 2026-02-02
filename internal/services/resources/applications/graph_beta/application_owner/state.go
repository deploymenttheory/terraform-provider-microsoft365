package graphBetaApplicationOwner

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote state from Microsoft Graph API to the Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data *ApplicationOwnerResourceModel, ownerObject graphmodels.DirectoryObjectable) {
	tflog.Debug(ctx, fmt.Sprintf("Mapping %s remote state to Terraform state", ResourceName))

	// If no owner object is found, the owner assignment doesn't exist
	if ownerObject == nil {
		tflog.Debug(ctx, "No owner object found - setting all computed fields to empty")
		data.OwnerType = types.StringValue("")
		data.OwnerDisplayName = types.StringValue("")
		return
	}

	// Get the @odata.type to determine the actual object type
	odataType := ownerObject.GetOdataType()
	if odataType != nil {
		switch *odataType {
		case "#microsoft.graph.user":
			data.OwnerType = types.StringValue("User")
			if userObj, ok := ownerObject.(graphmodels.Userable); ok {
				if displayName := userObj.GetDisplayName(); displayName != nil {
					data.OwnerDisplayName = convert.GraphToFrameworkString(displayName)
				} else {
					data.OwnerDisplayName = types.StringValue("")
				}
			}
		case "#microsoft.graph.servicePrincipal":
			data.OwnerType = types.StringValue("ServicePrincipal")
			if spObj, ok := ownerObject.(graphmodels.ServicePrincipalable); ok {
				if displayName := spObj.GetDisplayName(); displayName != nil {
					data.OwnerDisplayName = convert.GraphToFrameworkString(displayName)
				} else {
					data.OwnerDisplayName = types.StringValue("")
				}
			}
		default:
			tflog.Warn(ctx, fmt.Sprintf("Unknown owner object type: %s", *odataType))
			data.OwnerType = types.StringValue("Unknown")
			data.OwnerDisplayName = types.StringValue("")
		}
	} else {
		data.OwnerType = types.StringValue("Unknown")
		data.OwnerDisplayName = types.StringValue("")
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s remote state to Terraform state", ResourceName))
}

package graphBetaApplicationIdentifierUri

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the application response to the Terraform state
// For application identifier URIs, this is minimal since we only verify the application exists
func MapRemoteResourceStateToTerraform(ctx context.Context, data *ApplicationIdentifierUriResourceModel, application graphmodels.Applicationable) {
	if application == nil {
		tflog.Warn(ctx, "Received nil application in MapRemoteResourceStateToTerraform")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting to map remote state to Terraform state for %s", ResourceName))

	// For identifier URIs, we just maintain the state as-is
	// The identifier URI presence is validated during read

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s with identifier_uri: %s", ResourceName, data.IdentifierUri.ValueString()))
}

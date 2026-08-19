package graphBetaServicePrincipalPreferredTokenSigningKeyThumbprint

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the service principal response to the Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *ServicePrincipalPreferredTokenSigningKeyThumbprintResourceModel, servicePrincipal graphmodels.ServicePrincipalable) {
	if servicePrincipal == nil {
		tflog.Warn(ctx, "Received nil service principal in MapRemoteResourceStateToTerraform")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting to map remote state to Terraform state for %s", ResourceName))

	// Graph may return the thumbprint with different casing than configured; keep the
	// configured casing when the values match case-insensitively to avoid a perpetual diff.
	remote := servicePrincipal.GetPreferredTokenSigningKeyThumbprint()
	if remote == nil || !strings.EqualFold(*remote, data.Thumbprint.ValueString()) {
		data.Thumbprint = types.StringPointerValue(remote)
	}

	data.Id = data.ServicePrincipalID

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s for service principal: %s", ResourceName, data.ServicePrincipalID.ValueString()))
}

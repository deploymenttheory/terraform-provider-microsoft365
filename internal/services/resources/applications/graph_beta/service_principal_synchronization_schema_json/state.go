package graphBetaServicePrincipalSynchronizationSchemaJson

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the synchronization schema returned by Graph to the Terraform model.
//
// The collections are stored in the form Terraform's jsonencode() produces. The SchemaJSONValue semantic equality then
// keeps the configured value in state when Graph only differs in formatting or the properties it manages itself.
func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data ServicePrincipalSynchronizationSchemaJsonResourceModel,
	remoteResource *synchronizationSchemaResponse,
) (ServicePrincipalSynchronizationSchemaJsonResourceModel, error) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data, errNoSynchronizationSchema
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting to map remote state to Terraform state for %s", ResourceName))

	synchronizationRules, err := canonicalJSONArray(remoteResource.SynchronizationRules)
	if err != nil {
		return data, fmt.Errorf("failed to map synchronizationRules: %w", err)
	}

	directories, err := canonicalJSONArray(remoteResource.Directories)
	if err != nil {
		return data, fmt.Errorf("failed to map directories: %w", err)
	}

	data.ID = types.StringValue(resourceId(data.ServicePrincipalId.ValueString(), data.JobId.ValueString()))
	data.SynchronizationRulesJson = NewSchemaJSONValue(synchronizationRules)
	data.DirectoriesJson = NewSchemaJSONValue(directories)
	data.Version = convert.GraphToFrameworkString(remoteResource.Version)

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform state for %s", ResourceName))

	return data, nil
}

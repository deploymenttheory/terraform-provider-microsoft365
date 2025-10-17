package graphBetaCustomSecurityAttributeDefinition

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteResourceStateToTerraform(ctx context.Context, data *CustomSecurityAttributeDefinitionResourceModel, remoteResource graphmodels.CustomSecurityAttributeDefinitionable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	attributeSet := convert.GraphToFrameworkString(remoteResource.GetAttributeSet()).ValueString()
	name := convert.GraphToFrameworkString(remoteResource.GetName()).ValueString()
	data.ID = types.StringValue(fmt.Sprintf("%s_%s", attributeSet, name))

	data.AttributeSet = convert.GraphToFrameworkString(remoteResource.GetAttributeSet())
	data.Name = convert.GraphToFrameworkString(remoteResource.GetName())
	data.IsCollection = convert.GraphToFrameworkBool(remoteResource.GetIsCollection())
	data.IsSearchable = convert.GraphToFrameworkBool(remoteResource.GetIsSearchable())
	data.Status = convert.GraphToFrameworkString(remoteResource.GetStatus())
	data.Type = convert.GraphToFrameworkString(remoteResource.GetTypeEscaped())
	data.UsePreDefinedValuesOnly = convert.GraphToFrameworkBool(remoteResource.GetUsePreDefinedValuesOnly())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

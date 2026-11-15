package graphBetaIdentityAndAccessLifecycleWorkflowSettings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	identitygovernance "github.com/microsoftgraph/msgraph-beta-sdk-go/models/identitygovernance"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the remote lifecycle workflow settings from Kiota SDK to Terraform state.
// Graph does not return an id for this singleton, so the fixed singleton ID is always used.
func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data LifecycleWorkflowSettingsResourceModel,
	remoteResource identitygovernance.LifecycleManagementSettingsable,
) LifecycleWorkflowSettingsResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": singletonID,
	})

	data.ID = types.StringValue(singletonID)
	data.WorkflowScheduleIntervalInHours = convert.GraphToFrameworkInt32(remoteResource.GetWorkflowScheduleIntervalInHours())
	data.EmailSettings = mapEmailSettings(remoteResource.GetEmailSettings())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
	return data
}

// mapEmailSettings maps the email notification settings.
func mapEmailSettings(remote graphmodels.EmailSettingsable) types.Object {
	if remote == nil {
		return types.ObjectNull(emailSettingsAttrTypes)
	}
	return types.ObjectValueMust(emailSettingsAttrTypes, map[string]attr.Value{
		"sender_domain":        convert.GraphToFrameworkString(remote.GetSenderDomain()),
		"use_company_branding": convert.GraphToFrameworkBool(remote.GetUseCompanyBranding()),
	})
}

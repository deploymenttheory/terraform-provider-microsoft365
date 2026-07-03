package graphBetaMacOSDeviceEnrollmentPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// mapResourceToState maps the base DeviceManagementConfigurationPolicy fields to Terraform state.
// dep_onboarding_settings_id is intentionally left untouched here: Graph does not return the
// creationSource used to build it on GET, so the value already present in state/plan is preserved.
func mapResourceToState(
	ctx context.Context,
	stateModel *MacOSDeviceEnrollmentPolicyResourceModel,
	resource models.DeviceManagementConfigurationPolicyable,
) {
	if resource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(resource.GetId()),
	})

	stateModel.ID = convert.GraphToFrameworkString(resource.GetId())
	stateModel.Name = convert.GraphToFrameworkString(resource.GetName())
	stateModel.Description = convert.GraphToFrameworkString(resource.GetDescription())
	stateModel.IsAssigned = convert.GraphToFrameworkBool(resource.GetIsAssigned())
	stateModel.CreatedDateTime = convert.GraphToFrameworkTime(resource.GetCreatedDateTime())
	stateModel.LastModifiedDateTime = convert.GraphToFrameworkTime(resource.GetLastModifiedDateTime())
	stateModel.SettingsCount = convert.GraphToFrameworkInt32(resource.GetSettingCount())
	stateModel.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, resource.GetRoleScopeTagIds())

	if platforms := resource.GetPlatforms(); platforms != nil {
		stateModel.Platforms = types.StringValue(platforms.String())
	}

	if technologies := resource.GetTechnologies(); technologies != nil {
		stateModel.Technologies = types.StringValue(technologies.String())
	}

	if templateRef := resource.GetTemplateReference(); templateRef != nil {
		if templateId := templateRef.GetTemplateId(); templateId != nil {
			stateModel.TemplateId = convert.GraphToFrameworkString(templateId)
		}
		if templateFamily := templateRef.GetTemplateFamily(); templateFamily != nil {
			stateModel.TemplateFamily = types.StringValue(templateFamily.String())
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource state with id %s", stateModel.ID.ValueString()))
}

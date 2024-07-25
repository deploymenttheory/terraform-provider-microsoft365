package graphBetaAssignmentFilter

import (
	"context"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapRemoteStateToTerraform maps the remote state from the Graph API to the Terraform resource model.
// It populates the AssignmentFilterResourceModel with data from the DeviceAndAppManagementAssignmentFilterable.
func mapRemoteStateToTerraform(ctx context.Context, data *AssignmentFilterResourceModel, remoteResource models.DeviceAndAppManagementAssignmentFilterable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Mapping ID")
	data.ID = types.StringValue(helpers.StringPtrToString(remoteResource.GetId()))

	tflog.Debug(ctx, "Mapping DisplayName")
	data.DisplayName = types.StringValue(helpers.StringPtrToString(remoteResource.GetDisplayName()))

	tflog.Debug(ctx, "Mapping Description")
	data.Description = types.StringValue(helpers.StringPtrToString(remoteResource.GetDescription()))

	tflog.Debug(ctx, "Mapping Platform")
	if platform := remoteResource.GetPlatform(); platform != nil {
		data.Platform = types.StringValue(platform.String())
	} else {
		data.Platform = types.StringNull()
	}

	tflog.Debug(ctx, "Mapping Rule")
	data.Rule = types.StringValue(helpers.StringPtrToString(remoteResource.GetRule()))

	tflog.Debug(ctx, "Mapping AssignmentFilterManagementType")
	if managementType := remoteResource.GetAssignmentFilterManagementType(); managementType != nil {
		data.AssignmentFilterManagementType = types.StringValue(managementType.String())
	} else {
		data.AssignmentFilterManagementType = types.StringNull()
	}

	tflog.Debug(ctx, "Mapping CreatedDateTime")
	if createdDateTime := remoteResource.GetCreatedDateTime(); createdDateTime != nil {
		data.CreatedDateTime = types.StringValue(createdDateTime.Format(time.RFC3339))
	} else {
		data.CreatedDateTime = types.StringNull()
	}

	tflog.Debug(ctx, "Mapping LastModifiedDateTime")
	if lastModifiedDateTime := remoteResource.GetLastModifiedDateTime(); lastModifiedDateTime != nil {
		data.LastModifiedDateTime = types.StringValue(lastModifiedDateTime.Format(time.RFC3339))
	} else {
		data.LastModifiedDateTime = types.StringNull()
	}

	tflog.Debug(ctx, "Mapping RoleScopeTags")
	if roleScopeTags := remoteResource.GetRoleScopeTags(); len(roleScopeTags) > 0 {
		data.RoleScopeTags = types.ListValueMust(types.StringType, roleScopeTagsToValueSlice(roleScopeTags))
	} else {
		data.RoleScopeTags = types.ListNull(types.StringType)
	}

	tflog.Debug(ctx, "Mapping Payloads")
	if payloads := remoteResource.GetPayloads(); len(payloads) > 0 {
		data.Payloads = types.ListValueMust(types.ObjectType{AttrTypes: payloadAttributeTypes()}, payloadsToValueSlice(payloads))
	} else {
		data.Payloads = types.ListNull(types.ObjectType{AttrTypes: payloadAttributeTypes()})
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform")
}

// roleScopeTagsToValueSlice converts a slice of role scope tag strings to a slice of Terraform attr.Value.
// This is used to populate the RoleScopeTags field in the Terraform resource model.
func roleScopeTagsToValueSlice(roleScopeTags []string) []attr.Value {
	if roleScopeTags == nil {
		return []attr.Value{}
	}
	values := make([]attr.Value, len(roleScopeTags))
	for i, tag := range roleScopeTags {
		values[i] = types.StringValue(tag)
	}
	return values
}

// payloadAttributeTypes returns a map of attribute names to their Terraform types for the Payload object.
// This defines the structure of the Payload object in the Terraform resource model.
func payloadAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"payload_id":             types.StringType,
		"payload_type":           types.StringType,
		"group_id":               types.StringType,
		"assignment_filter_type": types.StringType,
	}
}

// payloadsToValueSlice converts a slice of PayloadByFilterable to a slice of Terraform attr.Value.
// This is used to populate the Payloads field in the Terraform resource model.
func payloadsToValueSlice(payloads []models.PayloadByFilterable) []attr.Value {
	values := make([]attr.Value, len(payloads))
	for i, payload := range payloads {
		payloadMap := map[string]attr.Value{
			"payload_id":             types.StringValue(helpers.StringPtrToString(payload.GetPayloadId())),
			"payload_type":           types.StringValue(payloadTypeToString(payload.GetPayloadType())),
			"group_id":               types.StringValue(helpers.StringPtrToString(payload.GetGroupId())),
			"assignment_filter_type": types.StringValue(assignmentFilterTypeToString(payload.GetAssignmentFilterType())),
		}
		values[i] = types.ObjectValueMust(payloadAttributeTypes(), payloadMap)
	}
	return values
}

// payloadTypeToString converts AssociatedAssignmentPayloadType to its string representation.
func payloadTypeToString(payloadType *models.AssociatedAssignmentPayloadType) string {
	if payloadType == nil {
		return ""
	}
	return (*payloadType).String()
}

// assignmentFilterTypeToString converts DeviceAndAppManagementAssignmentFilterType to its string representation.
func assignmentFilterTypeToString(filterType *models.DeviceAndAppManagementAssignmentFilterType) string {
	if filterType == nil {
		return ""
	}
	return (*filterType).String()
}

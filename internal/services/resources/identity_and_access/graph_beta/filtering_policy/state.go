package graphBetaFilteringPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapRemoteResourceStateToTerraform maps the remote filtering policy to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *FilteringPolicyResourceModel, remoteResource map[string]any) {
	tflog.Debug(ctx, "Starting MapRemoteResourceStateToTerraform", map[string]any{
		"remoteResource": remoteResource,
	})

	if id, ok := remoteResource["id"].(string); ok {
		tflog.Debug(ctx, "Mapping ID", map[string]any{"id": id})
		data.ID = types.StringValue(id)
	} else {
		tflog.Debug(ctx, "ID not found or not a string")
		data.ID = types.StringNull()
	}

	if name, ok := remoteResource["name"].(string); ok {
		tflog.Debug(ctx, "Mapping name", map[string]any{"name": name})
		data.Name = types.StringValue(name)
	} else {
		tflog.Debug(ctx, "name not found or not a string")
		data.Name = types.StringNull()
	}

	if description, ok := remoteResource["description"].(string); ok {
		tflog.Debug(ctx, "Mapping description", map[string]any{"description": description})
		data.Description = types.StringValue(description)
	} else {
		tflog.Debug(ctx, "description not found or not a string")
		data.Description = types.StringNull()
	}

	if action, ok := remoteResource["action"].(string); ok {
		tflog.Debug(ctx, "Mapping action", map[string]any{"action": action})
		data.Action = types.StringValue(action)
	} else {
		tflog.Debug(ctx, "action not found or not a string")
		data.Action = types.StringNull()
	}

	if createdDateTime, ok := remoteResource["createdDateTime"].(string); ok {
		tflog.Debug(ctx, "Mapping createdDateTime", map[string]any{"createdDateTime": createdDateTime})
		data.CreatedDateTime = types.StringValue(createdDateTime)
	} else {
		tflog.Debug(ctx, "createdDateTime not found or not a string")
		data.CreatedDateTime = types.StringNull()
	}

	if lastModifiedDateTime, ok := remoteResource["lastModifiedDateTime"].(string); ok {
		tflog.Debug(ctx, "Mapping lastModifiedDateTime", map[string]any{"lastModifiedDateTime": lastModifiedDateTime})
		data.LastModifiedDateTime = types.StringValue(lastModifiedDateTime)
	} else {
		tflog.Debug(ctx, "lastModifiedDateTime not found or not a string")
		data.LastModifiedDateTime = types.StringNull()
	}

	// Note: 'version' is not documented in the Microsoft Graph API documentation but is actually
	// included in API responses. We map it here as a computed field.
	if version, ok := remoteResource["version"].(string); ok {
		tflog.Debug(ctx, "Mapping version", map[string]any{"version": version})
		data.Version = types.StringValue(version)
	} else {
		tflog.Debug(ctx, "version not found or not a string")
		data.Version = types.StringNull()
	}

	// Note: 'state' and 'priority' are not mapped here as they are properties used when linking
	// policies to security profiles, not direct properties of the filtering policy resource itself.
	// The Microsoft Graph API documentation incorrectly lists these in the update documentation.

	tflog.Debug(ctx, "Completed MapRemoteResourceStateToTerraform")
}

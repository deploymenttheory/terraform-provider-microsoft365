package graphBetaTargetedManagedAppConfigurations

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *TargetedManagedAppConfigurationResourceModel, remoteResource graphmodels.TargetedManagedAppConfigurationable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.Version = convert.GraphToFrameworkString(remoteResource.GetVersion())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	// Custom settings
	if customSettings := remoteResource.GetCustomSettings(); customSettings != nil {
		data.CustomSettings = mapCustomSettingsToTerraform(customSettings)
	} else {
		data.CustomSettings = []KeyValuePairResourceModel{}
	}

	data.AppGroupType = convert.GraphToFrameworkEnum(remoteResource.GetAppGroupType())

	// Apps - conditionally populate based on app_group_type
	// When app_group_type = "allApps", the API response includes all possible mobile app identifiers.
	// These identifiers are not required for a valid request but are set by the API automatically.
	// To prevent configuration drift, only populate state with apps when the app_group_type requires explicit app tracking.
	appGroupType := convert.GraphToFrameworkEnum(remoteResource.GetAppGroupType())
	if !appGroupType.IsNull() && appGroupType.ValueString() == "allApps" {
		tflog.Debug(ctx, "AppGroupType is 'allApps', setting apps to empty set to avoid drift from API-populated app list")
		// Keep apps as empty set - API auto-populates all apps, but they're not required for this group type
		data.Apps = types.SetValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"mobile_app_identifier": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":           types.StringType,
						"bundle_id":      types.StringType,
						"package_id":     types.StringType,
						"windows_app_id": types.StringType,
					},
				},
				"version": types.StringType,
			},
		}, []attr.Value{})
	} else {
		tflog.Debug(ctx, "AppGroupType is not 'allApps', populating apps from API response")
		if apps := remoteResource.GetApps(); apps != nil {
			data.Apps = mapAppsToTerraformSet(ctx, apps)
		} else {
			data.Apps = types.SetNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"mobile_app_identifier": types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"type":           types.StringType,
							"bundle_id":      types.StringType,
							"package_id":     types.StringType,
							"windows_app_id": types.StringType,
						},
					},
					"version": types.StringType,
				},
			})
		}
	}

	if assignments := remoteResource.GetAssignments(); assignments != nil {
		data.Assignments = mapAssignmentsToTerraform(ctx, assignments)
	} else {
		data.Assignments = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":     types.StringType,
				"group_id": types.StringType,
			},
		})
	}

	// Settings catalog mapping - use direct settings array processing
	if settings := remoteResource.GetSettings(); len(settings) > 0 {
		tflog.Debug(ctx, "Mapping settings catalog from remote response")

		data.SettingsCatalog = &DeviceConfigV2GraphServiceResourceModel{}

		if err := StateConfigurationPolicySettings(ctx, data, settings, nil); err != nil {
			tflog.Error(ctx, "Failed to map settings catalog", map[string]any{"error": err})
			data.SettingsCatalog = nil
		}
	} else {
		tflog.Debug(ctx, "No settings catalog data in remote response. skipping")
		data.SettingsCatalog = nil
	}

	data.DeployedAppCount = convert.GraphToFrameworkInt32(remoteResource.GetDeployedAppCount())
	data.IsAssigned = convert.GraphToFrameworkBool(remoteResource.GetIsAssigned())
	data.TargetedAppManagementLevels = convert.GraphToFrameworkEnum(remoteResource.GetTargetedAppManagementLevels())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform for resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// mapCustomSettingsToTerraform maps the custom settings to the Terraform state
func mapCustomSettingsToTerraform(customSettings []graphmodels.KeyValuePairable) []KeyValuePairResourceModel {
	if len(customSettings) == 0 {
		return []KeyValuePairResourceModel{}
	}

	result := make([]KeyValuePairResourceModel, len(customSettings))
	for i, setting := range customSettings {
		result[i] = KeyValuePairResourceModel{
			Name:  convert.GraphToFrameworkString(setting.GetName()),
			Value: convert.GraphToFrameworkString(setting.GetValue()),
		}
	}

	return result
}

func mapAppsToTerraformSet(ctx context.Context, apps []graphmodels.ManagedMobileAppable) types.Set {
	if len(apps) == 0 {
		return types.SetValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"mobile_app_identifier": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":           types.StringType,
						"bundle_id":      types.StringType,
						"package_id":     types.StringType,
						"windows_app_id": types.StringType,
					},
				},
				"version": types.StringType,
			},
		}, []attr.Value{})
	}

	appValues := make([]attr.Value, len(apps))
	for i, app := range apps {
		appAttrs := map[string]attr.Value{
			"version": convert.GraphToFrameworkString(app.GetVersion()),
		}

		if identifier := app.GetMobileAppIdentifier(); identifier != nil {
			identifierAttrs := map[string]attr.Value{
				"type":           types.StringNull(),
				"bundle_id":      types.StringNull(),
				"package_id":     types.StringNull(),
				"windows_app_id": types.StringNull(),
			}

			// Determine the type based on the actual type
			switch typedIdentifier := identifier.(type) {
			case graphmodels.AndroidMobileAppIdentifierable:
				identifierAttrs["type"] = types.StringValue("android_mobile_app")
				identifierAttrs["package_id"] = convert.GraphToFrameworkString(typedIdentifier.GetPackageId())
			case graphmodels.IosMobileAppIdentifierable:
				identifierAttrs["type"] = types.StringValue("ios_mobile_app")
				identifierAttrs["bundle_id"] = convert.GraphToFrameworkString(typedIdentifier.GetBundleId())
			case graphmodels.WindowsAppIdentifierable:
				identifierAttrs["type"] = types.StringValue("windows_app")
				identifierAttrs["windows_app_id"] = convert.GraphToFrameworkString(typedIdentifier.GetWindowsAppId())
			default:
				// Fallback to OData type if we can't determine the specific type
				if odataType := identifier.GetOdataType(); odataType != nil {
					switch *odataType {
					case "#microsoft.graph.androidMobileAppIdentifier":
						identifierAttrs["type"] = types.StringValue("android_mobile_app")
					case "#microsoft.graph.iosMobileAppIdentifier":
						identifierAttrs["type"] = types.StringValue("ios_mobile_app")
					case "#microsoft.graph.windowsAppIdentifier":
						identifierAttrs["type"] = types.StringValue("windows_app")
					default:
						identifierAttrs["type"] = types.StringValue("unknown")
					}
				} else {
					identifierAttrs["type"] = types.StringValue("unknown")
				}
			}

			appAttrs["mobile_app_identifier"] = types.ObjectValueMust(
				map[string]attr.Type{
					"type":           types.StringType,
					"bundle_id":      types.StringType,
					"package_id":     types.StringType,
					"windows_app_id": types.StringType,
				},
				identifierAttrs,
			)
		} else {
			appAttrs["mobile_app_identifier"] = types.ObjectNull(map[string]attr.Type{
				"type":           types.StringType,
				"bundle_id":      types.StringType,
				"package_id":     types.StringType,
				"windows_app_id": types.StringType,
			})
		}

		appValues[i] = types.ObjectValueMust(
			map[string]attr.Type{
				"mobile_app_identifier": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"type":           types.StringType,
						"bundle_id":      types.StringType,
						"package_id":     types.StringType,
						"windows_app_id": types.StringType,
					},
				},
				"version": types.StringType,
			},
			appAttrs,
		)
	}

	return types.SetValueMust(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"mobile_app_identifier": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"type":           types.StringType,
					"bundle_id":      types.StringType,
					"package_id":     types.StringType,
					"windows_app_id": types.StringType,
				},
			},
			"version": types.StringType,
		},
	}, appValues)
}

func mapAssignmentsToTerraform(_ context.Context, assignments []graphmodels.TargetedManagedAppPolicyAssignmentable) types.Set {
	if len(assignments) == 0 {
		return types.SetValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":     types.StringType,
				"group_id": types.StringType,
			},
		}, []attr.Value{})
	}

	assignmentValues := make([]attr.Value, len(assignments))
	for i, assignment := range assignments {
		assignmentAttrs := map[string]attr.Value{
			"type":     types.StringNull(),
			"group_id": types.StringNull(),
		}

		if target := assignment.GetTarget(); target != nil {
			if odataType := target.GetOdataType(); odataType != nil {
				switch *odataType {
				case "#microsoft.graph.groupAssignmentTarget":
					assignmentAttrs["type"] = types.StringValue("groupAssignmentTarget")
				case "#microsoft.graph.exclusionGroupAssignmentTarget":
					assignmentAttrs["type"] = types.StringValue("exclusionGroupAssignmentTarget")
				}
			}

			if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
				if groupId := groupTarget.GetGroupId(); groupId != nil {
					assignmentAttrs["group_id"] = types.StringValue(*groupId)
				}
			} else if exclusionGroupTarget, ok := target.(graphmodels.ExclusionGroupAssignmentTargetable); ok {
				if groupId := exclusionGroupTarget.GetGroupId(); groupId != nil {
					assignmentAttrs["group_id"] = types.StringValue(*groupId)
				}
			}
		}

		assignmentValues[i] = types.ObjectValueMust(
			map[string]attr.Type{
				"type":     types.StringType,
				"group_id": types.StringType,
			},
			assignmentAttrs,
		)
	}

	return types.SetValueMust(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":     types.StringType,
			"group_id": types.StringType,
		},
	}, assignmentValues)
}

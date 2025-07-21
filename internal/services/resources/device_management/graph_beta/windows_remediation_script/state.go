package graphBetaWindowsRemediationScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote DeviceHealthScript resource state to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *DeviceHealthScriptResourceModel, remoteResource graphmodels.DeviceHealthScriptable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Publisher = convert.GraphToFrameworkString(remoteResource.GetPublisher())
	data.RunAs32Bit = convert.GraphToFrameworkBool(remoteResource.GetRunAs32Bit())
	data.EnforceSignatureCheck = convert.GraphToFrameworkBool(remoteResource.GetEnforceSignatureCheck())
	data.Version = convert.GraphToFrameworkString(remoteResource.GetVersion())
	data.IsGlobalScript = convert.GraphToFrameworkBool(remoteResource.GetIsGlobalScript())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.HighestAvailableVersion = convert.GraphToFrameworkString(remoteResource.GetHighestAvailableVersion())
	data.RunAsAccount = convert.GraphToFrameworkEnum(remoteResource.GetRunAsAccount())
	data.DeviceHealthScriptType = convert.GraphToFrameworkEnum(remoteResource.GetDeviceHealthScriptType())
	data.DetectionScriptContent = convert.GraphToFrameworkBytes(remoteResource.GetDetectionScriptContent())
	data.RemediationScriptContent = convert.GraphToFrameworkBytes(remoteResource.GetRemediationScriptContent())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	tflog.Debug(ctx, "Successfully mapped basic resource properties", map[string]interface{}{
		"resourceId":             data.ID.ValueString(),
		"displayName":            data.DisplayName.ValueString(),
		"runAsAccount":           data.RunAsAccount.ValueString(),
		"deviceHealthScriptType": data.DeviceHealthScriptType.ValueString(),
		"version":                data.Version.ValueString(),
		"isGlobalScript":         data.IsGlobalScript.ValueBool(),
	})

	assignments := remoteResource.GetAssignments()
	tflog.Debug(ctx, "Retrieved assignments from remote resource", map[string]interface{}{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments found, setting assignments to null", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
		data.Assignments = types.SetNull(WindowsRemediationScriptAssignmentType())
	} else {
		tflog.Debug(ctx, "Starting assignment mapping process", map[string]interface{}{
			"resourceId":      data.ID.ValueString(),
			"assignmentCount": len(assignments),
		})
		MapAssignmentsToTerraform(ctx, data, assignments)
		tflog.Debug(ctx, "Completed assignment mapping process", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// WindowsRemediationScriptAssignmentType returns the object type for WindowsRemediationScriptAssignmentModel
func WindowsRemediationScriptAssignmentType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":              types.StringType,
			"group_id":          types.StringType,
			"filter_id":         types.StringType,
			"filter_type":       types.StringType,
			"daily_schedule":    types.ObjectType{AttrTypes: dailyScheduleAttrTypes()},
			"hourly_schedule":   types.ObjectType{AttrTypes: hourlyScheduleAttrTypes()},
			"run_once_schedule": types.ObjectType{AttrTypes: runOnceScheduleAttrTypes()},
		},
	}
}

func dailyScheduleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interval": types.Int32Type,
		"time":     types.StringType,
		"use_utc":  types.BoolType,
	}
}

func hourlyScheduleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interval": types.Int32Type,
	}
}

func runOnceScheduleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"date":    types.StringType,
		"time":    types.StringType,
		"use_utc": types.BoolType,
	}
}

// MapAssignmentsToTerraform maps the remote DeviceHealthScript assignments to Terraform state
func MapAssignmentsToTerraform(ctx context.Context, data *DeviceHealthScriptResourceModel, assignments []graphmodels.DeviceHealthScriptAssignmentable) {
	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments to process")
		data.Assignments = types.SetNull(WindowsRemediationScriptAssignmentType())
		return
	}

	tflog.Debug(ctx, "Starting assignment mapping process", map[string]interface{}{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	assignmentValues := []attr.Value{}

	for i, assignment := range assignments {
		tflog.Debug(ctx, "Processing assignment", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		target := assignment.GetTarget()
		if target == nil {
			tflog.Warn(ctx, "Assignment target is nil, skipping assignment", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			continue
		}

		odataType := target.GetOdataType()
		if odataType == nil {
			tflog.Warn(ctx, "Assignment target OData type is nil, skipping assignment", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			continue
		}

		tflog.Debug(ctx, "Processing assignment target", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"targetType":      *odataType,
			"resourceId":      data.ID.ValueString(),
		})

		assignmentObj := map[string]attr.Value{
			"type":              types.StringNull(),
			"group_id":          types.StringNull(),
			"filter_id":         types.StringNull(),
			"filter_type":       types.StringNull(),
			"daily_schedule":    types.ObjectNull(dailyScheduleAttrTypes()),
			"hourly_schedule":   types.ObjectNull(hourlyScheduleAttrTypes()),
			"run_once_schedule": types.ObjectNull(runOnceScheduleAttrTypes()),
		}

		switch *odataType {
		case "#microsoft.graph.allDevicesAssignmentTarget":
			tflog.Debug(ctx, "Mapping allDevicesAssignmentTarget", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["type"] = types.StringValue("allDevicesAssignmentTarget")
			assignmentObj["group_id"] = types.StringNull()

		case "#microsoft.graph.allLicensedUsersAssignmentTarget":
			tflog.Debug(ctx, "Mapping allLicensedUsersAssignmentTarget", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["type"] = types.StringValue("allLicensedUsersAssignmentTarget")
			assignmentObj["group_id"] = types.StringNull()

		case "#microsoft.graph.groupAssignmentTarget":
			tflog.Debug(ctx, "Mapping groupAssignmentTarget", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["type"] = types.StringValue("groupAssignmentTarget")

			if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
				groupId := groupTarget.GetGroupId()
				if groupId != nil && *groupId != "" {
					tflog.Debug(ctx, "Setting group ID for group assignment target", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"groupId":         *groupId,
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				} else {
					tflog.Warn(ctx, "Group ID is nil/empty for group assignment target", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = types.StringNull()
				}
			} else {
				tflog.Error(ctx, "Failed to cast target to GroupAssignmentTargetable", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["group_id"] = types.StringNull()
			}

		case "#microsoft.graph.exclusionGroupAssignmentTarget":
			tflog.Debug(ctx, "Mapping exclusionGroupAssignmentTarget", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["type"] = types.StringValue("exclusionGroupAssignmentTarget")

			if groupTarget, ok := target.(graphmodels.ExclusionGroupAssignmentTargetable); ok {
				groupId := groupTarget.GetGroupId()
				if groupId != nil && *groupId != "" {
					tflog.Debug(ctx, "Setting group ID for exclusion group assignment target", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"groupId":         *groupId,
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = convert.GraphToFrameworkString(groupId)
				} else {
					tflog.Warn(ctx, "Group ID is nil/empty for exclusion group assignment target", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"resourceId":      data.ID.ValueString(),
					})
					assignmentObj["group_id"] = types.StringNull()
				}
			} else {
				tflog.Error(ctx, "Failed to cast target to ExclusionGroupAssignmentTargetable", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["group_id"] = types.StringNull()
			}

		default:
			tflog.Warn(ctx, "Unknown target type encountered", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"targetType":      *odataType,
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["group_id"] = types.StringNull()
		}

		tflog.Debug(ctx, "Processing assignment filters", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		filterID := target.GetDeviceAndAppManagementAssignmentFilterId()
		if filterID != nil && *filterID != "" && *filterID != "00000000-0000-0000-0000-000000000000" {
			tflog.Debug(ctx, "Assignment has meaningful filter ID", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"filterId":        *filterID,
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["filter_id"] = convert.GraphToFrameworkString(filterID)
		} else {
			tflog.Debug(ctx, "Assignment has no meaningful filter ID, using schema default", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["filter_id"] = types.StringValue("00000000-0000-0000-0000-000000000000")
		}

		filterType := target.GetDeviceAndAppManagementAssignmentFilterType()
		if filterType != nil {
			tflog.Debug(ctx, "Processing filter type", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"filterType":      *filterType,
				"resourceId":      data.ID.ValueString(),
			})

			switch *filterType {
			case graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				tflog.Debug(ctx, "Setting filter type to include", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["filter_type"] = types.StringValue("include")
			case graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				tflog.Debug(ctx, "Setting filter type to exclude", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["filter_type"] = types.StringValue("exclude")
			case graphmodels.NONE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				tflog.Debug(ctx, "Setting filter type to none", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["filter_type"] = types.StringValue("none")
			default:
				tflog.Debug(ctx, "Unknown filter type, using schema default", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"filterType":      *filterType,
					"resourceId":      data.ID.ValueString(),
				})
				assignmentObj["filter_type"] = types.StringValue("none")
			}
		} else {
			tflog.Debug(ctx, "No filter type specified, using schema default", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentObj["filter_type"] = types.StringValue("none")
		}

		tflog.Debug(ctx, "Processing assignment schedule", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		runSchedule := assignment.GetRunSchedule()
		if runSchedule != nil {
			scheduleType := runSchedule.GetOdataType()
			if scheduleType != nil {
				tflog.Debug(ctx, "Assignment has schedule", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"scheduleType":    *scheduleType,
					"resourceId":      data.ID.ValueString(),
				})

				switch *scheduleType {
				case "#microsoft.graph.deviceHealthScriptDailySchedule":
					tflog.Debug(ctx, "Processing daily schedule", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"resourceId":      data.ID.ValueString(),
					})

					if dailySchedule, ok := runSchedule.(graphmodels.DeviceHealthScriptDailyScheduleable); ok {
						dailyScheduleObj := map[string]attr.Value{
							"interval": convert.GraphToFrameworkInt32(dailySchedule.GetInterval()),
							"use_utc":  convert.GraphToFrameworkBool(dailySchedule.GetUseUtc()),
							"time":     types.StringNull(),
						}

						if interval := dailySchedule.GetInterval(); interval != nil {
							tflog.Debug(ctx, "Daily schedule interval", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"interval":        *interval,
								"resourceId":      data.ID.ValueString(),
							})
						}

						if useUtc := dailySchedule.GetUseUtc(); useUtc != nil {
							tflog.Debug(ctx, "Daily schedule UTC setting", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"useUtc":          *useUtc,
								"resourceId":      data.ID.ValueString(),
							})
						}

						if timeValue := dailySchedule.GetTime(); timeValue != nil {
							tflog.Debug(ctx, "Processing daily schedule time", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"timeValue":       timeValue.String(),
								"resourceId":      data.ID.ValueString(),
							})
							dailyScheduleObj["time"] = convert.GraphToFrameworkTimeOnlyWithPrecision(timeValue, 0)
						} else {
							tflog.Warn(ctx, "Daily schedule time is nil", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"resourceId":      data.ID.ValueString(),
							})
						}

						dailyObj, diags := types.ObjectValue(dailyScheduleAttrTypes(), dailyScheduleObj)
						if !diags.HasError() {
							tflog.Debug(ctx, "Successfully created daily schedule object", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"resourceId":      data.ID.ValueString(),
							})
							assignmentObj["daily_schedule"] = dailyObj
						} else {
							tflog.Error(ctx, "Failed to create daily schedule object", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"errors":          diags.Errors(),
								"resourceId":      data.ID.ValueString(),
							})
						}
					} else {
						tflog.Error(ctx, "Failed to cast run schedule to DeviceHealthScriptDailyScheduleable", map[string]interface{}{
							"assignmentIndex": i,
							"assignmentId":    assignment.GetId(),
							"resourceId":      data.ID.ValueString(),
						})
					}

				case "#microsoft.graph.deviceHealthScriptHourlySchedule":
					tflog.Debug(ctx, "Processing hourly schedule", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"resourceId":      data.ID.ValueString(),
					})

					if hourlySchedule, ok := runSchedule.(graphmodels.DeviceHealthScriptHourlyScheduleable); ok {
						hourlyScheduleObj := map[string]attr.Value{
							"interval": convert.GraphToFrameworkInt32(hourlySchedule.GetInterval()),
						}

						if interval := hourlySchedule.GetInterval(); interval != nil {
							tflog.Debug(ctx, "Hourly schedule interval", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"interval":        *interval,
								"resourceId":      data.ID.ValueString(),
							})
						}

						hourlyObj, diags := types.ObjectValue(hourlyScheduleAttrTypes(), hourlyScheduleObj)
						if !diags.HasError() {
							tflog.Debug(ctx, "Successfully created hourly schedule object", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"resourceId":      data.ID.ValueString(),
							})
							assignmentObj["hourly_schedule"] = hourlyObj
						} else {
							tflog.Error(ctx, "Failed to create hourly schedule object", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"errors":          diags.Errors(),
								"resourceId":      data.ID.ValueString(),
							})
						}
					} else {
						tflog.Error(ctx, "Failed to cast run schedule to DeviceHealthScriptHourlyScheduleable", map[string]interface{}{
							"assignmentIndex": i,
							"assignmentId":    assignment.GetId(),
							"resourceId":      data.ID.ValueString(),
						})
					}

				case "#microsoft.graph.deviceHealthScriptRunOnceSchedule":
					tflog.Debug(ctx, "Processing run once schedule", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"resourceId":      data.ID.ValueString(),
					})

					if runOnceSchedule, ok := runSchedule.(graphmodels.DeviceHealthScriptRunOnceScheduleable); ok {
						runOnceScheduleObj := map[string]attr.Value{
							"date":    types.StringNull(),
							"time":    types.StringNull(),
							"use_utc": convert.GraphToFrameworkBool(runOnceSchedule.GetUseUtc()),
						}

						if useUtc := runOnceSchedule.GetUseUtc(); useUtc != nil {
							tflog.Debug(ctx, "Run once schedule UTC setting", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"useUtc":          *useUtc,
								"resourceId":      data.ID.ValueString(),
							})
						}

						if dateValue := runOnceSchedule.GetDate(); dateValue != nil {
							tflog.Debug(ctx, "Processing run once schedule date", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"dateValue":       dateValue.String(),
								"resourceId":      data.ID.ValueString(),
							})
							runOnceScheduleObj["date"] = convert.GraphToFrameworkDateOnly(dateValue)
						} else {
							tflog.Warn(ctx, "Run once schedule date is nil", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"resourceId":      data.ID.ValueString(),
							})
						}

						if timeValue := runOnceSchedule.GetTime(); timeValue != nil {
							tflog.Debug(ctx, "Processing run once schedule time", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"timeValue":       timeValue.String(),
								"resourceId":      data.ID.ValueString(),
							})
							runOnceScheduleObj["time"] = convert.GraphToFrameworkTimeOnlyWithPrecision(timeValue, 0)
						} else {
							tflog.Warn(ctx, "Run once schedule time is nil", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"resourceId":      data.ID.ValueString(),
							})
						}

						runOnceObj, diags := types.ObjectValue(runOnceScheduleAttrTypes(), runOnceScheduleObj)
						if !diags.HasError() {
							tflog.Debug(ctx, "Successfully created run once schedule object", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"resourceId":      data.ID.ValueString(),
							})
							assignmentObj["run_once_schedule"] = runOnceObj
						} else {
							tflog.Error(ctx, "Failed to create run once schedule object", map[string]interface{}{
								"assignmentIndex": i,
								"assignmentId":    assignment.GetId(),
								"errors":          diags.Errors(),
								"resourceId":      data.ID.ValueString(),
							})
						}
					} else {
						tflog.Error(ctx, "Failed to cast run schedule to DeviceHealthScriptRunOnceScheduleable", map[string]interface{}{
							"assignmentIndex": i,
							"assignmentId":    assignment.GetId(),
							"resourceId":      data.ID.ValueString(),
						})
					}

				default:
					tflog.Warn(ctx, "Unknown schedule type encountered", map[string]interface{}{
						"assignmentIndex": i,
						"assignmentId":    assignment.GetId(),
						"scheduleType":    *scheduleType,
						"resourceId":      data.ID.ValueString(),
					})
				}
			} else {
				tflog.Warn(ctx, "Schedule OData type is nil", map[string]interface{}{
					"assignmentIndex": i,
					"assignmentId":    assignment.GetId(),
					"resourceId":      data.ID.ValueString(),
				})
			}
		} else {
			tflog.Debug(ctx, "No schedule found for assignment", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
		}

		tflog.Debug(ctx, "Creating assignment object value", map[string]interface{}{
			"assignmentIndex": i,
			"assignmentId":    assignment.GetId(),
			"resourceId":      data.ID.ValueString(),
		})

		objValue, diags := types.ObjectValue(WindowsRemediationScriptAssignmentType().(types.ObjectType).AttrTypes, assignmentObj)
		if !diags.HasError() {
			tflog.Debug(ctx, "Successfully created assignment object", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"resourceId":      data.ID.ValueString(),
			})
			assignmentValues = append(assignmentValues, objValue)
		} else {
			tflog.Error(ctx, "Failed to create assignment object value", map[string]interface{}{
				"assignmentIndex": i,
				"assignmentId":    assignment.GetId(),
				"errors":          diags.Errors(),
				"resourceId":      data.ID.ValueString(),
			})
		}
	}

	tflog.Debug(ctx, "Creating assignments set", map[string]interface{}{
		"processedAssignments": len(assignmentValues),
		"originalAssignments":  len(assignments),
		"resourceId":           data.ID.ValueString(),
	})

	if len(assignmentValues) > 0 {
		setVal, diags := types.SetValue(WindowsRemediationScriptAssignmentType(), assignmentValues)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to create assignments set", map[string]interface{}{
				"errors":     diags.Errors(),
				"resourceId": data.ID.ValueString(),
			})
			data.Assignments = types.SetNull(WindowsRemediationScriptAssignmentType())
		} else {
			tflog.Debug(ctx, "Successfully created assignments set", map[string]interface{}{
				"assignmentCount": len(assignmentValues),
				"resourceId":      data.ID.ValueString(),
			})
			data.Assignments = setVal
		}
	} else {
		tflog.Debug(ctx, "No valid assignments processed, setting assignments to null", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
		data.Assignments = types.SetNull(WindowsRemediationScriptAssignmentType())
	}

	tflog.Debug(ctx, "Finished mapping assignments to Terraform state", map[string]interface{}{
		"finalAssignmentCount": len(assignmentValues),
		"originalAssignments":  len(assignments),
		"resourceId":           data.ID.ValueString(),
	})
}

package schema

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// DeviceAndAppManagementAssignments returns a schema attribute for device and app management resource assignments that can be reused across different app types.
func DeviceAndAppManagementAssignments() schema.Attribute {
	return schema.ListNestedAttribute{
		Optional:    true,
		Description: "The list of group assignments for this device and app management resource.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"target": schema.SingleNestedAttribute{
					Required: true,
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:    true,
							Description: "The type of target. Possible values are: allLicensedUsers, allDevices, group.",
							Validators: []validator.String{
								stringvalidator.OneOf("allLicensedUsers", "allDevices", "group"),
							},
						},
						"group_id": schema.StringAttribute{
							Optional:    true,
							Description: "The ID of the group to assign the app to. Required when type is 'group'.",
						},
					},
					Description: "The target for this assignment.",
				},
				"intent": schema.StringAttribute{
					Required:    true,
					Description: "The intent of the assignment. Possible values are: available, required, uninstall, availableWithoutEnrollment.",
					Validators: []validator.String{
						stringvalidator.OneOf("available", "required", "uninstall", "availableWithoutEnrollment"),
					},
				},
				"settings": schema.SingleNestedAttribute{
					Optional: true,
					Attributes: map[string]schema.Attribute{
						"notifications": schema.StringAttribute{
							Optional:    true,
							Description: "The notification setting for the assignment. Possible values are: showAll, showReboot, hideAll.",
							Validators: []validator.String{
								stringvalidator.OneOf("showAll", "showReboot", "hideAll"),
							},
						},
					},
					Description: "The settings for this assignment.",
				},
			},
		},
	}
}

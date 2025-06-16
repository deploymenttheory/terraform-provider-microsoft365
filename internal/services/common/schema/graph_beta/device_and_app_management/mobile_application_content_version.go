package schema

import (
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// MobileAppContentVersionSchema returns the schema definition for content versions
func MobileAppContentVersionSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The committed content version of the app, including its files. Only the currently committed version is shown.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The unique identifier for this content version. This ID is assigned during creation of the content version. Read-only.",
				},
				"files": schema.SetNestedAttribute{
					Computed:            true,
					MarkdownDescription: "The files associated with this content version.",
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates the name of the file.",
								PlanModifiers: []planmodifier.String{
									planmodifiers.UseStateForUnknownString(),
								},
							},
							"size": schema.Int64Attribute{
								Computed:            true,
								MarkdownDescription: "Indicates the original size of the file, in bytes.",
								PlanModifiers: []planmodifier.Int64{
									planmodifiers.UseStateForUnknownInt64(),
								},
							},
							"size_encrypted": schema.Int64Attribute{
								Computed:            true,
								MarkdownDescription: "Indicates the size of the file after encryption, in bytes.",
								PlanModifiers: []planmodifier.Int64{
									planmodifiers.UseStateForUnknownInt64(),
								},
							},
							"upload_state": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates the state of the current upload request. This property is read-only.",
								PlanModifiers: []planmodifier.String{
									planmodifiers.UseStateForUnknownString(),
								},
							},
							"is_committed": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "A value indicating whether the file is committed. A committed app content file has been fully uploaded and validated by the Intune service. Read-only.",
								PlanModifiers: []planmodifier.Bool{
									planmodifiers.UseStateForUnknownBool(),
								},
							},
							"is_dependency": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates whether this content file is a dependency for the main content file.",
								PlanModifiers: []planmodifier.Bool{
									planmodifiers.UseStateForUnknownBool(),
								},
							},
							"is_framework_file": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates whether this content file is a framework file.",
								PlanModifiers: []planmodifier.Bool{
									planmodifiers.UseStateForUnknownBool(),
								},
							},
							"azure_storage_uri": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates the Azure Storage URI that the file is uploaded to. Read-only.",
								PlanModifiers: []planmodifier.String{
									planmodifiers.UseStateForUnknownString(),
								},
							},
							"azure_storage_uri_expiration": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates the date and time when the Azure storage URI expires, in ISO 8601 format. Read-only.",
								PlanModifiers: []planmodifier.String{
									planmodifiers.UseStateForUnknownString(),
								},
							},
							"created_date_time": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates created date and time associated with app content file, in ISO 8601 format. Read-only.",
								PlanModifiers: []planmodifier.String{
									planmodifiers.UseStateForUnknownString(),
								},
							},
						},
					},
				},
			},
		},
	}
}

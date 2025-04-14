package schema

import (
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func MobileAppContentVersionSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true,
		Computed: true,
		PlanModifiers: []planmodifier.List{
			listplanmodifier.UseStateForUnknown(),
		},
		MarkdownDescription: "The content versions of the app, including their files.",
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Computed:            true,
					MarkdownDescription: "The unique identifier for this mobileAppContentFile. This id is assigned during creation of the mobileAppContentFile. Read-only. This property is read-only.",
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
								MarkdownDescription: "Indicates the state of the current upload request. Possible values are: success, transientError, error, unknown, azureStorageUriRequestSuccess, azureStorageUriRequestPending, azureStorageUriRequestFailed, azureStorageUriRequestTimedOut, azureStorageUriRenewalSuccess, azureStorageUriRenewalPending, azureStorageUriRenewalFailed, azureStorageUriRenewalTimedOut, commitFileSuccess, commitFilePending, commitFileFailed, commitFileTimedOut. Default value is success. This property is read-only.",
								PlanModifiers: []planmodifier.String{
									planmodifiers.UseStateForUnknownString(),
								},
							},
							"is_committed": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "A value indicating whether the file is committed. A committed app content file has been fully uploaded and validated by the Intune service. TRUE means that app content file is committed, FALSE means that app content file is not committed. Defaults to FALSE. Read-only.",
								PlanModifiers: []planmodifier.Bool{
									planmodifiers.UseStateForUnknownBool(),
								},
							},
							"is_dependency": schema.BoolAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates whether this content file is a dependency for the main content file. TRUE means that the content file is a dependency, FALSE means that the content file is not a dependency and is the main content file. Defaults to FALSE.",
								PlanModifiers: []planmodifier.Bool{
									planmodifiers.UseStateForUnknownBool(),
								},
							},
							"azure_storage_uri": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates the Azure Storage URI that the file is uploaded to. Created by the service upon receiving a valid mobileAppContentFile. Read-only.",
								PlanModifiers: []planmodifier.String{
									planmodifiers.UseStateForUnknownString(),
								},
							},
							"azure_storage_uri_expiration": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates the date and time when the Azure storage URI expires, in ISO 8601 format. For example, midnight UTC on Jan 1, 2014 would look like this: '2014-01-01T00:00:00Z'. Read-only.",
								PlanModifiers: []planmodifier.String{
									planmodifiers.UseStateForUnknownString(),
								},
							},
							"created_date_time": schema.StringAttribute{
								Computed:            true,
								MarkdownDescription: "Indicates created date and time associated with app content file, in ISO 8601 format. For example, midnight UTC on Jan 1, 2014 would look like this: '2014-01-01T00:00:00Z'. Read-only. This property is read-only.",
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

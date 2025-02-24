package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func MobileAppContentVersionSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed:            true,
		MarkdownDescription: "The content version details for the macOS PKG app.",
		Attributes: map[string]schema.Attribute{
			"content_version_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The ID of the content version.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"file_count": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The number of files in the content version.",
			},
			"files": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "The list of files for this content version.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The file ID.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The file name.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"is_dependency": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates if the file is a dependency.",
						},
						"is_committed": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates if the file has been committed.",
						},
						"size": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The size of the file.",
						},
						"size_encrypted": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The encrypted size of the file.",
						},
						"upload_state": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The upload state of the file.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"created_date_time": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The creation date and time of the file.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"azure_storage_uri": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The Azure Storage URI for the file.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"azure_storage_uri_expiration": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The expiration time of the Azure Storage URI.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"is_framework_file": schema.BoolAttribute{
							Computed:            true,
							MarkdownDescription: "Indicates if the file is a framework file.",
						},
						"manifest": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "The manifest content of the file.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"size_encrypted_in_bytes": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The encrypted size in bytes of the file.",
						},
						"size_in_bytes": schema.Int64Attribute{
							Computed:            true,
							MarkdownDescription: "The size in bytes of the file.",
						},
					},
				},
			},
		},
	}
}

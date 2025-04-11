package schema

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func MobileAppMetadataSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "Metadata related to the installer file, such as size and checksums. This is automatically computed during app creation and updates.",
		Attributes: map[string]schema.Attribute{
			"installer_size_in_bytes": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The size of the installer file in bytes. Used to detect changes in content.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"installer_md5_checksum": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The MD5 checksum of the installer file. Used as an additional verification of file integrity and to detect changes.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"installer_sha256_checksum": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The SHA256 checksum of the installer file. Used as a more secure verification of file integrity and to detect changes.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

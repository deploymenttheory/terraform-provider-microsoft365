package schema

import (
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func MobileAppMetadataSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed:            true,
		Optional:            true,
		MarkdownDescription: "Metadata related to the installer file, such as size and checksums. This is automatically computed during app creation and updates.",
		PlanModifiers: []planmodifier.Object{
			planmodifiers.UseStateForUnknownObject(),
		},
		Attributes: map[string]schema.Attribute{
			"installer_size_in_bytes": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "The size of the installer file in bytes. Used to detect changes in content.",
			},
			"installer_md5_checksum": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The MD5 checksum of the installer file. Used as an additional verification of file integrity and to detect changes.",
			},
			"installer_sha256_checksum": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "The SHA256 checksum of the installer file. Used as a more secure verification of file integrity and to detect changes.",
			},
		},
	}
}

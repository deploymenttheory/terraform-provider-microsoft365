package schema

import (
	"regexp"

	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func MobileAppMetadataSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		MarkdownDescription: "Metadata related to the installer file, such as size and checksums. This is automatically computed during app creation and updates.",
		PlanModifiers: []planmodifier.Object{
			planmodifiers.UseStateForUnknownObject(),
		},
		Attributes: map[string]schema.Attribute{
			"installer_file_path_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The path to the PKG file to be uploaded. The file must be a valid `.pkg` file. Value is not returned by API call.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`.*\.pkg$`),
						"File path must point to a valid .pkg file.",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
			"installer_url_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The web location of the PKG file, can be a http(s) URL. The file must be a valid `.pkg` file. Value is not returned by API call.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(http|https|file)://.*$|^(/|./|../).*$`),
						"Must be a valid URL.",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
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

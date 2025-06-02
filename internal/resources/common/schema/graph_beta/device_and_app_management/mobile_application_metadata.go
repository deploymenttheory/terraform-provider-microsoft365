package schema

import (
	"regexp"

	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// MobileAppMacOSPkgInstallerMetadataSchema returns schema for macOS PKG app installer metadata
func MobileAppMacOSPkgInstallerMetadataSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		MarkdownDescription: "Metadata related to the PKG installer file, such as size and checksums. This is automatically computed during app creation and updates.",
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
					stringplanmodifier.RequiresReplace(),
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
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// MobileAppMacOSLobInstallerMetadataSchema returns schema for macOS LOB app installer metadata
func MobileAppMacOSLobInstallerMetadataSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		MarkdownDescription: "Metadata related to the LOB installer file, such as size and checksums. This is automatically computed during app creation and updates.",
		PlanModifiers: []planmodifier.Object{
			planmodifiers.UseStateForUnknownObject(),
		},
		Attributes: map[string]schema.Attribute{
			"installer_file_path_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The path to the LOB installer file to be uploaded. The file must be a valid `.pkg` file. Value is not returned by API call.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`.*\.pkg$`),
						"File path must point to a valid .pkg file.",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"installer_url_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The web location of the LOB installer file, can be a http(s) URL. The file must be a valid `.pkg` file. Value is not returned by API call.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(http|https|file)://.*$|^(/|./|../).*$`),
						"Must be a valid URL.",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

// MobileAppDmgInstallerMetadataSchema returns schema for macOS DMG app installer metadata
func MobileAppDmgInstallerMetadataSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		MarkdownDescription: "Metadata related to the DMG installer file, such as size and checksums. This is automatically computed during app creation and updates.",
		PlanModifiers: []planmodifier.Object{
			planmodifiers.UseStateForUnknownObject(),
		},
		Attributes: map[string]schema.Attribute{
			"installer_file_path_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The path to the DMG installer file to be uploaded. The file must be a valid `.dmg` file. Value is not returned by API call.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`.*\.dmg$`),
						"File path must point to a valid .dmg file.",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"installer_url_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The web location of the DMG installer file, can be a http(s) URL. The file must be a valid `.dmg` file. Value is not returned by API call.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(http|https|file)://.*$|^(/|./|../).*$`),
						"Must be a valid URL.",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

package schema

import (
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
						regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
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
						regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
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
						regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
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

// MobileAppWin32LobInstallerMetadataSchema describes prepackaged Win32 content.
func MobileAppWin32LobInstallerMetadataSchema() schema.SingleNestedAttribute {
	return win32InstallerSourceSchema(false)
}

// MobileAppWin32ZipInstallerMetadataSchema describes unencrypted installer ZIPs.
func MobileAppWin32ZipInstallerMetadataSchema() schema.SingleNestedAttribute {
	return win32InstallerSourceSchema(true)
}

func win32InstallerSourceSchema(plainZip bool) schema.SingleNestedAttribute {
	format := "a prepackaged `.intunewin` file containing Detection.xml and its encrypted payload"
	conflict := "app_installer_zip"
	extension := `(?i).*\.intunewin$`
	if plainZip {
		format = "an unencrypted installer ZIP (a `.zip` file, or a legacy ZIP renamed `.intunewin`); the provider encrypts it before upload"
		conflict = "app_installer"
		extension = `(?i).*\.(zip|intunewin)$`
	}
	return schema.SingleNestedAttribute{
		Optional:            true,
		MarkdownDescription: "Source for " + format + ". Mutually exclusive with `" + conflict + "`. Set exactly one local path or URL. Source changes publish a new content version without replacing the application. Omit both source blocks when managing an imported application's metadata only.",
		Validators: []validator.Object{
			objectvalidator.ConflictsWith(path.MatchRoot(conflict)),
		},
		Attributes: map[string]schema.Attribute{
			"installer_file_path_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Local path to " + format + ". Not returned by the API. Use a different path for each package version.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(extension), "File path must have the documented package extension."),
					stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("installer_url_source")),
				},
			},
			"installer_url_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "HTTP(S) URL for " + format + ". Not returned by the API. Use a different URL for each package version.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(constants.HttpOrHttpsUrlRegex), "Must be a valid URL."),
				},
			},
		},
	}
}

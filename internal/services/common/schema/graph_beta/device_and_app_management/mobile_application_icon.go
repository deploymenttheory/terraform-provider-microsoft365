package schema

import (
	"regexp"

	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func MobileAppIconSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		MarkdownDescription: "The source information for the app icon. Supports various image formats (JPEG, PNG, GIF, etc.) which will be automatically converted to PNG as required by Microsoft Intune.",
		PlanModifiers: []planmodifier.Object{
			planmodifiers.UseStateForUnknownObject(),
		},
		Attributes: map[string]schema.Attribute{
			"icon_file_path_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The file path to the icon file to be uploaded. Supports various image formats which will be automatically converted to PNG.",
			},
			"icon_url_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The web location of the icon file, can be a http(s) URL. Supports various image formats which will be automatically converted to PNG.",
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
		},
	}
}

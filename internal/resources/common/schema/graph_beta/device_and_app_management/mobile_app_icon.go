package schema

import (
	"regexp"

	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/plan_modifiers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func MobileAppIconSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		MarkdownDescription: "The source information for the app icon (PNG). Supports local file paths or URLs.",
		PlanModifiers: []planmodifier.Object{
			planmodifiers.UseStateForUnknownObject(),
		},
		Attributes: map[string]schema.Attribute{
			"icon_file_path_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The file path to the icon file (PNG) to be uploaded.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`\.png$`),
						"must end with .png file extension",
					),
				},
			},
			"icon_url_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The web location of the icon file (PNG), can be a http(s) URL.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(http|https|file)://.*$|^(/|./|../).*$`),
						"Must be a valid URL.",
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`\.png$`),
						"must end with .png file extension",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
		},
	}
}

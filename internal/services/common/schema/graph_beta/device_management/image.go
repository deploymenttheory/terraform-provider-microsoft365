package schema

import (
	"regexp"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	planmodifiers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/plan_modifiers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ImageSchema returns a schema for image handling that supports both file and URL sources
func ImageSchema(description string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		MarkdownDescription: description + " Supports various image formats which will be automatically converted as required by Microsoft Intune.",
		PlanModifiers: []planmodifier.Object{
			planmodifiers.UseStateForUnknownObject(),
		},
		Validators: []validator.Object{
			validators.ExactlyOneOf("image_file_path_source", "image_url_source"),
		},
		Attributes: map[string]schema.Attribute{
			"image_file_path_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The file path to the image file to be uploaded. Supports various image formats.",
			},
			"image_url_source": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "The web location of the image file, can be a http(s) URL. Supports various image formats.",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(constants.HttpOrHttpsUrlRegex),
						"must be a valid URL starting with http:// or https://",
					),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.UseStateForUnknownString(),
				},
			},
		},
	}
}

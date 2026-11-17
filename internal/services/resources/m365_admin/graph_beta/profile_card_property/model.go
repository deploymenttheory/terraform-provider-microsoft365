// REF: https://learn.microsoft.com/en-us/graph/api/resources/profilecardproperty?view=graph-rest-beta
package graphBetaM365AdminProfileCardProperty

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ProfileCardPropertyResourceModel represents the schema for the profile card property resource
type ProfileCardPropertyResourceModel struct {
	ID                    types.String                 `tfsdk:"id"`
	DirectoryPropertyName types.String                 `tfsdk:"directory_property_name"`
	IsVisible             types.Bool                   `tfsdk:"is_visible"`
	Annotations           []ProfileCardAnnotationModel `tfsdk:"annotations"`
	Timeouts              timeouts.Value               `tfsdk:"timeouts"`
}

// ProfileCardAnnotationModel represents one profileCardAnnotation (a custom label for the attribute)
type ProfileCardAnnotationModel struct {
	DisplayName   types.String                   `tfsdk:"display_name"`
	Localizations []DisplayNameLocalizationModel `tfsdk:"localizations"`
}

// DisplayNameLocalizationModel represents one displayNameLocalization (a label for a language)
type DisplayNameLocalizationModel struct {
	LanguageTag types.String `tfsdk:"language_tag"`
	DisplayName types.String `tfsdk:"display_name"`
}

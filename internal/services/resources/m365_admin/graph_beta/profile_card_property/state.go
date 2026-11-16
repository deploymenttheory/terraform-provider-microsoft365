package graphBetaM365AdminProfileCardProperty

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the remote profile card property from Kiota SDK to Terraform state
func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data ProfileCardPropertyResourceModel,
	remoteResource graphmodels.ProfileCardPropertyable,
) ProfileCardPropertyResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"directoryPropertyName": remoteResource.GetDirectoryPropertyName(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetDirectoryPropertyName())
	data.DirectoryPropertyName = convert.GraphToFrameworkString(remoteResource.GetDirectoryPropertyName())
	data.IsVisible = convert.GraphToFrameworkBool(remoteResource.GetIsVisible())
	data.Annotations = mapAnnotations(remoteResource.GetAnnotations())

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]any{
		"resourceId": data.ID.ValueString(),
	})

	return data
}

// mapAnnotations maps the remote annotations collection to the Terraform model.
// Graph returns an empty collection when no annotations or localizations are set; empty maps to null to match an omitted configuration value.
func mapAnnotations(annotations []graphmodels.ProfileCardAnnotationable) []ProfileCardAnnotationModel {
	if len(annotations) == 0 {
		return nil
	}

	result := make([]ProfileCardAnnotationModel, 0, len(annotations))

	for _, annotation := range annotations {
		if annotation == nil {
			continue
		}
		result = append(result, ProfileCardAnnotationModel{
			DisplayName:   convert.GraphToFrameworkString(annotation.GetDisplayName()),
			Localizations: mapLocalizations(annotation.GetLocalizations()),
		})
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

// mapLocalizations maps the remote localizations collection of an annotation to the Terraform model
func mapLocalizations(localizations []graphmodels.DisplayNameLocalizationable) []DisplayNameLocalizationModel {
	if len(localizations) == 0 {
		return nil
	}

	result := make([]DisplayNameLocalizationModel, 0, len(localizations))

	for _, localization := range localizations {
		if localization == nil {
			continue
		}
		result = append(result, DisplayNameLocalizationModel{
			LanguageTag: convert.GraphToFrameworkString(localization.GetLanguageTag()),
			DisplayName: convert.GraphToFrameworkString(localization.GetDisplayName()),
		})
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

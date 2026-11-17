package graphBetaM365AdminProfileCardProperty

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructResource converts the Terraform resource model to a Kiota SDK model
// Returns a ProfileCardProperty that can be serialized by Kiota.
// directoryPropertyName is the key, so it is only sent on create; it cannot change afterwards.
func constructResource(ctx context.Context, data *ProfileCardPropertyResourceModel) (graphmodels.ProfileCardPropertyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	isCreate := data.ID.IsNull() || data.ID.IsUnknown()
	requestBody := constructRequestBody(data, isCreate)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructRequestBody builds the profile card property request body from the model
func constructRequestBody(data *ProfileCardPropertyResourceModel, isCreate bool) graphmodels.ProfileCardPropertyable {
	requestBody := graphmodels.NewProfileCardProperty()

	if isCreate {
		convert.FrameworkToGraphString(data.DirectoryPropertyName, requestBody.SetDirectoryPropertyName)
	}

	convert.FrameworkToGraphBool(data.IsVisible, requestBody.SetIsVisible)

	// Unset annotations are omitted on create. PATCH keeps omitted properties, so an update
	// sends an empty collection to clear annotations that were removed from configuration.
	if len(data.Annotations) > 0 || !isCreate {
		requestBody.SetAnnotations(constructAnnotations(data.Annotations))
	}

	return requestBody
}

// constructAnnotations builds the annotations collection using SDK models
func constructAnnotations(data []ProfileCardAnnotationModel) []graphmodels.ProfileCardAnnotationable {
	annotations := make([]graphmodels.ProfileCardAnnotationable, 0, len(data))

	for _, item := range data {
		annotation := graphmodels.NewProfileCardAnnotation()
		convert.FrameworkToGraphString(item.DisplayName, annotation.SetDisplayName)

		localizations := make([]graphmodels.DisplayNameLocalizationable, 0, len(item.Localizations))
		for _, localizationItem := range item.Localizations {
			localization := graphmodels.NewDisplayNameLocalization()
			convert.FrameworkToGraphString(localizationItem.LanguageTag, localization.SetLanguageTag)
			convert.FrameworkToGraphString(localizationItem.DisplayName, localization.SetDisplayName)
			localizations = append(localizations, localization)
		}
		annotation.SetLocalizations(localizations)

		annotations = append(annotations, annotation)
	}

	return annotations
}

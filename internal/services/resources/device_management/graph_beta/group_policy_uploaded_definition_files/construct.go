package graphBetaGroupPolicyUploadedDefinitionFiles

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *GroupPolicyUploadedDefinitionFileResourceModel, isCreate bool) (graphmodels.GroupPolicyUploadedDefinitionFileable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	definitionFile := graphmodels.NewGroupPolicyUploadedDefinitionFile()
	if err := constructGroupPolicyUploadedDefinitionFile(ctx, data, definitionFile); err != nil {
		return nil, fmt.Errorf("failed to construct group policy uploaded definition file: %s", err)
	}
	requestBody := definitionFile

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructGroupPolicyUploadedDefinitionFile handles specific settings using SDK setters
func constructGroupPolicyUploadedDefinitionFile(ctx context.Context, data *GroupPolicyUploadedDefinitionFileResourceModel, definitionFile *graphmodels.GroupPolicyUploadedDefinitionFile) error {
	convert.FrameworkToGraphString(data.FileName, definitionFile.SetFileName)
	convert.FrameworkToGraphBytes(data.Content, definitionFile.SetContent)
	convert.FrameworkToGraphString(data.DefaultLanguageCode, definitionFile.SetDefaultLanguageCode)

	if !data.GroupPolicyUploadedLanguageFiles.IsNull() && !data.GroupPolicyUploadedLanguageFiles.IsUnknown() {
		var languageFilesModels []GroupPolicyUploadedLanguageFileModel
		diags := data.GroupPolicyUploadedLanguageFiles.ElementsAs(ctx, &languageFilesModels, false)
		if diags.HasError() {
			return fmt.Errorf("failed to parse group policy uploaded language files: %v", diags.Errors())
		}

		if len(languageFilesModels) > 0 {
			languageFiles, err := constructGroupPolicyUploadedLanguageFiles(ctx, languageFilesModels)
			if err != nil {
				return fmt.Errorf("failed to construct group policy uploaded language files: %s", err)
			}
			definitionFile.SetGroupPolicyUploadedLanguageFiles(languageFiles)
		}
	}

	return nil
}

// constructGroupPolicyUploadedLanguageFiles constructs the language files for the definition file
func constructGroupPolicyUploadedLanguageFiles(ctx context.Context, languageFilesModels []GroupPolicyUploadedLanguageFileModel) ([]graphmodels.GroupPolicyUploadedLanguageFileable, error) {
	languageFiles := make([]graphmodels.GroupPolicyUploadedLanguageFileable, 0, len(languageFilesModels))
	for _, languageFile := range languageFilesModels {
		languageFileModel := graphmodels.NewGroupPolicyUploadedLanguageFile()

		convert.FrameworkToGraphString(languageFile.FileName, languageFileModel.SetFileName)
		convert.FrameworkToGraphString(languageFile.LanguageCode, languageFileModel.SetLanguageCode)
		convert.FrameworkToGraphBytes(languageFile.Content, languageFileModel.SetContent)
		languageFiles = append(languageFiles, languageFileModel)
	}
	return languageFiles, nil
}

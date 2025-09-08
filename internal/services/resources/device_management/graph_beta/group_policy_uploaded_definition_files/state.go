package graphBetaGroupPolicyUploadedDefinitionFiles

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote GraphServiceClient object to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *GroupPolicyUploadedDefinitionFileResourceModel, remoteResource graphmodels.GroupPolicyUploadedDefinitionFileable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.FileName = convert.GraphToFrameworkString(remoteResource.GetFileName())
	data.DefaultLanguageCode = convert.GraphToFrameworkString(remoteResource.GetDefaultLanguageCode())
	data.TargetPrefix = convert.GraphToFrameworkString(remoteResource.GetTargetPrefix())
	data.TargetNamespace = convert.GraphToFrameworkString(remoteResource.GetTargetNamespace())
	data.PolicyType = convert.GraphToFrameworkEnum(remoteResource.GetPolicyType())
	data.Revision = convert.GraphToFrameworkString(remoteResource.GetRevision())
	data.Status = convert.GraphToFrameworkEnum(remoteResource.GetStatus())

	// Handle time.Time for dates
	if uploadDateTime := remoteResource.GetUploadDateTime(); uploadDateTime != nil {
		data.UploadDateTime = types.StringValue(uploadDateTime.Format(time.RFC3339))
	} else {
		data.UploadDateTime = types.StringNull()
	}

	if lastModifiedDateTime := remoteResource.GetLastModifiedDateTime(); lastModifiedDateTime != nil {
		data.LastModifiedDateTime = types.StringValue(lastModifiedDateTime.Format(time.RFC3339))
	} else {
		data.LastModifiedDateTime = types.StringNull()
	}

	// Content is not returned in GET responses
	// data.Content = convert.GraphToFrameworkString(remoteResource.GetContent())

	// Map language codes
	if languageCodes := remoteResource.GetLanguageCodes(); languageCodes != nil {
		data.LanguageCodes = convert.GraphToFrameworkStringList(languageCodes)
	} else {
		data.LanguageCodes = types.ListNull(types.StringType)
	}

	// Map language files - preserve existing language files if they're not returned by the API
	// This is important because the API doesn't return language files in the GET response
	languageFiles := remoteResource.GetGroupPolicyUploadedLanguageFiles()

	if len(languageFiles) > 0 {
		// If we have language files in the response, map them
		mappedLanguageFiles, err := mapLanguageFilesToState(ctx, languageFiles)
		if err != nil {
			tflog.Error(ctx, "Failed to map language files", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			data.GroupPolicyUploadedLanguageFiles = mappedLanguageFiles
		}
	} else if data.GroupPolicyUploadedLanguageFiles.IsNull() {
		// If we don't have language files in the response and the current state is null,
		// set it to an empty set
		data.GroupPolicyUploadedLanguageFiles = types.SetValueMust(
			GroupPolicyUploadedLanguageFileType(),
			[]attr.Value{},
		)
	}
	// Otherwise, keep the existing language files in the state

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// GroupPolicyUploadedLanguageFileType returns the object type for GroupPolicyUploadedLanguageFileModel
func GroupPolicyUploadedLanguageFileType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"file_name":     types.StringType,
			"language_code": types.StringType,
			"content":       types.StringType,
		},
	}
}

// mapLanguageFilesToState maps language files from SDK to state
func mapLanguageFilesToState(ctx context.Context, languageFiles []graphmodels.GroupPolicyUploadedLanguageFileable) (types.Set, error) {
	languageFileType := GroupPolicyUploadedLanguageFileType()

	if len(languageFiles) == 0 {
		// Return an empty set rather than null
		return types.SetValueMust(languageFileType, []attr.Value{}), nil
	}

	languageFileValues := make([]attr.Value, 0, len(languageFiles))

	for _, languageFile := range languageFiles {
		// Ensure we have valid values for all fields
		fileName := convert.GraphToFrameworkString(languageFile.GetFileName())
		if fileName.IsNull() || fileName.IsUnknown() {
			fileName = types.StringValue("")
		}

		languageCode := convert.GraphToFrameworkString(languageFile.GetLanguageCode())
		if languageCode.IsNull() || languageCode.IsUnknown() {
			languageCode = types.StringValue("")
		}

		content := convert.GraphToFrameworkBytes(languageFile.GetContent())
		if content.IsNull() || content.IsUnknown() {
			// Content is not returned in GET responses, but we need to preserve it
			content = types.StringValue("")
		}

		languageFileAttrs := map[string]attr.Value{
			"file_name":     fileName,
			"language_code": languageCode,
			"content":       content,
		}

		languageFileValue, _ := types.ObjectValue(languageFileType.(types.ObjectType).AttrTypes, languageFileAttrs)
		languageFileValues = append(languageFileValues, languageFileValue)
	}

	set, diags := types.SetValue(languageFileType, languageFileValues)
	if diags.HasError() {
		return types.SetValueMust(languageFileType, []attr.Value{}), fmt.Errorf("failed to create language files set")
	}
	return set, nil
}

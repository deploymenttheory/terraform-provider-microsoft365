package graphBetaMobileApp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a mobile app to a model
func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.MobileAppable) MobileAppModel {
	tflog.Debug(ctx, "Starting to map remote resource state to Terraform state", map[string]any{
		"resourceName": data.GetDisplayName(),
		"resourceId":   data.GetId(),
	})

	model := MobileAppModel{
		ID:                    convert.GraphToFrameworkString(data.GetId()),
		DisplayName:           convert.GraphToFrameworkString(data.GetDisplayName()),
		Description:           convert.GraphToFrameworkString(data.GetDescription()),
		Publisher:             convert.GraphToFrameworkString(data.GetPublisher()),
		Developer:             convert.GraphToFrameworkString(data.GetDeveloper()),
		Owner:                 convert.GraphToFrameworkString(data.GetOwner()),
		Notes:                 convert.GraphToFrameworkString(data.GetNotes()),
		CreatedDateTime:       convert.GraphToFrameworkTime(data.GetCreatedDateTime()),
		LastModifiedDateTime:  convert.GraphToFrameworkTime(data.GetLastModifiedDateTime()),
		InformationUrl:        convert.GraphToFrameworkString(data.GetInformationUrl()),
		PrivacyInformationUrl: convert.GraphToFrameworkString(data.GetPrivacyInformationUrl()),
		IsAssigned:            convert.GraphToFrameworkBool(data.GetIsAssigned()),
		IsFeatured:            convert.GraphToFrameworkBool(data.GetIsFeatured()),
		UploadState:           convert.GraphToFrameworkInt32(data.GetUploadState()),
		PublishingState:       convert.GraphToFrameworkEnum(data.GetPublishingState()),
		DependentAppCount:     convert.GraphToFrameworkInt32(data.GetDependentAppCount()),
		SupersededAppCount:    convert.GraphToFrameworkInt32(data.GetSupersededAppCount()),
		SupersedingAppCount:   convert.GraphToFrameworkInt32(data.GetSupersedingAppCount()),
	}

	// Handle role scope tag IDs
	model.RoleScopeTagIds = convert.GraphToFrameworkStringSlice(data.GetRoleScopeTagIds())

	// Handle categories
	categories := data.GetCategories()
	if categories != nil {
		categoryValues := make([]types.String, 0, len(categories))
		for _, category := range categories {
			categoryValues = append(categoryValues, convert.GraphToFrameworkString(category.GetDisplayName()))
		}
		model.Categories = categoryValues
	} else {
		model.Categories = []types.String{}
	}

	return model
}

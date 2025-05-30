package graphBetaMobileApp

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a mobile app to a model
func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.MobileAppable) MobileAppModel {
	model := MobileAppModel{
		ID:                    state.StringPointerValue(data.GetId()),
		DisplayName:           state.StringPointerValue(data.GetDisplayName()),
		Description:           state.StringPointerValue(data.GetDescription()),
		Publisher:             state.StringPointerValue(data.GetPublisher()),
		Developer:             state.StringPointerValue(data.GetDeveloper()),
		Owner:                 state.StringPointerValue(data.GetOwner()),
		Notes:                 state.StringPointerValue(data.GetNotes()),
		CreatedDateTime:       state.TimeToString(data.GetCreatedDateTime()),
		LastModifiedDateTime:  state.TimeToString(data.GetLastModifiedDateTime()),
		InformationUrl:        state.StringPointerValue(data.GetInformationUrl()),
		PrivacyInformationUrl: state.StringPointerValue(data.GetPrivacyInformationUrl()),
		IsAssigned:            state.BoolPointerValue(data.GetIsAssigned()),
		IsFeatured:            state.BoolPointerValue(data.GetIsFeatured()),
		UploadState:           state.Int32PointerValue(data.GetUploadState()),
		PublishingState:       state.EnumPtrToTypeString(data.GetPublishingState()),
		DependentAppCount:     state.Int32PointerValue(data.GetDependentAppCount()),
		SupersededAppCount:    state.Int32PointerValue(data.GetSupersededAppCount()),
		SupersedingAppCount:   state.Int32PtrToTypeInt32(data.GetSupersedingAppCount()),
	}

	// Handle role scope tag IDs
	roleScopeTagIds := data.GetRoleScopeTagIds()
	if roleScopeTagIds != nil {
		model.RoleScopeTagIds = state.SliceToTypeStringSlice(roleScopeTagIds)
	} else {
		model.RoleScopeTagIds = []types.String{}
	}

	// Handle categories
	categories := data.GetCategories()
	if categories != nil {
		categoryValues := make([]types.String, 0, len(categories))
		for _, category := range categories {
			if category.GetDisplayName() != nil {
				categoryValues = append(categoryValues, types.StringValue(*category.GetDisplayName()))
			}
		}
		model.Categories = categoryValues
	} else {
		model.Categories = []types.String{}
	}

	return model
}

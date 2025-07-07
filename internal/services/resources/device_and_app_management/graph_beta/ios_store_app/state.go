package graphBetaIOSStoreApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapResourceToState maps the Graph API response to the Terraform state.
func mapResourceToState(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, graphResponse graphmodels.MobileAppable, data *IOSStoreAppResourceModel) error {
	iosStoreApp, ok := graphResponse.(graphmodels.IosStoreAppable)
	if !ok {
		return fmt.Errorf("expected IosStoreAppable but got %T", graphResponse)
	}

	// Map base properties
	data.ID = convert.GraphToFrameworkString(iosStoreApp.GetId())
	data.DisplayName = convert.GraphToFrameworkString(iosStoreApp.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(iosStoreApp.GetDescription())
	data.Publisher = convert.GraphToFrameworkString(iosStoreApp.GetPublisher())
	data.InformationUrl = convert.GraphToFrameworkString(iosStoreApp.GetInformationUrl())
	data.PrivacyInformationUrl = convert.GraphToFrameworkString(iosStoreApp.GetPrivacyInformationUrl())
	data.Owner = convert.GraphToFrameworkString(iosStoreApp.GetOwner())
	data.Developer = convert.GraphToFrameworkString(iosStoreApp.GetDeveloper())
	data.Notes = convert.GraphToFrameworkString(iosStoreApp.GetNotes())
	data.IsFeatured = convert.GraphToFrameworkBool(iosStoreApp.GetIsFeatured())
	data.CreatedDateTime = convert.GraphToFrameworkTime(iosStoreApp.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(iosStoreApp.GetLastModifiedDateTime())
	data.PublishingState = convert.GraphToFrameworkEnum(iosStoreApp.GetPublishingState())
	data.DependentAppCount = convert.GraphToFrameworkInt32(iosStoreApp.GetDependentAppCount())
	data.IsAssigned = convert.GraphToFrameworkBool(iosStoreApp.GetIsAssigned())
	data.SupersededAppCount = convert.GraphToFrameworkInt32(iosStoreApp.GetSupersededAppCount())
	data.SupersedingAppCount = convert.GraphToFrameworkInt32(iosStoreApp.GetSupersedingAppCount())
	data.UploadState = convert.GraphToFrameworkInt32(iosStoreApp.GetUploadState())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, iosStoreApp.GetRoleScopeTagIds())

	// Map iOS Store app specific properties
	data.AppStoreUrl = convert.GraphToFrameworkString(iosStoreApp.GetAppStoreUrl())
	//data.BundleId = convert.GraphToFrameworkString(iosStoreApp.GetBundleId())

	// Map applicable device type
	if applicableDeviceType := iosStoreApp.GetApplicableDeviceType(); applicableDeviceType != nil {
		if data.ApplicableDeviceType == nil {
			data.ApplicableDeviceType = &IOSDeviceTypeResourceModel{}
		}
		data.ApplicableDeviceType.IPad = convert.GraphToFrameworkBool(applicableDeviceType.GetIPad())
		data.ApplicableDeviceType.IPhoneAndIPod = convert.GraphToFrameworkBool(applicableDeviceType.GetIPhoneAndIPod())
	}

	// Map minimum supported OS
	if minOS := iosStoreApp.GetMinimumSupportedOperatingSystem(); minOS != nil {
		if data.MinimumSupportedOperatingSystem == nil {
			data.MinimumSupportedOperatingSystem = &IOSMinimumOperatingSystemResourceModel{}
		}
		data.MinimumSupportedOperatingSystem.V8_0 = convert.GraphToFrameworkBool(minOS.GetV80())
		data.MinimumSupportedOperatingSystem.V9_0 = convert.GraphToFrameworkBool(minOS.GetV90())
		data.MinimumSupportedOperatingSystem.V10_0 = convert.GraphToFrameworkBool(minOS.GetV100())
		data.MinimumSupportedOperatingSystem.V11_0 = convert.GraphToFrameworkBool(minOS.GetV110())
		data.MinimumSupportedOperatingSystem.V12_0 = convert.GraphToFrameworkBool(minOS.GetV120())
		data.MinimumSupportedOperatingSystem.V13_0 = convert.GraphToFrameworkBool(minOS.GetV130())
		data.MinimumSupportedOperatingSystem.V14_0 = convert.GraphToFrameworkBool(minOS.GetV140())
		data.MinimumSupportedOperatingSystem.V15_0 = convert.GraphToFrameworkBool(minOS.GetV150())
		data.MinimumSupportedOperatingSystem.V16_0 = convert.GraphToFrameworkBool(minOS.GetV160())
		data.MinimumSupportedOperatingSystem.V17_0 = convert.GraphToFrameworkBool(minOS.GetV170())
		data.MinimumSupportedOperatingSystem.V18_0 = convert.GraphToFrameworkBool(minOS.GetV180())
	}

	if data.AppIcon != nil {
		tflog.Debug(ctx, "Preserving original app_icon values from configuration")
	} else if largeIcon := iosStoreApp.GetLargeIcon(); largeIcon != nil {
		data.AppIcon = &sharedmodels.MobileAppIconResourceModel{
			IconFilePathSource: types.StringNull(),
			IconURLSource:      types.StringNull(),
		}
	} else {
		data.AppIcon = nil
	}

	// Map categories
	data.Categories = sharedstater.MapMobileAppCategoriesStateToTerraform(ctx, iosStoreApp.GetCategories())

	// Map relationships
	if relationships := iosStoreApp.GetRelationships(); len(relationships) > 0 {
		relationshipElements := make([]attr.Value, 0, len(relationships))
		for _, relationship := range relationships {
			relationshipAttrs := map[string]attr.Value{
				"id":                            convert.GraphToFrameworkString(relationship.GetId()),
				"source_display_name":           convert.GraphToFrameworkString(relationship.GetSourceDisplayName()),
				"source_display_version":        convert.GraphToFrameworkString(relationship.GetSourceDisplayVersion()),
				"source_id":                     convert.GraphToFrameworkString(relationship.GetSourceId()),
				"source_publisher_display_name": convert.GraphToFrameworkString(relationship.GetSourcePublisherDisplayName()),
				"target_display_name":           convert.GraphToFrameworkString(relationship.GetTargetDisplayName()),
				"target_display_version":        convert.GraphToFrameworkString(relationship.GetTargetDisplayVersion()),
				"target_id":                     convert.GraphToFrameworkString(relationship.GetTargetId()),
				"target_publisher":              convert.GraphToFrameworkString(relationship.GetTargetPublisher()),
				"target_publisher_display_name": convert.GraphToFrameworkString(relationship.GetTargetPublisherDisplayName()),
				"target_type":                   convert.GraphToFrameworkEnum(relationship.GetTargetType()),
			}

			element, diags := types.ObjectValue(
				map[string]attr.Type{
					"id":                            types.StringType,
					"source_display_name":           types.StringType,
					"source_display_version":        types.StringType,
					"source_id":                     types.StringType,
					"source_publisher_display_name": types.StringType,
					"target_display_name":           types.StringType,
					"target_display_version":        types.StringType,
					"target_id":                     types.StringType,
					"target_publisher":              types.StringType,
					"target_publisher_display_name": types.StringType,
					"target_type":                   types.StringType,
				},
				relationshipAttrs,
			)

			if diags.HasError() {
				return fmt.Errorf("error creating relationship object: %v", diags)
			}

			relationshipElements = append(relationshipElements, element)
		}

		relationshipsList, diags := types.ListValue(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":                            types.StringType,
					"source_display_name":           types.StringType,
					"source_display_version":        types.StringType,
					"source_id":                     types.StringType,
					"source_publisher_display_name": types.StringType,
					"target_display_name":           types.StringType,
					"target_display_version":        types.StringType,
					"target_id":                     types.StringType,
					"target_publisher":              types.StringType,
					"target_publisher_display_name": types.StringType,
					"target_type":                   types.StringType,
				},
			},
			relationshipElements,
		)

		if diags.HasError() {
			return fmt.Errorf("error creating relationships list: %v", diags)
		}

		data.Relationships = relationshipsList
	} else {
		// Initialize as empty list if no relationships
		data.Relationships = basetypes.NewListNull(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":                            types.StringType,
					"source_display_name":           types.StringType,
					"source_display_version":        types.StringType,
					"source_id":                     types.StringType,
					"source_publisher_display_name": types.StringType,
					"target_display_name":           types.StringType,
					"target_display_version":        types.StringType,
					"target_id":                     types.StringType,
					"target_publisher":              types.StringType,
					"target_publisher_display_name": types.StringType,
					"target_type":                   types.StringType,
				},
			},
		)
	}

	return nil
}

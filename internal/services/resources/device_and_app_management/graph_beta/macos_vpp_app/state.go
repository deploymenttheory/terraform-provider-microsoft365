package graphBetaMacOSVppApp

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
func mapResourceToState(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, graphResponse graphmodels.MobileAppable, data *MacOSVppAppResourceModel) error {
	macOsVppApp, ok := graphResponse.(graphmodels.MacOsVppAppable)
	if !ok {
		return fmt.Errorf("expected MacOsVppAppable but got %T", graphResponse)
	}

	// Map base properties
	data.ID = convert.GraphToFrameworkString(macOsVppApp.GetId())
	data.DisplayName = convert.GraphToFrameworkString(macOsVppApp.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(macOsVppApp.GetDescription())
	data.Publisher = convert.GraphToFrameworkString(macOsVppApp.GetPublisher())
	data.InformationUrl = convert.GraphToFrameworkString(macOsVppApp.GetInformationUrl())
	data.PrivacyInformationUrl = convert.GraphToFrameworkString(macOsVppApp.GetPrivacyInformationUrl())
	data.Owner = convert.GraphToFrameworkString(macOsVppApp.GetOwner())
	data.Developer = convert.GraphToFrameworkString(macOsVppApp.GetDeveloper())
	data.Notes = convert.GraphToFrameworkString(macOsVppApp.GetNotes())
	data.IsFeatured = convert.GraphToFrameworkBool(macOsVppApp.GetIsFeatured())
	data.CreatedDateTime = convert.GraphToFrameworkTime(macOsVppApp.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(macOsVppApp.GetLastModifiedDateTime())
	data.PublishingState = convert.GraphToFrameworkEnum(macOsVppApp.GetPublishingState())
	data.DependentAppCount = convert.GraphToFrameworkInt32(macOsVppApp.GetDependentAppCount())
	data.IsAssigned = convert.GraphToFrameworkBool(macOsVppApp.GetIsAssigned())
	data.SupersededAppCount = convert.GraphToFrameworkInt32(macOsVppApp.GetSupersededAppCount())
	data.SupersedingAppCount = convert.GraphToFrameworkInt32(macOsVppApp.GetSupersedingAppCount())
	data.UploadState = convert.GraphToFrameworkInt32(macOsVppApp.GetUploadState())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, macOsVppApp.GetRoleScopeTagIds())

	// Map VPP app specific properties
	data.UsedLicenseCount = convert.GraphToFrameworkInt32(macOsVppApp.GetUsedLicenseCount())
	data.TotalLicenseCount = convert.GraphToFrameworkInt32(macOsVppApp.GetTotalLicenseCount())
	data.ReleaseDateTime = convert.GraphToFrameworkTime(macOsVppApp.GetReleaseDateTime())
	data.AppStoreUrl = convert.GraphToFrameworkString(macOsVppApp.GetAppStoreUrl())
	data.BundleId = convert.GraphToFrameworkString(macOsVppApp.GetBundleId())
	data.VppTokenId = convert.GraphToFrameworkString(macOsVppApp.GetVppTokenId())
	data.VppTokenDisplayName = convert.GraphToFrameworkString(macOsVppApp.GetVppTokenDisplayName())
	data.VppTokenOrganizationName = convert.GraphToFrameworkString(macOsVppApp.GetVppTokenOrganizationName())
	data.VppTokenAppleId = convert.GraphToFrameworkString(macOsVppApp.GetVppTokenAppleId())

	if vppTokenAccountType := macOsVppApp.GetVppTokenAccountType(); vppTokenAccountType != nil {
		data.VppTokenAccountType = convert.GraphToFrameworkEnum(vppTokenAccountType)
	}

	if data.AppIcon != nil {
		tflog.Debug(ctx, "Preserving original app_icon values from configuration")
	} else if largeIcon := macOsVppApp.GetLargeIcon(); largeIcon != nil {
		data.AppIcon = &sharedmodels.MobileAppIconResourceModel{
			IconFilePathSource: types.StringNull(),
			IconURLSource:      types.StringNull(),
		}
	} else {
		data.AppIcon = nil
	}

	// Map licensing type
	if licensingType := macOsVppApp.GetLicensingType(); licensingType != nil {
		if data.LicensingType == nil {
			data.LicensingType = &VppLicensingTypeResourceModel{}
		}
		data.LicensingType.SupportUserLicensing = convert.GraphToFrameworkBool(licensingType.GetSupportUserLicensing())
		data.LicensingType.SupportDeviceLicensing = convert.GraphToFrameworkBool(licensingType.GetSupportDeviceLicensing())
		data.LicensingType.SupportsUserLicensing = convert.GraphToFrameworkBool(licensingType.GetSupportsUserLicensing())
		data.LicensingType.SupportsDeviceLicensing = convert.GraphToFrameworkBool(licensingType.GetSupportsDeviceLicensing())
	}

	// Map categories
	data.Categories = sharedstater.MapMobileAppCategoriesStateToTerraform(ctx, macOsVppApp.GetCategories())

	// Map relationships
	if relationships := macOsVppApp.GetRelationships(); len(relationships) > 0 {
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

	// Map assigned licenses
	if assignedLicenses := macOsVppApp.GetAssignedLicenses(); len(assignedLicenses) > 0 {
		licenseElements := make([]attr.Value, 0, len(assignedLicenses))
		for _, license := range assignedLicenses {
			licenseAttrs := map[string]attr.Value{
				"user_id":             convert.GraphToFrameworkString(license.GetUserId()),
				"user_email_address":  convert.GraphToFrameworkString(license.GetUserEmailAddress()),
				"user_name":           convert.GraphToFrameworkString(license.GetUserName()),
				"user_principal_name": convert.GraphToFrameworkString(license.GetUserPrincipalName()),
			}

			element, diags := types.ObjectValue(
				map[string]attr.Type{
					"user_id":             types.StringType,
					"user_email_address":  types.StringType,
					"user_name":           types.StringType,
					"user_principal_name": types.StringType,
				},
				licenseAttrs,
			)

			if diags.HasError() {
				return fmt.Errorf("error creating assigned license object: %v", diags)
			}

			licenseElements = append(licenseElements, element)
		}

		licensesList, diags := types.ListValue(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"user_id":             types.StringType,
					"user_email_address":  types.StringType,
					"user_name":           types.StringType,
					"user_principal_name": types.StringType,
				},
			},
			licenseElements,
		)

		if diags.HasError() {
			return fmt.Errorf("error creating assigned licenses list: %v", diags)
		}

		data.AssignedLicenses = licensesList
	} else {
		// Initialize as empty list if no assigned licenses
		data.AssignedLicenses = basetypes.NewListNull(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"user_id":             types.StringType,
					"user_email_address":  types.StringType,
					"user_name":           types.StringType,
					"user_principal_name": types.StringType,
				},
			},
		)
	}

	// Map revoke license action results
	if revokeLicenseResults := macOsVppApp.GetRevokeLicenseActionResults(); len(revokeLicenseResults) > 0 {
		resultElements := make([]attr.Value, 0, len(revokeLicenseResults))
		for _, result := range revokeLicenseResults {
			resultAttrs := map[string]attr.Value{
				"user_id":                convert.GraphToFrameworkString(result.GetUserId()),
				"managed_device_id":      convert.GraphToFrameworkString(result.GetManagedDeviceId()),
				"total_licenses_count":   convert.GraphToFrameworkInt32(result.GetTotalLicensesCount()),
				"failed_licenses_count":  convert.GraphToFrameworkInt32(result.GetFailedLicensesCount()),
				"action_failure_reason":  convert.GraphToFrameworkEnum(result.GetActionFailureReason()),
				"action_name":            convert.GraphToFrameworkString(result.GetActionName()),
				"action_state":           convert.GraphToFrameworkEnum(result.GetActionState()),
				"start_date_time":        convert.GraphToFrameworkTime(result.GetStartDateTime()),
				"last_updated_date_time": convert.GraphToFrameworkTime(result.GetLastUpdatedDateTime()),
			}

			element, diags := types.ObjectValue(
				map[string]attr.Type{
					"user_id":                types.StringType,
					"managed_device_id":      types.StringType,
					"total_licenses_count":   types.Int32Type,
					"failed_licenses_count":  types.Int32Type,
					"action_failure_reason":  types.StringType,
					"action_name":            types.StringType,
					"action_state":           types.StringType,
					"start_date_time":        types.StringType,
					"last_updated_date_time": types.StringType,
				},
				resultAttrs,
			)

			if diags.HasError() {
				return fmt.Errorf("error creating revoke license result object: %v", diags)
			}

			resultElements = append(resultElements, element)
		}

		resultsList, diags := types.ListValue(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"user_id":                types.StringType,
					"managed_device_id":      types.StringType,
					"total_licenses_count":   types.Int32Type,
					"failed_licenses_count":  types.Int32Type,
					"action_failure_reason":  types.StringType,
					"action_name":            types.StringType,
					"action_state":           types.StringType,
					"start_date_time":        types.StringType,
					"last_updated_date_time": types.StringType,
				},
			},
			resultElements,
		)

		if diags.HasError() {
			return fmt.Errorf("error creating revoke license results list: %v", diags)
		}

		data.RevokeLicenseActionResults = resultsList
	} else {
		// Initialize as empty list if no revoke license results
		data.RevokeLicenseActionResults = basetypes.NewListNull(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"user_id":                types.StringType,
					"managed_device_id":      types.StringType,
					"total_licenses_count":   types.Int32Type,
					"failed_licenses_count":  types.Int32Type,
					"action_failure_reason":  types.StringType,
					"action_name":            types.StringType,
					"action_state":           types.StringType,
					"start_date_time":        types.StringType,
					"last_updated_date_time": types.StringType,
				},
			},
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s to state", ResourceName, data.ID.ValueString()))
	return nil
}

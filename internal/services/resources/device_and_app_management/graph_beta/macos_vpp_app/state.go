package graphBetaMacOSVppApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_and_app_management"
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
		data.VppTokenAccountType = convert.GraphToFrameworkString(vppTokenAccountType)
	}

	// Map large icon
	if largeIcon := macOsVppApp.GetLargeIcon(); largeIcon != nil {
		if data.LargeIcon == nil {
			data.LargeIcon = &sharedmodels.MimeContentResourceModel{}
		}
		data.LargeIcon.Type = convert.GraphToFrameworkString(largeIcon.GetTypeEscaped())
		if value := largeIcon.GetValue(); value != nil {
			data.LargeIcon.Value = convert.GraphToFrameworkStringFromBytes(value)
		}
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

	// Map assignments
	assignments, err := sharedstater.MapMobileAppAssignmentsStateToTerraform(ctx, client, data.ID.ValueString())
	if err != nil {
		return fmt.Errorf("error mapping mobile app assignments: %v", err)
	}
	data.Assignments = assignments

	// Map relationships
	if relationships := macOsVppApp.GetRelationships(); len(relationships) > 0 {
		var mappedRelationships []MobileAppRelationshipResourceModel
		for _, relationship := range relationships {
			mappedRelationship := MobileAppRelationshipResourceModel{
				ID:                         convert.GraphToFrameworkString(relationship.GetId()),
				SourceDisplayName:          convert.GraphToFrameworkString(relationship.GetSourceDisplayName()),
				SourceDisplayVersion:       convert.GraphToFrameworkString(relationship.GetSourceDisplayVersion()),
				SourceId:                   convert.GraphToFrameworkString(relationship.GetSourceId()),
				SourcePublisherDisplayName: convert.GraphToFrameworkString(relationship.GetSourcePublisherDisplayName()),
				TargetDisplayName:          convert.GraphToFrameworkString(relationship.GetTargetDisplayName()),
				TargetDisplayVersion:       convert.GraphToFrameworkString(relationship.GetTargetDisplayVersion()),
				TargetId:                   convert.GraphToFrameworkString(relationship.GetTargetId()),
				TargetPublisher:            convert.GraphToFrameworkString(relationship.GetTargetPublisher()),
				TargetPublisherDisplayName: convert.GraphToFrameworkString(relationship.GetTargetPublisherDisplayName()),
				TargetType:                 convert.GraphToFrameworkEnum(relationship.GetTargetType()),
			}
			mappedRelationships = append(mappedRelationships, mappedRelationship)
		}
		data.Relationships = mappedRelationships
	}

	// Map assigned licenses
	if assignedLicenses := macOsVppApp.GetAssignedLicenses(); len(assignedLicenses) > 0 {
		var mappedLicenses []MacOSVppAppAssignedLicenseResourceModel
		for _, license := range assignedLicenses {
			mappedLicense := MacOSVppAppAssignedLicenseResourceModel{
				UserId:           convert.GraphToFrameworkString(license.GetUserId()),
				DeviceId:         convert.GraphToFrameworkString(license.GetDeviceId()),
				LicenseType:      convert.GraphToFrameworkEnum(license.GetLicenseType()),
				UserEmailAddress: convert.GraphToFrameworkString(license.GetUserEmailAddress()),
			}
			mappedLicenses = append(mappedLicenses, mappedLicense)
		}
		data.AssignedLicenses = mappedLicenses
	}

	// Map revoke license action results
	if revokeLicenseResults := macOsVppApp.GetRevokeLicenseActionResults(); len(revokeLicenseResults) > 0 {
		var mappedResults []MacOSVppAppRevokeLicensesActionResultResourceModel
		for _, result := range revokeLicenseResults {
			mappedResult := MacOSVppAppRevokeLicensesActionResultResourceModel{
				UserId:              convert.GraphToFrameworkString(result.GetUserId()),
				ManagedDeviceId:     convert.GraphToFrameworkString(result.GetManagedDeviceId()),
				TotalLicensesCount:  convert.GraphToFrameworkInt32(result.GetTotalLicensesCount()),
				FailedLicensesCount: convert.GraphToFrameworkInt32(result.GetFailedLicensesCount()),
				ActionFailureReason: convert.GraphToFrameworkString(result.GetActionFailureReason()),
				ActionName:          convert.GraphToFrameworkString(result.GetActionName()),
				ActionState:         convert.GraphToFrameworkString(result.GetActionState()),
				StartDateTime:       convert.GraphToFrameworkTime(result.GetStartDateTime()),
				LastUpdatedDateTime: convert.GraphToFrameworkTime(result.GetLastUpdatedDateTime()),
			}
			mappedResults = append(mappedResults, mappedResult)
		}
		data.RevokeLicenseActionResults = mappedResults
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s to state", ResourceName, data.ID.ValueString()))
	return nil
}

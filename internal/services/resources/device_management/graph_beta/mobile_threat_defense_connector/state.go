package graphBetaMobileThreatDefenseConnector

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the remote mobile threat defense connector from Kiota SDK to Terraform state
func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data MobileThreatDefenseConnectorResourceModel,
	remoteResource graphmodels.MobileThreatDefenseConnectorable,
) MobileThreatDefenseConnectorResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.AndroidEnabled = convert.GraphToFrameworkBool(remoteResource.GetAndroidEnabled())
	data.IosEnabled = convert.GraphToFrameworkBool(remoteResource.GetIosEnabled())
	data.MacEnabled = convert.GraphToFrameworkBool(remoteResource.GetMacEnabled())
	data.WindowsEnabled = convert.GraphToFrameworkBool(remoteResource.GetWindowsEnabled())
	data.AndroidMobileApplicationManagementEnabled = convert.GraphToFrameworkBool(remoteResource.GetAndroidMobileApplicationManagementEnabled())
	data.IosMobileApplicationManagementEnabled = convert.GraphToFrameworkBool(remoteResource.GetIosMobileApplicationManagementEnabled())
	data.WindowsMobileApplicationManagementEnabled = convert.GraphToFrameworkBool(remoteResource.GetWindowsMobileApplicationManagementEnabled())
	data.PartnerUnsupportedOsVersionBlocked = convert.GraphToFrameworkBool(remoteResource.GetPartnerUnsupportedOsVersionBlocked())
	data.PartnerUnresponsivenessThresholdInDays = convert.GraphToFrameworkInt32(remoteResource.GetPartnerUnresponsivenessThresholdInDays())
	data.AllowPartnerToCollectIosApplicationMetadata = convert.GraphToFrameworkBool(remoteResource.GetAllowPartnerToCollectIOSApplicationMetadata())
	data.AllowPartnerToCollectIosPersonalApplicationMetadata = convert.GraphToFrameworkBool(remoteResource.GetAllowPartnerToCollectIOSPersonalApplicationMetadata())
	data.AllowPartnerToCollectIosCertificateMetadata = convert.GraphToFrameworkBool(remoteResource.GetAllowPartnerToCollectIosCertificateMetadata())
	data.AllowPartnerToCollectIosPersonalCertificateMetadata = convert.GraphToFrameworkBool(remoteResource.GetAllowPartnerToCollectIosPersonalCertificateMetadata())
	data.MicrosoftDefenderForEndpointAttachEnabled = convert.GraphToFrameworkBool(remoteResource.GetMicrosoftDefenderForEndpointAttachEnabled())
	data.GrantMobileThreatDefensePartnerRole = convert.GraphToFrameworkBool(remoteResource.GetGrantMobileThreatDefensePartnerRole())
	data.LaunchMobileThreatDefensePartnerOnSetupEnabled = convert.GraphToFrameworkBool(remoteResource.GetLaunchMobileThreatDefensePartnerOnSetupEnabled())
	data.AndroidDeviceBlockedOnMissingPartnerData = convert.GraphToFrameworkBool(remoteResource.GetAndroidDeviceBlockedOnMissingPartnerData())
	data.IosDeviceBlockedOnMissingPartnerData = convert.GraphToFrameworkBool(remoteResource.GetIosDeviceBlockedOnMissingPartnerData())
	data.MacDeviceBlockedOnMissingPartnerData = convert.GraphToFrameworkBool(remoteResource.GetMacDeviceBlockedOnMissingPartnerData())
	data.WindowsDeviceBlockedOnMissingPartnerData = convert.GraphToFrameworkBool(remoteResource.GetWindowsDeviceBlockedOnMissingPartnerData())
	data.PartnerState = convert.GraphToFrameworkEnum(remoteResource.GetPartnerState())
	data.LastHeartbeatDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastHeartbeatDateTime())
	data.PartnerDisplayName = convert.GraphToFrameworkString(remoteResource.GetPartnerDisplayName())
	data.PartnerSupportsAndroidMam = convert.GraphToFrameworkBool(remoteResource.GetPartnerSupportsAndroidMam())
	data.PartnerSupportsIosMam = convert.GraphToFrameworkBool(remoteResource.GetPartnerSupportsIosMam())
	data.PartnerSupportsWindowsMam = convert.GraphToFrameworkBool(remoteResource.GetPartnerSupportsWindowsMam())
	data.PartnerSupportsIosCertificateSync = convert.GraphToFrameworkBool(remoteResource.GetPartnerSupportsIosCertificateSync())
	data.PartnerMinimumSupportedAndroidVersion = convert.GraphToFrameworkString(remoteResource.GetPartnerMinimumSupportedAndroidVersion())
	data.PartnerMinimumSupportedIosVersion = convert.GraphToFrameworkString(remoteResource.GetPartnerMinimumSupportedIosVersion())
	data.PartnerMinimumSupportedMacVersion = convert.GraphToFrameworkString(remoteResource.GetPartnerMinimumSupportedMacVersion())
	data.PartnerMinimumSupportedWindowsVersion = convert.GraphToFrameworkString(remoteResource.GetPartnerMinimumSupportedWindowsVersion())

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state")

	return data
}

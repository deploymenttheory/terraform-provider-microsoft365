package graphBetaWindowsAutopilotDeviceIdentity

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote Windows Autopilot Device Identity resource state to the Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsAutopilotDeviceIdentityResourceModel, remoteResource graphmodels.WindowsAutopilotDeviceIdentityable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map Windows Autopilot Device Identity from API to Terraform state")

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.GroupTag = convert.GraphToFrameworkString(remoteResource.GetGroupTag())
	data.PurchaseOrderIdentifier = convert.GraphToFrameworkString(remoteResource.GetPurchaseOrderIdentifier())
	data.SerialNumber = convert.GraphToFrameworkString(remoteResource.GetSerialNumber())
	data.ProductKey = convert.GraphToFrameworkString(remoteResource.GetProductKey())
	data.Manufacturer = convert.GraphToFrameworkString(remoteResource.GetManufacturer())
	data.Model = convert.GraphToFrameworkString(remoteResource.GetModel())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.UserPrincipalName = convert.GraphToFrameworkString(remoteResource.GetUserPrincipalName())

	// Computed properties
	data.AddressableUserName = convert.GraphToFrameworkString(remoteResource.GetAddressableUserName())
	data.ResourceName = convert.GraphToFrameworkString(remoteResource.GetResourceName())
	data.SkuNumber = convert.GraphToFrameworkString(remoteResource.GetSkuNumber())
	data.SystemFamily = convert.GraphToFrameworkString(remoteResource.GetSystemFamily())
	data.AzureActiveDirectoryDeviceId = convert.GraphToFrameworkString(remoteResource.GetAzureActiveDirectoryDeviceId())
	data.AzureAdDeviceId = convert.GraphToFrameworkString(remoteResource.GetAzureAdDeviceId())
	data.ManagedDeviceId = convert.GraphToFrameworkString(remoteResource.GetManagedDeviceId())

	// Enrollment and profile status
	if enrollmentState := remoteResource.GetEnrollmentState(); enrollmentState != nil {
		data.EnrollmentState = convert.GraphToFrameworkEnum(enrollmentState)
	}

	data.LastContactedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastContactedDateTime())
	data.DeploymentProfileAssignedDateTime = convert.GraphToFrameworkTime(remoteResource.GetDeploymentProfileAssignedDateTime())

	if deploymentProfileAssignmentStatus := remoteResource.GetDeploymentProfileAssignmentStatus(); deploymentProfileAssignmentStatus != nil {
		data.DeploymentProfileAssignmentStatus = convert.GraphToFrameworkEnum(deploymentProfileAssignmentStatus)
	}

	if deploymentProfileAssignmentDetailedStatus := remoteResource.GetDeploymentProfileAssignmentDetailedStatus(); deploymentProfileAssignmentDetailedStatus != nil {
		data.DeploymentProfileAssignmentDetailedStatus = convert.GraphToFrameworkEnum(deploymentProfileAssignmentDetailedStatus)
	}

	// Remediation state
	if remediationState := remoteResource.GetRemediationState(); remediationState != nil {
		data.RemediationState = convert.GraphToFrameworkEnum(remediationState)
	}

	data.RemediationStateLastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetRemediationStateLastModifiedDateTime())

	// Userless enrollment status
	if userlessEnrollmentStatus := remoteResource.GetUserlessEnrollmentStatus(); userlessEnrollmentStatus != nil {
		data.UserlessEnrollmentStatus = convert.GraphToFrameworkEnum(userlessEnrollmentStatus)
	}

	tflog.Debug(ctx, "Finished mapping Windows Autopilot Device Identity from API to Terraform state")
}

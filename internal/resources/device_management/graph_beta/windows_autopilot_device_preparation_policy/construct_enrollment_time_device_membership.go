package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructEnrollmentTimeDeviceMembershipResource constructs the request body for setting the enrollment time device membership target
func constructEnrollmentTimeDeviceMembershipResource(ctx context.Context, deviceSecurityGroupID string) (devicemanagement.ConfigurationPoliciesItemSetEnrollmentTimeDeviceMembershipTargetPostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing enrollment time device membership target with security group: %s", deviceSecurityGroupID))

	// Create an EnrollmentTimeDeviceMembershipTarget object
	target := models.NewEnrollmentTimeDeviceMembershipTarget()

	// Set the target ID to the security group ID
	targetID := deviceSecurityGroupID
	target.SetTargetId(&targetID)

	// Set the target type to static security group
	targetType := models.STATICSECURITYGROUP_ENROLLMENTTIMEDEVICEMEMBERSHIPTARGETTYPE
	target.SetTargetType(&targetType)

	// Create the request body with the targets
	requestBody := devicemanagement.NewConfigurationPoliciesItemSetEnrollmentTimeDeviceMembershipTargetPostRequestBody()
	targets := []models.EnrollmentTimeDeviceMembershipTargetable{target}
	requestBody.SetEnrollmentTimeDeviceMembershipTargets(targets)

	return requestBody, nil
}

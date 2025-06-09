package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructJustInTimeAssignmentBody constructs the request body for assigning an enrollment time device membership target
func constructJustInTimeAssignmentBody(ctx context.Context, deviceSecurityGroupID string) (*devicemanagement.ConfigurationPoliciesItemSetEnrollmentTimeDeviceMembershipTargetPostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing enrollment time device membership target with security group: %s", deviceSecurityGroupID))

	targetType := models.STATICSECURITYGROUP_ENROLLMENTTIMEDEVICEMEMBERSHIPTARGETTYPE
	target := models.NewEnrollmentTimeDeviceMembershipTarget()
	target.SetTargetType(&targetType)
	target.SetTargetId(&deviceSecurityGroupID)

	body := devicemanagement.NewConfigurationPoliciesItemSetEnrollmentTimeDeviceMembershipTargetPostRequestBody()
	body.SetEnrollmentTimeDeviceMembershipTargets([]models.EnrollmentTimeDeviceMembershipTargetable{target})

	tflog.Debug(ctx, "Finished constructing enrollment time device membership target body")
	return body, nil
}

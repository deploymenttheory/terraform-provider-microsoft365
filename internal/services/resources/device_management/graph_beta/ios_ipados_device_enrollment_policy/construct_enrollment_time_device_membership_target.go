package graphBetaIOSiPadOSDeviceEnrollmentPolicy

import (
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructEnrollmentTimeDeviceMembershipTargetBody builds the request body for
// setEnrollmentTimeDeviceMembershipTarget, targeting the given static security group.
func constructEnrollmentTimeDeviceMembershipTargetBody(
	deviceSecurityGroupID string,
) *devicemanagement.ConfigurationPoliciesItemSetEnrollmentTimeDeviceMembershipTargetPostRequestBody {
	targetType := models.STATICSECURITYGROUP_ENROLLMENTTIMEDEVICEMEMBERSHIPTARGETTYPE
	target := models.NewEnrollmentTimeDeviceMembershipTarget()
	target.SetTargetType(&targetType)
	target.SetTargetId(&deviceSecurityGroupID)

	body := devicemanagement.NewConfigurationPoliciesItemSetEnrollmentTimeDeviceMembershipTargetPostRequestBody()
	body.SetEnrollmentTimeDeviceMembershipTargets(
		[]models.EnrollmentTimeDeviceMembershipTargetable{target},
	)

	return body
}

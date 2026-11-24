package graphBetaDeviceEnrollmentNotification_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaDeviceEnrollmentNotification "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/device_enrollment_notification"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// TestUnitResourceDeviceEnrollmentNotification_ImportRejectsInvalidID checks that ImportState rejects an ID that
// doesn't match the resource's ID format, before the ID can reach a Graph request URL.
func TestUnitResourceDeviceEnrollmentNotification_ImportRejectsInvalidID(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.NewMocks().AuthMocks.RegisterMocks()

	tests := []struct {
		name     string
		importID string
	}{
		{name: "plain GUID", importID: "00000000-0000-0000-0000-000000000000"},
		{name: "query string", importID: "00000000-0000-0000-0000-000000000000_EnrollmentNotificationsConfiguration?$select=id"},
		{name: "path traversal", importID: "../00000000-0000-0000-0000-000000000000_EnrollmentNotificationsConfiguration"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:        `resource "` + graphBetaDeviceEnrollmentNotification.ResourceName + `" "test" {}`,
						ResourceName:  graphBetaDeviceEnrollmentNotification.ResourceName + ".test",
						ImportState:   true,
						ImportStateId: tt.importID,
						ExpectError:   regexp.MustCompile(`Invalid Import ID`),
					},
				},
			})
		})
	}
}

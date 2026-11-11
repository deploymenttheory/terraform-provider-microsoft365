package graphBetaWindowsEnrollmentStatusPage

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func TestUnitWindowsEnrollmentStatusPageAssignmentFilterState(t *testing.T) {
	for _, tc := range []struct {
		name string
		id   *string
		want types.String
	}{
		{"omitted", nil, types.StringNull()},
		{"empty", new(""), types.StringNull()},
		{"unfiltered", new("00000000-0000-0000-0000-000000000000"), types.StringNull()},
		{"filtered", new("11111111-1111-4111-8111-111111111111"), types.StringValue("11111111-1111-4111-8111-111111111111")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			target := graphmodels.NewGroupAssignmentTarget()
			target.SetGroupId(new("22222222-2222-4222-8222-222222222222"))
			target.SetDeviceAndAppManagementAssignmentFilterId(tc.id)
			assignment := graphmodels.NewEnrollmentConfigurationAssignment()
			assignment.SetTarget(target)
			var data WindowsEnrollmentStatusPageResourceModel
			MapAssignmentsToTerraform(context.Background(), &data, []graphmodels.EnrollmentConfigurationAssignmentable{assignment})
			got := data.Assignments.Elements()[0].(types.Object).Attributes()["filter_id"]
			if !got.Equal(tc.want) {
				t.Fatalf("filter_id = %s, want %s", got, tc.want)
			}
		})
	}
}

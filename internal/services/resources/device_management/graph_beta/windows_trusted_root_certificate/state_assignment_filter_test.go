package graphBetaWindowsTrustedRootCertificate

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func TestUnitAssignmentFilterState(t *testing.T) {
	for _, tc := range []struct {
		name string
		id   *string
		want types.String
	}{
		{"omitted", nil, types.StringNull()},
		{"empty", new(""), types.StringNull()},
		{"zero GUID", new("00000000-0000-0000-0000-000000000000"), types.StringNull()},
		{"real filter", new("11111111-1111-4111-8111-111111111111"), types.StringValue("11111111-1111-4111-8111-111111111111")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			target := graphmodels.NewGroupAssignmentTarget()
			target.SetGroupId(new("22222222-2222-4222-8222-222222222222"))
			target.SetDeviceAndAppManagementAssignmentFilterId(tc.id)
			value, err := mapAssignment(target)
			require.NoError(t, err)
			got := value.(types.Object).Attributes()["filter_id"]
			require.True(t, got.Equal(tc.want), "filter_id = %s, want %s", got, tc.want)
		})
	}
}

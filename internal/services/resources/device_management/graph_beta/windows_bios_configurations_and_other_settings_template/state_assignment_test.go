package graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func groupAssignment(groupID string, filterID *string, filterType *graphmodels.DeviceAndAppManagementAssignmentFilterType) graphmodels.HardwareConfigurationAssignmentable {
	target := graphmodels.NewGroupAssignmentTarget()
	target.SetGroupId(&groupID)
	target.SetDeviceAndAppManagementAssignmentFilterId(filterID)
	target.SetDeviceAndAppManagementAssignmentFilterType(filterType)

	a := graphmodels.NewHardwareConfigurationAssignment()
	a.SetTarget(target)
	return a
}

func exclusionAssignment(groupID string) graphmodels.HardwareConfigurationAssignmentable {
	target := graphmodels.NewExclusionGroupAssignmentTarget()
	target.SetGroupId(&groupID)

	a := graphmodels.NewHardwareConfigurationAssignment()
	a.SetTarget(target)
	return a
}

func assignmentAttrs(t *testing.T, data *WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel, i int) map[string]string {
	t.Helper()

	elems := data.Assignments.Elements()
	if i >= len(elems) {
		t.Fatalf("assignment index %d out of range (%d elements)", i, len(elems))
	}

	obj, ok := elems[i].(types.Object)
	if !ok {
		t.Fatalf("assignment %d is %T, want types.Object", i, elems[i])
	}

	out := map[string]string{}
	for k, v := range obj.Attributes() {
		s, ok := v.(types.String)
		if !ok {
			t.Fatalf("attribute %q is %T, want types.String", k, v)
		}
		out[k] = s.ValueString()
	}
	return out
}

func TestUnitMapAssignmentsToTerraformWithNoAssignmentsIsNull(t *testing.T) {
	t.Parallel()

	data := &WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel{}
	MapAssignmentsToTerraform(context.Background(), data, nil)

	if !data.Assignments.IsNull() {
		t.Error("assignments should be null when there are none")
	}
}

func TestUnitMapAssignmentsToTerraformAppliesSchemaDefaultsForUnfilteredTargets(t *testing.T) {
	t.Parallel()

	include := graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
	exclude := graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
	none := graphmodels.NONE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
	filterID := "44444444-4444-4444-4444-444444444444"
	sentinel := "00000000-0000-0000-0000-000000000000"

	data := &WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel{}
	MapAssignmentsToTerraform(context.Background(), data, []graphmodels.HardwareConfigurationAssignmentable{
		groupAssignment("11111111-1111-1111-1111-111111111111", &filterID, &include),
		groupAssignment("22222222-2222-2222-2222-222222222222", &filterID, &exclude),
		// The API echoes the all-zero sentinel and "none"; both must map back to the schema defaults.
		groupAssignment("33333333-3333-3333-3333-333333333333", &sentinel, &none),
		// A target with no filter information at all must also land on the defaults.
		groupAssignment("44444444-4444-4444-4444-444444444444", nil, nil),
		exclusionAssignment("66666666-6666-6666-6666-666666666666"),
	})

	if got := len(data.Assignments.Elements()); got != 5 {
		t.Fatalf("assignments length = %d, want 5", got)
	}

	for i, want := range []map[string]string{
		{"type": "groupAssignmentTarget", "group_id": "11111111-1111-1111-1111-111111111111", "filter_id": filterID, "filter_type": "include"},
		{"type": "groupAssignmentTarget", "group_id": "22222222-2222-2222-2222-222222222222", "filter_id": filterID, "filter_type": "exclude"},
		{"type": "groupAssignmentTarget", "group_id": "33333333-3333-3333-3333-333333333333", "filter_id": sentinel, "filter_type": "none"},
		{"type": "groupAssignmentTarget", "group_id": "44444444-4444-4444-4444-444444444444", "filter_id": sentinel, "filter_type": "none"},
		{"type": "exclusionGroupAssignmentTarget", "group_id": "66666666-6666-6666-6666-666666666666", "filter_id": sentinel, "filter_type": "none"},
	} {
		got := assignmentAttrs(t, data, i)
		for k, v := range want {
			if got[k] != v {
				t.Errorf("assignment %d attribute %q = %q, want %q", i, k, got[k], v)
			}
		}
	}
}

func TestUnitMapAssignmentsToTerraformSkipsUnusableAssignments(t *testing.T) {
	t.Parallel()

	unsupported := graphmodels.NewHardwareConfigurationAssignment()
	unsupported.SetTarget(graphmodels.NewAllDevicesAssignmentTarget())

	nilTarget := graphmodels.NewHardwareConfigurationAssignment()

	noOdataType := graphmodels.NewHardwareConfigurationAssignment()
	bareTarget := graphmodels.NewDeviceAndAppManagementAssignmentTarget()
	bareTarget.SetOdataType(nil)
	noOdataType.SetTarget(bareTarget)

	data := &WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel{}
	MapAssignmentsToTerraform(context.Background(), data, []graphmodels.HardwareConfigurationAssignmentable{
		nilTarget,
		noOdataType,
		unsupported,
		exclusionAssignment("66666666-6666-6666-6666-666666666666"),
	})

	if got := len(data.Assignments.Elements()); got != 1 {
		t.Fatalf("assignments length = %d, want 1 (the three unusable entries should be skipped)", got)
	}
	if got := assignmentAttrs(t, data, 0)["type"]; got != "exclusionGroupAssignmentTarget" {
		t.Errorf("assignment 0 type = %q, want exclusionGroupAssignmentTarget", got)
	}
}

// An assignment whose target carries no group id still maps, with group_id left null.
func TestUnitMapAssignmentsToTerraformHandlesMissingGroupId(t *testing.T) {
	t.Parallel()

	target := graphmodels.NewGroupAssignmentTarget()
	target.SetGroupId(nil)
	a := graphmodels.NewHardwareConfigurationAssignment()
	a.SetTarget(target)

	data := &WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel{}
	MapAssignmentsToTerraform(context.Background(), data, []graphmodels.HardwareConfigurationAssignmentable{a})

	if got := len(data.Assignments.Elements()); got != 1 {
		t.Fatalf("assignments length = %d, want 1", got)
	}
	if got := assignmentAttrs(t, data, 0)["group_id"]; got != "" {
		t.Errorf("group_id = %q, want it to be null", got)
	}
}

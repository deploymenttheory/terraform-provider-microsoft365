package graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate

import (
	"context"
	"errors"
	"testing"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

const noFilter = "00000000-0000-0000-0000-000000000000"

// assignmentSet builds an assignments set matching the resource schema.
func assignmentSet(t *testing.T, objects ...map[string]attr.Value) types.Set {
	t.Helper()

	objType := WindowsBiosConfigurationsAndOtherSettingsTemplateAssignmentType()
	elems := make([]attr.Value, len(objects))
	for i, o := range objects {
		obj, diags := types.ObjectValue(objType.AttrTypes, o)
		if diags.HasError() {
			t.Fatalf("failed to build assignment object: %v", diags.Errors())
		}
		elems[i] = obj
	}

	set, diags := types.SetValue(objType, elems)
	if diags.HasError() {
		t.Fatalf("failed to build assignments set: %v", diags.Errors())
	}
	return set
}

func assignment(targetType, groupID, filterID, filterType string) map[string]attr.Value {
	return map[string]attr.Value{
		"type":        types.StringValue(targetType),
		"group_id":    types.StringValue(groupID),
		"filter_id":   types.StringValue(filterID),
		"filter_type": types.StringValue(filterType),
	}
}

func TestUnitConstructAssignmentWithNoAssignmentsSendsEmptyCollection(t *testing.T) {
	t.Parallel()

	for name, set := range map[string]types.Set{
		"null":    types.SetNull(WindowsBiosConfigurationsAndOtherSettingsTemplateAssignmentType()),
		"unknown": types.SetUnknown(WindowsBiosConfigurationsAndOtherSettingsTemplateAssignmentType()),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			body, err := constructAssignment(context.Background(),
				&WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel{Assignments: set})
			if err != nil {
				t.Fatalf("constructAssignment() returned an unexpected error: %v", err)
			}

			got := body.GetHardwareConfigurationAssignments()
			if got == nil {
				t.Fatal("hardwareConfigurationAssignments was nil, want an empty collection")
			}
			if len(got) != 0 {
				t.Errorf("hardwareConfigurationAssignments length = %d, want 0", len(got))
			}
		})
	}
}

func TestUnitConstructAssignmentBuildsGroupAndExclusionTargets(t *testing.T) {
	t.Parallel()

	data := &WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel{
		Assignments: assignmentSet(t,
			assignment("groupAssignmentTarget", "11111111-1111-1111-1111-111111111111", "44444444-4444-4444-4444-444444444444", "include"),
			assignment("groupAssignmentTarget", "22222222-2222-2222-2222-222222222222", "55555555-5555-5555-5555-555555555555", "exclude"),
			assignment("groupAssignmentTarget", "33333333-3333-3333-3333-333333333333", noFilter, "none"),
			assignment("exclusionGroupAssignmentTarget", "66666666-6666-6666-6666-666666666666", noFilter, "none"),
		),
	}

	body, err := constructAssignment(context.Background(), data)
	if err != nil {
		t.Fatalf("constructAssignment() returned an unexpected error: %v", err)
	}

	assignments := body.GetHardwareConfigurationAssignments()
	if len(assignments) != 4 {
		t.Fatalf("hardwareConfigurationAssignments length = %d, want 4", len(assignments))
	}

	includeTarget := assignments[0].GetTarget()
	if _, ok := includeTarget.(graphmodels.GroupAssignmentTargetable); !ok {
		t.Errorf("assignment 0 target type = %T, want GroupAssignmentTargetable", includeTarget)
	}
	if got := includeTarget.GetDeviceAndAppManagementAssignmentFilterId(); got == nil || *got != "44444444-4444-4444-4444-444444444444" {
		t.Errorf("assignment 0 filter id = %v, want the configured filter", got)
	}
	if got := includeTarget.GetDeviceAndAppManagementAssignmentFilterType(); got == nil ||
		*got != graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE {
		t.Errorf("assignment 0 filter type = %v, want include", got)
	}

	if got := assignments[1].GetTarget().GetDeviceAndAppManagementAssignmentFilterType(); got == nil ||
		*got != graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE {
		t.Errorf("assignment 1 filter type = %v, want exclude", got)
	}

	// The sentinel filter id means "no filter" and must not reach the API.
	if got := assignments[2].GetTarget().GetDeviceAndAppManagementAssignmentFilterId(); got != nil {
		t.Errorf("assignment 2 filter id = %q, want it to be omitted", *got)
	}

	if _, ok := assignments[3].GetTarget().(graphmodels.ExclusionGroupAssignmentTargetable); !ok {
		t.Errorf("assignment 3 target type = %T, want ExclusionGroupAssignmentTargetable", assignments[3].GetTarget())
	}
}

func TestUnitConstructAssignmentRejectsInvalidAssignments(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		object  map[string]attr.Value
		wantErr error
	}{
		"missing type": {
			object:  assignment("", "11111111-1111-1111-1111-111111111111", noFilter, "none"),
			wantErr: errUnsupportedTargetType,
		},
		"null type": {
			object: map[string]attr.Value{
				"type":        types.StringNull(),
				"group_id":    types.StringValue("11111111-1111-1111-1111-111111111111"),
				"filter_id":   types.StringValue(noFilter),
				"filter_type": types.StringValue("none"),
			},
			wantErr: errAssignmentMissingType,
		},
		"group target without group id": {
			object:  assignment("groupAssignmentTarget", "", noFilter, "none"),
			wantErr: errAssignmentMissingGroupId,
		},
		"exclusion target without group id": {
			object:  assignment("exclusionGroupAssignmentTarget", "", noFilter, "none"),
			wantErr: errAssignmentMissingGroupId,
		},
		"unsupported target": {
			object:  assignment("allDevicesAssignmentTarget", "", noFilter, "none"),
			wantErr: errUnsupportedTargetType,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data := &WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel{
				Assignments: assignmentSet(t, tc.object),
			}

			_, err := constructAssignment(context.Background(), data)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("constructAssignment() error = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestUnitConstructTargetIgnoresUnknownFilterType(t *testing.T) {
	t.Parallel()

	target, err := constructTarget(context.Background(), "groupAssignmentTarget",
		sharedmodels.DeviceManagementDeviceConfigurationAssignmentWithGroupFilterModel{
			Type:       types.StringValue("groupAssignmentTarget"),
			GroupId:    types.StringValue("11111111-1111-1111-1111-111111111111"),
			FilterId:   types.StringValue("44444444-4444-4444-4444-444444444444"),
			FilterType: types.StringValue("somethingElse"),
		})
	if err != nil {
		t.Fatalf("constructTarget() returned an unexpected error: %v", err)
	}

	// The filter id is still sent, but an unrecognised filter type is left unset rather than guessed.
	if got := target.GetDeviceAndAppManagementAssignmentFilterId(); got == nil {
		t.Error("filter id was not set")
	}
	if got := target.GetDeviceAndAppManagementAssignmentFilterType(); got != nil {
		t.Errorf("filter type = %v, want it to be left unset", got)
	}
}

package mocks

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/microsoft/kiota-abstractions-go/authentication"
	msgraph "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

func newMockContractClient(t *testing.T) *msgraph.GraphServiceClient {
	t.Helper()
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	registrar := &WindowsBiosConfigurationsAndOtherSettingsTemplateMock{}
	registrar.RegisterMocks()
	t.Cleanup(registrar.CleanupMockState)
	adapter, err := msgraph.NewGraphRequestAdapterWithParseNodeFactoryAndSerializationWriterFactoryAndHttpClient(
		&authentication.AnonymousAuthenticationProvider{},
		nil,
		nil,
		http.DefaultClient,
	)
	if err != nil {
		t.Fatal(err)
	}
	adapter.SetBaseUrl("https://graph.microsoft.com/beta")
	return msgraph.NewGraphServiceClient(adapter)
}

func TestUnitMockFilterMetadataExposesCreationTimesThroughSDK(t *testing.T) {
	client := newMockContractClient(t)
	for filterID, want := range map[string]string{
		"44444444-4444-4444-4444-444444444444": "2026-01-02T00:00:00Z",
		"55555555-5555-5555-5555-555555555555": "2026-01-01T00:00:00Z",
	} {
		filter, err := client.DeviceManagement().AssignmentFilters().
			ByDeviceAndAppManagementAssignmentFilterId(filterID).Get(context.Background(), nil)
		if err != nil {
			t.Fatal(err)
		}
		if filter == nil || filter.GetId() == nil || *filter.GetId() != filterID {
			t.Fatalf("metadata response does not identify filter %s", filterID)
		}
		created, err := time.Parse(time.RFC3339, want)
		if err != nil {
			t.Fatal(err)
		}
		if filter.GetCreatedDateTime() == nil || !filter.GetCreatedDateTime().Equal(created) {
			t.Errorf(
				"filter %s createdDateTime = %v, want %s",
				filterID,
				filter.GetCreatedDateTime(),
				want,
			)
		}
	}
}

func TestUnitMockAssignmentsRejectReverseOrderAndRetainAllTargets(t *testing.T) {
	ctx := context.Background()
	client := newMockContractClient(t)
	name, fileName := "unit-test-006-mock-contract", "contract.cctk"
	format := graphmodels.DELL_HARDWARECONFIGURATIONFORMAT
	profile := graphmodels.NewHardwareConfiguration()
	profile.SetDisplayName(&name)
	profile.SetFileName(&fileName)
	profile.SetConfigurationFileContent([]byte("[cctk]\nSecureBoot=Enabled\nTpm=On\n"))
	profile.SetHardwareConfigurationFormat(&format)
	profile.SetRoleScopeTagIds([]string{"0"})
	created, err := client.DeviceManagement().HardwareConfigurations().Post(ctx, profile, nil)
	if err != nil {
		t.Fatal(err)
	}
	item := client.DeviceManagement().
		HardwareConfigurations().
		ByHardwareConfigurationId(*created.GetId())
	assignments := make([]graphmodels.HardwareConfigurationAssignmentable, 0, 5)
	for index := 0; index < 5; index++ {
		groupID := fmt.Sprintf("00000000-0000-0000-0000-%012d", index+1)
		var target graphmodels.DeviceAndAppManagementAssignmentTargetable
		if index < 3 {
			included := graphmodels.NewGroupAssignmentTarget()
			included.SetGroupId(&groupID)
			target = included
		} else {
			excluded := graphmodels.NewExclusionGroupAssignmentTarget()
			excluded.SetGroupId(&groupID)
			target = excluded
		}
		if index < 2 {
			filterID := "44444444-4444-4444-4444-444444444444"
			filterType := graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
			if index == 1 {
				filterID = "55555555-5555-5555-5555-555555555555"
				filterType = graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
			}
			target.SetDeviceAndAppManagementAssignmentFilterId(&filterID)
			target.SetDeviceAndAppManagementAssignmentFilterType(&filterType)
		}
		assignment := graphmodels.NewHardwareConfigurationAssignment()
		assignment.SetTarget(target)
		assignments = append(assignments, assignment)
	}
	request := devicemanagement.NewHardwareConfigurationsItemAssignPostRequestBody()
	request.SetHardwareConfigurationAssignments(assignments)
	_, err = item.Assign().Post(ctx, request, nil)
	requireScopeValidationError(t, err)
	read, err := item.Get(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(read.GetAssignments()) != 0 {
		t.Fatal("rejected assignments mutated the empty policy")
	}

	assignments[0], assignments[1] = assignments[1], assignments[0]
	request.SetHardwareConfigurationAssignments(assignments)
	if _, err := item.Assign().Post(ctx, request, nil); err != nil {
		t.Fatalf("creation-ordered assignments were rejected: %v", err)
	}
	want := assignmentKeys(assignments)
	for _, rejectUpdate := range []bool{false, true} {
		if rejectUpdate {
			assignments[0], assignments[1] = assignments[1], assignments[0]
			request.SetHardwareConfigurationAssignments(assignments)
			_, err = item.Assign().Post(ctx, request, nil)
			requireScopeValidationError(t, err)
		}
		read, err = item.Get(ctx, nil)
		if err != nil {
			t.Fatal(err)
		}
		if got := assignmentKeys(read.GetAssignments()); !reflect.DeepEqual(got, want) {
			t.Errorf(
				"assignment targets changed (rejected update %t): got %v, want %v",
				rejectUpdate,
				got,
				want,
			)
		}
	}
}

func assignmentKeys(assignments []graphmodels.HardwareConfigurationAssignmentable) []string {
	keys := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		target := assignment.GetTarget()
		groupID := ""
		switch typed := target.(type) {
		case graphmodels.GroupAssignmentTargetable:
			groupID = stringPointer(typed.GetGroupId())
		case graphmodels.ExclusionGroupAssignmentTargetable:
			groupID = stringPointer(typed.GetGroupId())
		}
		filterType := "none"
		if target.GetDeviceAndAppManagementAssignmentFilterType() != nil {
			filterType = target.GetDeviceAndAppManagementAssignmentFilterType().String()
		}
		keys = append(
			keys,
			fmt.Sprintf("%s|%s|%s|%s", stringPointer(target.GetOdataType()), groupID,
				stringPointer(target.GetDeviceAndAppManagementAssignmentFilterId()), filterType),
		)
	}
	return keys
}

func stringPointer(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func requireScopeValidationError(t *testing.T, err error) {
	t.Helper()
	var graphError *odataerrors.ODataError
	if !errors.As(err, &graphError) || graphError.ResponseStatusCode != http.StatusBadRequest {
		t.Fatalf("reverse filter order must return HTTP400, got %v", err)
	}
	if !strings.Contains(err.Error(), "Scope Tags validation failed") {
		t.Fatalf("scope validation error was lost: %v", err)
	}
}

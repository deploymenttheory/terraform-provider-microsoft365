package graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/microsoft/kiota-abstractions-go/authentication"
	msgraph "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func TestUnitOrderAssignmentsByFilterCreationTime(t *testing.T) {
	const (
		newer = "11111111-1111-1111-1111-111111111111"
		older = "99999999-9999-9999-9999-999999999999"
	)
	for _, test := range []struct {
		name       string
		filterIDs  []string
		wantOrder  []string
		wantReads  int
		status     int
		missing    bool
		equalTimes bool
		wantError  bool
	}{
		{name: "empty"},
		{name: "unfiltered", filterIDs: []string{"", noFilter}, wantOrder: []string{"0", "1"}},
		{name: "one filter", filterIDs: []string{newer}, wantOrder: []string{"0"}},
		{name: "duplicate single filter", filterIDs: []string{newer, "", newer}, wantOrder: []string{"0", "1", "2"}},
		{
			name:      "reverse creation order with duplicate and unfiltered targets",
			filterIDs: []string{"", newer, noFilter, older, older, ""},
			wantOrder: []string{"0", "3", "2", "4", "1", "5"}, wantReads: 2,
		},
		{
			name: "equal timestamps retain stable order", filterIDs: []string{newer, older},
			wantOrder: []string{"0", "1"}, wantReads: 2, equalTimes: true,
		},
		{
			name: "missing metadata fails without changing targets", filterIDs: []string{newer, older},
			wantOrder: []string{"0", "1"}, wantReads: 1, missing: true, wantError: true,
		},
		{
			name: "inaccessible metadata fails without changing targets", filterIDs: []string{newer, older},
			wantOrder: []string{"0", "1"}, wantReads: 1, status: http.StatusForbidden, wantError: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			reads := make(map[string]int)
			var readsMu sync.Mutex
			server := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					if req.Method != http.MethodGet ||
						!strings.HasPrefix(
							req.URL.Path,
							"/beta/deviceManagement/assignmentFilters/",
						) {
						t.Errorf("unexpected API request: %s %s", req.Method, req.URL.Path)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					filterID := strings.TrimPrefix(
						req.URL.Path,
						"/beta/deviceManagement/assignmentFilters/",
					)
					readsMu.Lock()
					reads[filterID]++
					readsMu.Unlock()
					w.Header().Set("Content-Type", "application/json")
					if test.status != 0 {
						w.WriteHeader(test.status)
						fmt.Fprint(
							w,
							`{"error":{"code":"Forbidden","message":"Filter metadata is inaccessible"}}`,
						)
						return
					}
					if test.missing {
						fmt.Fprintf(w, `{"id":%q}`, filterID)
						return
					}
					created := "2026-01-02T00:00:00Z"
					if filterID == older && !test.equalTimes {
						created = "2026-01-01T00:00:00Z"
					}
					fmt.Fprintf(w, `{"id":%q,"createdDateTime":%q}`, filterID, created)
				}),
			)
			t.Cleanup(server.Close)
			adapter, err := msgraph.NewGraphRequestAdapter(
				&authentication.AnonymousAuthenticationProvider{},
			)
			if err != nil {
				t.Fatal(err)
			}
			adapter.SetBaseUrl(server.URL + "/beta")
			r := &WindowsBiosConfigurationsAndOtherSettingsTemplateResource{
				client: msgraph.NewGraphServiceClient(adapter),
			}
			body := devicemanagement.NewHardwareConfigurationsItemAssignPostRequestBody()
			assignments := make(
				[]graphmodels.HardwareConfigurationAssignmentable,
				0,
				len(test.filterIDs),
			)
			originals := make(map[string]graphmodels.HardwareConfigurationAssignmentable)
			for index, filterID := range test.filterIDs {
				id := fmt.Sprint(index)
				var target graphmodels.DeviceAndAppManagementAssignmentTargetable
				if index%3 == 2 {
					excluded := graphmodels.NewExclusionGroupAssignmentTarget()
					excluded.SetGroupId(&id)
					target = excluded
				} else {
					included := graphmodels.NewGroupAssignmentTarget()
					included.SetGroupId(&id)
					target = included
				}
				if filterID != "" {
					target.SetDeviceAndAppManagementAssignmentFilterId(&filterID)
					filterType := graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
					if filterID == newer {
						filterType = graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE
					}
					target.SetDeviceAndAppManagementAssignmentFilterType(&filterType)
				}
				assignment := graphmodels.NewHardwareConfigurationAssignment()
				assignment.SetId(&id)
				assignment.SetTarget(target)
				assignments = append(assignments, assignment)
				originals[id] = assignment
			}
			body.SetHardwareConfigurationAssignments(assignments)
			err = r.orderAssignmentsByFilterCreationTime(context.Background(), body)
			if (err != nil) != test.wantError {
				t.Fatalf(
					"orderAssignmentsByFilterCreationTime error = %v, want error %t",
					err,
					test.wantError,
				)
			}
			if test.missing && !errors.Is(err, errMissingFilterCreatedDateTime) {
				t.Errorf("missing metadata error = %v, want errMissingFilterCreatedDateTime", err)
			}
			var actual []string
			for _, assignment := range body.GetHardwareConfigurationAssignments() {
				id := *assignment.GetId()
				actual = append(actual, id)
				if assignment != originals[id] {
					t.Errorf("assignment %s changed instead of being reordered", id)
				}
			}
			if !reflect.DeepEqual(actual, test.wantOrder) {
				t.Errorf("assignment order = %v, want %v", actual, test.wantOrder)
			}
			count := 0
			readsMu.Lock()
			for id, readCount := range reads {
				if readCount != 1 {
					t.Errorf("filter %s read %d times, want once", id, readCount)
				}
				count += readCount
			}
			readsMu.Unlock()
			if count != test.wantReads {
				t.Errorf("metadata reads = %d, want %d", count, test.wantReads)
			}
		})
	}
}

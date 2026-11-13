package mocks_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"

	policymocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/windows_autopilot_device_preparation_policy/mocks"
)

type membershipResponse struct {
	ID      string `json:"id"`
	Targets []struct {
		TargetID string `json:"targetId"`
	} `json:"enrollmentTimeDeviceMembershipTargetValidationStatuses"`
}

func membershipRequest(t *testing.T, method, endpoint, body string, status int) membershipResponse {
	t.Helper()
	request, err := http.NewRequestWithContext(
		t.Context(),
		method,
		endpoint,
		strings.NewReader(body),
	)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != status {
		t.Fatalf("%s returned HTTP %d, want %d", method, response.StatusCode, status)
	}
	var result membershipResponse
	if status != http.StatusNoContent {
		if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
			t.Fatal(err)
		}
	}
	return result
}

func TestUnitAutopilotMembershipTargetMock_RoundTrip(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	mock := &policymocks.WindowsAutopilotDevicePreparationPolicyMock{}
	mock.RegisterMocks()
	t.Cleanup(mock.CleanupMockState)

	const endpoint = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies"
	created := membershipRequest(
		t,
		http.MethodPost,
		endpoint,
		`{"name":"unit-membership-policy"}`,
		http.StatusCreated,
	)
	if created.ID == "" {
		t.Fatal("create did not return an ID")
	}
	readEndpoint := endpoint + "('" + created.ID + "')/retrieveEnrollmentTimeDeviceMembershipTarget"
	if result := membershipRequest(
		t,
		http.MethodGet,
		readEndpoint,
		"",
		http.StatusOK,
	); len(
		result.Targets,
	) != 0 {
		t.Fatal("unconfigured policy returned a membership target")
	}

	for _, id := range []string{
		"00000000-0000-0000-0000-000000000001",
		"00000000-0000-0000-0000-000000000002",
	} {
		body := `{"enrollmentTimeDeviceMembershipTargets":[{"targetId":"` + id + `","targetType":"staticSecurityGroup"}]}`
		membershipRequest(
			t,
			http.MethodPost,
			endpoint+"/"+created.ID+"/setEnrollmentTimeDeviceMembershipTarget",
			body,
			http.StatusNoContent,
		)
		result := membershipRequest(t, http.MethodGet, readEndpoint, "", http.StatusOK)
		if len(result.Targets) != 1 || result.Targets[0].TargetID != id {
			t.Fatalf("membership target was not preserved: %+v", result.Targets)
		}
	}

	membershipRequest(t, http.MethodDelete, endpoint+"/"+created.ID, "", http.StatusNoContent)
	membershipRequest(t, http.MethodGet, readEndpoint, "", http.StatusNotFound)
}

func TestUnitAutopilotMembershipTargetMock_InvalidSetBody(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	mock := &policymocks.WindowsAutopilotDevicePreparationPolicyMock{}
	mock.RegisterMocks()
	t.Cleanup(mock.CleanupMockState)

	const endpoint = "https://graph.microsoft.com/beta/deviceManagement/configurationPolicies/00000000-0000-0000-0000-000000000001/setEnrollmentTimeDeviceMembershipTarget"
	membershipRequest(t, http.MethodPost, endpoint, "{", http.StatusBadRequest)
}

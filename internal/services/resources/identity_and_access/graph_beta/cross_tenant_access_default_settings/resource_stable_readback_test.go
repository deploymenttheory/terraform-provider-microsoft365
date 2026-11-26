package graphBetaCrossTenantAccessDefaultSettings_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

// A successful replica read can be followed by the pre-write policy. Both write
// operations must wait for a stable sequence before Terraform refreshes the plan.
func TestUnitResourceCrossTenantAccessDefaultSettings_19_StableWriteReadback(t *testing.T) {
	for _, operation := range []string{"create", "update"} {
		t.Run(operation, func(t *testing.T) {
			mocks.SetupUnitTestEnvironment(t)
			_, mock := setupMockEnvironment()
			defer httpmock.DeactivateAndReset()
			defer mock.CleanupMockState()
			content, err := os.ReadFile("tests/responses/validate_create/get_cross_tenant_access_default_settings_system_default.json")
			if err != nil {
				t.Fatal(err)
			}
			var current map[string]any
			if err := json.Unmarshal(content, &current); err != nil {
				t.Fatal(err)
			}
			var mu sync.Mutex
			var previous map[string]any
			writes, reads := 0, 0
			targetWrite := 1
			if operation == "update" {
				targetWrite = 2
			}
			endpoint := "https://graph.microsoft.com/beta/policies/crossTenantAccessPolicy/default"
			httpmock.RegisterResponder("PATCH", endpoint, func(req *http.Request) (*http.Response, error) {
				mu.Lock()
				defer mu.Unlock()
				var patch map[string]any
				if err := json.NewDecoder(req.Body).Decode(&patch); err != nil {
					return nil, err
				}
				previous = current
				current = make(map[string]any, len(previous))
				for key, value := range previous {
					current[key] = value
				}
				for key, value := range patch {
					current[key] = value
				}
				current["isServiceDefault"] = false
				writes++
				reads = 0
				return httpmock.NewStringResponse(204, ""), nil
			})
			httpmock.RegisterResponder("GET", endpoint, func(_ *http.Request) (*http.Response, error) {
				mu.Lock()
				defer mu.Unlock()
				if writes == targetWrite {
					reads++
					if reads == 2 {
						return httpmock.NewJsonResponse(200, previous)
					}
				}
				return httpmock.NewJsonResponse(200, current)
			})
			config := strings.ReplaceAll(loadUnitTestTerraform("resource_01_minimal.tf"), "restore_defaults_on_destroy = true", "restore_defaults_on_destroy = false")
			end := strings.LastIndex(config, "}")
			configured := config[:end] + "inbound_trust = { is_mfa_accepted = true, is_compliant_device_accepted = true, is_hybrid_azure_ad_joined_device_accepted = true }\n" + config[end:]
			checkStable := func(_ *terraform.State) error {
				mu.Lock()
				defer mu.Unlock()
				if reads < 5 {
					return fmt.Errorf("%s accepted an isolated matching replica after %d reads; want fresh, stale, then three fresh reads", operation, reads)
				}
				return nil
			}
			steps := []resource.TestStep{{Config: configured, Check: checkStable}}
			if operation == "update" {
				steps = []resource.TestStep{{Config: config}, {Config: configured, Check: checkStable}}
			}
			resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories, Steps: steps})
		})
	}
}

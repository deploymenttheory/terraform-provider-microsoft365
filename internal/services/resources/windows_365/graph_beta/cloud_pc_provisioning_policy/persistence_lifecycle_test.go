package graphBetaCloudPcProvisioningPolicy_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	policyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_365/graph_beta/cloud_pc_provisioning_policy/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestUnitResourceCloudPcProvisioningPolicy_PersistenceLifecycle(t *testing.T) {
	setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	(&policyMocks.CloudPcProvisioningPolicyMock{}).RegisterMocks()
	name := "microsoft365_graph_beta_windows_365_cloud_pc_provisioning_policy.lifecycle"
	base := testConfigMinimalWithResourceName("lifecycle")
	var policyID string
	expectedEnabled := "true"
	var steps []resource.TestStep
	for i, phase := range []struct{ config, value string }{{`provisioning_type = "sharedByEntraGroup"
user_settings_persistence_configuration = {enabled = true}`, "fourGB"},
		{`provisioning_type = "sharedByEntraGroup"
user_settings_persistence_configuration = {enabled = true, storage_size_category = "eightGB"}`, "eightGB"},
		{`provisioning_type = "sharedByEntraGroup"
user_settings_persistence_configuration = {enabled = true, storage_size_category = "sixteenGB"}`, "sixteenGB"},
		{`provisioning_type = "sharedByEntraGroup"
user_settings_persistence_configuration = {enabled = true, storage_size_category = "thirtyTwoGB"}`, "thirtyTwoGB"},
		{`provisioning_type = "sharedByEntraGroup"
user_settings_persistence_configuration = {enabled = true, storage_size_category = "sixtyFourGB"}`, "sixtyFourGB"},
		{`provisioning_type = "sharedByEntraGroup"
user_settings_persistence_configuration = {enabled = false, storage_size_category = "eightGB"}`, "eightGB"},
		{`provisioning_type = "sharedByEntraGroup"
user_settings_persistence_configuration = {enabled = false}`, "eightGB"},
		{`provisioning_type = "sharedByEntraGroup"`, "eightGB"}} {
		config := strings.Replace(base, "display_name =", phase.config+"\ndisplay_name =", 1)
		checks := []resource.TestCheckFunc{resource.TestCheckResourceAttr(name, "user_settings_persistence_configuration.storage_size_category", phase.value)}
		if i == 0 {
			checks = append(checks, resource.TestCheckResourceAttrWith(name, "id", func(value string) error { policyID = value; return nil }))
		} else {
			checks = append(checks, resource.TestCheckResourceAttrPtr(name, "id", &policyID))
		}
		if strings.Contains(phase.config, "enabled = false") {
			expectedEnabled = "false"
		}
		checks = append(checks, resource.TestCheckResourceAttr(name, "user_settings_persistence_configuration.enabled", expectedEnabled))

		if i > 0 {
			config = strings.Replace(config, "Test Minimal Provisioning Policy", fmt.Sprintf("Updated policy %d", i), 1)
		}
		steps = append(steps, resource.TestStep{Config: config, Check: resource.ComposeTestCheckFunc(checks...)}, resource.TestStep{ResourceName: name, ImportState: true, ImportStateVerify: true}, resource.TestStep{Config: config, PlanOnly: true})
	}
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories, Steps: steps})
}
func TestUnitResourceCloudPcProvisioningPolicy_PersistenceLifecycleInvalidProvisioningType(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories, Steps: []resource.TestStep{{Config: strings.Replace(testConfigMinimalWithResourceName("invalid"), "display_name =", `provisioning_type = "dedicated"
user_settings_persistence_configuration = {enabled = true}`+"\ndisplay_name =", 1), ExpectError: regexp.MustCompile("persistence requires")}}})
}

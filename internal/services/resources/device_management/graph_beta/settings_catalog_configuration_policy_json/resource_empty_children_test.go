package graphBetaSettingsCatalogConfigurationPolicyJson_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaSettingsCatalogConfigurationPolicyJson "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_configuration_policy_json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
)

// policyRequestURL matches the create (POST) and update (PUT) requests for a configuration policy, but not its assign action.
var policyRequestURL = regexp.MustCompile(`^/beta/deviceManagement/configurationPolicies(/\{uuid-\d+\}|\(('|%27)\{uuid-\d+\}('|%27)\))?$`)

// assertEmptyChildrenSent checks that the captured method request for the policy sent "children": [] for each setting
// value without children (#4956). Graph rejects "children": null with 400 ModelValidationFailure.
func assertEmptyChildrenSent(capture *mocks.RequestCapture, method string) error {
	requests, err := capture.Requests()
	if err != nil {
		return err
	}
	var body map[string]any
	urls := make([]string, 0, len(requests))
	for _, request := range requests {
		urls = append(urls, request.Method+" "+request.URL)
		if request.Method == method && policyRequestURL.MatchString(request.URL) {
			body, _ = request.Body.(map[string]any)
		}
	}
	if body == nil {
		return fmt.Errorf("no %s request with a JSON body was sent for the policy; requests: %v", method, urls)
	}
	settings, _ := body["settings"].([]any)
	if len(settings) != 3 {
		return fmt.Errorf("%s sent %d settings, want 3", method, len(settings))
	}
	for index, value := range []struct {
		path string
		get  func(instance map[string]any) any
	}{
		{"choiceSettingCollectionValue[0]", func(instance map[string]any) any { return firstElement(instance["choiceSettingCollectionValue"]) }},
		{"groupSettingCollectionValue[0]", func(instance map[string]any) any { return firstElement(instance["groupSettingCollectionValue"]) }},
		{"choiceSettingValue", func(instance map[string]any) any { return instance["choiceSettingValue"] }},
	} {
		setting, _ := settings[index].(map[string]any)
		instance, _ := setting["settingInstance"].(map[string]any)
		settingValue, _ := value.get(instance).(map[string]any)
		children, present := settingValue["children"]
		list, isList := children.([]any)
		if !present || !isList || list == nil || len(list) != 0 {
			return fmt.Errorf("%s settings[%d].settingInstance.%s.children = %#v (present %t), want []", method, index, value.path, children, present)
		}
	}
	return nil
}

// firstElement returns the first element of a JSON array, or nil when value is not a non-empty array.
func firstElement(value any) any {
	if list, ok := value.([]any); ok && len(list) > 0 {
		return list[0]
	}
	return nil
}

// Test 21: Settings without children are sent as "children": [] on create and update (#4956)
func TestUnitResourceSettingsCatalogConfigurationPolicyJson_21_EmptyChildren(t *testing.T) {
	resourceType := graphBetaSettingsCatalogConfigurationPolicyJson.ResourceName
	mocks.SetupUnitTestEnvironment(t)
	_, policyMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer policyMock.CleanupMockState()
	capture := mocks.CaptureGraphRequests(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_21_empty_children.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".empty_children").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".empty_children").Key("settings_count").HasValue("3"),
					check.That(resourceType+".empty_children").Key("assignments.#").HasValue("1"),
					func(_ *terraform.State) error { return assertEmptyChildrenSent(capture, "POST") },
				),
			},
			{
				Config: loadUnitTestTerraform("resource_21_empty_children_updated.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".empty_children").Key("description").HasValue("Test policy for settings without children - updated"),
					check.That(resourceType+".empty_children").Key("assignments.0.group_id").HasValue("22222222-2222-2222-2222-222222222222"),
					func(_ *terraform.State) error { return assertEmptyChildrenSent(capture, "PUT") },
				),
			},
		},
	})
}

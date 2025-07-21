package graphIOSMobileAppConfiguration_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// Unit Tests

func TestUnitIOSMobileAppConfigurationResource_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Read test fixtures
	postResp := readTestFixture(t, "Validate_Create/post_ios_mobile_app_configuration.json")
	getResp := readTestFixture(t, "Validate_Create/get_ios_mobile_app_configuration.json")
	getAssignmentsResp := readTestFixture(t, "Validate_Create/get_ios_mobile_app_configuration_assignments.json")
	getUpdatedResp := readTestFixture(t, "Validate_Update/get_ios_mobile_app_configuration_updated.json")
	notFoundResp := readTestFixture(t, "Validate_Delete/get_ios_mobile_app_configuration_not_found.json")

	// Mock Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/v1.0/deviceAppManagement/mobileAppConfigurations",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(201, postResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Track GET request count and delete state
	var getRequestCount int
	var isDeleted bool

	// Mock Read - handles all GET requests for the resource
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/v1.0/deviceAppManagement/mobileAppConfigurations/00000000-0000-0000-0000-000000000001",
		func(req *http.Request) (*http.Response, error) {
			getRequestCount++
			
			// After delete, always return 404
			if isDeleted {
				resp := httpmock.NewBytesResponse(404, notFoundResp)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			}
			
			// For update step (requests 3-4), return updated response
			if getRequestCount >= 3 && getRequestCount <= 4 {
				resp := httpmock.NewBytesResponse(200, getUpdatedResp)
				resp.Header.Set("Content-Type", "application/json")
				return resp, nil
			}
			
			// For create step (requests 1-2), return initial response
			resp := httpmock.NewBytesResponse(200, getResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Mock Get Assignments
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/v1.0/deviceAppManagement/mobileAppConfigurations/00000000-0000-0000-0000-000000000001/assignments",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(200, getAssignmentsResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Mock Update
	httpmock.RegisterResponder("PATCH", "https://graph.microsoft.com/v1.0/deviceAppManagement/mobileAppConfigurations/00000000-0000-0000-0000-000000000001",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(200, getUpdatedResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Mock Delete
	httpmock.RegisterResponder("DELETE", "https://graph.microsoft.com/v1.0/deviceAppManagement/mobileAppConfigurations/00000000-0000-0000-0000-000000000001",
		func(req *http.Request) (*http.Response, error) {
			isDeleted = true
			return httpmock.NewBytesResponse(204, nil), nil
		})

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIOSMobileAppConfigurationResource_basic("Test iOS Config"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "display_name", "Test iOS Config"),
					resource.TestCheckResourceAttr("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "description", "Test iOS mobile app configuration"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "id"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "created_date_time"),
					resource.TestCheckResourceAttrSet("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "last_modified_date_time"),
					resource.TestCheckResourceAttr("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "version", "1"),
				),
			},
			// Update and Read testing
			{
				Config: testAccIOSMobileAppConfigurationResource_basic("Updated iOS Config"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "display_name", "Updated iOS Config"),
				),
			},
			// Delete testing is implicit
		},
	})
}

func TestUnitIOSMobileAppConfigurationResource_Complete(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()

	// Read test fixtures
	postResp := readTestFixture(t, "Validate_Create/post_ios_mobile_app_configuration_complete.json")
	getResp := readTestFixture(t, "Validate_Create/get_ios_mobile_app_configuration_complete.json")
	getAssignmentsResp := readTestFixture(t, "Validate_Create/get_ios_mobile_app_configuration_assignments_complete.json")
	patchAssignmentsResp := readTestFixture(t, "Validate_Create/patch_ios_mobile_app_configuration_assignments.json")

	// Mock Create
	httpmock.RegisterResponder("POST", "https://graph.microsoft.com/v1.0/deviceAppManagement/mobileAppConfigurations",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(201, postResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Mock Read after Create
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/v1.0/deviceAppManagement/mobileAppConfigurations/00000000-0000-0000-0000-000000000002",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(200, getResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Mock Patch for assignments
	httpmock.RegisterResponder("PATCH", "https://graph.microsoft.com/v1.0/deviceAppManagement/mobileAppConfigurations/00000000-0000-0000-0000-000000000002",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(200, patchAssignmentsResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Mock Get Assignments
	httpmock.RegisterResponder("GET", "https://graph.microsoft.com/v1.0/deviceAppManagement/mobileAppConfigurations/00000000-0000-0000-0000-000000000002/assignments",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewBytesResponse(200, getAssignmentsResp)
			resp.Header.Set("Content-Type", "application/json")
			return resp, nil
		})

	// Mock Delete
	httpmock.RegisterResponder("DELETE", "https://graph.microsoft.com/v1.0/deviceAppManagement/mobileAppConfigurations/00000000-0000-0000-0000-000000000002",
		httpmock.NewBytesResponder(204, nil))

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with all fields
			{
				Config: testAccIOSMobileAppConfigurationResource_complete(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "display_name", "Complete iOS Config"),
					resource.TestCheckResourceAttr("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "description", "Complete iOS mobile app configuration with all features"),
					resource.TestCheckResourceAttr("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "targeted_mobile_apps.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "settings.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "settings.0.app_config_key", "setting1"),
					resource.TestCheckResourceAttr("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "settings.0.app_config_key_type", "stringType"),
					resource.TestCheckResourceAttr("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "settings.0.app_config_key_value", "value1"),
					resource.TestCheckResourceAttr("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "assignments.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "assignments.0.target.odata_type", "#microsoft.graph.groupAssignmentTarget"),
					resource.TestCheckResourceAttr("microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test", "assignments.0.target.group_id", "00000000-0000-0000-0000-000000000100"),
				),
			},
		},
	})
}

// Acceptance Tests

func TestAccIOSMobileAppConfigurationResource_Basic(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("Set TF_ACC=1 to run acceptance tests")
	}

	ctx := context.Background()
	_ = ctx
	resourceName := "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration.test"
	displayName := fmt.Sprintf("tftest-ios-config-%d", 12345)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIOSMobileAppConfigurationResource_basic(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "display_name", displayName),
					resource.TestCheckResourceAttr(resourceName, "description", "Test iOS mobile app configuration"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_date_time"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modified_date_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIOSMobileAppConfigurationResource_basic("Updated " + displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "display_name", "Updated "+displayName),
				),
			},
		},
	})
}

// Test configurations

func testAccIOSMobileAppConfigurationResource_basic(displayName string) string {
	return fmt.Sprintf(`
resource "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name = %[1]q
  description  = "Test iOS mobile app configuration"
}
`, displayName)
}

func testAccIOSMobileAppConfigurationResource_complete() string {
	return `
resource "microsoft365_graph_v1_device_and_app_management_ios_mobile_app_configuration" "test" {
  display_name = "Complete iOS Config"
  description  = "Complete iOS mobile app configuration with all features"
  
  targeted_mobile_apps = [
    "00000000-0000-0000-0000-000000000010",
    "00000000-0000-0000-0000-000000000011"
  ]
  
  encoded_setting_xml = base64encode("<configuration><setting>test</setting></configuration>")
  
  settings {
    app_config_key       = "setting1"
    app_config_key_type  = "stringType"
    app_config_key_value = "value1"
  }
  
  settings {
    app_config_key       = "setting2"
    app_config_key_type  = "integerType"
    app_config_key_value = "42"
  }
  
  assignments {
    target {
      odata_type = "#microsoft.graph.groupAssignmentTarget"
      group_id   = "00000000-0000-0000-0000-000000000100"
    }
  }
}
`
}

// Helper functions

func readTestFixture(t *testing.T, filename string) []byte {
	path := filepath.Join("tests", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test fixture %s: %v", filename, err)
	}
	return data
}

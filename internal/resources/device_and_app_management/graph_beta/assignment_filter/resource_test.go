package graphBetaAssignmentFilter

// import (
// 	"context"
// 	"fmt"
// 	"regexp"
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
// 	"github.com/hashicorp/terraform-plugin-testing/terraform"

// 	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider"
// 	"github.com/hashicorp/terraform-plugin-framework/providerserver"
// 	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
// )

// const (
// 	// providerConfig is a shared configuration to combine with the actual
// 	// test configuration so the microsoft365 client is properly configured.
// 	// It is also possible to use the microsoft365_ environment variables instead,
// 	// such as updating the Makefile and running the testing through that tool.
// 	providerConfig = `
// provider "microsoft365" {
// }
// `
// )

// var (
// 	// testAccProtoV6ProviderFactories are used to instantiate a provider during
// 	// acceptance testing. The factory function will be invoked for every Terraform
// 	// CLI command executed to create a provider server to which the CLI can
// 	// reattach.
// 	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
// 		"microsoft365": providerserver.NewProtocol6WithError(provider.New("test")()),
// 	}
// )

// func TestAccAssignmentFilterResource_basic(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Create and Read testing
// 			{
// 				Config: testAccAssignmentFilterResourceConfig_basic(),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					testAccCheckAssignmentFilterExists("microsoft365_graph_beta_device_and_app_management_assignment_filter.test"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_assignment_filter.test", "display_name", "Test Assignment Filter"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_assignment_filter.test", "description", "Test description"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_assignment_filter.test", "platform", "windows10AndLater"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_assignment_filter.test", "rule", "(device.deviceOwnership -eq \"Company\")"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_assignment_filter.test", "assignment_filter_management_type", "devices"),
// 					resource.TestCheckResourceAttrSet("microsoft365_graph_beta_device_and_app_management_assignment_filter.test", "id"),
// 				),
// 			},
// 			// ImportState testing
// 			{
// 				ResourceName:      "microsoft365_graph_beta_device_and_app_management_assignment_filter.test",
// 				ImportState:       true,
// 				ImportStateVerify: true,
// 			},
// 			// Update and Read testing
// 			{
// 				Config: testAccAssignmentFilterResourceConfig_update(),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					testAccCheckAssignmentFilterExists("microsoft365_graph_beta_device_and_app_management_assignment_filter.test"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_assignment_filter.test", "display_name", "Updated Test Assignment Filter"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_assignment_filter.test", "description", "Updated test description"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_assignment_filter.test", "platform", "android"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_assignment_filter.test", "rule", "(device.deviceOwnership -eq \"Personal\")"),
// 					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_and_app_management_assignment_filter.test", "assignment_filter_management_type", "apps"),
// 				),
// 			},
// 		},
// 	})
// }

// func TestAccAssignmentFilterResource_invalidPlatform(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config:      testAccAssignmentFilterResourceConfig_invalidPlatform(),
// 				ExpectError: regexp.MustCompile(`expected platform to be one of`),
// 			},
// 		},
// 	})
// }

// func TestAccAssignmentFilterResource_invalidRule(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config:      testAccAssignmentFilterResourceConfig_invalidRule(),
// 				ExpectError: regexp.MustCompile(`invalid rule syntax`),
// 			},
// 		},
// 	})
// }

// func TestAccAssignmentFilterResource_requiredFields(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			{
// 				Config:      testAccAssignmentFilterResourceConfig_missingRequiredFields(),
// 				ExpectError: regexp.MustCompile(`The argument "display_name" is required`),
// 			},
// 		},
// 	})
// }

// func testAccAssignmentFilterResourceConfig_basic() string {
// 	return `
// resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "test" {
//   display_name = "Test Assignment Filter"
//   description = "Test description"
//   platform = "windows10AndLater"
//   rule = "(device.deviceOwnership -eq \"Company\")"
//   assignment_filter_management_type = "devices"
// }
// `
// }

// func testAccAssignmentFilterResourceConfig_update() string {
// 	return `
// resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "test" {
//   display_name = "Updated Test Assignment Filter"
//   description = "Updated test description"
//   platform = "android"
//   rule = "(device.deviceOwnership -eq \"Personal\")"
//   assignment_filter_management_type = "apps"
// }
// `
// }

// func testAccAssignmentFilterResourceConfig_invalidPlatform() string {
// 	return `
// resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "test" {
//   display_name = "Invalid Platform Test"
//   description = "Test with invalid platform"
//   platform = "invalidPlatform"
//   rule = "(device.deviceOwnership -eq \"Company\")"
//   assignment_filter_management_type = "devices"
// }
// `
// }

// func testAccAssignmentFilterResourceConfig_invalidRule() string {
// 	return `
// resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "test" {
//   display_name = "Invalid Rule Test"
//   description = "Test with invalid rule"
//   platform = "windows10AndLater"
//   rule = "This is not a valid rule"
//   assignment_filter_management_type = "devices"
// }
// `
// }

// func testAccAssignmentFilterResourceConfig_missingRequiredFields() string {
// 	return `
// resource "microsoft365_graph_beta_device_and_app_management_assignment_filter" "test" {
//   description = "Missing required fields"
//   platform = "windows10AndLater"
//   rule = "(device.deviceOwnership -eq \"Company\")"
// }
// `
// }

// func testAccCheckAssignmentFilterExists(resourceName string) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources[resourceName]
// 		if !ok {
// 			return fmt.Errorf("Assignment Filter not found: %s", resourceName)
// 		}

// 		if rs.Primary.ID == "" {
// 			return fmt.Errorf("Assignment Filter ID is not set")
// 		}

// 		providerServer := testAccProtoV6ProviderFactories["microsoft365"]()
// 		if providerServer == nil {
// 			return fmt.Errorf("provider not initialized")
// 		}

// 		diags := providerServer.ConfigureProvider(context.Background(), &tfprotov6.ConfigureProviderRequest{})
// 		if diags.HasError() {
// 			return fmt.Errorf("error configuring provider: %v", diags)
// 		}

// 		// Get the provider's configured client
// 		providerFactoryRes := provider.New("test")()
// 		providerAdapter, ok := providerFactoryRes.(provider.Microsoft365Provider)
// 		if !ok {
// 			return fmt.Errorf("failed to convert to Microsoft365Provider")
// 		}

// 		client := providerAdapter.GetMSGraphClient()
// 		if client == nil {
// 			return fmt.Errorf("failed to get MS Graph client")
// 		}

// 		// Now use the client to check if the resource exists
// 		filter, err := client.DeviceManagement().AssignmentFilters().ByDeviceAndAppManagementAssignmentFilterId(rs.Primary.ID).Get(context.Background(), nil)
// 		if err != nil {
// 			return fmt.Errorf("error fetching Assignment Filter with resource name %s and id %s, %v", resourceName, rs.Primary.ID, err)
// 		}

// 		if filter == nil {
// 			return fmt.Errorf("Assignment Filter with ID %s not found", rs.Primary.ID)
// 		}

// 		return nil
// 	}
// }

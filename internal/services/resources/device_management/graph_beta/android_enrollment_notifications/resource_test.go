package graphBetaAndroidEnrollmentNotifications_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	notificationMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/android_enrollment_notifications/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *notificationMocks.AndroidEnrollmentNotificationsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	notificationMock := &notificationMocks.AndroidEnrollmentNotificationsMock{}
	notificationMock.RegisterMocks()
	return mockClient, notificationMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *notificationMocks.AndroidEnrollmentNotificationsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	notificationMock := &notificationMocks.AndroidEnrollmentNotificationsMock{}
	notificationMock.RegisterErrorMocks()
	return mockClient, notificationMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestAndroidEnrollmentNotificationsResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, notificationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer notificationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: func() string {
					config, _ := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
					return config
				}(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.unit", "display_name", "Unit Test - Android Enrollment Notifications"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.unit", "platform_type", "androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.unit", "device_enrollment_configuration_type", "deviceEnrollmentNotificationConfiguration"),
				),
			},
		},
	})
}

func TestAndroidEnrollmentNotificationsResource_CompleteConfiguration(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, notificationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer notificationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: func() string {
					config, _ := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
					return config
				}(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.complete"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.complete", "display_name", "Complete Test - Android Enrollment Notifications"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.complete", "platform_type", "androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.complete", "branding_options.#", "1"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.complete", "default_locale", "en-US"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.complete", "notification_templates.#", "2"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.complete", "localized_notification_messages.#", "2"),
				),
			},
		},
	})
}

func TestAndroidEnrollmentNotificationsResource_WithAssignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, notificationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer notificationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: func() string {
					config, _ := helpers.ParseHCLFile("tests/terraform/unit/resource_with_assignments.tf")
					return config
				}(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.with_assignments"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.with_assignments", "assignments.#", "2"),
				),
			},
		},
	})
}

func TestAndroidEnrollmentNotificationsResource_CreateError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, notificationMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer notificationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: func() string {
					config, _ := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
					return config
				}(),
				ExpectError: regexp.MustCompile(`Error creating resource`),
			},
		},
	})
}

func TestAndroidEnrollmentNotificationsResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, notificationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer notificationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create the resource
			{
				Config: func() string {
					config, _ := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
					return config
				}(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.unit", "display_name", "Unit Test - Android Enrollment Notifications"),
				),
			},
			// Step 2: Update the resource
			{
				Config: func() string {
					config, _ := helpers.ParseHCLFile("tests/terraform/unit/resource_maximal.tf")
					// Modify the resource name to match the first step
					config = strings.ReplaceAll(config, "microsoft365_graph_beta_device_management_android_enrollment_notifications \"complete\"", "microsoft365_graph_beta_device_management_android_enrollment_notifications \"unit\"")
					config = strings.ReplaceAll(config, "Complete Test - Android Enterprise Notifications", "Updated Test - Android Enterprise Notifications")
					return config
				}(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.unit"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.unit", "display_name", "Updated Test - Android Enterprise Notifications"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.unit", "localized_notification_messages.#", "2"),
				),
			},
		},
	})
}

func TestAndroidEnrollmentNotificationsResource_DeleteError(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, notificationMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer notificationMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: func() string {
					config, _ := helpers.ParseHCLFile("tests/terraform/unit/resource_minimal.tf")
					return config
				}(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.unit"),
				),
			},
		},
	})
}

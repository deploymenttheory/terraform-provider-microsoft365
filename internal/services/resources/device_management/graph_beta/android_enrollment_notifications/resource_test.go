package graphBetaAndroidEnrollmentNotifications_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	enrollmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/android_enrollment_notifications/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *enrollmentMocks.AndroidEnrollmentNotificationsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	enrollmentMock := &enrollmentMocks.AndroidEnrollmentNotificationsMock{}
	enrollmentMock.RegisterMocks()
	return mockClient, enrollmentMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *enrollmentMocks.AndroidEnrollmentNotificationsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	enrollmentMock := &enrollmentMocks.AndroidEnrollmentNotificationsMock{}
	enrollmentMock.RegisterErrorMocks()
	return mockClient, enrollmentMock
}

func testCheckExists(resourceName string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, "id")
}

func TestAndroidEnrollmentNotificationsResource_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAndroidEmailMinimal(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_android", "display_name", "email minimal android"),
					resource.TestMatchResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_android", "id", regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_android", "platform_type", "android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_android", "default_locale", "en-US"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_android", "role_scope_tag_ids.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_android", "role_scope_tag_ids.*", "0"),
				),
			},
		},
	})
}

func TestAndroidEnrollmentNotificationsResource_PlatformTypes(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAndroidEmailMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "display_name", "email maximal android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "platform_type", "android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "notification_templates.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "notification_templates.*", "email"),
				),
			},
			{
				Config: testConfigAndroidForWorkEmailMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_androidforwork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_androidforwork", "display_name", "email maximal androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_androidforwork", "platform_type", "androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_androidforwork", "notification_templates.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_androidforwork", "notification_templates.*", "email"),
				),
			},
		},
	})
}

func TestAndroidEnrollmentNotificationsResource_AllTerraformConfigurations(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Android Platform Tests
			{
				Config: testConfigAndroidEmailMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_android", "platform_type", "android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_android", "notification_templates.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_android", "notification_templates.*", "email"),
				),
			},
			{
				Config: testConfigAndroidEmailMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "platform_type", "android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "notification_templates.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "notification_templates.*", "email"),
				),
			},
			{
				Config: testConfigAndroidPushMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.push_maximal_android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.push_maximal_android", "platform_type", "android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.push_maximal_android", "notification_templates.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.push_maximal_android", "notification_templates.*", "push"),
				),
			},
			{
				Config: testConfigAndroidAllMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.all_android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.all_android", "platform_type", "android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.all_android", "notification_templates.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.all_android", "notification_templates.*", "email"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.all_android", "notification_templates.*", "push"),
				),
			},
			// AndroidForWork Platform Tests
			{
				Config: testConfigAndroidForWorkEmailMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_androidforwork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_androidforwork", "platform_type", "androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_androidforwork", "notification_templates.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_androidforwork", "notification_templates.*", "email"),
				),
			},
			{
				Config: testConfigAndroidForWorkEmailMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_androidforwork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_androidforwork", "platform_type", "androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_androidforwork", "notification_templates.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_androidforwork", "notification_templates.*", "email"),
				),
			},
			{
				Config: testConfigAndroidForWorkPushMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.push_maximal_androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.push_maximal_androidForWork", "platform_type", "androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.push_maximal_androidForWork", "notification_templates.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.push_maximal_androidForWork", "notification_templates.*", "push"),
				),
			},
			{
				Config: testConfigAndroidForWorkAllMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.all_androidforwork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.all_androidforwork", "platform_type", "androidForWork"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.all_androidforwork", "notification_templates.#", "2"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.all_androidforwork", "notification_templates.*", "email"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.all_androidforwork", "notification_templates.*", "push"),
				),
			},
		},
	})
}

func TestAndroidEnrollmentNotificationsResource_BrandingOptions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAndroidEmailMinimal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_android", "branding_options.#", "1"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_minimal_android", "branding_options.*", "none"),
				),
			},
			{
				Config: testConfigAndroidEmailMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "branding_options.#", "5"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "branding_options.*", "includeCompanyLogo"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "branding_options.*", "includeCompanyName"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "branding_options.*", "includeCompanyPortalLink"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "branding_options.*", "includeContactInformation"),
					resource.TestCheckTypeSetElemAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "branding_options.*", "includeDeviceDetails"),
				),
			},
		},
	})
}

func TestAndroidEnrollmentNotificationsResource_LocalizedMessages(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAndroidEmailMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "localized_notification_messages.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "localized_notification_messages.*", map[string]string{
						"locale":        "en-us",
						"subject":       "Device Enrollment Required",
						"template_type": "email",
						"is_default":    "true",
					}),
				),
			},
		},
	})
}

func TestAndroidEnrollmentNotificationsResource_Assignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigAndroidEmailMaximal(),
				Check: resource.ComposeTestCheckFunc(
					testCheckExists("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android"),
					resource.TestCheckResourceAttr("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "assignments.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("microsoft365_graph_beta_device_management_android_enrollment_notifications.email_maximal_android", "assignments.*", map[string]string{
						"type": "groupAssignmentTarget",
					}),
				),
			},
		},
	})
}

func TestAndroidEnrollmentNotificationsResource_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testConfigAndroidEmailMinimal(),
				ExpectError: regexp.MustCompile("Invalid Android Enrollment Notification data"),
			},
		},
	})
}

// Android Platform Configuration Functions
func testConfigAndroidEmailMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_android_email_minimal.tf")
	if err != nil {
		panic("failed to load android email minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigAndroidEmailMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_android_email_maximal.tf")
	if err != nil {
		panic("failed to load android email maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigAndroidPushMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_android_push_maximal.tf")
	if err != nil {
		panic("failed to load android push maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigAndroidAllMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_android_all_maximal.tf")
	if err != nil {
		panic("failed to load android all maximal config: " + err.Error())
	}
	return unitTestConfig
}

// AndroidForWork Platform Configuration Functions
func testConfigAndroidForWorkEmailMinimal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_androidForWork_email_minimal.tf")
	if err != nil {
		panic("failed to load androidForWork email minimal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigAndroidForWorkEmailMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_androidForWork_email_maximal.tf")
	if err != nil {
		panic("failed to load androidForWork email maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigAndroidForWorkPushMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_androidForWork_push_maximal.tf")
	if err != nil {
		panic("failed to load androidForWork push maximal config: " + err.Error())
	}
	return unitTestConfig
}

func testConfigAndroidForWorkAllMaximal() string {
	unitTestConfig, err := helpers.ParseHCLFile("tests/terraform/unit/resource_androidForWork_all_maximal.tf")
	if err != nil {
		panic("failed to load androidForWork all maximal config: " + err.Error())
	}
	return unitTestConfig
}

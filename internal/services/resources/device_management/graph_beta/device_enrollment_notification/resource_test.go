package graphBetaDeviceEnrollmentNotification_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	enrollmentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/device_enrollment_notification/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

// Helper function to load test configs from unit directory
func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

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

func TestUnitResourceDeviceEnrollmentNotification_01_Schema(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01_android_email_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".email_minimal_android").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_minimal_android").Key("display_name").HasValue("email minimal android"),
					check.That(resourceType+".email_minimal_android").Key("platform_type").HasValue("android"),
					check.That(resourceType+".email_minimal_android").Key("default_locale").HasValue("en-US"),
					check.That(resourceType+".email_minimal_android").Key("role_scope_tag_ids.#").HasValue("1"),
					check.That(resourceType+".email_minimal_android").Key("role_scope_tag_ids.*").ContainsTypeSetElement("0"),
				),
			},
		},
	})
}

func TestUnitResourceDeviceEnrollmentNotification_02_PlatformTypes(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_02_android_email_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".email_maximal_android").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_maximal_android").Key("display_name").HasValue("email maximal android"),
					check.That(resourceType+".email_maximal_android").Key("platform_type").HasValue("android"),
					check.That(resourceType+".email_maximal_android").Key("notification_templates.#").HasValue("1"),
					check.That(resourceType+".email_maximal_android").Key("notification_templates.*").ContainsTypeSetElement("email"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_06_androidForWork_email_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".email_maximal_androidforwork").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_maximal_androidforwork").Key("display_name").HasValue("email maximal androidForWork"),
					check.That(resourceType+".email_maximal_androidforwork").Key("platform_type").HasValue("androidForWork"),
					check.That(resourceType+".email_maximal_androidforwork").Key("notification_templates.#").HasValue("1"),
					check.That(resourceType+".email_maximal_androidforwork").Key("notification_templates.*").ContainsTypeSetElement("email"),
				),
			},
		},
	})
}

func TestUnitResourceDeviceEnrollmentNotification_03_AllTerraformConfigurations(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Android Platform Tests
			{
				Config: loadUnitTestTerraform("resource_01_android_email_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".email_minimal_android").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_minimal_android").Key("platform_type").HasValue("android"),
					check.That(resourceType+".email_minimal_android").Key("notification_templates.#").HasValue("1"),
					check.That(resourceType+".email_minimal_android").Key("notification_templates.*").ContainsTypeSetElement("email"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_02_android_email_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".email_maximal_android").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_maximal_android").Key("platform_type").HasValue("android"),
					check.That(resourceType+".email_maximal_android").Key("notification_templates.#").HasValue("1"),
					check.That(resourceType+".email_maximal_android").Key("notification_templates.*").ContainsTypeSetElement("email"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_03_android_push_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".push_maximal_android").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".push_maximal_android").Key("platform_type").HasValue("android"),
					check.That(resourceType+".push_maximal_android").Key("notification_templates.#").HasValue("1"),
					check.That(resourceType+".push_maximal_android").Key("notification_templates.*").ContainsTypeSetElement("push"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_04_android_all_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".all_android").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".all_android").Key("platform_type").HasValue("android"),
					check.That(resourceType+".all_android").Key("notification_templates.#").HasValue("2"),
					check.That(resourceType+".all_android").Key("notification_templates.*").ContainsTypeSetElement("email"),
					check.That(resourceType+".all_android").Key("notification_templates.*").ContainsTypeSetElement("push"),
				),
			},
			// AndroidForWork Platform Tests
			{
				Config: loadUnitTestTerraform("resource_05_androidForWork_email_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".email_minimal_androidforwork").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_minimal_androidforwork").Key("platform_type").HasValue("androidForWork"),
					check.That(resourceType+".email_minimal_androidforwork").Key("notification_templates.#").HasValue("1"),
					check.That(resourceType+".email_minimal_androidforwork").Key("notification_templates.*").ContainsTypeSetElement("email"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_06_androidForWork_email_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".email_maximal_androidforwork").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_maximal_androidforwork").Key("platform_type").HasValue("androidForWork"),
					check.That(resourceType+".email_maximal_androidforwork").Key("notification_templates.#").HasValue("1"),
					check.That(resourceType+".email_maximal_androidforwork").Key("notification_templates.*").ContainsTypeSetElement("email"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_07_androidForWork_push_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".push_maximal_androidForWork").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".push_maximal_androidForWork").Key("platform_type").HasValue("androidForWork"),
					check.That(resourceType+".push_maximal_androidForWork").Key("notification_templates.#").HasValue("1"),
					check.That(resourceType+".push_maximal_androidForWork").Key("notification_templates.*").ContainsTypeSetElement("push"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_08_androidForWork_all_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".all_androidforwork").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".all_androidforwork").Key("platform_type").HasValue("androidForWork"),
					check.That(resourceType+".all_androidforwork").Key("notification_templates.#").HasValue("2"),
					check.That(resourceType+".all_androidforwork").Key("notification_templates.*").ContainsTypeSetElement("email"),
					check.That(resourceType+".all_androidforwork").Key("notification_templates.*").ContainsTypeSetElement("push"),
				),
			},
		},
	})
}

func TestUnitResourceDeviceEnrollmentNotification_04_BrandingOptions(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_01_android_email_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".email_minimal_android").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_minimal_android").Key("branding_options.#").HasValue("1"),
					check.That(resourceType+".email_minimal_android").Key("branding_options.*").ContainsTypeSetElement("none"),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_02_android_email_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".email_maximal_android").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_maximal_android").Key("branding_options.#").HasValue("5"),
					check.That(resourceType+".email_maximal_android").Key("branding_options.*").ContainsTypeSetElement("includeCompanyLogo"),
					check.That(resourceType+".email_maximal_android").Key("branding_options.*").ContainsTypeSetElement("includeCompanyName"),
					check.That(resourceType+".email_maximal_android").Key("branding_options.*").ContainsTypeSetElement("includeCompanyPortalLink"),
					check.That(resourceType+".email_maximal_android").Key("branding_options.*").ContainsTypeSetElement("includeContactInformation"),
					check.That(resourceType+".email_maximal_android").Key("branding_options.*").ContainsTypeSetElement("includeDeviceDetails"),
				),
			},
		},
	})
}

func TestUnitResourceDeviceEnrollmentNotification_05_LocalizedMessages(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_02_android_email_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".email_maximal_android").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_maximal_android").Key("localized_notification_messages.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".email_maximal_android", "localized_notification_messages.*", map[string]string{
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

func TestUnitResourceDeviceEnrollmentNotification_06_Assignments(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_02_android_email_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".email_maximal_android").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_maximal_android").Key("assignments.#").HasValue("1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceType+".email_maximal_android", "assignments.*", map[string]string{
						"type": "groupAssignmentTarget",
					}),
				),
			},
		},
	})
}

func TestUnitResourceDeviceEnrollmentNotification_07_ErrorHandling(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, enrollmentMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer enrollmentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("resource_01_android_email_minimal.tf"),
				ExpectError: regexp.MustCompile("Invalid Android Enrollment Notification data"),
			},
		},
	})
}

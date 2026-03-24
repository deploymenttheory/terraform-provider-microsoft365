package graphBetaChangeNotificationsSubscription_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	subscriptionMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/change_notifications/graph_beta/subscription/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func setupMockEnvironment() (*mocks.Mocks, *subscriptionMocks.SubscriptionMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	subscriptionMock := &subscriptionMocks.SubscriptionMock{}
	subscriptionMock.RegisterMocks()
	return mockClient, subscriptionMock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

// SUB001: minimal create + import
func TestUnitResourceChangeNotificationsSubscription_01_SUB001(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, subscriptionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer subscriptionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_sub001-minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".sub001_minimal").Key("id").HasValue("11111111-1111-1111-1111-111111111111"),
					check.That(resourceType+".sub001_minimal").Key("change_type").HasValue("updated"),
					check.That(resourceType+".sub001_minimal").Key("resource").HasValue("users"),
					check.That(resourceType+".sub001_minimal").Key("notification_url").HasValue("https://example.com/webhook"),
				),
			},
			{
				ResourceName:            resourceType + ".sub001_minimal",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts", "client_state"},
			},
		},
	})
}

// SUB002: update expiration
func TestUnitResourceChangeNotificationsSubscription_02_SUB002(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, subscriptionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer subscriptionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("resource_sub001-minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".sub001_minimal").Key("id").Exists(),
				),
			},
			{
				Config: loadUnitTestTerraform("resource_sub002-update_expiration.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".sub001_minimal").Key("expiration_date_time").HasValue("2030-06-01T12:00:00Z"),
				),
			},
		},
	})
}

// SUB003: create error path (invalid notification URL)
func TestUnitResourceChangeNotificationsSubscription_03_SUB003(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, subscriptionMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer subscriptionMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "microsoft365_graph_beta_change_notifications_subscription" "sub003_error" {
  change_type          = "updated"
  notification_url     = "https://invalid-webhook.example/not-valid"
  resource             = "users"
  expiration_date_time = "2030-01-01T12:00:00Z"
}
`,
				ExpectError: regexp.MustCompile(`(?i)validation|error|failed`),
			},
		},
	})
}

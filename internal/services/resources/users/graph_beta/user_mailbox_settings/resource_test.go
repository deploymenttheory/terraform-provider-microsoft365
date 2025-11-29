package graphBetaUsersUserMailboxSettings_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	mailboxSettings "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/user_mailbox_settings"
	mailboxSettingsMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/user_mailbox_settings/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

const (
	resourceType = mailboxSettings.ResourceName
)

func setupMockEnvironment() (*mocks.Mocks, *mailboxSettingsMocks.UserMailboxSettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mailboxSettingsMock := &mailboxSettingsMocks.UserMailboxSettingsMock{}
	mailboxSettingsMock.RegisterMocks()
	return mockClient, mailboxSettingsMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *mailboxSettingsMocks.UserMailboxSettingsMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	mailboxSettingsMock := &mailboxSettingsMocks.UserMailboxSettingsMock{}
	mailboxSettingsMock.RegisterErrorMocks()
	return mockClient, mailboxSettingsMock
}

func TestUserMailboxSettingsResource_Basic(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mailboxSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mailboxSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("user_id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(resourceType+".minimal").Key("time_zone").HasValue("UTC"),
					check.That(resourceType+".minimal").Key("date_format").HasValue("MM/dd/yyyy"),
					check.That(resourceType+".minimal").Key("time_format").HasValue("hh:mm tt"),
					check.That(resourceType+".minimal").Key("id").Exists(),
				),
			},
			{
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUserMailboxSettingsResource_Update(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mailboxSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mailboxSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".minimal").Key("user_id").HasValue("00000000-0000-0000-0000-000000000001"),
					check.That(resourceType+".minimal").Key("time_zone").HasValue("UTC"),
				),
			},
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("user_id").HasValue("00000000-0000-0000-0000-000000000002"),
					check.That(resourceType+".maximal").Key("time_zone").HasValue("UTC"),
					check.That(resourceType+".maximal").Key("date_format").HasValue("MM/dd/yyyy"),
					check.That(resourceType+".maximal").Key("time_format").HasValue("hh:mm tt"),
					check.That(resourceType+".maximal").Key("delegate_meeting_message_delivery_options").HasValue("sendToDelegateOnly"),
					check.That(resourceType+".maximal").Key("automatic_replies_setting.0.status").HasValue("scheduled"),
					check.That(resourceType+".maximal").Key("automatic_replies_setting.0.external_audience").HasValue("all"),
					check.That(resourceType+".maximal").Key("automatic_replies_setting.0.scheduled_start_date_time.0.date_time").HasValue("2016-03-14T07:00:00.0000000"),
					check.That(resourceType+".maximal").Key("automatic_replies_setting.0.scheduled_start_date_time.0.time_zone").HasValue("UTC"),
					check.That(resourceType+".maximal").Key("automatic_replies_setting.0.scheduled_end_date_time.0.date_time").HasValue("2016-03-28T07:00:00.0000000"),
					check.That(resourceType+".maximal").Key("automatic_replies_setting.0.scheduled_end_date_time.0.time_zone").HasValue("UTC"),
					check.That(resourceType+".maximal").Key("language.0.locale").HasValue("en-US"),
					check.That(resourceType+".maximal").Key("language.0.display_name").Exists(),
					check.That(resourceType+".maximal").Key("working_hours.0.start_time").HasValue("08:00:00.0000000"),
					check.That(resourceType+".maximal").Key("working_hours.0.end_time").HasValue("17:00:00.0000000"),
					check.That(resourceType+".maximal").Key("working_hours.0.time_zone.0.name").HasValue("Pacific Standard Time"),
					check.That(resourceType+".maximal").Key("working_hours.0.days_of_week.#").HasValue("5"),
				),
			},
			{
				ResourceName:      resourceType + ".maximal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestUserMailboxSettingsResource_AutomaticReplies(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mailboxSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mailboxSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("automatic_replies_setting.0.status").HasValue("scheduled"),
					check.That(resourceType+".maximal").Key("automatic_replies_setting.0.external_audience").HasValue("all"),
					check.That(resourceType+".maximal").Key("automatic_replies_setting.0.internal_reply_message").Exists(),
					check.That(resourceType+".maximal").Key("automatic_replies_setting.0.external_reply_message").Exists(),
				),
			},
		},
	})
}

func TestUserMailboxSettingsResource_WorkingHours(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mailboxSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mailboxSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".maximal").Key("working_hours.0.days_of_week.#").HasValue("5"),
					check.That(resourceType+".maximal").Key("working_hours.0.start_time").HasValue("08:00:00.0000000"),
					check.That(resourceType+".maximal").Key("working_hours.0.end_time").HasValue("17:00:00.0000000"),
					check.That(resourceType+".maximal").Key("working_hours.0.time_zone.0.name").HasValue("Pacific Standard Time"),
				),
			},
		},
	})
}

func TestUserMailboxSettingsResource_RequiresImport(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, mailboxSettingsMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer mailboxSettingsMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType + ".minimal").Key("id").Exists(),
				),
			},
			{
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testConfigBasic returns a minimal Terraform configuration for testing
func testConfigBasic() string {
	return mocks.LoadUnitTerraformConfig("resource_minimal.tf")
}

// testConfigUpdate returns a maximal Terraform configuration for testing
func testConfigUpdate() string {
	return mocks.LoadUnitTerraformConfig("resource_maximal.tf")
}

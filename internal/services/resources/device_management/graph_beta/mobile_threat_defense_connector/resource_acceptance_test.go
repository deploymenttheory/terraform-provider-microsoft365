package graphBetaMobileThreatDefenseConnector_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	graphBetaMobileThreatDefenseConnector "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/mobile_threat_defense_connector"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaMobileThreatDefenseConnector.ResourceName

	// testResource is the test resource implementation for mobile threat defense connectors
	testResource = graphBetaMobileThreatDefenseConnector.MobileThreatDefenseConnectorTestResource{}
)

const (
	// Partner IDs used by the acceptance tests. Neither partner is connected in the test tenant,
	// so the tests own the connectors they create and delete them by exact ID.
	accPradeoPartnerId  = "8d0ed095-8191-4bd3-8a41-953b22d51ff7"
	accIVerifyPartnerId = "fc06ff9b-e2f9-4221-9a68-465684ee8662"
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// preCheckConnectorAbsent stops the test when the partner already has a connector in the tenant.
// Graph treats create as an upsert, so running against an existing connector would overwrite and then delete it.
func preCheckConnectorAbsent(t *testing.T, partnerId string) {
	t.Helper()
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		t.Fatalf("failed to create Graph client: %v", err)
	}
	_, err = graphClient.DeviceManagement().MobileThreatDefenseConnectors().ByMobileThreatDefenseConnectorId(partnerId).Get(context.Background(), nil)
	if err == nil {
		t.Fatalf("mobile threat defense connector %s already exists in the tenant; refusing to run because the test would overwrite and delete it", partnerId)
	}
	if errorInfo := errors.GraphError(context.Background(), err); errorInfo.StatusCode != 404 {
		t.Fatalf("failed to check whether mobile threat defense connector %s exists: %v", partnerId, err)
	}
}

func TestAccResourceMobileThreatDefenseConnector_01_Lifecycle(t *testing.T) {
	address := resourceType + ".test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			preCheckConnectorAbsent(t, accPradeoPartnerId)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 10*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal Pradeo connector")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(address).ExistsInGraph(testResource),
					check.That(address).Key("id").HasValue(accPradeoPartnerId),
					check.That(address).Key("android_enabled").HasValue("false"),
					check.That(address).Key("ios_enabled").HasValue("false"),
					check.That(address).Key("mac_enabled").HasValue("false"),
					check.That(address).Key("windows_enabled").HasValue("false"),
					check.That(address).Key("partner_unresponsiveness_threshold_in_days").HasValue("7"),
					check.That(address).Key("android_device_blocked_on_missing_partner_data").HasValue("false"),
					check.That(address).Key("partner_state").MatchesRegex(regexp.MustCompile(`^[a-zA-Z]+$`)),
					check.That(address).Key("last_heartbeat_date_time").DoesNotExist(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal Pradeo connector")
				},
				ResourceName:            address,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating Pradeo connector to maximal")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(address).ExistsInGraph(testResource),
					check.That(address).Key("id").HasValue(accPradeoPartnerId),
					check.That(address).Key("android_enabled").HasValue("true"),
					check.That(address).Key("ios_enabled").HasValue("true"),
					check.That(address).Key("mac_enabled").HasValue("true"),
					check.That(address).Key("windows_enabled").HasValue("true"),
					check.That(address).Key("android_mobile_application_management_enabled").HasValue("false"),
					check.That(address).Key("ios_mobile_application_management_enabled").HasValue("false"),
					check.That(address).Key("windows_mobile_application_management_enabled").HasValue("true"),
					check.That(address).Key("partner_unsupported_os_version_blocked").HasValue("true"),
					check.That(address).Key("partner_unresponsiveness_threshold_in_days").HasValue("14"),
					check.That(address).Key("allow_partner_to_collect_ios_application_metadata").HasValue("true"),
					check.That(address).Key("allow_partner_to_collect_ios_personal_application_metadata").HasValue("true"),
					check.That(address).Key("allow_partner_to_collect_ios_certificate_metadata").HasValue("true"),
					check.That(address).Key("allow_partner_to_collect_ios_personal_certificate_metadata").HasValue("true"),
					check.That(address).Key("grant_mobile_threat_defense_partner_role").HasValue("true"),
					check.That(address).Key("android_device_blocked_on_missing_partner_data").HasValue("true"),
					check.That(address).Key("ios_device_blocked_on_missing_partner_data").HasValue("true"),
					check.That(address).Key("mac_device_blocked_on_missing_partner_data").HasValue("true"),
					check.That(address).Key("windows_device_blocked_on_missing_partner_data").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing maximal Pradeo connector")
				},
				ResourceName:            address,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating Pradeo connector back to minimal")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(address).ExistsInGraph(testResource),
					check.That(address).Key("android_enabled").HasValue("false"),
					check.That(address).Key("ios_enabled").HasValue("false"),
					check.That(address).Key("mac_enabled").HasValue("false"),
					check.That(address).Key("windows_enabled").HasValue("false"),
					check.That(address).Key("windows_mobile_application_management_enabled").HasValue("false"),
					check.That(address).Key("partner_unresponsiveness_threshold_in_days").HasValue("7"),
					check.That(address).Key("grant_mobile_threat_defense_partner_role").HasValue("false"),
					check.That(address).Key("windows_device_blocked_on_missing_partner_data").HasValue("false"),
				),
			},
		},
	})
}

func TestAccResourceMobileThreatDefenseConnector_02_CreateWithSettings(t *testing.T) {
	address := resourceType + ".iverify"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			preCheckConnectorAbsent(t, accIVerifyPartnerId)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 10*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating iVerify connector with settings")
				},
				Config: loadAcceptanceTestTerraform("resource_iverify.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(address).ExistsInGraph(testResource),
					check.That(address).Key("id").HasValue(accIVerifyPartnerId),
					check.That(address).Key("android_enabled").HasValue("true"),
					check.That(address).Key("ios_enabled").HasValue("true"),
					check.That(address).Key("mac_enabled").HasValue("false"),
					check.That(address).Key("windows_enabled").HasValue("true"),
					check.That(address).Key("android_mobile_application_management_enabled").HasValue("false"),
					check.That(address).Key("ios_mobile_application_management_enabled").HasValue("false"),
					check.That(address).Key("windows_mobile_application_management_enabled").HasValue("false"),
					check.That(address).Key("partner_unsupported_os_version_blocked").HasValue("true"),
					check.That(address).Key("partner_unresponsiveness_threshold_in_days").HasValue("30"),
					check.That(address).Key("allow_partner_to_collect_ios_application_metadata").HasValue("true"),
					check.That(address).Key("allow_partner_to_collect_ios_personal_application_metadata").HasValue("false"),
					check.That(address).Key("allow_partner_to_collect_ios_certificate_metadata").HasValue("true"),
					check.That(address).Key("allow_partner_to_collect_ios_personal_certificate_metadata").HasValue("false"),
					check.That(address).Key("grant_mobile_threat_defense_partner_role").HasValue("false"),
					check.That(address).Key("android_device_blocked_on_missing_partner_data").HasValue("true"),
					check.That(address).Key("mac_device_blocked_on_missing_partner_data").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing iVerify connector")
				},
				ResourceName:            address,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

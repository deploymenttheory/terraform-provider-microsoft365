package graphBetaChangeNotificationsSubscription_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaChangeNotificationsSubscription "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/change_notifications/graph_beta/subscription"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	resourceType = graphBetaChangeNotificationsSubscription.ResourceName
	testResource = graphBetaChangeNotificationsSubscription.SubscriptionTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// SUB001: create, update expiration, import (payload shape from Graph subscription examples:
// created + Inbox messages + sample notificationUrl / clientState / TLS).
func TestAccResourceChangeNotificationsSubscription_01_SUB001(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating SUB001 mail Inbox messages subscription")
				},
				Config: loadAcceptanceTestTerraform("resource_sub001-acc_users_subscription.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("change notifications subscription", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".sub001_mail").ExistsInGraph(testResource),
					check.That(resourceType+".sub001_mail").Key("id").Exists(),
					check.That(resourceType+".sub001_mail").Key("resource").HasValue(`me/mailFolders('Inbox')/messages`),
					check.That(resourceType+".sub001_mail").Key("change_type").HasValue("created"),
					check.That(resourceType+".sub001_mail").Key("notification_url").HasValue("https://webhook.azurewebsites.net/api/send/myNotifyClient"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating SUB002 subscription expiration")
				},
				Config: loadAcceptanceTestTerraform("resource_sub002-acc_users_subscription.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("change notifications subscription", 30*time.Second)
						time.Sleep(30 * time.Second)
						return nil
					},
					check.That(resourceType+".sub001_mail").ExistsInGraph(testResource),
					check.That(resourceType+".sub001_mail").Key("id").Exists(),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing SUB001 subscription")
				},
				ResourceName:      resourceType + ".sub001_mail",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
					"client_state",
					"expiration_date_time",
				},
			},
		},
	})
}

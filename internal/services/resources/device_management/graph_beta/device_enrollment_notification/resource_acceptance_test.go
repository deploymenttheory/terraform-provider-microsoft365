package graphBetaDeviceEnrollmentNotification_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaDeviceEnrollmentNotification "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/device_enrollment_notification"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaDeviceEnrollmentNotification.ResourceName

	// testResource is the test resource implementation for device enrollment notifications
	testResource = graphBetaDeviceEnrollmentNotification.DeviceEnrollmentNotificationTestResource{}
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// Android Platform Tests
func TestAccResourceDeviceEnrollmentNotification_01_AndroidEmailMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating Android email minimal notification")
				},
				Config: loadAcceptanceTestTerraform("resource_01_android_email_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("device enrollment notification", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".email_minimal_android").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_minimal_android").Key("display_name").HasValue("email minimal android"),
					check.That(resourceType+".email_minimal_android").Key("platform_type").HasValue("android"),
					check.That(resourceType+".email_minimal_android").Key("notification_templates.#").HasValue("1"),
					check.That(resourceType+".email_minimal_android").Key("notification_templates.*").ContainsTypeSetElement("email"),
					check.That(resourceType+".email_minimal_android").Key("branding_options.#").HasValue("1"),
					check.That(resourceType+".email_minimal_android").Key("branding_options.*").ContainsTypeSetElement("none"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Android email minimal notification")
				},
				ResourceName:            resourceType + ".email_minimal_android",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"branding_options", "platform_type"},
			},
		},
	})
}

func TestAccResourceDeviceEnrollmentNotification_02_AndroidEmailMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating Android email maximal notification")
				},
				Config: loadAcceptanceTestTerraform("resource_02_android_email_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("device enrollment notification", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".email_maximal_android").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_maximal_android").Key("display_name").HasValue("email maximal android"),
					check.That(resourceType+".email_maximal_android").Key("platform_type").HasValue("android"),
					check.That(resourceType+".email_maximal_android").Key("notification_templates.#").HasValue("1"),
					check.That(resourceType+".email_maximal_android").Key("notification_templates.*").ContainsTypeSetElement("email"),
					check.That(resourceType+".email_maximal_android").Key("branding_options.#").HasValue("5"),
					check.That(resourceType+".email_maximal_android").Key("branding_options.*").ContainsTypeSetElement("includeCompanyLogo"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Android email maximal notification")
				},
				ResourceName:            resourceType + ".email_maximal_android",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"branding_options", "platform_type"},
			},
		},
	})
}

func TestAccResourceDeviceEnrollmentNotification_03_AndroidPushMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating Android push maximal notification")
				},
				Config: loadAcceptanceTestTerraform("resource_03_android_push_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("device enrollment notification", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".push_maximal_android").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".push_maximal_android").Key("display_name").HasValue("push maximal android"),
					check.That(resourceType+".push_maximal_android").Key("platform_type").HasValue("android"),
					check.That(resourceType+".push_maximal_android").Key("notification_templates.#").HasValue("1"),
					check.That(resourceType+".push_maximal_android").Key("notification_templates.*").ContainsTypeSetElement("push"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Android push maximal notification")
				},
				ResourceName:            resourceType + ".push_maximal_android",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"branding_options", "platform_type"},
			},
		},
	})
}

func TestAccResourceDeviceEnrollmentNotification_04_AndroidAllMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating Android all maximal notification")
				},
				Config: loadAcceptanceTestTerraform("resource_04_android_all_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("device enrollment notification", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".all_android").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".all_android").Key("display_name").HasValue("Complete Test - all android"),
					check.That(resourceType+".all_android").Key("platform_type").HasValue("android"),
					check.That(resourceType+".all_android").Key("notification_templates.#").HasValue("2"),
					check.That(resourceType+".all_android").Key("notification_templates.*").ContainsTypeSetElement("email"),
					check.That(resourceType+".all_android").Key("notification_templates.*").ContainsTypeSetElement("push"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Android all maximal notification")
				},
				ResourceName:            resourceType + ".all_android",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"branding_options", "platform_type"},
			},
		},
	})
}

// AndroidForWork Platform Tests
func TestAccResourceDeviceEnrollmentNotification_05_AndroidForWorkEmailMinimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AndroidForWork email minimal notification")
				},
				Config: loadAcceptanceTestTerraform("resource_05_androidForWork_email_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("device enrollment notification", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".email_minimal_androidforwork").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_minimal_androidforwork").Key("display_name").HasValue("email minimal androidForWork"),
					check.That(resourceType+".email_minimal_androidforwork").Key("platform_type").HasValue("androidForWork"),
					check.That(resourceType+".email_minimal_androidforwork").Key("notification_templates.#").HasValue("1"),
					check.That(resourceType+".email_minimal_androidforwork").Key("notification_templates.*").ContainsTypeSetElement("email"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AndroidForWork email minimal notification")
				},
				ResourceName:            resourceType + ".email_minimal_androidforwork",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"branding_options", "platform_type"},
			},
		},
	})
}

func TestAccResourceDeviceEnrollmentNotification_06_AndroidForWorkEmailMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AndroidForWork email maximal notification")
				},
				Config: loadAcceptanceTestTerraform("resource_06_androidForWork_email_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("device enrollment notification", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".email_maximal_androidforwork").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".email_maximal_androidforwork").Key("display_name").HasValue("email maximal androidForWork"),
					check.That(resourceType+".email_maximal_androidforwork").Key("platform_type").HasValue("androidForWork"),
					check.That(resourceType+".email_maximal_androidforwork").Key("notification_templates.#").HasValue("1"),
					check.That(resourceType+".email_maximal_androidforwork").Key("notification_templates.*").ContainsTypeSetElement("email"),
					check.That(resourceType+".email_maximal_androidforwork").Key("branding_options.#").HasValue("5"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AndroidForWork email maximal notification")
				},
				ResourceName:            resourceType + ".email_maximal_androidforwork",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"branding_options", "platform_type"},
			},
		},
	})
}

func TestAccResourceDeviceEnrollmentNotification_07_AndroidForWorkPushMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AndroidForWork push maximal notification")
				},
				Config: loadAcceptanceTestTerraform("resource_07_androidForWork_push_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("device enrollment notification", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".push_maximal_androidForWork").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".push_maximal_androidForWork").Key("display_name").HasValue("push maximal android"),
					check.That(resourceType+".push_maximal_androidForWork").Key("platform_type").HasValue("androidForWork"),
					check.That(resourceType+".push_maximal_androidForWork").Key("notification_templates.#").HasValue("1"),
					check.That(resourceType+".push_maximal_androidForWork").Key("notification_templates.*").ContainsTypeSetElement("push"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AndroidForWork push maximal notification")
				},
				ResourceName:            resourceType + ".push_maximal_androidForWork",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"branding_options", "platform_type"},
			},
		},
	})
}

func TestAccResourceDeviceEnrollmentNotification_08_AndroidForWorkAllMaximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			20*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: constants.ExternalProviderTimeVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating AndroidForWork all maximal notification")
				},
				Config: loadAcceptanceTestTerraform("resource_08_androidForWork_all_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("device enrollment notification", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".all_androidforwork").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+_EnrollmentNotificationsConfiguration$`)),
					check.That(resourceType+".all_androidforwork").Key("display_name").HasValue("Complete Test - all androidForWork"),
					check.That(resourceType+".all_androidforwork").Key("platform_type").HasValue("androidForWork"),
					check.That(resourceType+".all_androidforwork").Key("notification_templates.#").HasValue("2"),
					check.That(resourceType+".all_androidforwork").Key("notification_templates.*").ContainsTypeSetElement("email"),
					check.That(resourceType+".all_androidforwork").Key("notification_templates.*").ContainsTypeSetElement("push"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing AndroidForWork all maximal notification")
				},
				ResourceName:            resourceType + ".all_androidforwork",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"branding_options", "platform_type"},
			},
		},
	})
}

package graphBetaWin32App_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	win32 "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/win32_app"
	group "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/groups/graph_beta/group"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

const resourceType = win32.ResourceName

var testResource = win32.Win32AppTestResource{}

func loadAcceptanceTestTerraform(filename string, packagesURL string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	config = strings.ReplaceAll(config, "http://win32-packages.test", packagesURL)
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccResourceWin32App_01_Scenario_Minimal(t *testing.T) {
	packages := newAcceptancePackages(t)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			packages.prepare(t, "140.0.2")
		},
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("001_scenario_minimal.tf", packages.URL),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_001").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_001").Key("display_name").Exists(),
					check.That(resourceType+".test_001").Key("publisher").HasValue("Mozilla"),
					check.That(resourceType+".test_001").Key("content_version.#").HasValue("1"),
					check.That(resourceType+".test_001").Key("committed_content_version").Exists(),
				),
			},
			{Config: loadAcceptanceTestTerraform("001_scenario_minimal.tf", packages.URL), PlanOnly: true, ExpectNonEmptyPlan: false},
			importStep(resourceType + ".test_001"),
		},
	})
}

func TestAccResourceWin32App_02_Scenario_Maximal(t *testing.T) {
	packages := newAcceptancePackages(t)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			packages.prepare(t, "140.0.4")
		},
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		CheckDestroy: destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				Config: loadAcceptanceTestTerraform("002_scenario_maximal.tf", packages.URL),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_002").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_002").Key("display_version").HasValue("140.0.4"),
					check.That(resourceType+".test_002").Key("is_featured").HasValue("true"),
					check.That(resourceType+".test_002").Key("minimum_memory_in_mb").HasValue("2048"),
					check.That(resourceType+".test_002").Key("install_experience.run_as_account").HasValue("system"),
					check.That(resourceType+".test_002").Key("return_codes.#").HasValue("5"),
					check.That(resourceType+".test_002").Key("content_version.#").HasValue("1"),
				),
			},
			{Config: loadAcceptanceTestTerraform("002_scenario_maximal.tf", packages.URL), PlanOnly: true, ExpectNonEmptyPlan: false},
			importStep(resourceType + ".test_002"),
		},
	})
}

// Exercise two versions in each format and a format switch on the same app
// and group assignment. Minimal/maximal configurations are tested separately.
func TestAccResourceWin32App_03_Lifecycle_ContentUpdate(t *testing.T) {
	packages := newAcceptancePackages(t)
	appAddress := resourceType + ".test_003"
	checkContent := checkContentLifecycle(appAddress)
	var steps []resource.TestStep
	for _, filename := range []string{
		"003_lifecycle_content_update_step_1.tf",
		"003_lifecycle_content_update_step_2.tf",
		"003_lifecycle_content_update_step_3.tf",
		"003_lifecycle_content_update_step_4.tf",
	} {
		config := loadAcceptanceTestTerraform(filename, packages.URL)
		steps = append(steps,
			resource.TestStep{Config: config, Check: checkContent},
			resource.TestStep{Config: config, PlanOnly: true, ExpectNonEmptyPlan: false},
		)
	}
	steps = append(steps, importStep(appAddress))
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			packages.prepare(t, "140.0.2", "140.0.4")
		},
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		CheckDestroy: destroy.CheckDestroyedTypesFunc(30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: testResource},
			destroy.ResourceTypeMapping{ResourceType: group.ResourceName, TestResource: group.GroupTestResource{}},
		),
		Steps: steps,
	})
}

func importStep(address string) resource.TestStep {
	return resource.TestStep{
		ResourceName: address, ImportState: true, ImportStateVerify: true,
		// Graph cannot return the configured source or outer package filename.
		ImportStateVerifyIgnore: []string{"app_installer", "app_installer_zip", "file_name"},
	}
}

func checkContentLifecycle(appAddress string) resource.TestCheckFunc {
	var appID, assignmentID, lastVersion string
	return func(state *terraform.State) error {
		app, ok := state.RootModule().Resources[appAddress]
		if !ok || app.Primary == nil || app.Primary.ID == "" {
			return fmt.Errorf("app missing from state")
		}
		assignment, ok := state.RootModule().Resources["microsoft365_graph_beta_device_and_app_management_mobile_app_assignment.test"]
		if !ok || assignment.Primary == nil || assignment.Primary.ID == "" {
			return fmt.Errorf("assignment missing from state")
		}
		if appID == "" {
			appID, assignmentID = app.Primary.ID, assignment.Primary.ID
		} else if appID != app.Primary.ID || assignmentID != assignment.Primary.ID {
			return fmt.Errorf("content update replaced application or assignment")
		}
		version := app.Primary.Attributes["committed_content_version"]
		if version == "" || version == lastVersion {
			return fmt.Errorf("expected a new committed content version, got %q after %q", version, lastVersion)
		}
		if app.Primary.Attributes["content_version.#"] != "1" || app.Primary.Attributes["content_version.0.id"] != version {
			return fmt.Errorf("state must contain only the current committed content version")
		}
		client, err := acceptance.TestGraphClient()
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		remote, err := client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).Get(ctx, nil)
		if err != nil {
			return err
		}
		remoteApp, ok := remote.(graphmodels.Win32LobAppable)
		if !ok || remoteApp.GetCommittedContentVersion() == nil || *remoteApp.GetCommittedContentVersion() != version {
			return fmt.Errorf("Graph committed content version does not match state")
		}
		assignments, err := client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).Assignments().Get(ctx, nil)
		if err != nil {
			return err
		}
		groupID := assignment.Primary.Attributes["target.group_id"]
		for _, a := range assignments.GetValue() {
			if a.GetId() != nil && *a.GetId() == assignmentID {
				target, ok := a.GetTarget().(graphmodels.GroupAssignmentTargetable)
				if ok && target.GetGroupId() != nil && *target.GetGroupId() == groupID {
					lastVersion = version
					return nil
				}
			}
		}
		return fmt.Errorf("original group assignment is missing from Graph")
	}
}

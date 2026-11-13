package graphBetaDeviceAndAppManagementIOSManagedMobileApp_test

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaIosManagedAppProtection "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/ios_managed_app_protection"
	graphBetaIOSManagedMobileApp "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/ios_managed_mobile_app"
)

var (
	accTestResource = graphBetaIOSManagedMobileApp.IOSManagedMobileAppTestResource{}
	policyResource  = graphBetaIosManagedAppProtection.ResourceName + ".test"
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func accTestCase(t *testing.T, steps []resource.TestStep) resource.TestCase {
	return resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		CheckDestroy: destroy.CheckDestroyedTypesFunc(
			30*time.Second,
			destroy.ResourceTypeMapping{ResourceType: resourceType, TestResource: accTestResource},
			destroy.ResourceTypeMapping{
				ResourceType: graphBetaIosManagedAppProtection.ResourceName,
				TestResource: graphBetaIosManagedAppProtection.IosManagedAppProtectionTestResource{},
			},
		),
		Steps: steps,
	}
}

// checkGraphPolicyApps asserts the live policy targets exactly the given bundles,
// catching apps dropped or left behind by full-list targetApps writes.
func checkGraphPolicyApps(bundles ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		policy, ok := s.RootModule().Resources[policyResource]
		if !ok {
			return fmt.Errorf("%s not found in state", policyResource)
		}
		client, err := acceptance.TestGraphClient()
		if err != nil {
			return fmt.Errorf("create Graph client: %w", err)
		}
		response, err := client.
			DeviceAppManagement().
			IosManagedAppProtections().
			ByIosManagedAppProtectionId(policy.Primary.ID).
			Apps().
			Get(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("list policy apps: %w", err)
		}
		got := make([]string, 0, len(response.GetValue()))
		for _, app := range response.GetValue() {
			if app.GetId() != nil {
				got = append(got, *app.GetId())
			}
		}
		want := make([]string, 0, len(bundles))
		for _, pkg := range bundles {
			want = append(want, pkg+".ios")
		}
		slices.Sort(got)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			return fmt.Errorf("policy %s targets %v, want %v", policy.Primary.ID, got, want)
		}
		return nil
	}
}

func importID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("%s not found in state", resourceName)
		}
		return rs.Primary.Attributes["managed_app_protection_id"] + "/" + rs.Primary.ID, nil
	}
}

func TestAccResourceIOSManagedMobileApp_01_Lifecycle(t *testing.T) {
	r := check.That(resourceType + ".test")
	resource.Test(t, accTestCase(t, []resource.TestStep{
		{
			PreConfig: func() { testlog.StepAction(resourceType, "01: Targeting one app on a new policy") },
			Config:    loadAcceptanceTestTerraform("01_lifecycle_step1.tf"),
			Check: resource.ComposeTestCheckFunc(
				r.ExistsInGraph(accTestResource),
				r.Key("id").HasValue("com.example.acctest.lifecycle.ios"),
				r.Key("managed_app_protection_id").
					MatchesOtherKey(check.That(policyResource).Key("id")),
				r.Key("mobile_app_identifier.bundle_id").HasValue("com.example.acctest.lifecycle"),
				r.Key("version").MatchesRegex(regexp.MustCompile(`^-?\d+$`)),
				checkGraphPolicyApps("com.example.acctest.lifecycle"),
			),
		},
		{
			PreConfig:               func() { testlog.StepAction(resourceType, "01: Importing the app") },
			ResourceName:            resourceType + ".test",
			ImportState:             true,
			ImportStateIdFunc:       importID(resourceType + ".test"),
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{"timeouts"},
		},
		{
			PreConfig: func() { testlog.StepAction(resourceType, "01: Changing the bundle replaces the app") },
			Config:    loadAcceptanceTestTerraform("01_lifecycle_step2.tf"),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(
						resourceType+".test",
						plancheck.ResourceActionDestroyBeforeCreate,
					),
				},
			},
			Check: resource.ComposeTestCheckFunc(
				r.ExistsInGraph(accTestResource),
				r.Key("id").HasValue("com.example.acctest.renamed.ios"),
				checkGraphPolicyApps("com.example.acctest.renamed"),
			),
		},
	}))
}

func TestAccResourceIOSManagedMobileApp_02_SiblingsOnOnePolicy(t *testing.T) {
	resource.Test(t, accTestCase(t, []resource.TestStep{
		{
			PreConfig: func() { testlog.StepAction(resourceType, "02: Targeting two apps on one policy concurrently") },
			Config:    loadAcceptanceTestTerraform("02_siblings_step1.tf"),
			Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".first").ExistsInGraph(accTestResource),
				check.That(resourceType+".second").ExistsInGraph(accTestResource),
				checkGraphPolicyApps("com.example.acctest.first", "com.example.acctest.second"),
			),
		},
		{
			PreConfig: func() { testlog.StepAction(resourceType, "02: Removing one app keeps its sibling") },
			Config:    loadAcceptanceTestTerraform("02_siblings_step2.tf"),
			Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".first").ExistsInGraph(accTestResource),
				checkGraphPolicyApps("com.example.acctest.first"),
			),
		},
	}))
}

package graphBetaWindowsUpdatesAutopatchRing_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaWindowsUpdatesRing "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_ring"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	resourceType = graphBetaWindowsUpdatesRing.ResourceName
	testResource = graphBetaWindowsUpdatesRing.WindowsUpdatesAutopatchRingTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance test config " + filename + ": " + err.Error())
	}
	return config
}

func externalProviders() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"random": {
			Source:            "hashicorp/random",
			VersionConstraint: constants.ExternalProviderRandomVersion,
		},
		"time": {
			Source:            "hashicorp/time",
			VersionConstraint: constants.ExternalProviderTimeVersion,
		},
	}
}

// TestAccResourceWindowsUpdateRing_01_InclusionOnly creates a ring with
// included group assignments only, then verifies import.
func TestAccResourceWindowsUpdateRing_01_InclusionOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create ring with inclusion assignments only")
				},
				Config: loadAcceptanceTestTerraform("01_inclusion_only.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("policy_id").Exists(),
					check.That(resourceType+".test").Key("display_name").HasValue("Acc Test Ring 01 Inclusion Only"),
					check.That(resourceType+".test").Key("is_paused").HasValue("false"),
					check.That(resourceType+".test").Key("created_date_time").Exists(),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Import ring")
				},
				ResourceName: resourceType + ".test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".test"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".test")
					}
					return rs.Primary.Attributes["policy_id"] + "/" + rs.Primary.ID, nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

// TestAccResourceWindowsUpdateRing_02_ExclusionOnly creates a ring with
// excluded group assignments only and empty included assignments.
func TestAccResourceWindowsUpdateRing_02_ExclusionOnly(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create ring with exclusion assignments only")
				},
				Config: loadAcceptanceTestTerraform("02_exclusion_only.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("display_name").HasValue("Acc Test Ring 02 Exclusion Only"),
					check.That(resourceType+".test").Key("is_paused").HasValue("false"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
		},
	})
}

// TestAccResourceWindowsUpdateRing_03_NoAssignments creates a ring with no
// group assignments (empty included and excluded).
func TestAccResourceWindowsUpdateRing_03_NoAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create ring with no group assignments")
				},
				Config: loadAcceptanceTestTerraform("03_no_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("display_name").HasValue("Acc Test Ring 03 No Assignments"),
					check.That(resourceType+".test").Key("deferral_in_days").HasValue("0"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
		},
	})
}

// TestAccResourceWindowsUpdateRing_04_InclusionAndExclusion creates a ring with
// both included and excluded group assignments.
func TestAccResourceWindowsUpdateRing_04_InclusionAndExclusion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create ring with inclusion and exclusion assignments")
				},
				Config: loadAcceptanceTestTerraform("04_inclusion_and_exclusion.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("display_name").HasValue("Acc Test Ring 04 Inclusion And Exclusion"),
					check.That(resourceType+".test").Key("is_paused").HasValue("false"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
		},
	})
}

// TestAccResourceWindowsUpdateRing_05_MinToMaxAssignments starts with 1 included
// group and scales up to 2 included groups.
func TestAccResourceWindowsUpdateRing_05_MinToMaxAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create ring with 1 included group (min)")
				},
				Config: loadAcceptanceTestTerraform("05_min_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("display_name").HasValue("Acc Test Ring 05 Min Assignments"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Update ring to 2 included groups (max)")
					testlog.WaitForConsistency(fmt.Sprintf("%s (min to max)", resourceType), 5*time.Second)
					time.Sleep(5 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("05_max_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("display_name").HasValue("Acc Test Ring 05 Max Assignments"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
		},
	})
}

// TestAccResourceWindowsUpdateRing_06_MaxToMinAssignments starts with 2 included
// groups and scales down to 1 included group.
func TestAccResourceWindowsUpdateRing_06_MaxToMinAssignments(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders:        externalProviders(),
		CheckDestroy:             destroy.CheckDestroyedAllFunc(testResource, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Create ring with 2 included groups (max)")
				},
				Config: loadAcceptanceTestTerraform("05_max_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("display_name").HasValue("Acc Test Ring 05 Max Assignments"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Update ring to 1 included group (min)")
					testlog.WaitForConsistency(fmt.Sprintf("%s (max to min)", resourceType), 5*time.Second)
					time.Sleep(5 * time.Second)
				},
				Config: loadAcceptanceTestTerraform("05_min_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("display_name").HasValue("Acc Test Ring 05 Min Assignments"),
					check.That(resourceType+".test").ExistsInGraph(testResource),
				),
			},
		},
	})
}

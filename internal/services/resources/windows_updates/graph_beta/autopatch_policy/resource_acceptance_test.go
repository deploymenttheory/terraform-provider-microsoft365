package graphBetaWindowsUpdatesAutopatchPolicy_test

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
	graphBetaWindowsAutopatchPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/graph_beta/autopatch_policy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	testResource = graphBetaWindowsAutopatchPolicy.WindowsAutopatchPolicyTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return config
}

// WAP001: Minimal policy acceptance test
func TestAccResourceWindowsAutopatchPolicy_01_WAP001_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating WAP001 minimal Windows Autopatch policy")
				},
				Config: loadAcceptanceTestTerraform("resource_wap001-minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopatch policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".wap001_minimal").ExistsInGraph(testResource),
					check.That(resourceType+".wap001_minimal").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".wap001_minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-wap001-minimal-[a-z0-9]{8}$`)),
					check.That(resourceType+".wap001_minimal").Key("description").HasValue("Acceptance test - minimal policy"),
					check.That(resourceType+".wap001_minimal").Key("created_date_time").IsNotEmpty(),
					check.That(resourceType+".wap001_minimal").Key("last_modified_date_time").IsNotEmpty(),
					check.That(resourceType+".wap001_minimal").Key("approval_rules.#").HasValue("0"),
				),
			},
			{
				ResourceName:      resourceType + ".wap001_minimal",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// WAP002: Policy with approval rules acceptance test
func TestAccResourceWindowsAutopatchPolicy_02_WAP002_ApprovalRules(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating WAP002 Windows Autopatch policy with approval rules")
				},
				Config: loadAcceptanceTestTerraform("resource_wap002-approval-rules.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopatch policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".wap002_approval_rules").ExistsInGraph(testResource),
					check.That(resourceType+".wap002_approval_rules").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".wap002_approval_rules").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-wap002-approval-rules-[a-z0-9]{8}$`)),
					check.That(resourceType+".wap002_approval_rules").Key("description").HasValue("Acceptance test - policy with approval rules"),
					check.That(resourceType+".wap002_approval_rules").Key("approval_rules.#").HasValue("3"),
				),
			},
			{
				ResourceName:      resourceType + ".wap002_approval_rules",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// WAP003: Lifecycle test - minimal to with-rules
func TestAccResourceWindowsAutopatchPolicy_03_WAP003_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating WAP003 lifecycle step 1: minimal policy")
				},
				Config: loadAcceptanceTestTerraform("resource_wap003-lifecycle-step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopatch policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".wap003_lifecycle").ExistsInGraph(testResource),
					check.That(resourceType+".wap003_lifecycle").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".wap003_lifecycle").Key("description").HasValue("Acceptance test - lifecycle step 1: no approval rules"),
					check.That(resourceType+".wap003_lifecycle").Key("approval_rules.#").HasValue("0"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating WAP003 lifecycle step 2: adding approval rules")
				},
				Config: loadAcceptanceTestTerraform("resource_wap003-lifecycle-step2.tf"),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("windows autopatch policy", 15*time.Second)
						time.Sleep(15 * time.Second)
						return nil
					},
					check.That(resourceType+".wap003_lifecycle").ExistsInGraph(testResource),
					check.That(resourceType+".wap003_lifecycle").Key("description").HasValue("Acceptance test - lifecycle step 2: with approval rules"),
					check.That(resourceType+".wap003_lifecycle").Key("approval_rules.#").HasValue("2"),
				),
			},
		},
	})
}

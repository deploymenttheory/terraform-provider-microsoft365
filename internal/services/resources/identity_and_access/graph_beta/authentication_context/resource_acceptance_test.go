package graphBetaAuthenticationContext_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaAuthenticationContext "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/authentication_context"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testResourceName = "microsoft365_graph_beta_identity_and_access_authentication_context"
)

var (
	// Resource type name constructed from exported constants
	resourceType = constants.PROVIDER_NAME + "_" + graphBetaAuthenticationContext.ResourceName

	// testResource is the test resource implementation for authentication contexts
	testResource = graphBetaAuthenticationContext.AuthenticationContextTestResource{}
)

func TestAccAuthenticationContextResource_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			15*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating basic authentication context")
				},
				Config: testAccAuthenticationContextConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("authentication context", 10*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(testResourceName+".test").ExistsInGraph(testResource),
					check.That(testResourceName+".test").Key("id").HasValue("c90"),
					check.That(testResourceName+".test").Key("display_name").HasValue("Acceptance Test Context"),
					check.That(testResourceName+".test").Key("description").HasValue("Context for acceptance testing"),
					check.That(testResourceName+".test").Key("is_available").HasValue("true"),
				),
			},
		},
	})
}

func TestAccAuthenticationContextResource_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			15*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating initial authentication context")
				},
				Config: testAccAuthenticationContextConfigUpdate1(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("authentication context", 10*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(testResourceName+".test").ExistsInGraph(testResource),
					check.That(testResourceName+".test").Key("display_name").HasValue("Initial Context"),
					check.That(testResourceName+".test").Key("is_available").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating authentication context")
				},
				Config: testAccAuthenticationContextConfigUpdate2(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("authentication context", 10*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(testResourceName+".test").ExistsInGraph(testResource),
					check.That(testResourceName+".test").Key("display_name").HasValue("Updated Context"),
					check.That(testResourceName+".test").Key("description").HasValue("Updated description"),
					check.That(testResourceName+".test").Key("is_available").HasValue("false"),
				),
			},
		},
	})
}

func TestAccAuthenticationContextResource_Import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			15*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating authentication context for import")
				},
				Config: testAccAuthenticationContextConfigImport(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("authentication context", 10*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing authentication context")
				},
				ResourceName:      testResourceName + ".test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAuthenticationContextResource_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			15*time.Second,
		),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: ">= 3.7.2",
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal authentication context")
				},
				Config: testAccAuthenticationContextConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("authentication context", 10*time.Second)
						time.Sleep(5 * time.Second)
						return nil
					},
					check.That(testResourceName+".test").ExistsInGraph(testResource),
					check.That(testResourceName+".test").Key("id").HasValue("c93"),
					check.That(testResourceName+".test").Key("display_name").HasValue("Minimal Context"),
					check.That(testResourceName+".test").Key("is_available").HasValue("false"),
				),
			},
		},
	})
}

// Configuration helper functions
func testAccAuthenticationContextConfigBasic() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/01_basic.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load authentication context basic config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccAuthenticationContextConfigUpdate1() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_update_1.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load authentication context update 1 config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccAuthenticationContextConfigUpdate2() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/02_update_2.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load authentication context update 2 config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccAuthenticationContextConfigImport() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/03_import.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load authentication context import config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccAuthenticationContextConfigMinimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/04_minimal.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load authentication context minimal config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

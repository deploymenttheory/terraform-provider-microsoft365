package graphBetaApplicationsTokenLifetimePolicy_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaTokenLifetimePolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/applications/graph_beta/token_lifetime_policy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	testResource = graphBetaTokenLifetimePolicy.TokenLifetimePolicyTestResource{}
)

func TestAccResourceTokenLifetimePolicy_01_Minimal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkTokenLifetimePoliciesDestroyed,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating minimal token lifetime policy")
				},
				Config: testAccConfigMinimal(),
				Check: resource.ComposeTestCheckFunc(
					checkTokenLifetimePolicyExists(resourceType+".minimal"),
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-min-[0-9a-fA-F-]+$`)),
					check.That(resourceType+".minimal").Key("definition.#").HasValue("1"),
					check.That(resourceType+".minimal").Key("is_organization_default").HasValue("false"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing minimal token lifetime policy")
				},
				ResourceName:      resourceType + ".minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func TestAccResourceTokenLifetimePolicy_02_Maximal(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkTokenLifetimePoliciesDestroyed,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source:            "hashicorp/random",
				VersionConstraint: constants.ExternalProviderRandomVersion,
			},
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Creating maximal token lifetime policy")
				},
				Config: testAccConfigMaximal(),
				Check: resource.ComposeTestCheckFunc(
					checkTokenLifetimePolicyExists(resourceType+".maximal"),
					check.That(resourceType+".maximal").Key("id").Exists(),
					check.That(resourceType+".maximal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-max-[0-9a-fA-F-]+$`)),
					check.That(resourceType+".maximal").Key("definition.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing maximal token lifetime policy")
				},
				ResourceName:      resourceType + ".maximal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

func testAccConfigMinimal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_minimal.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load token lifetime policy minimal acceptance config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

func testAccConfigMaximal() string {
	accTestConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_maximal.tf")
	if err != nil {
		panic(fmt.Sprintf("failed to load token lifetime policy maximal acceptance config: %s", err.Error()))
	}
	return acceptance.ConfiguredM365ProviderBlock(accTestConfig)
}

// A fresh policy can read successfully through one replica and still return
// 404 through the independent Graph client used by acceptance checks.
func checkTokenLifetimePolicyExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		stored, ok := state.RootModule().Resources[resourceName]
		if !ok || stored.Primary == nil {
			return fmt.Errorf("%s is missing from Terraform state", resourceName)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		return retry.RetryContext(ctx, time.Minute, func() *retry.RetryError {
			exists, err := testResource.Exists(ctx, nil, stored.Primary)
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if exists == nil || !*exists {
				return retry.RetryableError(fmt.Errorf("policy %s is not yet visible to the Graph check client", stored.Primary.ID))
			}
			return nil
		})
	}
}

func checkTokenLifetimePoliciesDestroyed(state *terraform.State) error {
	for _, stored := range state.RootModule().Resources {
		if stored.Type != resourceType || stored.Primary == nil {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		err := retry.RetryContext(ctx, time.Minute, func() *retry.RetryError {
			exists, err := testResource.Exists(ctx, nil, stored.Primary)
			if err != nil {
				return retry.NonRetryableError(err)
			}
			if exists != nil && *exists {
				return retry.RetryableError(fmt.Errorf("policy %s is still visible after deletion", stored.Primary.ID))
			}
			return nil
		})
		cancel()
		if err != nil {
			return err
		}
	}
	return nil
}

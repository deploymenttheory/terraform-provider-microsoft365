package graphConditionalAccessTermsOfUse_test

import (
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphConditionalAccessTermsOfUse "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_v1.0/conditional_access_terms_of_use"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	// Resource type name from the resource package
	resourceType = graphConditionalAccessTermsOfUse.ResourceName

	// testResource is the test resource implementation for conditional access terms of use
	testResource = graphConditionalAccessTermsOfUse.ConditionalAccessTermsOfUseTestResource{}
)

func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccConditionalAccessTermsOfUseResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			5*time.Second,
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
					testlog.StepAction(resourceType, "Creating")
				},
				Config: loadAcceptanceTestTerraform("resource_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".acc_test_conditional_access_terms_of_use_minimal").ExistsInGraph(testResource),
					check.That(resourceType+".acc_test_conditional_access_terms_of_use_minimal").Key("id").Exists(),
					check.That(resourceType+".acc_test_conditional_access_terms_of_use_minimal").Key("display_name").IsNotEmpty(),
					check.That(resourceType+".acc_test_conditional_access_terms_of_use_minimal").Key("is_viewing_before_acceptance_required").HasValue("true"),
					check.That(resourceType+".acc_test_conditional_access_terms_of_use_minimal").Key("is_per_device_acceptance_required").HasValue("false"),
					check.That(resourceType+".acc_test_conditional_access_terms_of_use_minimal").Key("user_reaccept_required_frequency").HasValue("P10D"),
					check.That(resourceType+".acc_test_conditional_access_terms_of_use_minimal").Key("file.localizations.#").HasValue("1"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName:      resourceType + ".acc_test_conditional_access_terms_of_use_minimal",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"file.localizations.0.file_data",
					"file.localizations.0.file_data.data",
					"file.localizations.0.file_data.%",
				},
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating")
				},
				Config: loadAcceptanceTestTerraform("resource_maximal.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".acc_test_conditional_access_terms_of_use_maximal").ExistsInGraph(testResource),
					check.That(resourceType+".acc_test_conditional_access_terms_of_use_maximal").Key("id").Exists(),
					check.That(resourceType+".acc_test_conditional_access_terms_of_use_maximal").Key("display_name").IsNotEmpty(),
					check.That(resourceType+".acc_test_conditional_access_terms_of_use_maximal").Key("file.localizations.#").HasValue("30"),
				),
			},
		},
	})
}

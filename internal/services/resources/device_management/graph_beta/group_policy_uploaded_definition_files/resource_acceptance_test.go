package graphBetaGroupPolicyUploadedDefinitionFiles_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Helper function to load test configs from acceptance directory
func loadAcceptanceTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/" + filename)
	if err != nil {
		panic("failed to load acceptance config " + filename + ": " + err.Error())
	}
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func TestAccGroupPolicyUploadedDefinitionFilesResource_Mozilla(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("resource_group_policy_uploaded_definition_files_mozilla.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".mozilla").ExistsInGraph(testResource),
					check.That(resourceType+".mozilla").Key("file_name").HasValue("mozilla.admx"),
					check.That(resourceType+".mozilla").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".mozilla").Key("default_language_code").HasValue("en-US"),
					check.That(resourceType+".mozilla").Key("group_policy_uploaded_language_files.#").HasValue("1"),
				),
			},
		},
	})
}

func TestAccGroupPolicyUploadedDefinitionFilesResource_Google(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("resource_group_policy_uploaded_definition_files_google.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".google").ExistsInGraph(testResource),
					check.That(resourceType+".google").Key("file_name").HasValue("google.admx"),
					check.That(resourceType+".google").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".google").Key("default_language_code").HasValue("en-US"),
					check.That(resourceType+".google").Key("group_policy_uploaded_language_files.#").HasValue("1"),
				),
			},
		},
	})
}

func TestAccGroupPolicyUploadedDefinitionFilesResource_Update(t *testing.T) {
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
				Config: loadAcceptanceTestTerraform("resource_group_policy_uploaded_definition_files_chrome.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_update").ExistsInGraph(testResource),
					check.That(resourceType+".test_update").Key("file_name").HasValue("chrome.admx"),
					check.That(resourceType+".test_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_update").Key("default_language_code").HasValue("en-US"),
					check.That(resourceType+".test_update").Key("group_policy_uploaded_language_files.#").HasValue("1"),
				),
			},
			{
				Config: loadAcceptanceTestTerraform("resource_group_policy_uploaded_definition_files_mozilla_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test_update").ExistsInGraph(testResource),
					check.That(resourceType+".test_update").Key("file_name").HasValue("mozilla.admx"),
					check.That(resourceType+".test_update").Key("id").MatchesRegex(regexp.MustCompile(`^[0-9a-fA-F-]+$`)),
					check.That(resourceType+".test_update").Key("default_language_code").HasValue("en-US"),
					check.That(resourceType+".test_update").Key("group_policy_uploaded_language_files.#").HasValue("1"),
				),
			},
		},
	})
}

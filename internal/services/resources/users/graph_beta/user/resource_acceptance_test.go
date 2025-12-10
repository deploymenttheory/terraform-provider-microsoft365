package graphBetaUsersUser_test

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/testlog"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaUsersUser "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/users/graph_beta/user"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	// Resource type name from the resource package
	resourceType = graphBetaUsersUser.ResourceName

	// testResource is the test resource implementation for users
	testResource = graphBetaUsersUser.UserTestResource{}
)

func TestAccUserResource_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { mocks.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(
			testResource,
			resourceType,
			30*time.Second, // Increased wait time for dependency user with manager relationship
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
				Config: testAccUserConfig_minimal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".minimal").ExistsInGraph(testResource),
					check.That(resourceType+".minimal").Key("id").Exists(),
					check.That(resourceType+".minimal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-user-minimal-[a-z0-9]{8}$`)),
					check.That(resourceType+".minimal").Key("user_principal_name").MatchesRegex(regexp.MustCompile(`^acc-test-user-minimal-[a-z0-9]{8}@deploymenttheory\.com$`)),
					check.That(resourceType+".minimal").Key("mail_nickname").MatchesRegex(regexp.MustCompile(`^acc-test-user-minimal-[a-z0-9]{8}$`)),
					check.That(resourceType+".minimal").Key("account_enabled").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing")
				},
				ResourceName: resourceType + ".minimal",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".minimal"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".minimal")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password_profile",
					"password_profile.%",
					"password_profile.password",
					"password_profile.force_change_password_next_sign_in",
					"password_profile.force_change_password_next_sign_in_with_mfa",
				},
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Updating to Maximal")
				},
				Config: testAccUserConfig_maximal(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".maximal").ExistsInGraph(testResource),
					check.That(resourceType+".maximal").Key("id").Exists(),
					check.That(resourceType+".maximal").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-user-maximal-[a-z0-9]{8}$`)),
					check.That(resourceType+".maximal").Key("user_principal_name").MatchesRegex(regexp.MustCompile(`^acc-test-user-maximal-[a-z0-9]{8}@deploymenttheory\.com$`)),
					check.That(resourceType+".maximal").Key("mail_nickname").MatchesRegex(regexp.MustCompile(`^acc-test-user-maximal-[a-z0-9]{8}$`)),
					check.That(resourceType+".maximal").Key("given_name").HasValue("Maximal"),
					check.That(resourceType+".maximal").Key("surname").HasValue("User"),
					check.That(resourceType+".maximal").Key("job_title").HasValue("Senior Developer"),
					check.That(resourceType+".maximal").Key("department").HasValue("Engineering"),
					check.That(resourceType+".maximal").Key("company_name").HasValue("Deployment Theory"),
					check.That(resourceType+".maximal").Key("employee_id").HasValue("1234567890"),
					check.That(resourceType+".maximal").Key("employee_type").HasValue("full time"),
					check.That(resourceType+".maximal").Key("employee_hire_date").HasValue("2025-11-21T00:00:00Z"),
					check.That(resourceType+".maximal").Key("office_location").HasValue("Building A"),
					check.That(resourceType+".maximal").Key("city").HasValue("Redmond"),
					check.That(resourceType+".maximal").Key("state").HasValue("WA"),
					check.That(resourceType+".maximal").Key("country").HasValue("US"),
					check.That(resourceType+".maximal").Key("street_address").HasValue("123 street"),
					check.That(resourceType+".maximal").Key("postal_code").HasValue("98052"),
					check.That(resourceType+".maximal").Key("usage_location").HasValue("US"),
					check.That(resourceType+".maximal").Key("mobile_phone").HasValue("+1 425-555-0101"),
					check.That(resourceType+".maximal").Key("fax_number").HasValue("+1 425-555-0102"),
					check.That(resourceType+".maximal").Key("preferred_language").HasValue("en-US"),
					check.That(resourceType+".maximal").Key("password_policies").HasValue("DisablePasswordExpiration"),
					check.That(resourceType+".maximal").Key("age_group").HasValue("NotAdult"),
					check.That(resourceType+".maximal").Key("consent_provided_for_minor").HasValue("Granted"),
					check.That(resourceType+".maximal").Key("manager_id").Exists(),
					check.That(resourceType+".maximal").Key("account_enabled").HasValue("true"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Maximal")
				},
				ResourceName: resourceType + ".maximal",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".maximal"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".maximal")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password_profile",
					"password_profile.%",
					"password_profile.password",
					"password_profile.force_change_password_next_sign_in",
					"password_profile.force_change_password_next_sign_in_with_mfa",
					"manager_id", // Manager relationship may not be readable in all scenarios
				},
			},
		},
	})
}

func TestAccUserResource_CustomSecurityAttributes(t *testing.T) {
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
					testlog.StepAction(resourceType, "Creating with Custom Security Attributes")
				},
				Config: testAccUserConfig_customSecAtt(),
				Check: resource.ComposeTestCheckFunc(
					func(_ *terraform.State) error {
						testlog.WaitForConsistency("group", 20*time.Second)
						time.Sleep(20 * time.Second)
						return nil
					},
					check.That(resourceType+".with_custom_security_attributes").ExistsInGraph(testResource),
					check.That(resourceType+".with_custom_security_attributes").Key("id").Exists(),
					check.That(resourceType+".with_custom_security_attributes").Key("display_name").MatchesRegex(regexp.MustCompile(`^acc-test-user-custom-sec-att-[a-z0-9]{8}$`)),
					check.That(resourceType+".with_custom_security_attributes").Key("user_principal_name").MatchesRegex(regexp.MustCompile(`^acc-test-user-custom-sec-att-[a-z0-9]{8}@deploymenttheory\.com$`)),
					check.That(resourceType+".with_custom_security_attributes").Key("mail_nickname").MatchesRegex(regexp.MustCompile(`^acc-test-user-custom-sec-att-[a-z0-9]{8}$`)),
					check.That(resourceType+".with_custom_security_attributes").Key("account_enabled").HasValue("true"),
					check.That(resourceType+".with_custom_security_attributes").Key("custom_security_attributes.#").HasValue("2"),
					check.That(resourceType+".with_custom_security_attributes").Key("custom_security_attributes.0.attribute_set").HasValue("Engineering"),
					check.That(resourceType+".with_custom_security_attributes").Key("custom_security_attributes.1.attribute_set").HasValue("Marketing"),
				),
			},
			{
				PreConfig: func() {
					testlog.StepAction(resourceType, "Importing Custom Security Attributes User")
				},
				ResourceName: resourceType + ".with_custom_security_attributes",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceType+".with_custom_security_attributes"]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceType+".with_custom_security_attributes")
					}
					hardDelete := rs.Primary.Attributes["hard_delete"]
					return fmt.Sprintf("%s:hard_delete=%s", rs.Primary.ID, hardDelete), nil
				},
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password_profile",
					"password_profile.%",
					"password_profile.password",
					"password_profile.force_change_password_next_sign_in",
					"password_profile.force_change_password_next_sign_in_with_mfa",
				},
			},
		},
	})
}

// Test configuration functions
func testAccUserConfig_minimal() string {
	config := mocks.LoadTerraformConfigFile("resource_minimal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccUserConfig_maximal() string {
	config := mocks.LoadTerraformConfigFile("resource_maximal.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

func testAccUserConfig_customSecAtt() string {
	config := mocks.LoadTerraformConfigFile("resource_custom_sec_att.tf")
	return acceptance.ConfiguredM365ProviderBlock(config)
}

package graphBetaMobileAppSupersedence_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	supersedence "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/mobile_app_supersedence"
	supersedenceMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/mobile_app_supersedence/mocks"
)

const resourceType = supersedence.ResourceName

var testResource = supersedence.MobileAppSupersedenceTestResource{}

func setupMockEnvironment(t *testing.T) *supersedenceMocks.MobileAppSupersedenceMock {
	t.Helper()
	t.Setenv("TF_ACC", "0")
	t.Setenv("MS365_TEST_MODE", "true")
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	mock := &supersedenceMocks.MobileAppSupersedenceMock{}
	mock.RegisterMocks()
	t.Cleanup(func() { mock.CleanupMockState(); httpmock.DeactivateAndReset() })
	return mock
}

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit config " + filename + ": " + err.Error())
	}
	return config
}

func TestUnitResourceMobileAppSupersedence_01_Update(t *testing.T) {
	setupMockEnvironment(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: loadUnitTestTerraform("01_update.tf"), Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").IsSet(),
				check.That(resourceType+".test").
					Key("source_id").
					HasValue("11111111-1111-1111-1111-111111111111"),
				check.That(resourceType+".test").
					Key("target_id").
					HasValue("22222222-2222-2222-2222-222222222222"),
				check.That(resourceType+".test").Key("supersedence_type").HasValue("update"),
				check.That(resourceType+".test").Key("target_type").HasValue("child"),
			)},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true},
		},
	})
}

func TestUnitResourceMobileAppSupersedence_02_ReplaceToUpdate(t *testing.T) {
	setupMockEnvironment(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_replace.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("supersedence_type").HasValue("replace"),
					check.That(resourceType+".test").Key("timeouts.create").HasValue("5m"),
				),
			},
			{
				Config: loadUnitTestTerraform("01_update.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: check.That(resourceType + ".test").
					Key("supersedence_type").
					HasValue("update"),
			},
			{Config: loadUnitTestTerraform("01_update.tf"), PlanOnly: true},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true},
		},
	})
}

func TestUnitResourceMobileAppSupersedence_03_SiblingRemoval(t *testing.T) {
	setupMockEnvironment(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_siblings_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("supersedence_type").HasValue("update"),
					check.That(resourceType+".sibling").
						Key("supersedence_type").
						HasValue("replace"),
				),
			},
			{
				Config: loadUnitTestTerraform("05_siblings_step2.tf"),
				Check: check.That(resourceType + ".sibling").
					Key("supersedence_type").
					HasValue("replace"),
			},
			{Config: loadUnitTestTerraform("05_siblings_step2.tf"), PlanOnly: true},
			{Config: loadUnitTestTerraform("05_siblings_step3.tf")},
		},
	})
}

func TestUnitResourceMobileAppSupersedence_04_Validation(t *testing.T) {
	for _, tc := range []struct{ filename, message string }{
		{"08_invalid_supersedence_type.tf", "Attribute supersedence_type value must be one of"},
		{"09_invalid_source_id.tf", "must be a valid app GUID"},
		{"10_invalid_target_id.tf", "must be a valid app GUID"},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			setupMockEnvironment(t)
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      loadUnitTestTerraform(tc.filename),
						PlanOnly:    true,
						ExpectError: regexp.MustCompile(tc.message),
					},
				},
			})
		})
	}
}

func TestUnitResourceMobileAppSupersedence_05_ErrorHandling(t *testing.T) {
	mock := setupMockEnvironment(t)
	mock.RegisterErrorMocks()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_update.tf"),
				ExpectError: regexp.MustCompile(`Invalid mobile app supersedence\s+relationship`),
			},
		},
	})
}

func TestUnitResourceMobileAppSupersedence_06_PreservesDependency(t *testing.T) {
	setupMockEnvironment(t)
	const dependencyType = "microsoft365_graph_beta_device_and_app_management_mobile_app_dependency"
	const dependencyConfig = `
resource "microsoft365_graph_beta_device_and_app_management_mobile_app_dependency" "mixed" {
  source_id = "11111111-1111-1111-1111-111111111111"
  target_id = "33333333-3333-3333-3333-333333333333"
  dependency_type = "detect"
}
`
	config := loadUnitTestTerraform("01_update.tf")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config + dependencyConfig},
			{
				Config: strings.ReplaceAll(config, `"update"`, `"replace"`) + dependencyConfig,
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("supersedence_type").HasValue("replace"),
					check.That(dependencyType+".mixed").Key("dependency_type").HasValue("detect"),
				),
			},
			{
				Config: dependencyConfig,
				Check: check.That(dependencyType + ".mixed").
					Key("dependency_type").
					HasValue("detect"),
			},
			{Config: dependencyConfig, PlanOnly: true},
		},
	})
}

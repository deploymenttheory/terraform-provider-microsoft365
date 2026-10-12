package graphBetaMobileAppDependency_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	dependency "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/mobile_app_dependency"
	dependencyMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/mobile_app_dependency/mocks"
)

const resourceType = dependency.ResourceName

var testResource = dependency.MobileAppDependencyTestResource{}

func setupMockEnvironment(t *testing.T) *dependencyMocks.MobileAppDependencyMock {
	t.Helper()
	t.Setenv("TF_ACC", "0")
	t.Setenv("MS365_TEST_MODE", "true")
	httpmock.Activate()
	mocks.NewMocks().AuthMocks.RegisterMocks()
	mock := &dependencyMocks.MobileAppDependencyMock{}
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

func TestUnitResourceMobileAppDependency_01_Detect(t *testing.T) {
	setupMockEnvironment(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: loadUnitTestTerraform("01_detect.tf"), Check: resource.ComposeTestCheckFunc(
				check.That(resourceType+".test").Key("id").IsSet(),
				check.That(resourceType+".test").
					Key("source_id").
					HasValue("11111111-1111-1111-1111-111111111111"),
				check.That(resourceType+".test").
					Key("target_id").
					HasValue("22222222-2222-2222-2222-222222222222"),
				check.That(resourceType+".test").Key("dependency_type").HasValue("detect"),
				check.That(resourceType+".test").Key("target_type").HasValue("child"),
			)},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true},
		},
	})
}

func TestUnitResourceMobileAppDependency_02_AutoInstallToDetect(t *testing.T) {
	setupMockEnvironment(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_auto_install.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("dependency_type").HasValue("autoInstall"),
					check.That(resourceType+".test").Key("timeouts.create").HasValue("5m"),
				),
			},
			{
				Config: loadUnitTestTerraform("01_detect.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionUpdate,
						),
					},
				},
				Check: check.That(resourceType + ".test").
					Key("dependency_type").
					HasValue("detect"),
			},
			{Config: loadUnitTestTerraform("01_detect.tf"), PlanOnly: true},
			{ResourceName: resourceType + ".test", ImportState: true, ImportStateVerify: true},
		},
	})
}

func TestUnitResourceMobileAppDependency_03_SiblingRemoval(t *testing.T) {
	setupMockEnvironment(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_siblings_step1.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("dependency_type").HasValue("detect"),
					check.That(resourceType+".sibling").
						Key("dependency_type").
						HasValue("autoInstall"),
				),
			},
			{
				Config: loadUnitTestTerraform("05_siblings_step2.tf"),
				Check: check.That(resourceType + ".sibling").
					Key("dependency_type").
					HasValue("autoInstall"),
			},
			{Config: loadUnitTestTerraform("05_siblings_step2.tf"), PlanOnly: true},
			{Config: loadUnitTestTerraform("05_siblings_step3.tf")},
		},
	})
}

func TestUnitResourceMobileAppDependency_04_Validation(t *testing.T) {
	for _, tc := range []struct{ filename, message string }{
		{"08_invalid_dependency_type.tf", "Attribute dependency_type value must be one of"},
		{"09_invalid_source_id.tf", "must be a valid app GUID"},
		{"10_invalid_target_id.tf", "must be a valid app GUID"},
		{"11_self_dependency.tf", "An app cannot depend on itself"},
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

func TestUnitResourceMobileAppDependency_05_ErrorHandling(t *testing.T) {
	mock := setupMockEnvironment(t)
	mock.RegisterErrorMocks()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_detect.tf"),
				ExpectError: regexp.MustCompile("Invalid mobile app dependency relationship"),
			},
		},
	})
}

package graphBetaDeviceAndAppManagementIOSManagedMobileApp_test

import (
	"fmt"
	"regexp"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	graphBetaIOSManagedMobileApp "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/ios_managed_mobile_app"
	appMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/ios_managed_mobile_app/mocks"
)

var resourceType = graphBetaIOSManagedMobileApp.ResourceName

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment(t *testing.T, errorMocks bool) {
	t.Helper()
	mocks.SetupUnitTestEnvironment(t)
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	appMock := &appMocks.IOSManagedMobileAppMock{}
	if errorMocks {
		appMock.RegisterErrorMocks()
	} else {
		appMock.RegisterMocks()
	}
	t.Cleanup(appMock.CleanupMockState)
}

// checkPolicyApps asserts the mock policy targets exactly the given bundles.
func checkPolicyApps(bundles ...string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		want := make([]string, 0, len(bundles))
		for _, pkg := range bundles {
			want = append(want, appMocks.AppID(pkg))
		}
		slices.Sort(want)
		got := appMocks.PolicyAppIDs(appMocks.PolicyID)
		if !slices.Equal(got, want) {
			return fmt.Errorf("policy targets %v, want %v", got, want)
		}
		return nil
	}
}

func TestUnitResourceIOSManagedMobileApp_01_Lifecycle_Minimal(t *testing.T) {
	setupMockEnvironment(t, false)
	r := check.That(resourceType + ".test")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyApps(),
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_minimal.tf"),
				Check: resource.ComposeTestCheckFunc(
					r.Key("id").HasValue("com.example.testapp.ios"),
					r.Key("managed_app_protection_id").HasValue(appMocks.PolicyID),
					r.Key("mobile_app_identifier.bundle_id").HasValue("com.example.testapp"),
					r.Key("version").MatchesRegex(regexp.MustCompile(`^-?\d+$`)),
					checkPolicyApps("com.example.testapp"),
				),
			},
			{
				ResourceName:      resourceType + ".test",
				ImportState:       true,
				ImportStateId:     appMocks.PolicyID + "/com.example.testapp.ios",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"timeouts",
				},
			},
		},
	})
}

// Two resources on the same policy apply concurrently; the per-policy lock must
// stop either full-list targetApps write from dropping the other app.
func TestUnitResourceIOSManagedMobileApp_02_TwoAppsSamePolicy(t *testing.T) {
	setupMockEnvironment(t, false)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyApps(),
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_two_apps_same_policy.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".first").
						Key("id").
						HasValue("com.example.first.ios"),
					check.That(resourceType+".second").
						Key("id").
						HasValue("com.example.second.ios"),
					checkPolicyApps("com.example.first", "com.example.second"),
				),
			},
			{
				// Removing one resource must leave the other app targeted.
				Config: loadUnitTestTerraform("03_two_apps_remove_one.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".first").
						Key("id").
						HasValue("com.example.first.ios"),
					checkPolicyApps("com.example.first"),
				),
			},
		},
	})
}

func TestUnitResourceIOSManagedMobileApp_03_BundleChangeReplaces(t *testing.T) {
	setupMockEnvironment(t, false)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		CheckDestroy:             checkPolicyApps(),
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_bundle_change_step1.tf"),
				Check:  checkPolicyApps("com.example.testapp"),
			},
			{
				Config: loadUnitTestTerraform("04_bundle_change_step2.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							resourceType+".test",
							plancheck.ResourceActionDestroyBeforeCreate,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").
						Key("id").
						HasValue("com.example.renamed.ios"),
					checkPolicyApps("com.example.renamed"),
				),
			},
		},
	})
}

// Apps already targeted outside Terraform are left alone and must be imported.
func TestUnitResourceIOSManagedMobileApp_04_AlreadyTargeted(t *testing.T) {
	setupMockEnvironment(t, false)
	appMocks.SeedApps(appMocks.PolicyID, "com.example.existing")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("05_already_targeted.tf"),
				ExpectError: regexp.MustCompile(
					`already targeted by this policy[\s\S]*Import it with ID`,
				),
			},
		},
	})
	if err := checkPolicyApps("com.example.existing")(nil); err != nil {
		t.Fatal(err)
	}
}

func TestUnitResourceIOSManagedMobileApp_05_VersionNotConfigurable(t *testing.T) {
	setupMockEnvironment(t, false)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("06_version_not_configurable.tf"),
				ExpectError: regexp.MustCompile(`Invalid Configuration for Read-Only Attribute`),
			},
		},
	})
}

func TestUnitResourceIOSManagedMobileApp_06_MissingPolicy(t *testing.T) {
	setupMockEnvironment(t, false)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("07_missing_policy.tf"),
				ExpectError: regexp.MustCompile(`(?i)not found`),
			},
		},
	})
}

func TestUnitResourceIOSManagedMobileApp_07_GraphError(t *testing.T) {
	setupMockEnvironment(t, true)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_minimal.tf"),
				ExpectError: regexp.MustCompile(`(?i)forbidden|insufficient privileges`),
			},
		},
	})
}

func TestUnitResourceIOSManagedMobileApp_08_InvalidImportID(t *testing.T) {
	setupMockEnvironment(t, false)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        loadUnitTestTerraform("01_minimal.tf"),
				ResourceName:  resourceType + ".test",
				ImportState:   true,
				ImportStateId: "com.example.testapp.ios",
				ExpectError: regexp.MustCompile(
					`Expected import ID in format 'managedAppProtectionId/appId'`,
				),
			},
		},
	})
}

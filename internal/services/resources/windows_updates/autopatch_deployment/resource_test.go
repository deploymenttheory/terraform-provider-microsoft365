package graphBetaWindowsUpdatesAutopatchDeployment_test

import (
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/check"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	WindowsUpdatesAutopatchDeploymentResource "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_deployment"
	deploymentMocks "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/windows_updates/autopatch_deployment/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

var (
	resourceType = WindowsUpdatesAutopatchDeploymentResource.ResourceName
)

func loadUnitTestTerraform(filename string) string {
	config, err := helpers.ParseHCLFile("tests/terraform/unit/" + filename)
	if err != nil {
		panic("failed to load unit test config " + filename + ": " + err.Error())
	}
	return config
}

func setupMockEnvironment() (*mocks.Mocks, *deploymentMocks.WindowsUpdateDeploymentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	deploymentMock := &deploymentMocks.WindowsUpdateDeploymentMock{}
	deploymentMock.RegisterMocks()
	return mockClient, deploymentMock
}

func setupErrorMockEnvironment() (*mocks.Mocks, *deploymentMocks.WindowsUpdateDeploymentMock) {
	httpmock.Activate()
	mockClient := mocks.NewMocks()
	mockClient.AuthMocks.RegisterMocks()
	deploymentMock := &deploymentMocks.WindowsUpdateDeploymentMock{}
	deploymentMock.RegisterErrorMocks()
	return mockClient, deploymentMock
}

func TestUnitResourceWindowsUpdateDeployment_01_FeatureUpdateDeployment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deploymentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deploymentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("01_feature_update_deployment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("content.catalog_entry_id").HasValue("f341705b-0b15-4ce3-aaf2-6a1681d78606"),
					check.That(resourceType+".test").Key("content.catalog_entry_type").HasValue("featureUpdate"),
					check.That(resourceType+".test").Key("settings.schedule.start_date_time").HasValue("2024-01-15T10:00:00Z"),
					check.That(resourceType+".test").Key("settings.schedule.gradual_rollout.duration_between_offers").HasValue("P1W"),
					check.That(resourceType+".test").Key("settings.schedule.gradual_rollout.devices_per_offer").HasValue("100"),
					check.That(resourceType+".test").Key("settings.monitoring.monitoring_rules.#").HasValue("1"),
					check.That(resourceType+".test").Key("created_date_time").Exists(),
					check.That(resourceType+".test").Key("last_modified_date_time").Exists(),
				),
			},
			{
				ResourceName:            resourceType + ".test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"timeouts"},
			},
		},
	})
}

func TestUnitResourceWindowsUpdateDeployment_02_QualityUpdateDeployment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deploymentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deploymentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("02_quality_update_deployment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("content.catalog_entry_type").HasValue("qualityUpdate"),
					check.That(resourceType+".test").Key("settings.schedule.gradual_rollout.end_date_time").HasValue("2024-02-01T10:00:00Z"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdateDeployment_03_MinimalDeployment(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deploymentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deploymentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("03_minimal_deployment.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("content.catalog_entry_id").HasValue("minimal-catalog-entry-id"),
					check.That(resourceType+".test").Key("content.catalog_entry_type").HasValue("featureUpdate"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdateDeployment_04_IneligibleOfferFallback(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deploymentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deploymentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: loadUnitTestTerraform("04_ineligible_offer_fallback.tf"),
				Check: resource.ComposeTestCheckFunc(
					check.That(resourceType+".test").Key("id").Exists(),
					check.That(resourceType+".test").Key("content.catalog_entry_type").HasValue("featureUpdate"),
					check.That(resourceType+".test").Key("settings.monitoring.monitoring_rules.#").HasValue("1"),
				),
			},
		},
	})
}

func TestUnitResourceWindowsUpdateDeployment_06_ValidationIneligibleWrongAction(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deploymentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deploymentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("05_invalid_ineligible_wrong_action.tf"),
				ExpectError: regexp.MustCompile(`ineligible.*offerFallback|offerFallback.*ineligible|Invalid field combination`),
			},
		},
	})
}

func TestUnitResourceWindowsUpdateDeployment_07_ValidationRollbackOfferFallback(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deploymentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deploymentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("06_invalid_rollback_offer_fallback.tf"),
				ExpectError: regexp.MustCompile(`offerFallback.*ineligible|ineligible.*offerFallback|Invalid field combination`),
			},
		},
	})
}

func TestUnitResourceWindowsUpdateDeployment_08_ValidationOfferFallbackWithThreshold(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deploymentMock := setupMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deploymentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("07_invalid_offer_fallback_with_threshold.tf"),
				ExpectError: regexp.MustCompile(`threshold.*offerFallback|offerFallback.*threshold|Invalid field combination`),
			},
		},
	})
}

func TestUnitResourceWindowsUpdateDeployment_05_Error(t *testing.T) {
	mocks.SetupUnitTestEnvironment(t)
	_, deploymentMock := setupErrorMockEnvironment()
	defer httpmock.DeactivateAndReset()
	defer deploymentMock.CleanupMockState()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      loadUnitTestTerraform("01_feature_update_deployment.tf"),
				ExpectError: regexp.MustCompile("BadRequest|400|Invalid"),
			},
		},
	})
}

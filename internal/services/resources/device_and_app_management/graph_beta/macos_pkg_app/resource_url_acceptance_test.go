package graphBetaMacOSPKGApp_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	grapherrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	pkgmetadata "github.com/deploymenttheory/terraform-provider-microsoft365/internal/utilities/device_and_app_management/installers/macos_pkg"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func TestAccResourceMacOSPkgApp_04_ImportEmptyPrivacyURL(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("TF_ACC must be set to 1")
	}
	mocks.TestAccPreCheck(t)
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()
	client, err := acceptance.TestGraphClient() //nolint:contextcheck // The shared acceptance client factory creates its own context.
	if err != nil {
		t.Fatal(err)
	}
	packagePath := filepath.Join(t.TempDir(), "Nudge.pkg")
	downloadNudgePackage(t, packagePath)
	metadata, err := pkgmetadata.ExtractFieldsFromPKGFile(ctx, packagePath, "Info.plist", []pkgmetadata.Field{
		{Key: "CFBundleIdentifier", Required: true},
		{Key: "CFBundleShortVersionString", Required: true},
	})
	if err != nil {
		t.Fatal(err)
	}

	// Import only needs existing app metadata. Create the empty Graph value
	// directly, since Terraform correctly rejects an explicitly empty URL.
	// The other acceptance cases exercise real package upload and publication.
	displayName := "acc-test-macos-pkg-urls-" + uuid.NewString()[:8]
	description, publisher, empty := "Optional URL import regression", "MacAdmins", ""
	fileName := filepath.Base(packagePath)
	bundleID, bundleVersion := metadata[0].Values["CFBundleIdentifier"], metadata[0].Values["CFBundleShortVersionString"]
	body := graphmodels.NewMacOSPkgApp()
	body.SetDisplayName(&displayName)
	body.SetDescription(&description)
	body.SetPublisher(&publisher)
	body.SetPrivacyInformationUrl(&empty)
	body.SetFileName(&fileName)
	body.SetPrimaryBundleId(&bundleID)
	body.SetPrimaryBundleVersion(&bundleVersion)
	includedApp := graphmodels.NewMacOSIncludedApp()
	includedApp.SetBundleId(&bundleID)
	includedApp.SetBundleVersion(&bundleVersion)
	body.SetIncludedApps([]graphmodels.MacOSIncludedAppable{includedApp})
	minimumOS := graphmodels.NewMacOSMinimumOperatingSystem()
	enabled := true
	minimumOS.SetV120(&enabled)
	body.SetMinimumSupportedOperatingSystem(minimumOS)
	app, err := client.DeviceAppManagement().MobileApps().Post(ctx, body, nil)
	if err != nil {
		t.Fatal(err)
	}
	appID := *app.GetId()
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), time.Minute) //nolint:contextcheck // t.Context is cancelled before cleanup.
		defer cleanupCancel()
		current, readErr := client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).Get(cleanupCtx, nil)
		if readErr != nil && grapherrors.GraphError(cleanupCtx, readErr).StatusCode == 404 {
			return
		}
		if readErr != nil {
			t.Errorf("read owned import fixture for cleanup: %v", readErr)
			return
		}
		if current.GetDisplayName() == nil || *current.GetDisplayName() != displayName {
			t.Error("import fixture ownership check failed")
			return
		}
		if deleteErr := client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).Delete(cleanupCtx, nil); deleteErr != nil {
			t.Errorf("delete owned import fixture: %v", deleteErr)
			return
		}
		_, readErr = client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).Get(cleanupCtx, nil)
		if readErr == nil || grapherrors.GraphError(cleanupCtx, readErr).StatusCode != 404 {
			t.Error("owned import fixture deletion was not verified")
		}
	})
	for range 3 {
		current, readErr := client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).Get(ctx, nil)
		if readErr != nil {
			t.Fatal(readErr)
		}
		if current.GetPrivacyInformationUrl() == nil || *current.GetPrivacyInformationUrl() != "" {
			t.Fatal("Graph did not preserve an empty privacy URL on the import fixture")
		}
	}

	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_url_import.tf")
	if err != nil {
		t.Fatal(err)
	}
	config = strings.ReplaceAll(config, "APP_NAME", displayName)
	importConfig, err := helpers.ParseHCLFile("tests/terraform/acceptance/import_empty_privacy_url.tf")
	if err != nil {
		t.Fatal(err)
	}
	importConfig = strings.ReplaceAll(importConfig, "IMPORT_ID", appID)
	name := resourceType + ".test"
	resource.Test(t, resource.TestCase{
		TerraformVersionChecks:   []tfversion.TerraformVersionCheck{tfversion.SkipBelow(tfversion.Version1_5_0)},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(macOSPKGAppTestResource{}, resourceType, 30*time.Second),
		Steps: []resource.TestStep{
			{
				ResourceName: name, ImportStateId: appID,
				ImportState: true, GenerateConfig: true,
				ImportStateKind:        resource.ImportBlockWithID,
				ImportStateConfigExact: true,
				Config:                 importConfig,
				ImportPlanChecks: resource.ImportPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectKnownValue(name, tfjsonpath.New("privacy_information_url"), knownvalue.Null()),
				}},
			},
			{
				ResourceName: name, ImportStateId: appID,
				ImportState: true, ImportStatePersist: true,
				Config: config,
			},
			{Config: config, PlanOnly: true, ExpectNonEmptyPlan: false},
			{
				Config:   strings.Replace(config, "# PRIVACY_URL", "privacy_information_url = null", 1),
				PlanOnly: true, ExpectNonEmptyPlan: false,
			},
		},
	})
}

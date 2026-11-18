package graphBetaMacOSPKGApp_test

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const nudgePackageURL = "https://github.com/macadmins/nudge/releases/download/v2.1.4.81863/Nudge-2.1.4.81863.pkg"

func downloadNudgePackage(t *testing.T, destination string) {
	t.Helper()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, nudgePackageURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	client := &http.Client{Timeout: 2 * time.Minute}
	response, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("Nudge package download: HTTP %d", response.StatusCode)
	}
	content, err := io.ReadAll(io.LimitReader(response.Body, 2*1024*1024))
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%x", sha256.Sum256(content)) != "361ef3210b109100dbe4993be0acf7d54090645bb877c02f3b1cf105df3e1c2c" {
		t.Fatal("Nudge package SHA-256 does not match the official release")
	}
	if err := os.WriteFile(destination, content, 0600); err != nil {
		t.Fatal(err)
	}
}

func TestAccResourceMacOSPkgApp_02_MetadataUpdateWithoutLocalInstaller(t *testing.T) {
	packagePath := filepath.Join(t.TempDir(), "Nudge.pkg")
	testMetadataUpdateWithoutInstaller(t, packagePath, func() {
		downloadNudgePackage(t, packagePath)
	}, func() {
		if err := os.Remove(packagePath); err != nil {
			t.Fatal(err)
		}
	})
}

func TestAccResourceMacOSPkgApp_03_MetadataUpdateWithoutInstallerURL(t *testing.T) {
	packagePath := filepath.Join(t.TempDir(), "Nudge.pkg")
	var unavailable atomic.Bool
	var requests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		if unavailable.Load() {
			w.WriteHeader(http.StatusGone)
			return
		}
		http.ServeFile(w, r, packagePath)
	}))
	defer server.Close()
	var createRequests int64
	testMetadataUpdateWithoutInstaller(t, server.URL+"/Nudge.pkg", func() {
		downloadNudgePackage(t, packagePath)
	}, func() {
		createRequests = requests.Load()
		if createRequests == 0 {
			t.Fatal("Create did not download the installer")
		}
		unavailable.Store(true)
	}, func(_ *terraform.State) error {
		if requests.Load() != createRequests {
			return fmt.Errorf("metadata update downloaded the unchanged installer: requests before %d, after %d", createRequests, requests.Load())
		}
		return nil
	})
}

func testMetadataUpdateWithoutInstaller(t *testing.T, source string, prepare, removeSource func(), extraChecks ...resource.TestCheckFunc) {
	t.Helper()
	config, err := helpers.ParseHCLFile("tests/terraform/acceptance/resource_metadata_update.tf")
	if err != nil {
		t.Fatal(err)
	}
	sourceField := "installer_file_path_source"
	if strings.HasPrefix(source, "http") {
		sourceField = "installer_url_source"
	}
	config = strings.Replace(config, "installer_file_path_source = \"INSTALLER_SOURCE\"", sourceField+" = "+fmt.Sprintf("%q", filepath.ToSlash(source)), 1)
	updated := strings.Replace(config, "ignore_version_detection = false", "ignore_version_detection = true", 1)
	updated = strings.Replace(updated, "Original app metadata", "Updated app metadata", 1)
	name := resourceType + ".test"
	var appID, contentVersion string
	rememberContent := func(s *terraform.State) error {
		state := s.RootModule().Resources[name].Primary
		appID = state.ID
		contentVersion = state.Attributes["content_version.0.id"]
		if appID == "" || contentVersion == "" {
			return fmt.Errorf("created app has no app ID or committed content version")
		}
		return nil
	}
	checkContent := func(s *terraform.State) error {
		state := s.RootModule().Resources[name].Primary
		if state.ID != appID || state.Attributes["content_version.0.id"] != contentVersion {
			return fmt.Errorf("metadata update replaced the app or changed its committed content version")
		}
		return nil
	}
	checks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(name, "description", "Updated app metadata"),
		resource.TestCheckResourceAttr(name, "macos_pkg_app.ignore_version_detection", "true"),
		resource.TestCheckResourceAttr(name, "publishing_state", "published"),
		resource.TestCheckResourceAttr(name, "macos_pkg_app.primary_bundle_id", "com.github.macadmins.Nudge"),
		resource.TestCheckResourceAttr(name, "macos_pkg_app.primary_bundle_version", "2.1.4.81863"),
		checkContent,
	}
	checks = append(checks, extraChecks...)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			mocks.TestAccPreCheck(t)
			prepare()
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             destroy.CheckDestroyedAllFunc(macOSPKGAppTestResource{}, resourceType, 30*time.Second),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {Source: "hashicorp/random", VersionConstraint: constants.ExternalProviderRandomVersion},
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "publishing_state", "published"),
					rememberContent,
				),
			},
			{PreConfig: removeSource, Config: updated, Check: resource.ComposeTestCheckFunc(checks...)},
			{Config: updated, PlanOnly: true, ExpectNonEmptyPlan: false},
		},
	})
}

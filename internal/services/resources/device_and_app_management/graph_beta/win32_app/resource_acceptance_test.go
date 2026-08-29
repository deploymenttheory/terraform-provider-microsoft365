package graphBetaWin32App_test

import (
	"archive/zip"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/destroy"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/exists"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/mocks"
	win32 "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_and_app_management/graph_beta/win32_app"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

type win32TestResource struct{}

func (win32TestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	return exists.CheckResourceExists(ctx, state, func(client *msgraphbetasdk.GraphServiceClient, ctx context.Context, state *terraform.InstanceState) error {
		_, err := client.DeviceAppManagement().MobileApps().ByMobileAppId(state.ID).Get(ctx, nil)
		return err
	})
}

func TestAccResourceWin32App_ZipContentLifecycle(t *testing.T)       { testContentLifecycle(t, true) }
func TestAccResourceWin32App_IntuneWinContentLifecycle(t *testing.T) { testContentLifecycle(t, false) }

// Both tests create an app, publish a second version, switch formats, and
// publish another version. Each step is followed by an explicit empty plan.
func testContentLifecycle(t *testing.T, startWithZip bool) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("set TF_ACC=1 and the test-tenant/fixture variables described in tests/README.md")
	}
	mocks.TestAccPreCheck(t)
	groupID := os.Getenv("WIN32_APP_TEST_GROUP_ID")
	require.NotEmpty(t, groupID, "WIN32_APP_TEST_GROUP_ID must identify an empty, dedicated test group")
	packaged := []string{os.Getenv("WIN32_APP_INTUNEWIN_V1"), os.Getenv("WIN32_APP_INTUNEWIN_V2")}
	for i, p := range packaged {
		require.NotEmpty(t, p, "WIN32_APP_INTUNEWIN_V%d is required", i+1)
		absolute, err := filepath.Abs(p)
		require.NoError(t, err)
		packaged[i] = absolute
		_, err = os.Stat(absolute)
		require.NoError(t, err)
	}
	require.NotEqual(t, packaged[0], packaged[1], "use distinct artifact paths for each version")
	plain := []string{acceptanceZip(t, "v1"), acceptanceZip(t, "v2")}
	template, err := os.ReadFile("tests/terraform/acceptance/content_lifecycle.tf")
	require.NoError(t, err)
	name := fmt.Sprintf("acc-test-win32-content-%d", time.Now().UnixNano())
	appAddress := win32.ResourceName + ".test"
	assignmentAddress := "microsoft365_graph_beta_device_and_app_management_mobile_app_assignment.test"
	var appID, assignmentID, lastVersion string
	check := func(state *terraform.State) error {
		app, ok := state.RootModule().Resources[appAddress]
		if !ok {
			return fmt.Errorf("app missing from state")
		}
		assignment, ok := state.RootModule().Resources[assignmentAddress]
		if !ok {
			return fmt.Errorf("assignment missing from state")
		}
		if appID == "" {
			appID = app.Primary.ID
			assignmentID = assignment.Primary.ID
		} else if appID != app.Primary.ID || assignmentID != assignment.Primary.ID {
			return fmt.Errorf("content update replaced application or assignment")
		}
		version := app.Primary.Attributes["committed_content_version"]
		if version == "" || version == lastVersion {
			return fmt.Errorf("expected a new committed content version, got %q after %q", version, lastVersion)
		}
		if app.Primary.Attributes["content_version.#"] != "1" || app.Primary.Attributes["content_version.0.id"] != version {
			return fmt.Errorf("state must contain only the current committed content version")
		}
		client, err := acceptance.TestGraphClient()
		if err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		assignments, err := client.DeviceAppManagement().MobileApps().ByMobileAppId(appID).Assignments().Get(ctx, nil)
		if err != nil {
			return err
		}
		found := false
		for _, a := range assignments.GetValue() {
			if a.GetId() != nil && *a.GetId() == assignmentID {
				target, ok := a.GetTarget().(graphmodels.GroupAssignmentTargetable)
				found = ok && target.GetGroupId() != nil && *target.GetGroupId() == groupID
			}
		}
		if !found {
			return fmt.Errorf("original group assignment is missing from Graph")
		}
		lastVersion = version
		return nil
	}
	var steps []resource.TestStep
	for _, zipMode := range []bool{startWithZip, !startWithZip} {
		block, paths := "app_installer", packaged
		if zipMode {
			block, paths = "app_installer_zip", plain
		}
		for _, source := range paths {
			config := strings.NewReplacer("{{NAME}}", strconv.Quote(name), "{{GROUP_ID}}", strconv.Quote(groupID), "{{SOURCE_BLOCK}}", block, "{{SOURCE_PATH}}", strconv.Quote(source)).Replace(string(template))
			config = acceptance.ConfiguredM365ProviderBlock(config)
			steps = append(steps, resource.TestStep{Config: config, Check: check}, resource.TestStep{Config: config, PlanOnly: true, ExpectNonEmptyPlan: false})
		}
	}
	steps = append(steps, resource.TestStep{ResourceName: appAddress, ImportState: true, ImportStateVerify: true, ImportStateVerifyIgnore: []string{"app_installer", "app_installer_zip", "file_name"}})
	resource.Test(t, resource.TestCase{
		PreCheck: func() { mocks.TestAccPreCheck(t) }, ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		CheckDestroy: destroy.CheckDestroyedAllFunc(win32TestResource{}, win32.ResourceName, 30*time.Second), Steps: steps,
	})
}

func acceptanceZip(t *testing.T, version string) string {
	t.Helper()
	script, err := os.ReadFile(filepath.Join("tests", "fixtures", version, "setup.cmd"))
	require.NoError(t, err)
	target := filepath.Join(t.TempDir(), version+".zip")
	file, err := os.Create(target)
	require.NoError(t, err)
	archive := zip.NewWriter(file)
	entry, err := archive.Create("setup.cmd")
	require.NoError(t, err)
	_, err = entry.Write(script)
	require.NoError(t, err)
	require.NoError(t, archive.Close())
	require.NoError(t, file.Close())
	return target
}

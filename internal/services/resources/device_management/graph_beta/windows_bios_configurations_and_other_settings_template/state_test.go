package graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate

import (
	"context"
	"encoding/base64"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func ptr[T any](v T) *T { return &v }

func TestUnitMapRemoteResourceStateToTerraformIgnoresNilRemoteResource(t *testing.T) {
	t.Parallel()

	data := &WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel{
		DisplayName: types.StringValue("unchanged"),
	}

	MapRemoteResourceStateToTerraform(context.Background(), data, nil)

	if got := data.DisplayName.ValueString(); got != "unchanged" {
		t.Errorf("display_name = %q, want it left untouched", got)
	}
}

func TestUnitMapRemoteResourceStateToTerraformReEncodesConfigurationFileContent(t *testing.T) {
	t.Parallel()

	raw := "[cctk]\nSecureBoot=Enabled\n"
	created := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	format := graphmodels.DELL_HARDWARECONFIGURATIONFORMAT

	remote := graphmodels.NewHardwareConfiguration()
	remote.SetId(ptr("11111111-1111-1111-1111-111111111111"))
	remote.SetDisplayName(ptr("bios-template"))
	remote.SetDescription(ptr("a description"))
	remote.SetFileName(ptr("bios.cctk"))
	// Kiota hands back the decoded bytes, so the mapper must re-encode them.
	remote.SetConfigurationFileContent([]byte(raw))
	remote.SetHardwareConfigurationFormat(&format)
	remote.SetPerDevicePasswordDisabled(ptr(true))
	remote.SetRoleScopeTagIds([]string{"0", "1"})
	remote.SetVersion(ptr(int32(3)))
	remote.SetCreatedDateTime(&created)
	remote.SetLastModifiedDateTime(&created)

	data := &WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel{}
	MapRemoteResourceStateToTerraform(context.Background(), data, remote)

	want := base64.StdEncoding.EncodeToString([]byte(raw))
	if got := data.ConfigurationFileContent.ValueString(); got != want {
		t.Errorf("configuration_file_content = %q, want %q", got, want)
	}
	if got := data.HardwareConfigurationFormat.ValueString(); got != "dell" {
		t.Errorf("hardware_configuration_format = %q, want %q", got, "dell")
	}
	if got := data.Version.ValueInt32(); got != 3 {
		t.Errorf("version = %d, want 3", got)
	}
	if data.CreatedDateTime.IsNull() || data.LastModifiedDateTime.IsNull() {
		t.Error("created_date_time / last_modified_date_time were not mapped")
	}
	if !data.Assignments.IsNull() {
		t.Error("assignments should be null when the API returns none")
	}
}

// The API omits configurationFileContent on some responses; the configured value must survive.
func TestUnitMapRemoteResourceStateToTerraformRetainsConfigurationFileContentWhenOmitted(t *testing.T) {
	t.Parallel()

	remote := graphmodels.NewHardwareConfiguration()
	remote.SetId(ptr("11111111-1111-1111-1111-111111111111"))
	remote.SetDisplayName(ptr("bios-template"))

	configured := base64.StdEncoding.EncodeToString([]byte("[cctk]\n"))
	data := &WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel{
		ConfigurationFileContent: types.StringValue(configured),
	}

	MapRemoteResourceStateToTerraform(context.Background(), data, remote)

	if got := data.ConfigurationFileContent.ValueString(); got != configured {
		t.Errorf("configuration_file_content = %q, want the configured value %q", got, configured)
	}
}

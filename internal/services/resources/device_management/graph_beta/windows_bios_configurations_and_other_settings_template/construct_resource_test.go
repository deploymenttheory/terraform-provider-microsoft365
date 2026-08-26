package graphBetaWindowsBiosConfigurationsAndOtherSettingsTemplate

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func stringSet(t *testing.T, values ...string) types.Set {
	t.Helper()

	elems := make([]attr.Value, len(values))
	for i, v := range values {
		elems[i] = types.StringValue(v)
	}

	set, diags := types.SetValue(types.StringType, elems)
	if diags.HasError() {
		t.Fatalf("failed to build string set: %v", diags.Errors())
	}
	return set
}

func TestUnitConstructResourceEncodesConfigurationFileContentAsBase64(t *testing.T) {
	t.Parallel()

	raw := "[cctk]\nSecureBoot=Enabled\n"
	data := &WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel{
		DisplayName:                 types.StringValue("bios-template"),
		Description:                 types.StringValue("a description"),
		FileName:                    types.StringValue("bios.cctk"),
		ConfigurationFileContent:    types.StringValue(base64.StdEncoding.EncodeToString([]byte(raw))),
		HardwareConfigurationFormat: types.StringValue("dell"),
		PerDevicePasswordDisabled:   types.BoolValue(true),
		RoleScopeTagIds:             stringSet(t, "0", "1"),
	}

	body, err := constructResource(context.Background(), data)
	if err != nil {
		t.Fatalf("constructResource() returned an unexpected error: %v", err)
	}

	// The SDK holds the decoded bytes; Kiota re-encodes them to base64 on the wire, so the
	// value the Graph API receives is byte for byte what the practitioner supplied.
	if got := string(body.GetConfigurationFileContent()); got != raw {
		t.Errorf("configurationFileContent = %q, want %q", got, raw)
	}
	if got := body.GetHardwareConfigurationFormat().String(); got != "dell" {
		t.Errorf("hardwareConfigurationFormat = %q, want %q", got, "dell")
	}
	if got := body.GetPerDevicePasswordDisabled(); got == nil || !*got {
		t.Error("perDevicePasswordDisabled was not set to true")
	}
	if got := len(body.GetRoleScopeTagIds()); got != 2 {
		t.Errorf("roleScopeTagIds length = %d, want 2", got)
	}
}

func TestUnitConstructResourceRejectsNonBase64ConfigurationFileContent(t *testing.T) {
	t.Parallel()

	data := &WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel{
		DisplayName:              types.StringValue("bios-template"),
		FileName:                 types.StringValue("bios.cctk"),
		ConfigurationFileContent: types.StringValue("[cctk] raw text, not base64"),
		RoleScopeTagIds:          stringSet(t, "0"),
	}

	_, err := constructResource(context.Background(), data)
	if !errors.Is(err, errDecodeConfigurationFileContent) {
		t.Fatalf("constructResource() error = %v, want %v", err, errDecodeConfigurationFileContent)
	}
}

func TestUnitConstructResourceOmitsUnsetOptionalAttributes(t *testing.T) {
	t.Parallel()

	data := &WindowsBiosConfigurationsAndOtherSettingsTemplateResourceModel{
		DisplayName:                 types.StringValue("bios-template"),
		Description:                 types.StringNull(),
		FileName:                    types.StringValue("bios.cctk"),
		ConfigurationFileContent:    types.StringNull(),
		HardwareConfigurationFormat: types.StringNull(),
		PerDevicePasswordDisabled:   types.BoolNull(),
		RoleScopeTagIds:             types.SetNull(types.StringType),
	}

	body, err := constructResource(context.Background(), data)
	if err != nil {
		t.Fatalf("constructResource() returned an unexpected error: %v", err)
	}

	if body.GetConfigurationFileContent() != nil {
		t.Error("configurationFileContent was set from a null attribute")
	}
	if body.GetHardwareConfigurationFormat() != nil {
		t.Error("hardwareConfigurationFormat was set from a null attribute")
	}
	if body.GetDescription() != nil {
		t.Error("description was set from a null attribute")
	}
}

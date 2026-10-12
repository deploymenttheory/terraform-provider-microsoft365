package graphBetaDeviceManagementTemplateJson_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_management"
	templatejson "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_template_json"
)

func templatePointer[T any](value T) *T { return &value }

func TestUnitSettingsCatalogTemplateJsonStateResolvesTemplateType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		templateType   string
		templateID     string
		creationSource string
	}{
		{
			"windows_anti_virus_microsoft_defender_antivirus",
			"804339ad-1553-4478-a742-138fb5807418_1",
			"",
		},
		{"windows_firewall", "6078910e-d808-4a9f-a51d-1b8a7bacb7c0_1", ""},
		{"windows_disk_encryption_bitlocker", "46ddfc50-d10f-4867-b852-9434254b3bff_1", ""},
		{
			"windows_attack_surface_reduction_attack_surface_reduction_rules",
			"e8c053d6-9f95-42b1-a7f1-ebfd71c67a4b_1",
			"",
		},
		{
			"windows_local_admin_password_solution_(windows_LAPS)",
			"adc46e5a-f4aa-4ff6-aeff-4f27bc525796_1",
			"",
		},
		{"windows_endpoint_detection_and_response", "0385b795-0f2f-44ac-8602-9f65bf6adede_1", ""},
		{
			"security_baseline_for_microsoft_defender_for_endpoint_version_24H1",
			"49b8320f-e179-472e-8e2c-2fde00289ca2_1",
			"",
		},
		{
			"security_baseline_for_windows_10_and_later_version_24H2",
			"66df8dce-0166-4b82-92f7-1f74e3ca17a3_4",
			"",
		},
		{
			"linux_anti_virus_microsoft_defender_antivirus",
			"4cfd164c-5e8a-4ea9-b15d-9aa71e4ffff4_1",
			"",
		},
		{
			"macOS_anti_virus_microsoft_defender_antivirus",
			"2d345ec2-c817-49e5-9156-3ed416dc972a_1",
			"",
		},
		{"windows_(config_mgr)_firewall_profile", "c2791bb6-ad62-412d-99dc-cb179ef72dee_1", ""},
		{
			"windows_(config_mgr)_attack_surface_reduction_exploit_protection",
			"",
			"ASR_ExploitProtection",
		},
		{"windows_(config_mgr)_attack_surface_reduction_web_protection", "", "ASR_WebProtection"},
		{"windows_(config_mgr)_anti_virus_microsoft_defender_antivirus", "", "SccmAV"},
		{"windows_(config_mgr)_anti_virus_windows_security_experience", "", "WindowsSecurity"},
		{"windows_(config_mgr)_endpoint_detection_and_response", "", "SccmEDR"},
		{"windows_(config_mgr)_firewall", "", "Firewall"},
	}

	for _, test := range tests {
		t.Run(test.templateType, func(t *testing.T) {
			t.Parallel()
			remote := graphmodels.NewDeviceManagementConfigurationPolicy()
			reference := graphmodels.NewDeviceManagementConfigurationPolicyTemplateReference()
			reference.SetTemplateId(&test.templateID)
			remote.SetTemplateReference(reference)
			remote.SetCreationSource(&test.creationSource)

			// Import begins with null; refresh must also replace an outdated known value.
			for _, prior := range []types.String{types.StringNull(), types.StringValue("windows_account_protection")} {
				data := sharedmodels.SettingsCatalogJsonResourceModel{
					SettingsCatalogTemplateType: prior,
				}
				templatejson.MapRemoteResourceStateToTerraform(context.Background(), &data, remote)
				require.Equal(
					t,
					types.StringValue(test.templateType),
					data.SettingsCatalogTemplateType,
				)
			}
		})
	}
}

func TestUnitSettingsCatalogTemplateJsonStateHandlesMissingTemplateMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		nilResource    bool
		hasReference   bool
		templateID     *string
		creationSource *string
		resolvedType   string
	}{
		{name: "nil resource", nilResource: true},
		{name: "missing reference"},
		{name: "nil template ID", hasReference: true},
		{name: "empty template ID", hasReference: true, templateID: templatePointer("")},
		{name: "unknown template ID", hasReference: true, templateID: templatePointer("unknown_1")},
		{
			name:         "unknown template version",
			hasReference: true,
			templateID:   templatePointer("6078910e-d808-4a9f-a51d-1b8a7bacb7c0_99"),
		},
		{name: "unknown creation source", creationSource: templatePointer("unknown")},
		{name: "empty creation source", creationSource: templatePointer("")},
		{
			name:           "creation source without reference",
			creationSource: templatePointer("Firewall"),
			resolvedType:   "windows_(config_mgr)_firewall",
		},
		{
			name:           "creation source with nil ID",
			hasReference:   true,
			creationSource: templatePointer("SccmAV"),
			resolvedType:   "windows_(config_mgr)_anti_virus_microsoft_defender_antivirus",
		},
		{
			name:           "template ID takes precedence",
			hasReference:   true,
			templateID:     templatePointer("6078910e-d808-4a9f-a51d-1b8a7bacb7c0_1"),
			creationSource: templatePointer("Firewall"),
			resolvedType:   "windows_firewall",
		},
		{
			name:           "unknown ID prevents creation source fallback",
			hasReference:   true,
			templateID:     templatePointer("unknown_1"),
			creationSource: templatePointer("Firewall"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var remote graphmodels.DeviceManagementConfigurationPolicyable
			if !test.nilResource {
				remote = graphmodels.NewDeviceManagementConfigurationPolicy()
				remote.SetCreationSource(test.creationSource)
				if test.hasReference {
					reference := graphmodels.NewDeviceManagementConfigurationPolicyTemplateReference()
					reference.SetTemplateId(test.templateID)
					remote.SetTemplateReference(reference)
				}
			}

			for _, prior := range []types.String{types.StringNull(), types.StringUnknown(), types.StringValue("windows_firewall")} {
				data := sharedmodels.SettingsCatalogJsonResourceModel{
					SettingsCatalogTemplateType: prior,
				}
				templatejson.MapRemoteResourceStateToTerraform(context.Background(), &data, remote)
				want := prior
				if test.resolvedType != "" {
					want = types.StringValue(test.resolvedType)
				}
				require.Equal(t, want, data.SettingsCatalogTemplateType)
			}
		})
	}
}

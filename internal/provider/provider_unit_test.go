package provider_test

import (
	"context"
	"testing"

	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider"
)

func TestM365Provider_UnitTestMode(t *testing.T) {
	// Create provider in unit test mode
	providerFunc := provider.NewMicrosoft365Provider("test", true)
	p := providerFunc()

	// Simple test that provider is created successfully
	assert.NotNil(t, p)
}

func TestM365Provider_ValidAuthMethods(t *testing.T) {
	validAuthMethods := []string{
		"azure_developer_cli",
		"azure_cli",
		"client_secret",
		"client_certificate",
		"interactive_browser",
		"device_code",
		"workload_identity",
		"managed_identity",
		"oidc",
		"oidc_github",
		"oidc_azure_devops",
	}

	for _, method := range validAuthMethods {
		t.Run(method, func(t *testing.T) {
			testConfig := `
provider "microsoft365" {
  auth_method = "` + method + `"
}
`

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"microsoft365": providerserver.NewProtocol6WithError(
						provider.NewMicrosoft365Provider("test", true)(),
					),
				},
				Steps: []resource.TestStep{
					{
						Config: testConfig,
						Check:  resource.ComposeTestCheckFunc(),
					},
				},
			})
		})
	}
}

func TestM365Provider_ValidClouds(t *testing.T) {
	validClouds := []string{
		"public",
		"dod",
		"gcc",
		"gcchigh",
		"china",
		"ex",
		"rx",
	}

	for _, cloud := range validClouds {
		t.Run(cloud, func(t *testing.T) {
			testConfig := `
provider "microsoft365" {
  cloud = "` + cloud + `"
  auth_method = "device_code"
}
`

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"microsoft365": providerserver.NewProtocol6WithError(
						provider.NewMicrosoft365Provider("test", true)(),
					),
				},
				Steps: []resource.TestStep{
					{
						Config: testConfig,
						Check:  resource.ComposeTestCheckFunc(),
					},
				},
			})
		})
	}
}

func TestUnit_M365Provider_ClientCertificateSourceValidation(t *testing.T) {
	for _, tc := range []struct {
		name      string
		path      any
		base64    any
		wantError bool
	}{
		{name: "file_only", path: "/path/to/certificate.pfx"},
		{name: "base64_only", base64: "sensitive-certificate-data"},
		{name: "both_sources", path: "/path/to/certificate.pfx", base64: "sensitive-certificate-data", wantError: true},
		{name: "both_empty_strings", path: "", base64: "", wantError: true},
		{name: "neither_source"},
		{name: "unknown_base64", base64: tftypes.UnknownValue},
		{name: "file_with_unknown_base64", path: "/path/to/certificate.pfx", base64: tftypes.UnknownValue},
		{name: "base64_with_unknown_file", path: tftypes.UnknownValue, base64: "sensitive-certificate-data"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			p := provider.NewMicrosoft365Provider("test", true)()
			var schemaResponse frameworkprovider.SchemaResponse
			p.Schema(ctx, frameworkprovider.SchemaRequest{}, &schemaResponse)
			require.False(t, schemaResponse.Diagnostics.HasError())
			schemaType := schemaResponse.Schema.Type().TerraformType(ctx).(tftypes.Object)
			values := make(map[string]tftypes.Value)
			for name, attributeType := range schemaType.AttributeTypes {
				values[name] = tftypes.NewValue(attributeType, nil)
			}
			entraType := schemaType.AttributeTypes["entra_id_options"].(tftypes.Object)
			entraValues := make(map[string]tftypes.Value)
			for name, attributeType := range entraType.AttributeTypes {
				entraValues[name] = tftypes.NewValue(attributeType, nil)
			}
			entraValues["client_certificate"] = tftypes.NewValue(tftypes.String, tc.path)
			entraValues["client_certificate_base64"] = tftypes.NewValue(tftypes.String, tc.base64)
			values["entra_id_options"] = tftypes.NewValue(entraType, entraValues)
			dynamic, err := tfprotov6.NewDynamicValue(
				schemaType,
				tftypes.NewValue(schemaType, values),
			)
			require.NoError(t, err)
			server := providerserver.NewProtocol6(p)()
			response, err := server.ValidateProviderConfig(
				ctx,
				&tfprotov6.ValidateProviderConfigRequest{Config: &dynamic},
			)
			require.NoError(t, err)
			var errorCount int
			for _, diagnostic := range response.Diagnostics {
				if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
					errorCount++
					assert.Contains(t, diagnostic.Detail, "cannot be specified when")
					assert.NotContains(t, diagnostic.Detail, "sensitive-certificate-data")
				}
			}
			if tc.wantError {
				assert.Positive(t, errorCount)
			} else {
				assert.Zero(t, errorCount, "%v", response.Diagnostics)
			}
		})
	}
}

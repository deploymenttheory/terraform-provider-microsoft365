package provider_test

import (
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
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
					"microsoft365": providerserver.NewProtocol6WithError(provider.NewMicrosoft365Provider("test", true)()),
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
					"microsoft365": providerserver.NewProtocol6WithError(provider.NewMicrosoft365Provider("test", true)()),
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
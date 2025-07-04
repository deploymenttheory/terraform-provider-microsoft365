package powershell

import "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"

// ResourceData combines GraphClientInterface and ProviderData for PowerShell resources
type ResourceData struct {
	GraphClient  client.GraphClientInterface
	ProviderData *client.ProviderData
}

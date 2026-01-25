package provider

import (
	"context"

	graphBetaAuditEvents "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/ephemerals/multitenant_management/graph_beta/audit_events"
	//windowsAutopilotDeviceCSVImport "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/ephemeral/utility/graph_beta/windows_autopilot_device_csv_import"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/provider"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ProviderWithEphemeralResources = &M365Provider{}

// EphemeralResources defines the ephemeral resources implemented in the provider.
func (p *M365Provider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		graphBetaAuditEvents.NewAuditEventsEphemeralResource,
		//windowsAutopilotDeviceCSVImport.NewWindowsAutopilotDeviceCSVImportEphemeralResource,
	}
}

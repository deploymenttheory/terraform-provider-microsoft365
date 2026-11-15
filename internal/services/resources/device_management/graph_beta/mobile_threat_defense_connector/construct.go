package graphBetaMobileThreatDefenseConnector

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructResource converts the Terraform resource model to a Kiota SDK model
// Returns a MobileThreatDefenseConnector that can be serialized by Kiota.
// The connector body needs no lookups, so the client is accepted only to keep the CRUD call shape consistent.
func constructResource(
	ctx context.Context,
	_ *msgraphbetasdk.GraphServiceClient,
	data *MobileThreatDefenseConnectorResourceModel,
) (graphmodels.MobileThreatDefenseConnectorable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := constructRequestBody(data)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructRequestBody builds the connector request body from the Terraform model.
// The partner ID is sent as the entity id because Graph uses it to identify the partner on create.
// The *DeviceBlockedOnMissingPartnerData fields are derived by Graph from the *Enabled toggles,
// and the partner* capability fields are set by the service, so none of them are sent.
func constructRequestBody(data *MobileThreatDefenseConnectorResourceModel) graphmodels.MobileThreatDefenseConnectorable {
	requestBody := graphmodels.NewMobileThreatDefenseConnector()

	convert.FrameworkToGraphString(data.ID, requestBody.SetId)
	convert.FrameworkToGraphBool(data.AndroidEnabled, requestBody.SetAndroidEnabled)
	convert.FrameworkToGraphBool(data.IosEnabled, requestBody.SetIosEnabled)
	convert.FrameworkToGraphBool(data.MacEnabled, requestBody.SetMacEnabled)
	convert.FrameworkToGraphBool(data.WindowsEnabled, requestBody.SetWindowsEnabled)
	convert.FrameworkToGraphBool(data.AndroidMobileApplicationManagementEnabled, requestBody.SetAndroidMobileApplicationManagementEnabled)
	convert.FrameworkToGraphBool(data.IosMobileApplicationManagementEnabled, requestBody.SetIosMobileApplicationManagementEnabled)
	convert.FrameworkToGraphBool(data.WindowsMobileApplicationManagementEnabled, requestBody.SetWindowsMobileApplicationManagementEnabled)
	convert.FrameworkToGraphBool(data.PartnerUnsupportedOsVersionBlocked, requestBody.SetPartnerUnsupportedOsVersionBlocked)
	convert.FrameworkToGraphInt32(data.PartnerUnresponsivenessThresholdInDays, requestBody.SetPartnerUnresponsivenessThresholdInDays)
	convert.FrameworkToGraphBool(data.AllowPartnerToCollectIosApplicationMetadata, requestBody.SetAllowPartnerToCollectIOSApplicationMetadata)
	convert.FrameworkToGraphBool(data.AllowPartnerToCollectIosPersonalApplicationMetadata, requestBody.SetAllowPartnerToCollectIOSPersonalApplicationMetadata)
	convert.FrameworkToGraphBool(data.AllowPartnerToCollectIosCertificateMetadata, requestBody.SetAllowPartnerToCollectIosCertificateMetadata)
	convert.FrameworkToGraphBool(data.AllowPartnerToCollectIosPersonalCertificateMetadata, requestBody.SetAllowPartnerToCollectIosPersonalCertificateMetadata)
	convert.FrameworkToGraphBool(data.MicrosoftDefenderForEndpointAttachEnabled, requestBody.SetMicrosoftDefenderForEndpointAttachEnabled)
	convert.FrameworkToGraphBool(data.GrantMobileThreatDefensePartnerRole, requestBody.SetGrantMobileThreatDefensePartnerRole)
	convert.FrameworkToGraphBool(data.LaunchMobileThreatDefensePartnerOnSetupEnabled, requestBody.SetLaunchMobileThreatDefensePartnerOnSetupEnabled)

	return requestBody
}

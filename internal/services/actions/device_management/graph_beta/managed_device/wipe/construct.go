package graphBetaWipeManagedDevice

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructRequest(ctx context.Context, data *WipeManagedDeviceActionModel) (*devicemanagement.ManagedDevicesItemWipePostRequestBody, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s request", ActionName))

	requestBody := devicemanagement.NewManagedDevicesItemWipePostRequestBody()

	// Set keep enrollment data
	convert.FrameworkToGraphBool(data.KeepEnrollmentData, requestBody.SetKeepEnrollmentData)

	// Set keep user data
	convert.FrameworkToGraphBool(data.KeepUserData, requestBody.SetKeepUserData)

	// Set macOS unlock code
	convert.FrameworkToGraphString(data.MacOsUnlockCode, requestBody.SetMacOsUnlockCode)

	// Set obliteration behavior
	if !data.ObliterationBehavior.IsNull() && !data.ObliterationBehavior.IsUnknown() {
		obliterationBehaviorValue := data.ObliterationBehavior.ValueString()
		var obliterationBehavior graphmodels.ObliterationBehavior

		switch obliterationBehaviorValue {
		case "default":
			obliterationBehavior = graphmodels.DEFAULT_OBLITERATIONBEHAVIOR
		case "doNotObliterate":
			obliterationBehavior = graphmodels.DONOTOBLITERATE_OBLITERATIONBEHAVIOR
		case "obliterateWithWarning":
			obliterationBehavior = graphmodels.OBLITERATEWITHWARNING_OBLITERATIONBEHAVIOR
		case "always":
			obliterationBehavior = graphmodels.ALWAYS_OBLITERATIONBEHAVIOR
		default:
			return nil, fmt.Errorf("invalid obliteration_behavior value: %s", obliterationBehaviorValue)
		}

		requestBody.SetObliterationBehavior(&obliterationBehavior)
		tflog.Debug(ctx, fmt.Sprintf("Set obliteration behavior: %s", obliterationBehaviorValue))
	}

	// Set persist eSIM data plan
	convert.FrameworkToGraphBool(data.PersistEsimDataPlan, requestBody.SetPersistEsimDataPlan)

	// Set use protected wipe
	convert.FrameworkToGraphBool(data.UseProtectedWipe, requestBody.SetUseProtectedWipe)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for action %s", ActionName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s request", ActionName))
	return requestBody, nil
}

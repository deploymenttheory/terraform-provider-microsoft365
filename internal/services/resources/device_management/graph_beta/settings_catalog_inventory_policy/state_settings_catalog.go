package graphBetaSettingsCatalogInventoryPolicy

import (
	"context"

	configPolicy "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/device_management/graph_beta/settings_catalog_configuration_policy"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func StateInventoryPolicySettings(ctx context.Context, data *InventoryPolicyResourceModel, settingsResponse graphmodels.DeviceManagementConfigurationSettingCollectionResponseable, plan *InventoryPolicyResourceModel) error {
	tempModel := &configPolicy.SettingsCatalogProfileResourceModel{
		ID: data.ID,
	}

	var tempPlan *configPolicy.SettingsCatalogProfileResourceModel
	if plan != nil && plan.ConfigurationPolicy != nil {
		tempPlan = &configPolicy.SettingsCatalogProfileResourceModel{
			ID:                  plan.ID,
			ConfigurationPolicy: plan.ConfigurationPolicy,
		}
	}

	err := configPolicy.StateConfigurationPolicySettings(ctx, tempModel, settingsResponse, tempPlan)
	if err != nil {
		return err
	}

	data.ConfigurationPolicy = tempModel.ConfigurationPolicy

	if data.Description.IsNull() || data.Description.IsUnknown() {
		data.Description = types.StringValue("")
	}

	return nil
}

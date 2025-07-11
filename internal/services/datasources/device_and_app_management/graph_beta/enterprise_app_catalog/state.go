package enterpriseappcatalog

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapSinglePackageToModel maps a single mobile app catalog package to a model
func mapSinglePackageToModel(ctx context.Context, data models.MobileAppCatalogPackageable) MobileAppCatalogPackageModel {
	tflog.Debug(ctx, "Mapping mobile app catalog package to model", map[string]interface{}{
		"packageId": data.GetId(),
	})

	model := MobileAppCatalogPackageModel{
		Id:                       convert.GraphToFrameworkString(data.GetId()),
		ProductId:                convert.GraphToFrameworkString(data.GetProductId()),
		ProductDisplayName:       convert.GraphToFrameworkString(data.GetProductDisplayName()),
		PublisherDisplayName:     convert.GraphToFrameworkString(data.GetPublisherDisplayName()),
		VersionDisplayName:       convert.GraphToFrameworkString(data.GetVersionDisplayName()),
		BranchDisplayName:        convert.GraphToFrameworkString(data.GetBranchDisplayName()),
		ApplicableArchitectures:  convert.GraphToFrameworkString(data.GetApplicableArchitectures()),
		PackageAutoUpdateCapable: convert.GraphToFrameworkBool(data.GetPackageAutoUpdateCapable()),
	}

	// Handle locales
	locales := data.GetLocales()
	if locales != nil {
		localeValues := make([]types.String, 0, len(locales))
		for _, locale := range locales {
			localeValues = append(localeValues, types.StringValue(locale))
		}
		model.Locales = localeValues
	} else {
		model.Locales = []types.String{}
	}

	return model
}

// MapRemoteStateToDataSource maps a slice of mobile app catalog packages to models
func MapRemoteStateToDataSource(ctx context.Context, data []models.MobileAppCatalogPackageable) ([]MobileAppCatalogPackageModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := make([]MobileAppCatalogPackageModel, 0, len(data))

	for _, item := range data {
		model := mapSinglePackageToModel(ctx, item)
		result = append(result, model)
	}

	return result, diags
}

// MapAppConfigToModel maps app configuration data from API response to AppConfigModel
func MapAppConfigToModel(ctx context.Context, data map[string]interface{}) (*MobileAppCatalogPackageConfigurationModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Map the response to the AppConfigModel
	appConfig := &MobileAppCatalogPackageConfigurationModel{
		ODataType:               types.StringValue(getString(data, "@odata.type")),
		DisplayName:             types.StringValue(getString(data, "displayName")),
		Description:             types.StringValue(getString(data, "description")),
		Publisher:               types.StringValue(getString(data, "publisher")),
		Developer:               types.StringValue(getString(data, "developer")),
		PrivacyInformationUrl:   types.StringValue(getString(data, "privacyInformationUrl")),
		InformationUrl:          types.StringValue(getString(data, "informationUrl")),
		FileName:                types.StringValue(getString(data, "fileName")),
		Size:                    types.Int64Value(getInt64(data, "size")),
		InstallCommandLine:      types.StringValue(getString(data, "installCommandLine")),
		UninstallCommandLine:    types.StringValue(getString(data, "uninstallCommandLine")),
		ApplicableArchitectures: types.StringValue(getString(data, "applicableArchitectures")),
		AllowedArchitectures:    types.StringValue(getString(data, "allowedArchitectures")),
		SetupFilePath:           types.StringValue(getString(data, "setupFilePath")),
		MinSupportedWinRelease:  types.StringValue(getString(data, "minimumSupportedWindowsRelease")),
		DisplayVersion:          types.StringValue(getString(data, "displayVersion")),
		AllowAvailableUninstall: types.BoolValue(getBool(data, "allowAvailableUninstall")),
	}

	// Handle rules
	if rules, ok := data["rules"].([]interface{}); ok {
		appConfig.Rules = parseRules(rules)
	} else {
		appConfig.Rules = []RuleModel{}
	}

	// Handle install experience
	if installExp, ok := data["installExperience"].(map[string]interface{}); ok {
		appConfig.InstallExperience = &InstallExperienceModel{
			RunAsAccount:          types.StringValue(getString(installExp, "runAsAccount")),
			MaxRunTimeInMinutes:   types.Int64Value(getInt64(installExp, "maxRunTimeInMinutes")),
			DeviceRestartBehavior: types.StringValue(getString(installExp, "deviceRestartBehavior")),
		}
	}

	// Handle return codes
	if returnCodes, ok := data["returnCodes"].([]interface{}); ok {
		appConfig.ReturnCodes = parseReturnCodes(returnCodes)
	} else {
		appConfig.ReturnCodes = []ReturnCodeModel{}
	}

	return appConfig, diags
}

// parseRules converts rule data from API response to RuleModel slice
func parseRules(rules []interface{}) []RuleModel {
	result := make([]RuleModel, 0, len(rules))

	for _, r := range rules {
		if rule, ok := r.(map[string]interface{}); ok {
			model := RuleModel{
				ODataType:            types.StringValue(getString(rule, "@odata.type")),
				RuleType:             types.StringValue(getString(rule, "ruleType")),
				Path:                 types.StringValue(getString(rule, "path")),
				FileOrFolderName:     types.StringValue(getString(rule, "fileOrFolderName")),
				Check32BitOn64System: types.BoolValue(getBool(rule, "check32BitOn64System")),
				OperationType:        types.StringValue(getString(rule, "operationType")),
				Operator:             types.StringValue(getString(rule, "operator")),
				ComparisonValue:      types.StringValue(getString(rule, "comparisonValue")),
				KeyPath:              types.StringValue(getString(rule, "keyPath")),
				ValueName:            types.StringValue(getString(rule, "valueName")),
			}
			result = append(result, model)
		}
	}

	return result
}

// parseReturnCodes converts return code data from API response to ReturnCodeModel slice
func parseReturnCodes(codes []interface{}) []ReturnCodeModel {
	result := make([]ReturnCodeModel, 0, len(codes))

	for _, c := range codes {
		if code, ok := c.(map[string]interface{}); ok {
			model := ReturnCodeModel{
				ReturnCode: types.Int64Value(getInt64(code, "returnCode")),
				Type:       types.StringValue(getString(code, "type")),
			}
			result = append(result, model)
		}
	}

	return result
}

// Helper functions to safely extract values from the response map
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key].(string); ok {
		return val
	}
	return ""
}

func getInt64(data map[string]interface{}, key string) int64 {
	switch val := data[key].(type) {
	case int64:
		return val
	case int:
		return int64(val)
	case float64:
		return int64(val)
	default:
		return 0
	}
}

func getBool(data map[string]interface{}, key string) bool {
	if val, ok := data[key].(bool); ok {
		return val
	}
	return false
}

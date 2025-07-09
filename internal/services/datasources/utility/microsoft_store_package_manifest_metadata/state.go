package utilityMicrosoftStorePackageManifest

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mapRemoteStateToTerraformState converts the API response to Terraform model
func (d *MicrosoftStorePackageManifestDataSource) mapRemoteStateToTerraformState(ctx context.Context, manifests []interface{}) ([]PackageManifestDataSourceModel, diag.Diagnostics) {
	var terraformManifests []PackageManifestDataSourceModel
	var diags diag.Diagnostics

	for i, manifestInterface := range manifests {
		manifest, ok := manifestInterface.(map[string]interface{})
		if !ok {
			diags.AddError(
				"Invalid Manifest Data",
				fmt.Sprintf("Expected manifest at index %d to be a map, got %T", i, manifestInterface),
			)
			continue
		}

		terraformManifest := PackageManifestDataSourceModel{
			Type:              d.getStringValue(manifest, "$type"),
			PackageIdentifier: d.getStringValue(manifest, "PackageIdentifier"),
		}

		// Map versions
		if versionsInterface, exists := manifest["Versions"]; exists && versionsInterface != nil {
			if versionsList, ok := versionsInterface.([]interface{}); ok {
				versions, versionDiags := d.mapVersions(ctx, versionsList)
				diags.Append(versionDiags...)
				terraformManifest.Versions = versions
			}
		}

		terraformManifests = append(terraformManifests, terraformManifest)
	}

	return terraformManifests, diags
}

// mapVersions converts version data to Terraform models
func (d *MicrosoftStorePackageManifestDataSource) mapVersions(ctx context.Context, versionsList []interface{}) ([]PackageVersionDataSourceModel, diag.Diagnostics) {
	var terraformVersions []PackageVersionDataSourceModel
	var diags diag.Diagnostics

	for i, versionInterface := range versionsList {
		version, ok := versionInterface.(map[string]interface{})
		if !ok {
			diags.AddError(
				"Invalid Version Data",
				fmt.Sprintf("Expected version at index %d to be a map, got %T", i, versionInterface),
			)
			continue
		}

		terraformVersion := PackageVersionDataSourceModel{
			Type:           d.getStringValue(version, "$type"),
			PackageVersion: d.getStringValue(version, "PackageVersion"),
		}

		// Map default locale
		if defaultLocaleInterface, exists := version["DefaultLocale"]; exists && defaultLocaleInterface != nil {
			if defaultLocaleMap, ok := defaultLocaleInterface.(map[string]interface{}); ok {
				defaultLocale, localeDiags := d.mapDefaultLocale(ctx, defaultLocaleMap)
				diags.Append(localeDiags...)
				terraformVersion.DefaultLocale = defaultLocale
			}
		}

		// Map locales
		if localesInterface, exists := version["Locales"]; exists && localesInterface != nil {
			if localesList, ok := localesInterface.([]interface{}); ok {
				locales, localesDiags := d.mapLocales(ctx, localesList)
				diags.Append(localesDiags...)
				terraformVersion.Locales = locales
			}
		}

		// Map installers
		if installersInterface, exists := version["Installers"]; exists && installersInterface != nil {
			if installersList, ok := installersInterface.([]interface{}); ok {
				installers, installersDiags := d.mapInstallers(ctx, installersList)
				diags.Append(installersDiags...)
				terraformVersion.Installers = installers
			}
		}

		terraformVersions = append(terraformVersions, terraformVersion)
	}

	return terraformVersions, diags
}

// mapDefaultLocale converts default locale data to Terraform model
func (d *MicrosoftStorePackageManifestDataSource) mapDefaultLocale(ctx context.Context, localeMap map[string]interface{}) (*DefaultLocaleDataSourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	defaultLocale := &DefaultLocaleDataSourceModel{
		Type:                d.getStringValue(localeMap, "$type"),
		PackageLocale:       d.getStringValue(localeMap, "PackageLocale"),
		Publisher:           d.getStringValue(localeMap, "Publisher"),
		PublisherUrl:        d.getStringValue(localeMap, "PublisherUrl"),
		PrivacyUrl:          d.getStringValue(localeMap, "PrivacyUrl"),
		PublisherSupportUrl: d.getStringValue(localeMap, "PublisherSupportUrl"),
		PackageName:         d.getStringValue(localeMap, "PackageName"),
		License:             d.getStringValue(localeMap, "License"),
		Copyright:           d.getStringValue(localeMap, "Copyright"),
		ShortDescription:    d.getStringValue(localeMap, "ShortDescription"),
		Description:         d.getStringValue(localeMap, "Description"),
	}

	// Map tags
	if tagsInterface, exists := localeMap["Tags"]; exists && tagsInterface != nil {
		if tagsList, ok := tagsInterface.([]interface{}); ok {
			tags := d.mapStringList(tagsList)
			defaultLocale.Tags = tags
		}
	}

	// Map agreements
	if agreementsInterface, exists := localeMap["Agreements"]; exists && agreementsInterface != nil {
		if agreementsList, ok := agreementsInterface.([]interface{}); ok {
			agreements, agreementDiags := d.mapAgreements(ctx, agreementsList)
			diags.Append(agreementDiags...)
			defaultLocale.Agreements = agreements
		}
	}

	return defaultLocale, diags
}

// mapLocales converts locale data to Terraform models
func (d *MicrosoftStorePackageManifestDataSource) mapLocales(ctx context.Context, localesList []interface{}) ([]LocaleDataSourceModel, diag.Diagnostics) {
	var terraformLocales []LocaleDataSourceModel
	var diags diag.Diagnostics

	for i, localeInterface := range localesList {
		locale, ok := localeInterface.(map[string]interface{})
		if !ok {
			diags.AddError(
				"Invalid Locale Data",
				fmt.Sprintf("Expected locale at index %d to be a map, got %T", i, localeInterface),
			)
			continue
		}

		terraformLocale := LocaleDataSourceModel{
			Type:                d.getStringValue(locale, "$type"),
			PackageLocale:       d.getStringValue(locale, "PackageLocale"),
			Publisher:           d.getStringValue(locale, "Publisher"),
			PublisherUrl:        d.getStringValue(locale, "PublisherUrl"),
			PrivacyUrl:          d.getStringValue(locale, "PrivacyUrl"),
			PublisherSupportUrl: d.getStringValue(locale, "PublisherSupportUrl"),
			PackageName:         d.getStringValue(locale, "PackageName"),
			License:             d.getStringValue(locale, "License"),
			Copyright:           d.getStringValue(locale, "Copyright"),
			ShortDescription:    d.getStringValue(locale, "ShortDescription"),
			Description:         d.getStringValue(locale, "Description"),
		}

		// Map tags
		if tagsInterface, exists := locale["Tags"]; exists && tagsInterface != nil {
			if tagsList, ok := tagsInterface.([]interface{}); ok {
				tags := d.mapStringList(tagsList)
				terraformLocale.Tags = tags
			}
		}

		terraformLocales = append(terraformLocales, terraformLocale)
	}

	return terraformLocales, diags
}

// mapAgreements converts agreement data to Terraform models
func (d *MicrosoftStorePackageManifestDataSource) mapAgreements(ctx context.Context, agreementsList []interface{}) ([]AgreementDataSourceModel, diag.Diagnostics) {
	var terraformAgreements []AgreementDataSourceModel
	var diags diag.Diagnostics

	for i, agreementInterface := range agreementsList {
		agreement, ok := agreementInterface.(map[string]interface{})
		if !ok {
			diags.AddError(
				"Invalid Agreement Data",
				fmt.Sprintf("Expected agreement at index %d to be a map, got %T", i, agreementInterface),
			)
			continue
		}

		terraformAgreement := AgreementDataSourceModel{
			Type:           d.getStringValue(agreement, "$type"),
			AgreementLabel: d.getStringValue(agreement, "AgreementLabel"),
			Agreement:      d.getStringValue(agreement, "Agreement"),
			AgreementUrl:   d.getStringValue(agreement, "AgreementUrl"),
		}

		terraformAgreements = append(terraformAgreements, terraformAgreement)
	}

	return terraformAgreements, diags
}

// mapInstallers converts installer data to Terraform models
func (d *MicrosoftStorePackageManifestDataSource) mapInstallers(ctx context.Context, installersList []interface{}) ([]InstallerDataSourceModel, diag.Diagnostics) {
	var terraformInstallers []InstallerDataSourceModel
	var diags diag.Diagnostics

	for i, installerInterface := range installersList {
		installer, ok := installerInterface.(map[string]interface{})
		if !ok {
			diags.AddError(
				"Invalid Installer Data",
				fmt.Sprintf("Expected installer at index %d to be a map, got %T", i, installerInterface),
			)
			continue
		}

		terraformInstaller := InstallerDataSourceModel{
			Type:                      d.getStringValue(installer, "$type"),
			MSStoreProductIdentifier:  d.getStringValue(installer, "MSStoreProductIdentifier"),
			Architecture:              d.getStringValue(installer, "Architecture"),
			InstallerType:             d.getStringValue(installer, "InstallerType"),
			PackageFamilyName:         d.getStringValue(installer, "PackageFamilyName"),
			Scope:                     d.getStringValue(installer, "Scope"),
			DownloadCommandProhibited: d.getBoolValue(installer, "DownloadCommandProhibited"),
			InstallerSha256:           d.getStringValue(installer, "InstallerSha256"),
			InstallerUrl:              d.getStringValue(installer, "InstallerUrl"),
			InstallerLocale:           d.getStringValue(installer, "InstallerLocale"),
			MinimumOSVersion:          d.getStringValue(installer, "MinimumOSVersion"),
		}

		// Map installer success codes
		if successCodesInterface, exists := installer["InstallerSuccessCodes"]; exists && successCodesInterface != nil {
			if successCodesList, ok := successCodesInterface.([]interface{}); ok {
				successCodes := d.mapInt64List(successCodesList)
				terraformInstaller.InstallerSuccessCodes = successCodes
			}
		}

		// Map markets
		if marketsInterface, exists := installer["Markets"]; exists && marketsInterface != nil {
			if marketsMap, ok := marketsInterface.(map[string]interface{}); ok {
				markets := d.mapMarkets(marketsMap)
				terraformInstaller.Markets = markets
			}
		}

		// Map installer switches
		if switchesInterface, exists := installer["InstallerSwitches"]; exists && switchesInterface != nil {
			if switchesMap, ok := switchesInterface.(map[string]interface{}); ok {
				switches := d.mapInstallerSwitches(switchesMap)
				terraformInstaller.InstallerSwitches = switches
			}
		}

		// Map expected return codes
		if returnCodesInterface, exists := installer["ExpectedReturnCodes"]; exists && returnCodesInterface != nil {
			if returnCodesList, ok := returnCodesInterface.([]interface{}); ok {
				returnCodes, returnCodeDiags := d.mapExpectedReturnCodes(ctx, returnCodesList)
				diags.Append(returnCodeDiags...)
				terraformInstaller.ExpectedReturnCodes = returnCodes
			}
		}

		// Map apps and features entries
		if entriesInterface, exists := installer["AppsAndFeaturesEntries"]; exists && entriesInterface != nil {
			if entriesList, ok := entriesInterface.([]interface{}); ok {
				entries, entriesDiags := d.mapAppsAndFeaturesEntries(ctx, entriesList)
				diags.Append(entriesDiags...)
				terraformInstaller.AppsAndFeaturesEntries = entries
			}
		}

		terraformInstallers = append(terraformInstallers, terraformInstaller)
	}

	return terraformInstallers, diags
}

// mapMarkets converts markets data to Terraform model
func (d *MicrosoftStorePackageManifestDataSource) mapMarkets(marketsMap map[string]interface{}) *MarketsDataSourceModel {
	markets := &MarketsDataSourceModel{
		Type: d.getStringValue(marketsMap, "$type"),
	}

	if allowedMarketsInterface, exists := marketsMap["AllowedMarkets"]; exists && allowedMarketsInterface != nil {
		if allowedMarketsList, ok := allowedMarketsInterface.([]interface{}); ok {
			allowedMarkets := d.mapStringList(allowedMarketsList)
			markets.AllowedMarkets = allowedMarkets
		}
	}

	return markets
}

// mapInstallerSwitches converts installer switches data to Terraform model
func (d *MicrosoftStorePackageManifestDataSource) mapInstallerSwitches(switchesMap map[string]interface{}) *InstallerSwitchesDataSourceModel {
	return &InstallerSwitchesDataSourceModel{
		Type:   d.getStringValue(switchesMap, "$type"),
		Silent: d.getStringValue(switchesMap, "Silent"),
	}
}

// mapExpectedReturnCodes converts expected return codes data to Terraform models
func (d *MicrosoftStorePackageManifestDataSource) mapExpectedReturnCodes(ctx context.Context, returnCodesList []interface{}) ([]ExpectedReturnCodeDataSourceModel, diag.Diagnostics) {
	var terraformReturnCodes []ExpectedReturnCodeDataSourceModel
	var diags diag.Diagnostics

	for i, returnCodeInterface := range returnCodesList {
		returnCode, ok := returnCodeInterface.(map[string]interface{})
		if !ok {
			diags.AddError(
				"Invalid Return Code Data",
				fmt.Sprintf("Expected return code at index %d to be a map, got %T", i, returnCodeInterface),
			)
			continue
		}

		terraformReturnCode := ExpectedReturnCodeDataSourceModel{
			Type:                d.getStringValue(returnCode, "$type"),
			InstallerReturnCode: d.getInt64Value(returnCode, "InstallerReturnCode"),
			ReturnResponse:      d.getStringValue(returnCode, "ReturnResponse"),
		}

		terraformReturnCodes = append(terraformReturnCodes, terraformReturnCode)
	}

	return terraformReturnCodes, diags
}

// mapAppsAndFeaturesEntries converts apps and features entries data to Terraform models
func (d *MicrosoftStorePackageManifestDataSource) mapAppsAndFeaturesEntries(ctx context.Context, entriesList []interface{}) ([]AppsAndFeaturesEntryDataSourceModel, diag.Diagnostics) {
	var terraformEntries []AppsAndFeaturesEntryDataSourceModel
	var diags diag.Diagnostics

	for i, entryInterface := range entriesList {
		entry, ok := entryInterface.(map[string]interface{})
		if !ok {
			diags.AddError(
				"Invalid Apps and Features Entry Data",
				fmt.Sprintf("Expected entry at index %d to be a map, got %T", i, entryInterface),
			)
			continue
		}

		terraformEntry := AppsAndFeaturesEntryDataSourceModel{
			Type:           d.getStringValue(entry, "$type"),
			DisplayName:    d.getStringValue(entry, "DisplayName"),
			Publisher:      d.getStringValue(entry, "Publisher"),
			DisplayVersion: d.getStringValue(entry, "DisplayVersion"),
			ProductCode:    d.getStringValue(entry, "ProductCode"),
			InstallerType:  d.getStringValue(entry, "InstallerType"),
		}

		terraformEntries = append(terraformEntries, terraformEntry)
	}

	return terraformEntries, diags
}

// Helper functions for type conversion and null handling

// getStringValue safely extracts a string value from a map, returning types.StringNull() if not found or null
func (d *MicrosoftStorePackageManifestDataSource) getStringValue(data map[string]interface{}, key string) types.String {
	if value, exists := data[key]; exists && value != nil {
		if strValue, ok := value.(string); ok && strValue != "" {
			return types.StringValue(strValue)
		}
	}
	return types.StringNull()
}

// getBoolValue safely extracts a bool value from a map, returning types.BoolNull() if not found or null
func (d *MicrosoftStorePackageManifestDataSource) getBoolValue(data map[string]interface{}, key string) types.Bool {
	if value, exists := data[key]; exists && value != nil {
		if boolValue, ok := value.(bool); ok {
			return types.BoolValue(boolValue)
		}
	}
	return types.BoolNull()
}

// getInt64Value safely extracts an int64 value from a map, returning types.Int64Null() if not found or null
func (d *MicrosoftStorePackageManifestDataSource) getInt64Value(data map[string]interface{}, key string) types.Int64 {
	if value, exists := data[key]; exists && value != nil {
		switch v := value.(type) {
		case int:
			return types.Int64Value(int64(v))
		case int64:
			return types.Int64Value(v)
		case float64:
			return types.Int64Value(int64(v))
		}
	}
	return types.Int64Null()
}

// mapStringList converts a list of interfaces to a list of types.String
func (d *MicrosoftStorePackageManifestDataSource) mapStringList(list []interface{}) []types.String {
	var result []types.String
	for _, item := range list {
		if strItem, ok := item.(string); ok && strItem != "" {
			result = append(result, types.StringValue(strItem))
		}
	}
	return result
}

// mapInt64List converts a list of interfaces to a list of types.Int64
func (d *MicrosoftStorePackageManifestDataSource) mapInt64List(list []interface{}) []types.Int64 {
	var result []types.Int64
	for _, item := range list {
		switch v := item.(type) {
		case int:
			result = append(result, types.Int64Value(int64(v)))
		case int64:
			result = append(result, types.Int64Value(v))
		case float64:
			result = append(result, types.Int64Value(int64(v)))
		}
	}
	return result
}

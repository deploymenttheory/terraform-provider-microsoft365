package utilityMicrosoftStorePackageManifest

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MicrosoftStorePackageManifestDataSourceModel defines the data source model
type MicrosoftStorePackageManifestDataSourceModel struct {
	ID                types.String                     `tfsdk:"id"`
	PackageIdentifier types.String                     `tfsdk:"package_identifier"`
	SearchTerm        types.String                     `tfsdk:"search_term"`
	Manifests         []PackageManifestDataSourceModel `tfsdk:"manifests"`
	Timeouts          timeouts.Value                   `tfsdk:"timeouts"`
}

// PackageManifestDataSourceModel represents a single package manifest
type PackageManifestDataSourceModel struct {
	Type              types.String                    `tfsdk:"type"`
	PackageIdentifier types.String                    `tfsdk:"package_identifier"`
	Versions          []PackageVersionDataSourceModel `tfsdk:"versions"`
}

// PackageVersionDataSourceModel represents a package version
type PackageVersionDataSourceModel struct {
	Type           types.String                  `tfsdk:"type"`
	PackageVersion types.String                  `tfsdk:"package_version"`
	DefaultLocale  *DefaultLocaleDataSourceModel `tfsdk:"default_locale"`
	Locales        []LocaleDataSourceModel       `tfsdk:"locales"`
	Installers     []InstallerDataSourceModel    `tfsdk:"installers"`
}

// DefaultLocaleDataSourceModel represents the default locale information
type DefaultLocaleDataSourceModel struct {
	Type                types.String               `tfsdk:"type"`
	PackageLocale       types.String               `tfsdk:"package_locale"`
	Publisher           types.String               `tfsdk:"publisher"`
	PublisherUrl        types.String               `tfsdk:"publisher_url"`
	PrivacyUrl          types.String               `tfsdk:"privacy_url"`
	PublisherSupportUrl types.String               `tfsdk:"publisher_support_url"`
	PackageName         types.String               `tfsdk:"package_name"`
	License             types.String               `tfsdk:"license"`
	Copyright           types.String               `tfsdk:"copyright"`
	ShortDescription    types.String               `tfsdk:"short_description"`
	Description         types.String               `tfsdk:"description"`
	Tags                []types.String             `tfsdk:"tags"`
	Agreements          []AgreementDataSourceModel `tfsdk:"agreements"`
}

// LocaleDataSourceModel represents locale-specific information (same as DefaultLocale but without Agreements)
type LocaleDataSourceModel struct {
	Type                types.String   `tfsdk:"type"`
	PackageLocale       types.String   `tfsdk:"package_locale"`
	Publisher           types.String   `tfsdk:"publisher"`
	PublisherUrl        types.String   `tfsdk:"publisher_url"`
	PrivacyUrl          types.String   `tfsdk:"privacy_url"`
	PublisherSupportUrl types.String   `tfsdk:"publisher_support_url"`
	PackageName         types.String   `tfsdk:"package_name"`
	License             types.String   `tfsdk:"license"`
	Copyright           types.String   `tfsdk:"copyright"`
	ShortDescription    types.String   `tfsdk:"short_description"`
	Description         types.String   `tfsdk:"description"`
	Tags                []types.String `tfsdk:"tags"`
}

// AgreementDataSourceModel represents agreement information
type AgreementDataSourceModel struct {
	Type           types.String `tfsdk:"type"`
	AgreementLabel types.String `tfsdk:"agreement_label"`
	Agreement      types.String `tfsdk:"agreement"`
	AgreementUrl   types.String `tfsdk:"agreement_url"`
}

// InstallerDataSourceModel represents installer information
type InstallerDataSourceModel struct {
	Type                      types.String            `tfsdk:"type"`
	MSStoreProductIdentifier  types.String            `tfsdk:"ms_store_product_identifier"`
	Architecture              types.String            `tfsdk:"architecture"`
	InstallerType             types.String            `tfsdk:"installer_type"`
	Markets                   *MarketsDataSourceModel `tfsdk:"markets"`
	PackageFamilyName         types.String            `tfsdk:"package_family_name"`
	Scope                     types.String            `tfsdk:"scope"`
	DownloadCommandProhibited types.Bool              `tfsdk:"download_command_prohibited"`

	// Additional fields for SparkInstaller
	InstallerSha256        types.String                          `tfsdk:"installer_sha256"`
	InstallerUrl           types.String                          `tfsdk:"installer_url"`
	InstallerLocale        types.String                          `tfsdk:"installer_locale"`
	MinimumOSVersion       types.String                          `tfsdk:"minimum_os_version"`
	InstallerSwitches      *InstallerSwitchesDataSourceModel     `tfsdk:"installer_switches"`
	InstallerSuccessCodes  []types.Int64                         `tfsdk:"installer_success_codes"`
	ExpectedReturnCodes    []ExpectedReturnCodeDataSourceModel   `tfsdk:"expected_return_codes"`
	AppsAndFeaturesEntries []AppsAndFeaturesEntryDataSourceModel `tfsdk:"apps_and_features_entries"`
}

// MarketsDataSourceModel represents market information
type MarketsDataSourceModel struct {
	Type           types.String   `tfsdk:"type"`
	AllowedMarkets []types.String `tfsdk:"allowed_markets"`
}

// InstallerSwitchesDataSourceModel represents installer switches
type InstallerSwitchesDataSourceModel struct {
	Type   types.String `tfsdk:"type"`
	Silent types.String `tfsdk:"silent"`
}

// ExpectedReturnCodeDataSourceModel represents expected return codes
type ExpectedReturnCodeDataSourceModel struct {
	Type                types.String `tfsdk:"type"`
	InstallerReturnCode types.Int64  `tfsdk:"installer_return_code"`
	ReturnResponse      types.String `tfsdk:"return_response"`
}

// AppsAndFeaturesEntryDataSourceModel represents apps and features entry
type AppsAndFeaturesEntryDataSourceModel struct {
	Type           types.String `tfsdk:"type"`
	DisplayName    types.String `tfsdk:"display_name"`
	Publisher      types.String `tfsdk:"publisher"`
	DisplayVersion types.String `tfsdk:"display_version"`
	ProductCode    types.String `tfsdk:"product_code"`
	InstallerType  types.String `tfsdk:"installer_type"`
}

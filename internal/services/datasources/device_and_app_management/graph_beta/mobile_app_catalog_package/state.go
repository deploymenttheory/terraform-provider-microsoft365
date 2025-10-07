package graphBetaMobileAppCatalogPackage

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps the remote state to the data source model
func MapRemoteStateToDataSource(ctx context.Context, packageItem graphmodels.MobileAppCatalogPackageable) MobileAppCatalogPackageModel {
	model := MobileAppCatalogPackageModel{}

	// Map ID
	if packageItem.GetId() != nil {
		model.ID = types.StringValue(*packageItem.GetId())
	} else {
		model.ID = types.StringNull()
		tflog.Warn(ctx, "id field is missing in mobile app catalog package")
	}

	// Map product ID
	if packageItem.GetProductId() != nil {
		model.ProductID = types.StringValue(*packageItem.GetProductId())
	} else {
		model.ProductID = types.StringNull()
		tflog.Warn(ctx, "productId field is missing in mobile app catalog package")
	}

	// Map product display name
	if packageItem.GetProductDisplayName() != nil {
		model.ProductDisplayName = types.StringValue(*packageItem.GetProductDisplayName())
	} else {
		model.ProductDisplayName = types.StringNull()
		tflog.Warn(ctx, "productDisplayName field is missing in mobile app catalog package")
	}

	// Map publisher display name
	if packageItem.GetPublisherDisplayName() != nil {
		model.PublisherDisplayName = types.StringValue(*packageItem.GetPublisherDisplayName())
	} else {
		model.PublisherDisplayName = types.StringNull()
		tflog.Warn(ctx, "publisherDisplayName field is missing in mobile app catalog package")
	}

	// Map version display name
	if packageItem.GetVersionDisplayName() != nil {
		model.VersionDisplayName = types.StringValue(*packageItem.GetVersionDisplayName())
	} else {
		model.VersionDisplayName = types.StringNull()
		tflog.Warn(ctx, "versionDisplayName field is missing in mobile app catalog package")
	}

	// Map branch display name - check if it's a Win32MobileAppCatalogPackage
	if win32Package, ok := packageItem.(*graphmodels.Win32MobileAppCatalogPackage); ok {
		if win32Package.GetBranchDisplayName() != nil {
			model.BranchDisplayName = types.StringValue(*win32Package.GetBranchDisplayName())
		} else {
			model.BranchDisplayName = types.StringNull()
		}
	} else {
		model.BranchDisplayName = types.StringNull()
	}

	// Map applicable architectures - check if it's a Win32MobileAppCatalogPackage
	if win32Package, ok := packageItem.(*graphmodels.Win32MobileAppCatalogPackage); ok {
		if win32Package.GetApplicableArchitectures() != nil {
			// Convert enum to string
			archValue := win32Package.GetApplicableArchitectures().String()
			model.ApplicableArchitectures = types.StringValue(archValue)
		} else {
			model.ApplicableArchitectures = types.StringNull()
		}
	} else {
		model.ApplicableArchitectures = types.StringNull()
	}

	// Map locales - check if it's a Win32MobileAppCatalogPackage
	if win32Package, ok := packageItem.(*graphmodels.Win32MobileAppCatalogPackage); ok {
		if win32Package.GetLocales() != nil {
			var locales []types.String
			for _, locale := range win32Package.GetLocales() {
				locales = append(locales, types.StringValue(locale))
			}
			model.Locales = locales
		} else {
			model.Locales = []types.String{}
		}
	} else {
		model.Locales = []types.String{}
	}

	// Map package auto update capable - check if it's a Win32MobileAppCatalogPackage
	if win32Package, ok := packageItem.(*graphmodels.Win32MobileAppCatalogPackage); ok {
		if win32Package.GetPackageAutoUpdateCapable() != nil {
			model.PackageAutoUpdateCapable = types.BoolValue(*win32Package.GetPackageAutoUpdateCapable())
		} else {
			model.PackageAutoUpdateCapable = types.BoolNull()
		}
	} else {
		model.PackageAutoUpdateCapable = types.BoolNull()
	}

	return model
}

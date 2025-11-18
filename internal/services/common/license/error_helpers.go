package license

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// CheckAndAddLicenseError checks if a 403 error is due to missing licenses
// and adds an appropriate diagnostic error message if so.
// Returns true if a license error was detected and added to diagnostics.
func CheckAndAddLicenseError(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	diags *diag.Diagnostics,
	featureName string,
	originalError error,
) bool {
	tflog.Debug(ctx, "Checking if error is due to missing license", map[string]any{
		"feature_name": featureName,
	})

	if client == nil || diags == nil || featureName == "" {
		tflog.Warn(ctx, "Invalid parameters for license error check", map[string]any{
			"feature_name": featureName,
			"client_nil":   client == nil,
			"diags_nil":    diags == nil,
		})
		return false
	}

	requiredLicenses := GetRequiredLicensesForFeature(featureName)
	if len(requiredLicenses) == 0 {
		tflog.Debug(ctx, "No license requirements defined for feature", map[string]any{
			"feature_name": featureName,
		})
		return false
	}

	tflog.Debug(ctx, "Checking for required licenses after error", map[string]any{
		"feature_name":      featureName,
		"required_licenses": requiredLicenses,
	})

	hasLicense := false
	for _, licenseName := range requiredLicenses {
		found, err := CheckLicenseByName(ctx, client, licenseName)
		if err != nil {
			tflog.Warn(ctx, "Failed to check license during error handling", map[string]any{
				"feature_name": featureName,
				"license_name": licenseName,
				"error":        err.Error(),
			})
			continue
		}

		if found {
			tflog.Debug(ctx, "Found required license, error not due to licensing", map[string]any{
				"feature_name": featureName,
				"license_name": licenseName,
			})
			hasLicense = true
			break
		}
	}

	if !hasLicense {
		tflog.Warn(ctx, "Error appears to be due to missing license", map[string]any{
			"feature_name":      featureName,
			"required_licenses": requiredLicenses,
		})

		errorMessage := fmt.Sprintf(
			"This operation failed because the tenant is missing a required license.\n\n"+
				"%s\n\n"+
				"The operation failed with: %s\n\n"+
				"Please ensure your tenant has the appropriate license before using this resource.",
			FormatRequiredLicensesMessage(featureName),
			originalError.Error(),
		)

		diags.AddError("Missing Required License", errorMessage)
		return true
	}

	tflog.Debug(ctx, "Required license found, error not due to licensing", map[string]any{
		"feature_name": featureName,
	})
	return false
}

// GetLicenseDiagnosticInfo returns a formatted string with license information
// suitable for diagnostic messages (when licenses ARE found but something else is wrong)
func GetLicenseDiagnosticInfo(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, featureName string) string {
	tflog.Debug(ctx, "Retrieving license diagnostic info", map[string]any{
		"feature_name": featureName,
	})

	if client == nil {
		tflog.Warn(ctx, "Cannot retrieve license diagnostic info, client is nil", map[string]any{
			"feature_name": featureName,
		})
		return "Unable to retrieve license information (client is nil)"
	}

	requiredLicenses := GetRequiredLicensesForFeature(featureName)
	if len(requiredLicenses) == 0 {
		tflog.Debug(ctx, "No license requirements for feature", map[string]any{
			"feature_name": featureName,
		})
		return ""
	}

	foundLicenses := []string{}
	for _, licenseName := range requiredLicenses {
		found, err := CheckLicenseByName(ctx, client, licenseName)
		if err != nil {
			tflog.Warn(ctx, "Error checking license for diagnostic info", map[string]any{
				"feature_name": featureName,
				"license_name": licenseName,
				"error":        err.Error(),
			})
			continue
		}
		if found {
			foundLicenses = append(foundLicenses, licenseName)
		}
	}

	if len(foundLicenses) == 0 {
		tflog.Debug(ctx, "No required licenses found for diagnostic info", map[string]any{
			"feature_name": featureName,
		})
		return ""
	}

	result := "Tenant licenses found: "
	for i, lic := range foundLicenses {
		if i > 0 {
			result += ", "
		}
		result += lic
	}

	tflog.Debug(ctx, "Generated license diagnostic info", map[string]any{
		"feature_name":    featureName,
		"found_licenses":  foundLicenses,
		"diagnostic_info": result,
	})

	return result
}

// HasRequiredLicense checks if the tenant has at least one of the required licenses
// for a given feature. Returns true if found, false otherwise.
// This is a convenience function for simple boolean checks.
func HasRequiredLicense(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, featureName string) bool {
	tflog.Debug(ctx, "Checking if tenant has required license for feature", map[string]any{
		"feature_name": featureName,
	})

	if client == nil || featureName == "" {
		tflog.Warn(ctx, "Invalid parameters for license check", map[string]any{
			"feature_name": featureName,
			"client_nil":   client == nil,
		})
		return false
	}

	requiredLicenses := GetRequiredLicensesForFeature(featureName)
	if len(requiredLicenses) == 0 {
		tflog.Debug(ctx, "No license requirements defined for feature, assuming available", map[string]any{
			"feature_name": featureName,
		})
		// No license requirements defined - assume available
		return true
	}

	tflog.Debug(ctx, "Checking for required licenses", map[string]any{
		"feature_name":            featureName,
		"required_licenses_count": len(requiredLicenses),
		"required_licenses":       requiredLicenses,
	})

	for i, licenseName := range requiredLicenses {
		tflog.Debug(ctx, "Checking specific license requirement", map[string]any{
			"feature_name": featureName,
			"license_name": licenseName,
			"check_number": i + 1,
			"total_checks": len(requiredLicenses),
		})

		found, err := CheckLicenseByName(ctx, client, licenseName)
		if err != nil {
			tflog.Warn(ctx, "Failed to check license", map[string]any{
				"feature_name": featureName,
				"license_name": licenseName,
				"error":        err.Error(),
			})
			continue
		}

		if found {
			tflog.Info(ctx, "Found required license for feature", map[string]any{
				"feature_name": featureName,
				"license_name": licenseName,
			})
			return true
		}
	}

	tflog.Warn(ctx, "No required license found for feature", map[string]any{
		"feature_name":           featureName,
		"required_licenses":      requiredLicenses,
		"licenses_checked_count": len(requiredLicenses),
	})
	return false
}

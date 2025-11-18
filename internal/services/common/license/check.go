package license

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// LicenseInfo represents a simplified view of a tenant license
type LicenseInfo struct {
	SkuID           string
	SkuPartNumber   string
	ServicePlanName string
	Enabled         bool
	ConsumedUnits   int32
	PrepaidUnits    int32
}

// GetAllEnabledLicenses retrieves all enabled licenses in the tenant
// Returns a slice of LicenseInfo containing SKU and service plan details
func GetAllEnabledLicenses(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) ([]LicenseInfo, error) {
	tflog.Debug(ctx, "Retrieving all enabled licenses from tenant")

	if client == nil {
		return nil, fmt.Errorf("graph client cannot be nil")
	}

	subscribedSkus, err := client.
		SubscribedSkus().
		Get(ctx, nil)

	if err != nil {
		tflog.Error(ctx, "Failed to retrieve subscribed SKUs", map[string]any{
			"error": err.Error(),
		})
		return nil, fmt.Errorf("failed to retrieve tenant licenses: %w", err)
	}

	if subscribedSkus == nil {
		tflog.Warn(ctx, "No subscribed SKUs found in tenant")
		return []LicenseInfo{}, nil
	}

	skuCollection := subscribedSkus.GetValue()
	if len(skuCollection) == 0 {
		tflog.Warn(ctx, "Subscribed SKUs collection is empty")
		return []LicenseInfo{}, nil
	}

	var licenses []LicenseInfo

	for _, sku := range skuCollection {
		skuID := ""
		if sku.GetSkuId() != nil {
			skuID = sku.GetSkuId().String()
		}

		skuPartNumber := ""
		if sku.GetSkuPartNumber() != nil {
			skuPartNumber = *sku.GetSkuPartNumber()
		}

		capabilityStatus := ""
		if sku.GetCapabilityStatus() != nil {
			capabilityStatus = *sku.GetCapabilityStatus()
		}

		consumedUnits := int32(0)
		if sku.GetConsumedUnits() != nil {
			consumedUnits = *sku.GetConsumedUnits()
		}

		prepaidUnits := int32(0)
		if sku.GetPrepaidUnits() != nil && sku.GetPrepaidUnits().GetEnabled() != nil {
			prepaidUnits = *sku.GetPrepaidUnits().GetEnabled()
		}

		// Check if SKU is enabled (capabilityStatus should be "Enabled")
		skuEnabled := strings.EqualFold(capabilityStatus, "Enabled")

		tflog.Debug(ctx, "Found tenant SKU", map[string]any{
			"sku_part_number":   skuPartNumber,
			"sku_id":            skuID,
			"capability_status": capabilityStatus,
			"sku_enabled":       skuEnabled,
			"consumed_units":    consumedUnits,
			"prepaid_units":     prepaidUnits,
		})

		// Only add SKU-level license if the SKU is enabled
		if skuEnabled {
			licenses = append(licenses, LicenseInfo{
				SkuID:           skuID,
				SkuPartNumber:   skuPartNumber,
				ServicePlanName: "", // SKU-level, no specific service plan
				Enabled:         true,
				ConsumedUnits:   consumedUnits,
				PrepaidUnits:    prepaidUnits,
			})
		}

		// Add service plan details (only if parent SKU is enabled)
		if skuEnabled {
			servicePlans := sku.GetServicePlans()
			for _, plan := range servicePlans {
				planName := ""
				if plan.GetServicePlanName() != nil {
					planName = *plan.GetServicePlanName()
				}

				planID := ""
				if plan.GetServicePlanId() != nil {
					planID = plan.GetServicePlanId().String()
				}

				// Check provisioning status for service plan
				provisioningStatus := ""
				if plan.GetProvisioningStatus() != nil {
					provisioningStatus = *plan.GetProvisioningStatus()
				}

				// Service plan is considered enabled if provisioning status is Success, PendingActivation, or PendingInput
				planEnabled := strings.EqualFold(provisioningStatus, "Success") ||
					strings.EqualFold(provisioningStatus, "PendingActivation") ||
					strings.EqualFold(provisioningStatus, "PendingInput")

				tflog.Debug(ctx, "Found service plan in SKU", map[string]any{
					"service_plan_name":   planName,
					"service_plan_id":     planID,
					"sku_part_number":     skuPartNumber,
					"provisioning_status": provisioningStatus,
					"plan_enabled":        planEnabled,
				})

				// Only add enabled service plans from enabled SKUs
				if planEnabled {
					licenses = append(licenses, LicenseInfo{
						SkuID:           skuID,
						SkuPartNumber:   skuPartNumber,
						ServicePlanName: planName,
						Enabled:         true,
						ConsumedUnits:   consumedUnits,
						PrepaidUnits:    prepaidUnits,
					})
				}
			}
		}
	}

	// Count SKUs vs service plans
	skuCount := 0
	servicePlanCount := 0
	for _, lic := range licenses {
		if lic.ServicePlanName == "" {
			skuCount++
		} else {
			servicePlanCount++
		}
	}

	tflog.Debug(ctx, "License retrieval summary", map[string]any{
		"total_licenses":        len(licenses),
		"enabled_skus":          skuCount,
		"enabled_service_plans": servicePlanCount,
	})

	return licenses, nil
}

// CheckLicenseByName checks if a specific license (SKU or service plan) exists in the tenant
// The licenseName parameter can be either a SKU part number (e.g., "ENTERPRISEPACK", "ENTRA_SUITE")
// or a service plan name (e.g., "MICROSOFT_ENTRA_INTERNET_ACCESS", "AAD_PREMIUM_P2")
// Returns true if the license is found and enabled, false otherwise
func CheckLicenseByName(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, licenseName string) (bool, error) {
	tflog.Debug(ctx, "Starting license check", map[string]any{
		"license_name": licenseName,
	})

	if client == nil {
		return false, fmt.Errorf("graph client cannot be nil")
	}

	if licenseName == "" {
		return false, fmt.Errorf("licenseName cannot be empty")
	}

	licenses, err := GetAllEnabledLicenses(ctx, client)
	if err != nil {
		tflog.Error(ctx, "Failed to retrieve licenses for check", map[string]any{
			"license_name": licenseName,
			"error":        err.Error(),
		})
		return false, err
	}

	tflog.Debug(ctx, "Searching through licenses", map[string]any{
		"license_name":   licenseName,
		"total_licenses": len(licenses),
	})

	for _, license := range licenses {
		// Check SKU part number
		if license.SkuPartNumber != "" {
			tflog.Debug(ctx, "Comparing license with SKU", map[string]any{
				"search_term":   licenseName,
				"sku_in_tenant": license.SkuPartNumber,
				"match":         strings.EqualFold(license.SkuPartNumber, licenseName),
			})
			if strings.EqualFold(license.SkuPartNumber, licenseName) {
				tflog.Debug(ctx, "✓ Found matching SKU license", map[string]any{
					"searched_for":   licenseName,
					"found_sku":      license.SkuPartNumber,
					"consumed_units": license.ConsumedUnits,
					"prepaid_units":  license.PrepaidUnits,
				})
				return true, nil
			}
		}

		// Check service plan name
		if license.ServicePlanName != "" {
			tflog.Debug(ctx, "Comparing license with service plan", map[string]any{
				"search_term":            licenseName,
				"service_plan_in_tenant": license.ServicePlanName,
				"in_sku":                 license.SkuPartNumber,
				"match":                  strings.EqualFold(license.ServicePlanName, licenseName),
			})
			if strings.EqualFold(license.ServicePlanName, licenseName) {
				tflog.Debug(ctx, "✓ Found matching service plan license", map[string]any{
					"searched_for":       licenseName,
					"found_service_plan": license.ServicePlanName,
					"in_sku":             license.SkuPartNumber,
				})
				return true, nil
			}
		}
	}

	tflog.Warn(ctx, "License not found in tenant", map[string]any{
		"license_name":   licenseName,
		"total_searched": len(licenses),
	})
	return false, nil
}

// GetLicensesByName returns all licenses matching the given name (SKU or service plan)
// This is useful when you need detailed information about a specific license
func GetLicensesByName(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, licenseName string) ([]LicenseInfo, error) {
	tflog.Debug(ctx, "Retrieving detailed licenses", map[string]any{
		"license_name": licenseName,
	})

	if client == nil {
		return nil, fmt.Errorf("graph client cannot be nil")
	}

	if licenseName == "" {
		return nil, fmt.Errorf("licenseName cannot be empty")
	}

	licenses, err := GetAllEnabledLicenses(ctx, client)
	if err != nil {
		tflog.Error(ctx, "Failed to retrieve licenses", map[string]any{
			"license_name": licenseName,
			"error":        err.Error(),
		})
		return nil, err
	}

	normalizedSearchName := strings.ToLower(strings.TrimSpace(licenseName))
	var matchingLicenses []LicenseInfo

	for _, license := range licenses {
		if strings.EqualFold(license.SkuPartNumber, normalizedSearchName) ||
			strings.EqualFold(license.ServicePlanName, normalizedSearchName) {
			matchingLicenses = append(matchingLicenses, license)
			tflog.Debug(ctx, "Added matching license to results", map[string]any{
				"license_name":    licenseName,
				"sku_part_number": license.SkuPartNumber,
				"service_plan":    license.ServicePlanName,
			})
		}
	}

	tflog.Debug(ctx, "Completed license retrieval", map[string]any{
		"license_name":   licenseName,
		"matches_found":  len(matchingLicenses),
		"total_searched": len(licenses),
	})

	return matchingLicenses, nil
}

// FormatLicensesForError returns a formatted string of license info suitable for error messages
func FormatLicensesForError(licenses []LicenseInfo) string {
	if len(licenses) == 0 {
		return "No licenses found"
	}

	var sb strings.Builder
	sb.WriteString("Licenses found:\n")

	for _, lic := range licenses {
		if lic.ServicePlanName == "" {
			sb.WriteString(fmt.Sprintf("  - SKU: %s (%d/%d units)\n",
				lic.SkuPartNumber, lic.ConsumedUnits, lic.PrepaidUnits))
		} else {
			sb.WriteString(fmt.Sprintf("  - Service Plan: %s (in SKU: %s)\n",
				lic.ServicePlanName, lic.SkuPartNumber))
		}
	}

	return sb.String()
}

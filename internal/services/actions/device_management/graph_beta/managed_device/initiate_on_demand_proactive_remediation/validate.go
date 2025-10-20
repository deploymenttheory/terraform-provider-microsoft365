package graphBetaInitiateOnDemandProactiveRemediationManagedDevice

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (a *InitiateOnDemandProactiveRemediationManagedDeviceAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data InitiateOnDemandProactiveRemediationManagedDeviceActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that at least one device list is provided
	if len(data.ManagedDevices) == 0 && len(data.ComanagedDevices) == 0 {
		resp.Diagnostics.AddError(
			"No Devices Specified",
			"At least one of 'managed_devices' or 'comanaged_devices' must be provided with at least one device configuration.",
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Validating on-demand proactive remediation action for %d managed and %d co-managed device(s)",
		len(data.ManagedDevices), len(data.ComanagedDevices)))

	var nonExistentManagedDevices []string
	var nonExistentComanagedDevices []string
	var unsupportedManagedDevices []string
	var unsupportedComanagedDevices []string
	var invalidScriptPolicies []string

	// Track script policy IDs we've already validated to avoid duplicate checks
	validatedScriptPolicies := make(map[string]bool)

	// Validate managed devices
	for _, managedDevice := range data.ManagedDevices {
		deviceID := managedDevice.DeviceID.ValueString()
		scriptPolicyID := managedDevice.ScriptPolicyID.ValueString()

		// Validate script policy ID exists (if not already validated)
		if !validatedScriptPolicies[scriptPolicyID] {
			script, err := a.client.
				DeviceManagement().
				DeviceHealthScripts().
				ByDeviceHealthScriptId(scriptPolicyID).
				Get(ctx, nil)

			if err != nil {
				if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
					invalidScriptPolicies = append(invalidScriptPolicies, scriptPolicyID)
					validatedScriptPolicies[scriptPolicyID] = false
				} else {
					resp.Diagnostics.AddAttributeError(
						path.Root("managed_devices"),
						"Error Validating Script Policy",
						fmt.Sprintf("Failed to validate script policy %s: %s", scriptPolicyID, err.Error()),
					)
				}
			} else if script != nil {
				validatedScriptPolicies[scriptPolicyID] = true
				tflog.Debug(ctx, fmt.Sprintf("Script policy %s validated successfully", scriptPolicyID))
			}
		}

		// Skip device validation if script policy doesn't exist
		if exists, checked := validatedScriptPolicies[scriptPolicyID]; checked && !exists {
			continue
		}

		// Validate device exists
		device, err := a.client.
			DeviceManagement().
			ManagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				nonExistentManagedDevices = append(nonExistentManagedDevices, deviceID)
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("managed_devices"),
					"Error Validating Managed Device Existence",
					fmt.Sprintf("Failed to check existence of managed device %s: %s", deviceID, err.Error()),
				)
			}
		} else if device != nil {
			// Check if device is Windows (proactive remediations are Windows-specific)
			if device.GetOperatingSystem() != nil {
				osName := *device.GetOperatingSystem()
				if !strings.Contains(strings.ToLower(osName), "windows") {
					unsupportedManagedDevices = append(unsupportedManagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, osName))
					continue
				}
			}

			tflog.Debug(ctx, fmt.Sprintf("Managed device %s validated successfully for script policy %s",
				deviceID, managedDevice.ScriptPolicyID.ValueString()))
		}
	}

	// Validate co-managed devices
	for _, comanagedDevice := range data.ComanagedDevices {
		deviceID := comanagedDevice.DeviceID.ValueString()
		scriptPolicyID := comanagedDevice.ScriptPolicyID.ValueString()

		// Validate script policy ID exists (if not already validated)
		if !validatedScriptPolicies[scriptPolicyID] {
			script, err := a.client.
				DeviceManagement().
				DeviceHealthScripts().
				ByDeviceHealthScriptId(scriptPolicyID).
				Get(ctx, nil)

			if err != nil {
				if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
					invalidScriptPolicies = append(invalidScriptPolicies, scriptPolicyID)
					validatedScriptPolicies[scriptPolicyID] = false
				} else {
					resp.Diagnostics.AddAttributeError(
						path.Root("comanaged_devices"),
						"Error Validating Script Policy",
						fmt.Sprintf("Failed to validate script policy %s: %s", scriptPolicyID, err.Error()),
					)
				}
			} else if script != nil {
				validatedScriptPolicies[scriptPolicyID] = true
				tflog.Debug(ctx, fmt.Sprintf("Script policy %s validated successfully", scriptPolicyID))
			}
		}

		// Skip device validation if script policy doesn't exist
		if exists, checked := validatedScriptPolicies[scriptPolicyID]; checked && !exists {
			continue
		}

		// Validate device exists
		device, err := a.client.
			DeviceManagement().
			ComanagedDevices().
			ByManagedDeviceId(deviceID).
			Get(ctx, nil)

		if err != nil {
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				nonExistentComanagedDevices = append(nonExistentComanagedDevices, deviceID)
			} else {
				resp.Diagnostics.AddAttributeError(
					path.Root("comanaged_devices"),
					"Error Validating Co-Managed Device Existence",
					fmt.Sprintf("Failed to check existence of co-managed device %s: %s", deviceID, err.Error()),
				)
			}
		} else if device != nil {
			// Check if device is Windows
			if device.GetOperatingSystem() != nil {
				osName := *device.GetOperatingSystem()
				if !strings.Contains(strings.ToLower(osName), "windows") {
					unsupportedComanagedDevices = append(unsupportedComanagedDevices, fmt.Sprintf("%s (OS: %s)", deviceID, osName))
					continue
				}
			}

			tflog.Debug(ctx, fmt.Sprintf("Co-managed device %s validated successfully for script policy %s",
				deviceID, comanagedDevice.ScriptPolicyID.ValueString()))
		}
	}

	// Report invalid script policies first (most critical)
	if len(invalidScriptPolicies) > 0 {
		// Remove duplicates for cleaner error message
		uniqueScriptPolicies := make(map[string]bool)
		var uniqueList []string
		for _, id := range invalidScriptPolicies {
			if !uniqueScriptPolicies[id] {
				uniqueScriptPolicies[id] = true
				uniqueList = append(uniqueList, id)
			}
		}

		resp.Diagnostics.AddError(
			"Invalid Script Policy IDs",
			fmt.Sprintf("The following script policy IDs (proactive remediation scripts) do not exist in Intune: %s\n\n"+
				"Please verify:\n"+
				"- Script policy IDs are correct (GUIDs from Intune → Devices → Remediations)\n"+
				"- Scripts are published (not in draft state)\n"+
				"- You have permissions to access the scripts\n\n"+
				"To find script policy IDs:\n"+
				"1. Navigate to Intune → Devices → Remediations\n"+
				"2. Select the remediation script\n"+
				"3. Copy GUID from URL or use Graph API to list device health scripts",
				strings.Join(uniqueList, ", ")),
		)
	}

	if len(nonExistentManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("managed_devices"),
			"Non-Existent Managed Devices",
			fmt.Sprintf("The following managed device IDs do not exist or are not managed by Intune: %s. "+
				"Please ensure all device IDs are correct and refer to existing managed devices.",
				strings.Join(nonExistentManagedDevices, ", ")),
		)
	}

	if len(nonExistentComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("comanaged_devices"),
			"Non-Existent Co-Managed Devices",
			fmt.Sprintf("The following co-managed device IDs do not exist or are not managed by Intune: %s. "+
				"Please ensure all device IDs are correct and refer to existing co-managed devices.",
				strings.Join(nonExistentComanagedDevices, ", ")),
		)
	}

	if len(unsupportedManagedDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("managed_devices"),
			"Non-Windows Devices Detected",
			fmt.Sprintf("The following managed devices are not Windows devices: %s. "+
				"Proactive remediations (health scripts) are only supported on Windows 10/11 devices. "+
				"These devices will be skipped.",
				strings.Join(unsupportedManagedDevices, ", ")),
		)
	}

	if len(unsupportedComanagedDevices) > 0 {
		resp.Diagnostics.AddAttributeWarning(
			path.Root("comanaged_devices"),
			"Non-Windows Devices Detected",
			fmt.Sprintf("The following co-managed devices are not Windows devices: %s. "+
				"Proactive remediations (health scripts) are only supported on Windows 10/11 devices. "+
				"These devices will be skipped.",
				strings.Join(unsupportedComanagedDevices, ", ")),
		)
	}

	// Add informational note about the operation
	if len(data.ManagedDevices)+len(data.ComanagedDevices) > 0 {
		resp.Diagnostics.AddWarning(
			"On-Demand Proactive Remediation",
			fmt.Sprintf("This action will initiate on-demand proactive remediation for %d device(s).\n\n"+
				"Important notes:\n"+
				"- Remediation scripts will execute immediately on device check-in\n"+
				"- Scripts run with SYSTEM privileges\n"+
				"- The specified script policy must already be deployed to each device\n"+
				"- This does not create a new script deployment\n"+
				"- Results will be available in Intune portal and reports\n"+
				"- Scripts may take several minutes to complete execution",
				len(data.ManagedDevices)+len(data.ComanagedDevices)),
		)
	}
}


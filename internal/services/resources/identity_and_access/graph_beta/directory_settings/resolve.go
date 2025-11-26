package graphBetaDirectorySettings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// resolveInstantiatedDirectorySettingsID searches for an existing instantiated directory settings object
// that was created from the specified template ID.
//
// Directory settings work in two stages:
// 1. Template (directorySettingTemplate): A read-only blueprint with a static UUID (e.g., "62375ab9-6b52-47ed-826b-58e47e0e304b")
// 2. Instantiated Object (directorySetting): A modifiable settings object created from a template with its own unique ID
//
// This function:
// - Takes a template ID as input
// - Queries all instantiated directory settings objects in the tenant
// - Finds the one where templateId matches the provided template ID
// - Returns the instantiated object's unique ID (not the template ID)
//
// Returns:
// - The instantiated settings object's ID if found (e.g., "d39871af-cc19-4610-9107-9bbf2a95fe82")
// - Empty string if no instantiated object exists for this template
// - Error if the API call fails
func (r *DirectorySettingsResource) resolveInstantiatedDirectorySettingsID(ctx context.Context, templateID string) (string, error) {
	tflog.Debug(ctx, fmt.Sprintf("Searching for existing settings object for template: %s", templateID))

	settings, err := r.client.
		Settings().
		Get(ctx, nil)

	if err != nil {
		return "", err
	}

	if settings == nil || settings.GetValue() == nil {
		tflog.Debug(ctx, "No settings objects found at tenant level")
		return "", nil
	}

	for _, setting := range settings.GetValue() {
		settingTemplateId := setting.GetTemplateId()
		if settingTemplateId != nil && *settingTemplateId == templateID {
			id := setting.GetId()
			if id != nil {
				tflog.Debug(ctx, fmt.Sprintf("Found existing settings object for template %s", templateID), map[string]any{
					"settingsId": *id,
				})
				return *id, nil
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("No settings object found for template: %s", templateID))
	return "", nil
}

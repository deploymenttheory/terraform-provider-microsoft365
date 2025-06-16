package graphBetaDeviceEnrollmentNotificationConfiguration

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Helper function to resolve notification template ID for a given type
func (r *DeviceEnrollmentNotificationConfigurationResource) resolveNotificationTemplateID(ctx context.Context, configID, templateType string) (string, error) {
	tflog.Debug(ctx, fmt.Sprintf("Resolving template ID for configuration %s, type %s", configID, templateType))

	config, err := r.client.
		DeviceManagement().
		DeviceEnrollmentConfigurations().
		ByDeviceEnrollmentConfigurationId(configID).
		Get(ctx, nil)

	if err != nil {
		return "", fmt.Errorf("failed to get device enrollment configuration: %w", err)
	}

	if enrollmentNotificationConfig, ok := config.(msgraphbetamodels.DeviceEnrollmentNotificationConfigurationable); ok {
		if templates := enrollmentNotificationConfig.GetNotificationTemplates(); templates != nil {
			for _, template := range templates {
				capitalizedType := strings.ToUpper(templateType[:1]) + templateType[1:]
				if strings.HasPrefix(template, capitalizedType+"_") {
					guid := extractGUIDFromTemplateID(template)
					tflog.Debug(ctx, fmt.Sprintf("Resolved template ID for type '%s': %s (GUID: %s)", templateType, template, guid))
					return guid, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no template found for type %s in configuration %s", templateType, configID)
}

// Helper function to extract GUID from template ID
func extractGUIDFromTemplateID(templateID string) string {
	parts := strings.Split(templateID, "_")
	if len(parts) >= 2 {
		return parts[1]
	}
	return templateID
}

package graphBetaWindowsUpdatesApplicableContent

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphadmin "github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
)

func (d *ApplicableContentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ApplicableContentDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var object ApplicableContentDataSourceModel
	object.AudienceId = config.AudienceId
	object.CatalogEntryType = config.CatalogEntryType
	object.DriverClass = config.DriverClass
	object.Manufacturer = config.Manufacturer
	object.DeviceId = config.DeviceId
	object.IncludeNoMatches = config.IncludeNoMatches
	object.ODataFilter = config.ODataFilter
	object.Timeouts = config.Timeouts

	audienceId := object.AudienceId.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading %s for audience_id: %s", DataSourceName, audienceId))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestConfig := &graphadmin.WindowsUpdatesDeploymentAudiencesItemApplicableContentRequestBuilderGetRequestConfiguration{
		QueryParameters: &graphadmin.WindowsUpdatesDeploymentAudiencesItemApplicableContentRequestBuilderGetQueryParameters{
			Expand: []string{"catalogEntry", "matchedDevices"},
		},
	}

	respList, err := d.client.
		Admin().
		Windows().
		Updates().
		DeploymentAudiences().
		ByDeploymentAudienceId(audienceId).
		ApplicableContent().
		Get(ctx, requestConfig)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	var applicableContent []ApplicableContent

	if respList != nil && respList.GetValue() != nil {
		for _, content := range respList.GetValue() {
			mappedContent := MapRemoteStateToDataSource(ctx, content)
			
			// Apply client-side filters
			if shouldIncludeContent(ctx, mappedContent, &object) {
				applicableContent = append(applicableContent, mappedContent)
			}
		}
	}

	object.ApplicableContent = applicableContent

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d applicable content entries (after filtering)", DataSourceName, len(applicableContent)))
}

// shouldIncludeContent applies client-side filtering to determine if content should be included
func shouldIncludeContent(ctx context.Context, content ApplicableContent, config *ApplicableContentDataSourceModel) bool {
	// Filter by catalog entry type
	if !config.CatalogEntryType.IsNull() && !config.CatalogEntryType.IsUnknown() {
		catalogType := config.CatalogEntryType.ValueString()
		if content.CatalogEntry != nil {
			// Map the filter value to OData type check
			// This would need to check the actual @odata.type from the API response
			// For now, we'll use the display name as a heuristic
			displayName := content.CatalogEntry.DisplayName.ValueString()
			switch catalogType {
			case "driver":
				// Driver updates typically have manufacturer/provider info
				if content.CatalogEntry.Manufacturer.IsNull() && content.CatalogEntry.Provider.IsNull() {
					return false
				}
			case "quality":
				// Quality updates typically contain "SecurityUpdate" or "Update"
				if !contains(displayName, "SecurityUpdate") && !contains(displayName, "Update") {
					return false
				}
			case "feature":
				// Feature updates typically contain "Feature" or version numbers
				if !contains(displayName, "Feature") {
					return false
				}
			}
		}
	}

	// Filter by driver class
	if !config.DriverClass.IsNull() && !config.DriverClass.IsUnknown() {
		if content.CatalogEntry == nil || content.CatalogEntry.DriverClass.IsNull() {
			return false
		}
		if content.CatalogEntry.DriverClass.ValueString() != config.DriverClass.ValueString() {
			return false
		}
	}

	// Filter by manufacturer
	if !config.Manufacturer.IsNull() && !config.Manufacturer.IsUnknown() {
		if content.CatalogEntry == nil || content.CatalogEntry.Manufacturer.IsNull() {
			return false
		}
		if content.CatalogEntry.Manufacturer.ValueString() != config.Manufacturer.ValueString() {
			return false
		}
	}

	// Filter by device ID
	if !config.DeviceId.IsNull() && !config.DeviceId.IsUnknown() {
		deviceId := config.DeviceId.ValueString()
		found := false
		for _, device := range content.MatchedDevices {
			if device.DeviceId.ValueString() == deviceId {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Filter by include_no_matches
	if !config.IncludeNoMatches.IsNull() && !config.IncludeNoMatches.IsUnknown() {
		if !config.IncludeNoMatches.ValueBool() && len(content.MatchedDevices) == 0 {
			return false
		}
	}

	return true
}

// contains is a helper function for case-insensitive substring matching
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

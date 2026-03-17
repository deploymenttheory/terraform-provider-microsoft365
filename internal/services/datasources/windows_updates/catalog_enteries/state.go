package graphBetaWindowsUpdateCatalog

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

// MapRemoteStateToDataSource maps a Windows Update Catalog Entry to a model
func MapRemoteStateToDataSource(data graphmodels.CatalogEntryable) WindowsUpdateCatalogEntry {
	model := WindowsUpdateCatalogEntry{
		ID:                      convert.GraphToFrameworkString(data.GetId()),
		DisplayName:             convert.GraphToFrameworkString(data.GetDisplayName()),
		ReleaseDateTime:         convert.GraphToFrameworkTime(data.GetReleaseDateTime()),
		DeployableUntilDateTime: convert.GraphToFrameworkTime(data.GetDeployableUntilDateTime()),
	}

	// Determine the type of catalog entry and set type-specific fields
	switch entry := data.(type) {
	case *graphmodels.FeatureUpdateCatalogEntry:
		model.CatalogEntryType = types.StringValue("featureUpdate")
		model.Version = convert.GraphToFrameworkString(entry.GetVersion())

	case *graphmodels.QualityUpdateCatalogEntry:
		model.CatalogEntryType = types.StringValue("qualityUpdate")
		model.CatalogName = convert.GraphToFrameworkString(entry.GetCatalogName())
		model.ShortName = convert.GraphToFrameworkString(entry.GetShortName())
		model.IsExpeditable = convert.GraphToFrameworkBool(entry.GetIsExpeditable())

		if entry.GetQualityUpdateClassification() != nil {
			model.QualityUpdateClassification = types.StringValue(entry.GetQualityUpdateClassification().String())
		} else {
			model.QualityUpdateClassification = types.StringNull()
		}

		if entry.GetQualityUpdateCadence() != nil {
			model.QualityUpdateCadence = types.StringValue(entry.GetQualityUpdateCadence().String())
		} else {
			model.QualityUpdateCadence = types.StringNull()
		}

		// Map CVE severity information
		if cveSeverity := entry.GetCveSeverityInformation(); cveSeverity != nil {
			model.CveSeverityInformation = &CveSeverityInformation{
				MaxBaseScore: convert.GraphToFrameworkFloat64(cveSeverity.GetMaxBaseScore()),
			}

			// Handle MaxSeverity enum
			if maxSeverity := cveSeverity.GetMaxSeverity(); maxSeverity != nil {
				model.CveSeverityInformation.MaxSeverity = types.StringValue(maxSeverity.String())
			} else {
				model.CveSeverityInformation.MaxSeverity = types.StringNull()
			}

			// Map exploited CVEs
			if exploitedCves := cveSeverity.GetExploitedCves(); exploitedCves != nil {
				var cves []ExploitedCve
				for _, cve := range exploitedCves {
					cves = append(cves, ExploitedCve{
						Number: convert.GraphToFrameworkString(cve.GetNumber()),
						Url:    convert.GraphToFrameworkString(cve.GetUrl()),
					})
				}
				model.CveSeverityInformation.ExploitedCves = cves
			}
		}

	default:
		model.CatalogEntryType = types.StringValue("unknown")
	}

	return model
}

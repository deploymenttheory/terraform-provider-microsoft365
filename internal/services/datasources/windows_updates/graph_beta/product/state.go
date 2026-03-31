package graphBetaWindowsUpdateProduct

import (
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

func MapRemoteStateToDataSource(data graphmodels.Productable) WindowsUpdateProduct {
	model := WindowsUpdateProduct{
		ID:        convert.GraphToFrameworkString(data.GetId()),
		Name:      convert.GraphToFrameworkString(data.GetName()),
		GroupName: convert.GraphToFrameworkString(data.GetGroupName()),
	}

	if friendlyNames := data.GetFriendlyNames(); friendlyNames != nil {
		model.FriendlyNames = make([]types.String, 0, len(friendlyNames))
		for _, name := range friendlyNames {
			model.FriendlyNames = append(model.FriendlyNames, types.StringValue(name))
		}
	}

	if revisions := data.GetRevisions(); revisions != nil {
		model.Revisions = make([]ProductRevision, 0, len(revisions))
		for _, rev := range revisions {
			revision := ProductRevision{
				ID:              convert.GraphToFrameworkString(rev.GetId()),
				DisplayName:     convert.GraphToFrameworkString(rev.GetDisplayName()),
				ReleaseDateTime: convert.GraphToFrameworkTime(rev.GetReleaseDateTime()),
				Version:         convert.GraphToFrameworkString(rev.GetVersion()),
			}

			if osBuild := rev.GetOsBuild(); osBuild != nil {
				revision.OSBuild = &OSBuild{
					MajorVersion:        convert.GraphToFrameworkInt32(osBuild.GetMajorVersion()),
					MinorVersion:        convert.GraphToFrameworkInt32(osBuild.GetMinorVersion()),
					BuildNumber:         convert.GraphToFrameworkInt32(osBuild.GetBuildNumber()),
					UpdateBuildRevision: convert.GraphToFrameworkInt32(osBuild.GetUpdateBuildRevision()),
				}
			}

			if catalogEntry := rev.GetCatalogEntry(); catalogEntry != nil {
				revision.CatalogEntry = &CatalogEntry{
					ID:                      convert.GraphToFrameworkString(catalogEntry.GetId()),
					DisplayName:             convert.GraphToFrameworkString(catalogEntry.GetDisplayName()),
					ReleaseDateTime:         convert.GraphToFrameworkTime(catalogEntry.GetReleaseDateTime()),
					DeployableUntilDateTime: convert.GraphToFrameworkTime(catalogEntry.GetDeployableUntilDateTime()),
				}

				if qualityEntry, ok := catalogEntry.(*graphmodels.QualityUpdateCatalogEntry); ok {
					revision.CatalogEntry.CatalogName = convert.GraphToFrameworkString(qualityEntry.GetCatalogName())
					revision.CatalogEntry.ShortName = convert.GraphToFrameworkString(qualityEntry.GetShortName())
					revision.CatalogEntry.IsExpeditable = convert.GraphToFrameworkBool(qualityEntry.GetIsExpeditable())

					if qualityEntry.GetQualityUpdateClassification() != nil {
						revision.CatalogEntry.QualityUpdateClassification = types.StringValue(qualityEntry.GetQualityUpdateClassification().String())
					} else {
						revision.CatalogEntry.QualityUpdateClassification = types.StringNull()
					}

					if qualityEntry.GetQualityUpdateCadence() != nil {
						revision.CatalogEntry.QualityUpdateCadence = types.StringValue(qualityEntry.GetQualityUpdateCadence().String())
					} else {
						revision.CatalogEntry.QualityUpdateCadence = types.StringNull()
					}
				}
			}

			if kbArticle := rev.GetKnowledgeBaseArticle(); kbArticle != nil {
				revision.KnowledgeBaseArticle = &KnowledgeBaseArticle{
					ID:  convert.GraphToFrameworkString(kbArticle.GetId()),
					URL: convert.GraphToFrameworkString(kbArticle.GetUrl()),
				}
			}

			model.Revisions = append(model.Revisions, revision)
		}
	}

	if knownIssues := data.GetKnownIssues(); knownIssues != nil {
		model.KnownIssues = make([]KnownIssue, 0, len(knownIssues))
		for _, issue := range knownIssues {
			knownIssue := KnownIssue{
				ID:                  convert.GraphToFrameworkString(issue.GetId()),
				Title:               convert.GraphToFrameworkString(issue.GetTitle()),
				Description:         convert.GraphToFrameworkString(issue.GetDescription()),
				WebViewURL:          convert.GraphToFrameworkString(issue.GetWebViewUrl()),
				StartDateTime:       convert.GraphToFrameworkTime(issue.GetStartDateTime()),
				ResolvedDateTime:    convert.GraphToFrameworkTime(issue.GetResolvedDateTime()),
				LastUpdatedDateTime: convert.GraphToFrameworkTime(issue.GetLastUpdatedDateTime()),
			}

			if status := issue.GetStatus(); status != nil {
				knownIssue.Status = types.StringValue(status.String())
			} else {
				knownIssue.Status = types.StringNull()
			}

			if originatingKB := issue.GetOriginatingKnowledgeBaseArticle(); originatingKB != nil {
				knownIssue.OriginatingKnowledgeBaseArticle = &KnowledgeBaseArticle{
					ID:  convert.GraphToFrameworkString(originatingKB.GetId()),
					URL: convert.GraphToFrameworkString(originatingKB.GetUrl()),
				}
			}

			if resolvingKB := issue.GetResolvingKnowledgeBaseArticle(); resolvingKB != nil {
				knownIssue.ResolvingKnowledgeBaseArticle = &KnowledgeBaseArticle{
					ID:  convert.GraphToFrameworkString(resolvingKB.GetId()),
					URL: convert.GraphToFrameworkString(resolvingKB.GetUrl()),
				}
			}

			if safeguardHoldIDs := issue.GetSafeguardHoldIds(); safeguardHoldIDs != nil {
				knownIssue.SafeguardHoldIDs = make([]types.String, 0, len(safeguardHoldIDs))
				for _, holdID := range safeguardHoldIDs {
					knownIssue.SafeguardHoldIDs = append(knownIssue.SafeguardHoldIDs, types.StringValue(fmt.Sprintf("%d", holdID)))
				}
			}

			model.KnownIssues = append(model.KnownIssues, knownIssue)
		}
	}

	return model
}

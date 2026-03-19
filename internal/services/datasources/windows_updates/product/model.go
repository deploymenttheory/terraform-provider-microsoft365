// REF: https://learn.microsoft.com/en-us/graph/api/resources/windowsupdates-product?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/windowsupdates-product-findbykbnumber?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/windowsupdates-product-findbycatalogid?view=graph-rest-beta
package graphBetaWindowsUpdateProduct

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WindowsUpdateProductDataSourceModel struct {
	SearchType  types.String           `tfsdk:"search_type"`
	SearchValue types.String           `tfsdk:"search_value"`
	Products    []WindowsUpdateProduct `tfsdk:"products"`
	Timeouts    timeouts.Value         `tfsdk:"timeouts"`
}

type WindowsUpdateProduct struct {
	ID            types.String      `tfsdk:"id"`
	Name          types.String      `tfsdk:"name"`
	GroupName     types.String      `tfsdk:"group_name"`
	FriendlyNames []types.String    `tfsdk:"friendly_names"`
	Revisions     []ProductRevision `tfsdk:"revisions"`
	KnownIssues   []KnownIssue      `tfsdk:"known_issues"`
}

type ProductRevision struct {
	ID                   types.String          `tfsdk:"id"`
	DisplayName          types.String          `tfsdk:"display_name"`
	ReleaseDateTime      types.String          `tfsdk:"release_date_time"`
	Version              types.String          `tfsdk:"version"`
	OSBuild              *OSBuild              `tfsdk:"os_build"`
	CatalogEntry         *CatalogEntry         `tfsdk:"catalog_entry"`
	KnowledgeBaseArticle *KnowledgeBaseArticle `tfsdk:"knowledge_base_article"`
}

type OSBuild struct {
	MajorVersion        types.Int32 `tfsdk:"major_version"`
	MinorVersion        types.Int32 `tfsdk:"minor_version"`
	BuildNumber         types.Int32 `tfsdk:"build_number"`
	UpdateBuildRevision types.Int32 `tfsdk:"update_build_revision"`
}

type CatalogEntry struct {
	ID                          types.String `tfsdk:"id"`
	DisplayName                 types.String `tfsdk:"display_name"`
	ReleaseDateTime             types.String `tfsdk:"release_date_time"`
	DeployableUntilDateTime     types.String `tfsdk:"deployable_until_date_time"`
	CatalogName                 types.String `tfsdk:"catalog_name"`
	IsExpeditable               types.Bool   `tfsdk:"is_expeditable"`
	QualityUpdateClassification types.String `tfsdk:"quality_update_classification"`
	QualityUpdateCadence        types.String `tfsdk:"quality_update_cadence"`
	ShortName                   types.String `tfsdk:"short_name"`
}

type KnownIssue struct {
	ID                              types.String          `tfsdk:"id"`
	Title                           types.String          `tfsdk:"title"`
	Description                     types.String          `tfsdk:"description"`
	Status                          types.String          `tfsdk:"status"`
	WebViewURL                      types.String          `tfsdk:"web_view_url"`
	StartDateTime                   types.String          `tfsdk:"start_date_time"`
	ResolvedDateTime                types.String          `tfsdk:"resolved_date_time"`
	LastUpdatedDateTime             types.String          `tfsdk:"last_updated_date_time"`
	OriginatingKnowledgeBaseArticle *KnowledgeBaseArticle `tfsdk:"originating_knowledge_base_article"`
	ResolvingKnowledgeBaseArticle   *KnowledgeBaseArticle `tfsdk:"resolving_knowledge_base_article"`
	SafeguardHoldIDs                []types.String        `tfsdk:"safeguard_hold_ids"`
}

type KnowledgeBaseArticle struct {
	ID  types.String `tfsdk:"id"`
	URL types.String `tfsdk:"url"`
}

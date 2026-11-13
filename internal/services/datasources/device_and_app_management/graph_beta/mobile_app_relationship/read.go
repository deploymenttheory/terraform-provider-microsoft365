// REF: https://learn.microsoft.com/en-us/graph/api/intune-apps-mobileapprelationship-list?view=graph-rest-beta
package graphBetaMobileAppRelationship

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	apprelationships "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management/app_relationships"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// Read looks up a single relationship. The root /deviceAppManagement/mobileAppRelationships
// collection rejects $filter, $skip, $orderby and $count, so every list-based lookup reads
// the per-app /mobileApps/{source_id}/relationships navigation instead.
func (d *MobileAppRelationshipDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var object MobileAppRelationshipDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for datasource: %s", DataSourceName))

	ctx, cancel := crud.HandleTimeout(
		ctx,
		object.Timeouts.Read,
		ReadTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	sourceID := object.SourceID.ValueString()
	targetType := object.TargetType.ValueString()

	var (
		match    relationshipMatcher
		filter   string
		criteria string
	)

	switch determineLookupMethod(object) {
	case lookupByODataQuery:
		filter = object.ODataQuery.ValueString()
		match = matchAll
		criteria = fmt.Sprintf("odata_query %q on app %s", filter, sourceID)
	case lookupByTargetId:
		match = matchTargetId(object.TargetID.ValueString())
		criteria = fmt.Sprintf("target_id %q on app %s", object.TargetID.ValueString(), sourceID)
	case lookupByTargetDisplayName:
		match = matchTargetDisplayName(object.TargetDisplayName.ValueString())
		criteria = fmt.Sprintf(
			"target_display_name %q on app %s",
			object.TargetDisplayName.ValueString(),
			sourceID,
		)
	case lookupByRelationshipId:
		relationshipID := object.RelationshipID.ValueString()
		match = matchId(relationshipID)
		criteria = fmt.Sprintf("relationship_id %q", relationshipID)
		if sourceID == "" {
			sourceID = apprelationships.SourceID(relationshipID)
		}
		if sourceID == "" {
			d.readByOpaqueId(ctx, object, resp)
			return
		}
	}

	relationships, err := apprelationships.ListFiltered(ctx, d.client, sourceID, filter)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("App %s returned %d relationships", sourceID, len(relationships)))

	relationship, err := selectRelationship(relationships, match, targetType, criteria)
	if err != nil {
		resp.Diagnostics.AddError("Mobile App Relationship Lookup Failed", err.Error())
		return
	}

	d.setState(ctx, relationship, object, resp)
}

// readByOpaqueId fetches a relationship whose ID does not encode its source app.
func (d *MobileAppRelationshipDataSource) readByOpaqueId(
	ctx context.Context,
	object MobileAppRelationshipDataSourceModel,
	resp *datasource.ReadResponse,
) {
	relationship, err := d.client.
		DeviceAppManagement().
		MobileAppRelationships().
		ByMobileAppRelationshipId(object.RelationshipID.ValueString()).
		Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	relationship, err = selectRelationship(
		[]graphmodels.MobileAppRelationshipable{relationship},
		matchAll,
		object.TargetType.ValueString(),
		fmt.Sprintf("relationship_id %q", object.RelationshipID.ValueString()),
	)
	if err != nil {
		resp.Diagnostics.AddError("Mobile App Relationship Lookup Failed", err.Error())
		return
	}

	d.setState(ctx, relationship, object, resp)
}

func (d *MobileAppRelationshipDataSource) setState(
	ctx context.Context,
	relationship graphmodels.MobileAppRelationshipable,
	object MobileAppRelationshipDataSourceModel,
	resp *datasource.ReadResponse,
) {
	state := MapRemoteStateToDataSource(ctx, relationship, object)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(
		ctx,
		fmt.Sprintf(
			"Finished Datasource Read Method: %s (%s)",
			DataSourceName,
			state.ID.ValueString(),
		),
	)
}

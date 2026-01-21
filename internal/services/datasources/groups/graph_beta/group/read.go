// REF: https://learn.microsoft.com/en-us/graph/api/group-get?view=graph-rest-beta
package graphBetaGroup

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/groups"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// lookupMethod represents the different ways to look up a group
type lookupMethod int

const (
	lookupByODataQuery lookupMethod = iota
	lookupByObjectId
	lookupByDisplayName
	lookupByMailNickname
)

// Read handles the Read operation.
func (d *GroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object GroupDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for datasource: %s", DataSourceName))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var group graphmodels.Groupable
	var err error

	method := determineLookupMethod(object)
	switch method {
	case lookupByODataQuery:
		group, err = d.getGroupByODataQuery(ctx, object)
	case lookupByObjectId:
		group, err = d.getGroupByObjectId(ctx, object)
	case lookupByDisplayName:
		group, err = d.getGroupByDisplayName(ctx, object)
	case lookupByMailNickname:
		group, err = d.getGroupByMailNickname(ctx, object)
	}

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	if group == nil || group.GetId() == nil {
		resp.Diagnostics.AddError(
			"Group Not Found",
			"The group lookup did not return a valid group with an ID. The group may not exist or may not have fully propagated in the directory.",
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully found group with ID: %s", *group.GetId()))

	members, err := d.listAllGroupMembers(ctx, *group.GetId())
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to retrieve group members: %v", err))
	}

	owners, err := d.listAllGroupOwners(ctx, *group.GetId())
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to retrieve group owners: %v", err))
	}

	mappedState := MapRemoteStateToDataSource(ctx, group, members, owners, object)

	resp.Diagnostics.Append(resp.State.Set(ctx, &mappedState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s", DataSourceName))
}

// determineLookupMethod determines which lookup method to use based on provided attributes
func determineLookupMethod(object GroupDataSourceModel) lookupMethod {
	switch {
	case !object.ODataQuery.IsNull() && object.ODataQuery.ValueString() != "":
		return lookupByODataQuery
	case !object.ObjectId.IsNull() && object.ObjectId.ValueString() != "":
		return lookupByObjectId
	case !object.DisplayName.IsNull() && object.DisplayName.ValueString() != "":
		return lookupByDisplayName
	case !object.MailNickname.IsNull() && object.MailNickname.ValueString() != "":
		return lookupByMailNickname
	default:
		return lookupByObjectId // This should never happen due to schema validators
	}
}

// getGroupByObjectId retrieves a group by its object ID
// Includes retry logic because even direct GET can return 404 immediately after creation
func (d *GroupDataSource) getGroupByObjectId(ctx context.Context, object GroupDataSourceModel) (graphmodels.Groupable, error) {
	objectId := object.ObjectId.ValueString()

	maxRetries := 6
	retryDelay := 10 * time.Second
	startTime := time.Now()

	var group graphmodels.Groupable
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		group, err = d.client.Groups().ByGroupId(objectId).Get(ctx, nil)

		if err == nil && group != nil && group.GetId() != nil {
			return group, nil
		}

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errors.IsNonRetryableReadError(&errorInfo) {
				return nil, err
			}
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("group not found after %d attempts (%v total wait): %w", maxRetries, time.Since(startTime), err)
	}
	return nil, fmt.Errorf("group not found or invalid after %d attempts (%v total wait)", maxRetries, time.Since(startTime))
}

// getGroupByODataQuery retrieves a group using a custom OData query
func (d *GroupDataSource) getGroupByODataQuery(ctx context.Context, object GroupDataSourceModel) (graphmodels.Groupable, error) {
	filter := object.ODataQuery.ValueString()
	return d.executeOdataQueryWithRetry(ctx, filter, fmt.Sprintf("OData query: %s", filter))
}

// getGroupByDisplayName retrieves a group by display name with optional filters
func (d *GroupDataSource) getGroupByDisplayName(ctx context.Context, object GroupDataSourceModel) (graphmodels.Groupable, error) {
	filter := buildFilterQuery(
		fmt.Sprintf("displayName eq '%s'", object.DisplayName.ValueString()),
		object.MailEnabled,
		object.SecurityEnabled,
	)
	return d.executeOdataQueryWithRetry(ctx, filter, fmt.Sprintf("display_name: %s", object.DisplayName.ValueString()))
}

// getGroupByMailNickname retrieves a group by mail nickname with optional filters
func (d *GroupDataSource) getGroupByMailNickname(ctx context.Context, object GroupDataSourceModel) (graphmodels.Groupable, error) {
	filter := buildFilterQuery(
		fmt.Sprintf("mailNickname eq '%s'", object.MailNickname.ValueString()),
		object.MailEnabled,
		object.SecurityEnabled,
	)
	return d.executeOdataQueryWithRetry(ctx, filter, fmt.Sprintf("mail_nickname: %s", object.MailNickname.ValueString()))
}

// buildFilterQuery constructs an OData filter query with optional mail_enabled and security_enabled filters
func buildFilterQuery(baseFilter string, mailEnabled, securityEnabled interface {
	IsNull() bool
	ValueBool() bool
}) string {
	filter := baseFilter

	if !mailEnabled.IsNull() {
		filter += fmt.Sprintf(" and mailEnabled eq %t", mailEnabled.ValueBool())
	}

	if !securityEnabled.IsNull() {
		filter += fmt.Sprintf(" and securityEnabled eq %t", securityEnabled.ValueBool())
	}

	return filter
}

// executeOdataQueryWithRetry executes a filtered query with retry logic for eventual consistency
func (d *GroupDataSource) executeOdataQueryWithRetry(ctx context.Context, filter string, description string) (graphmodels.Groupable, error) {
	maxRetries := 6
	retryDelay := 10 * time.Second

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &groups.GroupsRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &groups.GroupsRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		groupsResponse, err := d.client.Groups().Get(ctx, requestConfig)
		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errors.IsNonRetryableReadError(&errorInfo) {
				return nil, err
			}
		} else if len(groupsResponse.GetValue()) > 0 {
			return validateSingleGroup(groupsResponse.GetValue(), description)
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	return validateSingleGroup([]graphmodels.Groupable{}, description)
}

// validateSingleGroup ensures exactly one group was returned
func validateSingleGroup(groupList []graphmodels.Groupable, criteria string) (graphmodels.Groupable, error) {
	switch len(groupList) {
	case 0:
		return nil, fmt.Errorf("no group found with %s", criteria)
	case 1:
		return groupList[0], nil
	default:
		return nil, fmt.Errorf("found %d groups with %s. The query must return exactly one group", len(groupList), criteria)
	}
}

// listAllGroupMembers retrieves all members of a group using pagination
func (d *GroupDataSource) listAllGroupMembers(ctx context.Context, groupId string) ([]string, error) {
	var memberIds []string

	membersResponse, err := d.client.
		Groups().
		ByGroupId(groupId).
		Members().
		Get(ctx, nil)

	if err != nil {
		return nil, err
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.DirectoryObjectable](
		membersResponse,
		d.client.GetAdapter(),
		graphmodels.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator for members: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(item graphmodels.DirectoryObjectable) bool {
		if item != nil && item.GetId() != nil {
			memberIds = append(memberIds, *item.GetId())
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate member pages: %w", err)
	}

	return memberIds, nil
}

// listAllGroupOwners retrieves all owners of a group using pagination
func (d *GroupDataSource) listAllGroupOwners(ctx context.Context, groupId string) ([]string, error) {
	var ownerIds []string

	ownersResponse, err := d.client.
		Groups().
		ByGroupId(groupId).
		Owners().
		Get(ctx, nil)

	if err != nil {
		return nil, err
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.DirectoryObjectable](
		ownersResponse,
		d.client.GetAdapter(),
		graphmodels.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator for owners: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(item graphmodels.DirectoryObjectable) bool {
		if item != nil && item.GetId() != nil {
			ownerIds = append(ownerIds, *item.GetId())
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate owner pages: %w", err)
	}

	return ownerIds, nil
}

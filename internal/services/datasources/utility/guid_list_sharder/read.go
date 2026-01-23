package utilityGuidListSharder

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
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devices"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/groups"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/users"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

func (d *guidListSharderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state GuidListSharderDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	resourceType := state.ResourceType.ValueString()
	strategy := state.Strategy.ValueString()

	var groupId string
	if !state.GroupId.IsNull() {
		groupId = state.GroupId.ValueString()
	}

	var odataQuery string
	if !state.ODataQuery.IsNull() {
		odataQuery = state.ODataQuery.ValueString()
	}

	var seed string
	if !state.Seed.IsNull() {
		seed = state.Seed.ValueString()
	}

	var shardCount int
	var percentages []int64

	if !state.ShardPercentages.IsNull() {
		diags = state.ShardPercentages.ElementsAs(ctx, &percentages, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		shardCount = len(percentages)
	} else {
		shardCount = int(state.ShardCount.ValueInt64())
	}

	// get source GUIDs based upon context
	var guids []string

	switch resourceType {
	case "users":
		guids = d.listAllUsers(ctx, resp, odataQuery)
	case "devices":
		guids = d.listAllDevices(ctx, resp, odataQuery)
	case "group_members":
		guids = d.listAllGroupMembers(ctx, resp, groupId, odataQuery)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d GUIDs for resource_type '%s'", len(guids), resourceType))

	// Apply the sharding strategy based upon context
	var shards [][]string

	switch strategy {
	case "hash":
		shards = shardByHash(guids, shardCount, seed)
	case "round-robin":
		shards = shardByRoundRobin(guids, shardCount, seed)
	case "percentage":
		shards = shardByPercentage(guids, percentages, seed)
	}

	if err := setStateToTerraform(ctx, &state, shards, resourceType, shardCount, strategy); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Set Computed State",
			fmt.Sprintf("Error setting computed state attributes: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", DataSourceName))
}

// listAllUsers retrieves all user GUIDs from Microsoft Graph API
func (d *guidListSharderDataSource) listAllUsers(ctx context.Context, resp *datasource.ReadResponse, filter string) []string {
	var guids []string

	var requestConfig *users.UsersRequestBuilderGetRequestConfiguration
	if filter != "" {
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")
		requestConfig = &users.UsersRequestBuilderGetRequestConfiguration{
			Headers: headers,
			QueryParameters: &users.UsersRequestBuilderGetQueryParameters{
				Filter: &filter,
			},
		}
	}

	usersResponse, err := d.client.
		Users().
		Get(ctx, requestConfig)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return nil
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.Userable](
		usersResponse,
		d.client.GetAdapter(),
		graphmodels.CreateUserCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Page Iterator",
			fmt.Sprintf("Failed to create page iterator: %s", err.Error()),
		)
		return nil
	}

	err = pageIterator.Iterate(ctx, func(item graphmodels.Userable) bool {
		if item != nil && item.GetId() != nil {
			guids = append(guids, *item.GetId())
		}
		return true
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Iterating Pages",
			fmt.Sprintf("Failed to iterate pages: %s", err.Error()),
		)
		return nil
	}

	return guids
}

// listAllDevices retrieves all device GUIDs from Microsoft Graph API
func (d *guidListSharderDataSource) listAllDevices(ctx context.Context, resp *datasource.ReadResponse, filter string) []string {
	var guids []string

	var requestConfig *devices.DevicesRequestBuilderGetRequestConfiguration
	if filter != "" {
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")
		requestConfig = &devices.DevicesRequestBuilderGetRequestConfiguration{
			Headers: headers,
			QueryParameters: &devices.DevicesRequestBuilderGetQueryParameters{
				Filter: &filter,
			},
		}
	}

	devicesResponse, err := d.client.
		Devices().
		Get(ctx, requestConfig)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return nil
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.Deviceable](
		devicesResponse,
		d.client.GetAdapter(),
		graphmodels.CreateDeviceCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Page Iterator",
			fmt.Sprintf("Failed to create page iterator: %s", err.Error()),
		)
		return nil
	}

	err = pageIterator.Iterate(ctx, func(item graphmodels.Deviceable) bool {
		if item != nil && item.GetId() != nil {
			guids = append(guids, *item.GetId())
		}
		return true
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Iterating Pages",
			fmt.Sprintf("Failed to iterate pages: %s", err.Error()),
		)
		return nil
	}

	return guids
}

// listAllGroupMembers retrieves all group member GUIDs from Microsoft Graph API
func (d *guidListSharderDataSource) listAllGroupMembers(ctx context.Context, resp *datasource.ReadResponse, groupId string, filter string) []string {
	var guids []string

	var requestConfig *groups.ItemMembersRequestBuilderGetRequestConfiguration
	if filter != "" {
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")
		requestConfig = &groups.ItemMembersRequestBuilderGetRequestConfiguration{
			Headers: headers,
			QueryParameters: &groups.ItemMembersRequestBuilderGetQueryParameters{
				Filter: &filter,
			},
		}
	}

	membersResponse, err := d.client.
		Groups().
		ByGroupId(groupId).
		Members().
		Get(ctx, requestConfig)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return nil
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.DirectoryObjectable](
		membersResponse,
		d.client.GetAdapter(),
		graphmodels.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Page Iterator",
			fmt.Sprintf("Failed to create page iterator: %s", err.Error()),
		)
		return nil
	}

	err = pageIterator.Iterate(ctx, func(item graphmodels.DirectoryObjectable) bool {
		if item != nil && item.GetId() != nil {
			guids = append(guids, *item.GetId())
		}
		return true
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Iterating Pages",
			fmt.Sprintf("Failed to iterate pages: %s", err.Error()),
		)
		return nil
	}

	return guids
}

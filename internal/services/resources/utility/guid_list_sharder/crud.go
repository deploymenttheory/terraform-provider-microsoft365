package utilityGuidListSharder

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/applications"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devices"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/groups"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/serviceprincipals"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/users"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// Create computes shard assignments from the current Graph API member list and stores them in state.
// recalculate_on_next_run is always honoured as true on first create — initial assignments must be computed.
//
// API Calls:
//   - GET /users                      (when resource_type = "users")
//   - GET /devices                    (when resource_type = "devices")
//   - GET /applications               (when resource_type = "applications")
//   - GET /servicePrincipals          (when resource_type = "service_principals")
//   - GET /groups/{id}/members        (when resource_type = "group_members")
func (r *GuidListSharderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan GuidListSharderResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Create of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	shards, shardCount := r.computeShards(ctx, &plan, resp, constants.TfOperationCreate)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := setStateToTerraform(ctx, &plan, shards, shardCount, plan.Strategy.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to Set State", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read returns shard assignments from state.
//
// When recalculate_on_next_run = false: returns the stored shard assignments unchanged,
// skipping all Graph API calls. This prevents membership churn from causing reassignments
// on every plan and apply.
//
// When recalculate_on_next_run = true: re-queries the Graph API and reruns the sharding
// algorithm, updating state with fresh assignments.
//
// API Calls (only when recalculate_on_next_run = true):
//   - GET /users                      (when resource_type = "users")
//   - GET /devices                    (when resource_type = "devices")
//   - GET /applications               (when resource_type = "applications")
//   - GET /servicePrincipals          (when resource_type = "service_principals")
//   - GET /groups/{id}/members        (when resource_type = "group_members")
func (r *GuidListSharderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GuidListSharderResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.RecalculateOnNextRun.ValueBool() {
		tflog.Debug(ctx, fmt.Sprintf("recalculate_on_next_run=false, returning cached shard state for %s", ResourceName))
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	shards, shardCount := r.computeShards(ctx, &state, resp, constants.TfOperationRead)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := setStateToTerraform(ctx, &state, shards, shardCount, state.Strategy.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to Set State", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update recomputes shard assignments whenever the configuration changes.
// An explicit terraform apply always triggers a fresh query and reshard,
// regardless of the recalculate_on_next_run value. This ensures that intentional
// configuration changes (e.g. changing shard_count or strategy) are always applied.
//
// API Calls:
//   - GET /users                      (when resource_type = "users")
//   - GET /devices                    (when resource_type = "devices")
//   - GET /applications               (when resource_type = "applications")
//   - GET /servicePrincipals          (when resource_type = "service_principals")
//   - GET /groups/{id}/members        (when resource_type = "group_members")
func (r *GuidListSharderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GuidListSharderResourceModel
	var state GuidListSharderResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Preserve the existing id — it was generated on Create and should remain stable.
	plan.Id = state.Id

	if !plan.RecalculateOnNextRun.ValueBool() {
		// Flag is false: accept config changes (e.g. toggling the flag itself, or changing
		// strategy/shard_count for a future reshard) but preserve existing shard assignments.
		// The caller explicitly does not want a reshard on this apply.
		tflog.Debug(ctx, fmt.Sprintf("recalculate_on_next_run=false, preserving shard assignments for %s", ResourceName))
		plan.Shards = state.Shards
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	shards, shardCount := r.computeShards(ctx, &plan, resp, constants.TfOperationUpdate)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := setStateToTerraform(ctx, &plan, shards, shardCount, plan.Strategy.ValueString()); err != nil {
		resp.Diagnostics.AddError("Failed to Set State", err.Error())
		return
	}

	// Restore the stable id after setStateToTerraform (which recomputes it from config).
	plan.Id = state.Id

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete removes the resource from Terraform state.
// There is no backing API resource to delete — shard assignments exist only in state.
func (r *GuidListSharderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))
	resp.State.RemoveResource(ctx)
	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

// =============================================================================
// Internal helpers
// =============================================================================

// computeShards queries the Graph API for the configured resource_type and runs
// the sharding algorithm. Returns (shards, shardCount).
// Any API errors are added directly to resp via errors.HandleKiotaGraphError.
func (r *GuidListSharderResource) computeShards(ctx context.Context, m *GuidListSharderResourceModel, resp any, operation string) ([][]string, int) {
	resourceType := m.ResourceType.ValueString()
	strategy := m.Strategy.ValueString()

	var groupId, odataQuery, seed string
	if !m.GroupId.IsNull() {
		groupId = m.GroupId.ValueString()
	}
	if !m.ODataFilter.IsNull() {
		odataQuery = m.ODataFilter.ValueString()
	}
	if !m.Seed.IsNull() {
		seed = m.Seed.ValueString()
	}

	var shardCount int
	var percentages []int64
	var sizes []int64

	if !m.ShardPercentages.IsNull() {
		diags := m.ShardPercentages.ElementsAs(ctx, &percentages, false)
		if diags.HasError() {
			// ElementsAs errors are framework-level, not API errors — add directly.
			switch r := resp.(type) {
			case *resource.CreateResponse:
				r.Diagnostics.Append(diags...)
			case *resource.ReadResponse:
				r.Diagnostics.Append(diags...)
			case *resource.UpdateResponse:
				r.Diagnostics.Append(diags...)
			}
			return nil, 0
		}
		shardCount = len(percentages)
	} else if !m.ShardSizes.IsNull() {
		diags := m.ShardSizes.ElementsAs(ctx, &sizes, false)
		if diags.HasError() {
			switch r := resp.(type) {
			case *resource.CreateResponse:
				r.Diagnostics.Append(diags...)
			case *resource.ReadResponse:
				r.Diagnostics.Append(diags...)
			case *resource.UpdateResponse:
				r.Diagnostics.Append(diags...)
			}
			return nil, 0
		}
		shardCount = len(sizes)
	} else {
		shardCount = int(m.ShardCount.ValueInt64())
	}

	var guids []string
	var hadError bool

	switch resourceType {
	case "users":
		guids, hadError = r.listAllUsers(ctx, resp, operation, odataQuery)
	case "devices":
		guids, hadError = r.listAllDevices(ctx, resp, operation, odataQuery)
	case "applications":
		guids, hadError = r.listAllApplications(ctx, resp, operation, odataQuery)
	case "service_principals":
		guids, hadError = r.listAllServicePrincipals(ctx, resp, operation, odataQuery)
	case "group_members":
		guids, hadError = r.listAllGroupMembers(ctx, resp, operation, groupId, odataQuery)
	}

	if hadError {
		return nil, 0
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d GUIDs for resource_type '%s'", len(guids), resourceType))

	var shards [][]string
	switch strategy {
	case "round-robin":
		shards = shardByRoundRobin(guids, shardCount, seed)
	case "percentage":
		shards = shardByPercentage(guids, percentages, seed)
	case "size":
		shards = shardBySize(guids, sizes, seed)
	case "rendezvous":
		shards = shardByRendezvous(guids, shardCount, seed)
	}

	return shards, shardCount
}

// listAllUsers retrieves all user GUIDs from Microsoft Graph API.
func (r *GuidListSharderResource) listAllUsers(ctx context.Context, resp any, operation string, filter string) ([]string, bool) {
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

	usersResponse, err := r.client.Users().Get(ctx, requestConfig)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.Userable](
		usersResponse,
		r.client.GetAdapter(),
		graphmodels.CreateUserCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	if err = pageIterator.Iterate(ctx, func(item graphmodels.Userable) bool {
		if item != nil && item.GetId() != nil {
			guids = append(guids, *item.GetId())
		}
		return true
	}); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	return guids, false
}

// listAllDevices retrieves all device GUIDs from Microsoft Graph API.
func (r *GuidListSharderResource) listAllDevices(ctx context.Context, resp any, operation string, filter string) ([]string, bool) {
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

	devicesResponse, err := r.client.Devices().Get(ctx, requestConfig)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.Deviceable](
		devicesResponse,
		r.client.GetAdapter(),
		graphmodels.CreateDeviceCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	if err = pageIterator.Iterate(ctx, func(item graphmodels.Deviceable) bool {
		if item != nil && item.GetId() != nil {
			guids = append(guids, *item.GetId())
		}
		return true
	}); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	return guids, false
}

// listAllApplications retrieves all application GUIDs from Microsoft Graph API.
func (r *GuidListSharderResource) listAllApplications(ctx context.Context, resp any, operation string, filter string) ([]string, bool) {
	var guids []string

	var requestConfig *applications.ApplicationsRequestBuilderGetRequestConfiguration
	if filter != "" {
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")
		requestConfig = &applications.ApplicationsRequestBuilderGetRequestConfiguration{
			Headers: headers,
			QueryParameters: &applications.ApplicationsRequestBuilderGetQueryParameters{
				Filter: &filter,
			},
		}
	}

	applicationsResponse, err := r.client.Applications().Get(ctx, requestConfig)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.Applicationable](
		applicationsResponse,
		r.client.GetAdapter(),
		graphmodels.CreateApplicationCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	if err = pageIterator.Iterate(ctx, func(item graphmodels.Applicationable) bool {
		if item != nil && item.GetId() != nil {
			guids = append(guids, *item.GetId())
		}
		return true
	}); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	return guids, false
}

// listAllServicePrincipals retrieves all service principal GUIDs from Microsoft Graph API.
func (r *GuidListSharderResource) listAllServicePrincipals(ctx context.Context, resp any, operation string, filter string) ([]string, bool) {
	var guids []string

	var requestConfig *serviceprincipals.ServicePrincipalsRequestBuilderGetRequestConfiguration
	if filter != "" {
		headers := abstractions.NewRequestHeaders()
		headers.Add("ConsistencyLevel", "eventual")
		requestConfig = &serviceprincipals.ServicePrincipalsRequestBuilderGetRequestConfiguration{
			Headers: headers,
			QueryParameters: &serviceprincipals.ServicePrincipalsRequestBuilderGetQueryParameters{
				Filter: &filter,
			},
		}
	}

	spResponse, err := r.client.ServicePrincipals().Get(ctx, requestConfig)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.ServicePrincipalable](
		spResponse,
		r.client.GetAdapter(),
		graphmodels.CreateServicePrincipalCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	if err = pageIterator.Iterate(ctx, func(item graphmodels.ServicePrincipalable) bool {
		if item != nil && item.GetId() != nil {
			guids = append(guids, *item.GetId())
		}
		return true
	}); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	return guids, false
}

// listAllGroupMembers retrieves all group member GUIDs from Microsoft Graph API.
func (r *GuidListSharderResource) listAllGroupMembers(ctx context.Context, resp any, operation string, groupId string, filter string) ([]string, bool) {
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

	membersResponse, err := r.client.Groups().ByGroupId(groupId).Members().Get(ctx, requestConfig)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.DirectoryObjectable](
		membersResponse,
		r.client.GetAdapter(),
		graphmodels.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	if err = pageIterator.Iterate(ctx, func(item graphmodels.DirectoryObjectable) bool {
		if item != nil && item.GetId() != nil {
			guids = append(guids, *item.GetId())
		}
		return true
	}); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return nil, true
	}

	return guids, false
}

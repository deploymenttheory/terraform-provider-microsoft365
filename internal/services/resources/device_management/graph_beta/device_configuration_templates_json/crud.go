package graphBetaDeviceConfigurationTemplatesJson

import (
	"context"
	"encoding/json"
	"maps"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	kiotaerrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	identitymodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

func (r *DeviceConfigurationTemplatesJsonResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data DeviceConfigurationTemplatesJsonResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		data.Timeouts.Create,
		CreateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()
	request, err := constructResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Invalid device configuration profile", err.Error())
		return
	}
	assignments, err := constructAssignment(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Invalid profile assignments", err.Error())
		return
	}
	remote, err := r.client.DeviceManagement().DeviceConfigurations().Post(ctx, request, nil)
	if err != nil {
		kiotaerrors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationCreate,
			r.WritePermissions,
		)
		return
	}
	if remote == nil || remote.GetId() == nil || *remote.GetId() == "" {
		resp.Diagnostics.AddError(
			"Invalid create response",
			"Graph did not return the created profile ID.",
		)
		return
	}
	data.ID = types.StringPointerValue(remote.GetId())
	if data.Description.IsUnknown() {
		data.Description = convert.GraphToFrameworkString(remote.GetDescription())
	}
	if data.RoleScopeTagIds.IsUnknown() {
		data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remote.GetRoleScopeTagIds())
	}
	// Preserve the remote ID before the second API write so an assignment failure cannot orphan the profile.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.createRelationships(ctx, data.ID.ValueString(), request); err != nil {
		kiotaerrors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationCreate,
			r.WritePermissions,
		)
		return
	}
	if len(assignments.GetAssignments()) > 0 {
		_, err = r.client.DeviceManagement().
			DeviceConfigurations().
			ByDeviceConfigurationId(data.ID.ValueString()).
			Assign().
			Post(ctx, assignments, nil)
		if err != nil {
			kiotaerrors.HandleKiotaGraphError(
				ctx,
				err,
				resp,
				constants.TfOperationCreate,
				r.WritePermissions,
			)
			return
		}
	}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State},
		&crud.CreateResponseContainer{CreateResponse: resp},
		opts,
	); err != nil {
		resp.Diagnostics.AddError("Cannot read created device configuration profile", err.Error())
	}
}

func (r *DeviceConfigurationTemplatesJsonResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data DeviceConfigurationTemplatesJsonResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		data.Timeouts.Read,
		ReadTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()
	tflog.Debug(
		ctx,
		"Reading device configuration profile",
		map[string]any{"id": data.ID.ValueString()},
	)
	operation := constants.TfOperationRead
	if retryOperation, ok := ctx.Value("retry_operation").(string); ok {
		operation = retryOperation
	}
	var response json.RawMessage
	err := customrequests.JSONRequest(ctx, r.client.GetAdapter(), abstractions.GET,
		r.ResourcePath+"/{id}", map[string]string{"id": data.ID.ValueString()}, nil, &response)
	if err != nil {
		kiotaerrors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, operation, r.ReadPermissions,
			kiotaerrors.GraphErrorOptions{PreserveStateOnReadBadRequest: true})
		return
	}
	if err := MapRemoteResourceStateToTerraform(ctx, &data, response); err != nil {
		resp.Diagnostics.AddError("Cannot read profile metadata", err.Error())
		return
	}
	if err := r.MapRemoteSettingsStateToTerraform(ctx, &data, response); err != nil {
		kiotaerrors.HandleKiotaGraphErrorWithOptions(
			ctx,
			err,
			resp,
			operation,
			r.ReadPermissions,
			kiotaerrors.GraphErrorOptions{
				PreserveStateOnReadBadRequest: true,
				PreserveStateOnReadNotFound:   true,
			},
		)
		return
	}
	var assignments []graphmodels.DeviceConfigurationAssignmentable
	builder := r.client.DeviceManagement().
		DeviceConfigurations().
		ByDeviceConfigurationId(data.ID.ValueString()).
		Assignments()
	for {
		page, err := builder.Get(ctx, nil)
		if err != nil {
			kiotaerrors.HandleKiotaGraphErrorWithOptions(
				ctx,
				err,
				resp,
				operation,
				r.ReadPermissions,
				kiotaerrors.GraphErrorOptions{
					PreserveStateOnReadBadRequest: true,
					PreserveStateOnReadNotFound:   true,
				},
			)
			return
		}
		if page == nil {
			resp.Diagnostics.AddError(
				"Invalid assignment response",
				"Graph returned an empty assignment response.",
			)
			return
		}
		assignments = append(assignments, page.GetValue()...)
		next := page.GetOdataNextLink()
		if next == nil || *next == "" {
			break
		}
		builder = builder.WithUrl(*next)
	}
	resp.Diagnostics.Append(MapAssignmentsToTerraform(ctx, &data, assignments)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if resp.Identity != nil {
		resp.Diagnostics.Append(
			resp.Identity.Set(ctx, identitymodels.ResourceIdentity{ID: data.ID.ValueString()})...)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeviceConfigurationTemplatesJsonResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data, prior DeviceConfigurationTemplatesJsonResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &prior)...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.ID = prior.ID
	if data.Description.IsUnknown() {
		data.Description = prior.Description
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		data.Timeouts.Update,
		UpdateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()
	request, err := constructResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Invalid device configuration profile", err.Error())
		return
	}
	previous, err := constructResource(ctx, &prior)
	if err != nil {
		resp.Diagnostics.AddError("Invalid previous device configuration profile", err.Error())
		return
	}
	if *request.GetOdataType() != *previous.GetOdataType() {
		resp.Diagnostics.AddError(
			"Profile type cannot be changed",
			"Recreate the profile with Terraform's -replace option to change the root @odata.type.",
		)
		return
	}
	assignments, err := constructAssignment(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Invalid profile assignments", err.Error())
		return
	}
	// Relationship bindings are updated through $ref: PATCH cannot remove them.
	settings, previousSettings := maps.Clone(
		request.GetAdditionalData(),
	), maps.Clone(
		previous.GetAdditionalData(),
	)
	for name := range deviceConfigurationRelationships[*request.GetOdataType()] {
		delete(settings, name+"@odata.bind")
		delete(previousSettings, name+"@odata.bind")
	}
	if !reflect.DeepEqual(settings, previousSettings) ||
		!data.DisplayName.Equal(prior.DisplayName) ||
		!data.Description.Equal(prior.Description) ||
		!data.RoleScopeTagIds.Equal(prior.RoleScopeTagIds) {
		bindings := request.GetAdditionalData()
		request.SetAdditionalData(settings)
		_, err = r.client.DeviceManagement().
			DeviceConfigurations().
			ByDeviceConfigurationId(data.ID.ValueString()).
			Patch(ctx, request, nil)
		request.SetAdditionalData(bindings)
		if err != nil {
			kiotaerrors.HandleKiotaGraphError(
				ctx,
				err,
				resp,
				constants.TfOperationUpdate,
				r.WritePermissions,
			)
			return
		}
	}
	if err := r.updateRelationships(ctx, data.ID.ValueString(), previous, request); err != nil {
		kiotaerrors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationUpdate,
			r.WritePermissions,
		)
		return
	}
	if !data.Assignments.Equal(prior.Assignments) {
		_, err = r.client.DeviceManagement().
			DeviceConfigurations().
			ByDeviceConfigurationId(data.ID.ValueString()).
			Assign().
			Post(ctx, assignments, nil)
		if err != nil {
			kiotaerrors.HandleKiotaGraphError(
				ctx,
				err,
				resp,
				constants.TfOperationUpdate,
				r.WritePermissions,
			)
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	if err := crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State},
		&crud.UpdateResponseContainer{UpdateResponse: resp},
		opts,
	); err != nil {
		resp.Diagnostics.AddError("Cannot read updated device configuration profile", err.Error())
	}
}

func (r *DeviceConfigurationTemplatesJsonResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var data DeviceConfigurationTemplatesJsonResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		data.Timeouts.Delete,
		DeleteTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()
	if err := r.client.DeviceManagement().
		DeviceConfigurations().
		ByDeviceConfigurationId(data.ID.ValueString()).
		Delete(ctx, nil); err != nil {
		kiotaerrors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationDelete,
			r.WritePermissions,
		)
		return
	}
	resp.State.RemoveResource(ctx)
}

package graphBetaBrowserSite

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation.
func (r *BrowserSiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object BrowserSiteResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing browser site",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	browserSiteListId := object.BrowserSiteListAssignmentID.ValueString()

	createdSite, err := r.client.
		Admin().
		Edge().
		InternetExplorerMode().
		SiteLists().
		ByBrowserSiteListId(browserSiteListId).
		Sites().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*createdSite.GetId())
	if resp.Diagnostics.HasError() {
		return
	}

	MapRemoteStateToTerraform(ctx, &object, createdSite)

	object.BrowserSiteListAssignmentID = types.StringValue(browserSiteListId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation.
func (r *BrowserSiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object BrowserSiteResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with ID: %s", r.ProviderTypeName, r.TypeName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	browserSiteListId := object.BrowserSiteListAssignmentID.ValueString()

	respResource, err := r.client.
		Admin().
		Edge().
		InternetExplorerMode().
		SiteLists().
		ByBrowserSiteListId(browserSiteListId).
		Sites().
		ByBrowserSiteId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, respResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation.
func (r *BrowserSiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object BrowserSiteResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if object.ID.IsNull() || object.ID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Error updating browser site",
			fmt.Sprintf("Resource ID is missing: %s_%s", r.ProviderTypeName, r.TypeName),
		)
		return
	}

	if object.BrowserSiteListAssignmentID.IsNull() || object.BrowserSiteListAssignmentID.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Error updating browser site",
			fmt.Sprintf("BrowserSiteListAssignmentID is missing: %s_%s", r.ProviderTypeName, r.TypeName),
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing browser site",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
		)
		return
	}

	browserSiteListId := object.BrowserSiteListAssignmentID.ValueString()
	browserSiteId := object.ID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Updating browser site with ID: %s in list: %s", browserSiteId, browserSiteListId))

	_, err = r.client.
		Admin().
		Edge().
		InternetExplorerMode().
		SiteLists().
		ByBrowserSiteListId(browserSiteListId).
		Sites().
		ByBrowserSiteId(browserSiteId).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	updatedResource, err := r.client.
		Admin().
		Edge().
		InternetExplorerMode().
		SiteLists().
		ByBrowserSiteListId(browserSiteListId).
		Sites().
		ByBrowserSiteId(browserSiteId).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read after Update", r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, updatedResource)

	object.BrowserSiteListAssignmentID = types.StringValue(browserSiteListId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation.
func (r *BrowserSiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object BrowserSiteResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	browserSiteListId := object.BrowserSiteListAssignmentID.ValueString()

	err := r.client.
		Admin().
		Edge().
		InternetExplorerMode().
		SiteLists().
		ByBrowserSiteListId(browserSiteListId).
		Sites().
		ByBrowserSiteId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.State.RemoveResource(ctx)
}

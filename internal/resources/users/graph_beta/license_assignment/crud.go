package graphBetaUserLicenseAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/users"
)

// Create handles the creation of a user license assignment.
func (r *UserLicenseAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object UserLicenseAssignmentResourceModel

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

	deadline, _ := ctx.Deadline()
	retryTimeout := time.Until(deadline) - time.Second

	object.ID = object.UserId

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing license assignment request",
			fmt.Sprintf("Could not construct license assignment request: %s", err.Error()),
		)
		return
	}

	_, err = r.client.
		Users().
		ByUserId(object.UserId.ValueString()).
		AssignLicense().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully assigned licenses to user: %s", object.UserId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, resource.ReadRequest{
			State:        resp.State,
			ProviderMeta: req.ProviderMeta,
		}, readResp)

		if readResp.Diagnostics.HasError() {
			return retry.NonRetryableError(fmt.Errorf("error reading resource state after Create Method: %s", readResp.Diagnostics.Errors()))
		}

		resp.State = readResp.State
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for resource creation",
			fmt.Sprintf("Failed to verify resource creation: %s", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", r.TypeName))
}

// Read retrieves the current state of a user's license assignments.
func (r *UserLicenseAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object UserLicenseAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading user license assignments for user ID: %s", object.UserId.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestParameters := &users.UserItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.UserItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "userPrincipalName", "assignedLicenses"},
		},
	}

	user, err := r.client.
		Users().
		ByUserId(object.UserId.ValueString()).
		Get(ctx, requestParameters)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, user)

	licenseDetails, err := r.client.
		Users().
		ByUserId(object.UserId.ValueString()).
		LicenseDetails().
		Get(ctx, nil)

	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to get license details: %s", err.Error()))
	} else {
		MapLicenseDetailsToTerraform(ctx, &object, licenseDetails)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", r.TypeName))
}

// Update handles updates to a user's license assignments.
func (r *UserLicenseAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object UserLicenseAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	deadline, _ := ctx.Deadline()
	retryTimeout := time.Until(deadline) - time.Second

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing license assignment request for update",
			fmt.Sprintf("Could not construct license assignment request: %s", err.Error()),
		)
		return
	}

	_, err = r.client.
		Users().
		ByUserId(object.UserId.ValueString()).
		AssignLicense().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully updated licenses for user: %s", object.UserId.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, resource.ReadRequest{
			State:        resp.State,
			ProviderMeta: req.ProviderMeta,
		}, readResp)

		if readResp.Diagnostics.HasError() {
			return retry.NonRetryableError(fmt.Errorf("error reading resource state after Update Method: %s", readResp.Diagnostics.Errors()))
		}

		resp.State = readResp.State
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for resource update",
			fmt.Sprintf("Failed to verify resource update: %s", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", r.TypeName))
}

// Delete handles the deletion of a user license assignment (removes all managed licenses).
func (r *UserLicenseAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object UserLicenseAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete method for: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Get current licenses to remove
	currentLicenses := make([]string, 0)
	for _, license := range object.AddLicenses {
		currentLicenses = append(currentLicenses, license.SkuId.ValueString())
	}

	// Also check the remove_licenses set in case there are licenses to remove
	removeLicensesSet := object.RemoveLicenses.Elements()
	for _, licenseVal := range removeLicensesSet {
		if strVal, ok := licenseVal.(types.String); ok {
			currentLicenses = append(currentLicenses, strVal.ValueString())
		}
	}

	if len(currentLicenses) > 0 {
		// Create request to remove all licenses managed by this resource
		requestBody := users.NewItemAssignLicensePostRequestBody()
		requestBody.SetAddLicenses([]graphmodels.AssignedLicenseable{})

		// Convert license IDs to UUIDs for removal
		removeLicenseGUIDs := make([]uuid.UUID, 0, len(currentLicenses))
		for _, licenseId := range currentLicenses {
			licenseUUID, err := uuid.Parse(licenseId)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error parsing license ID",
					fmt.Sprintf("Could not parse license ID %s as UUID: %s", licenseId, err.Error()),
				)
				return
			}
			removeLicenseGUIDs = append(removeLicenseGUIDs, licenseUUID)
		}
		requestBody.SetRemoveLicenses(removeLicenseGUIDs)

		_, err := r.client.
			Users().
			ByUserId(object.UserId.ValueString()).
			AssignLicense().
			Post(ctx, requestBody, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully removed licenses from user: %s", object.UserId.ValueString()))
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", r.TypeName))
}

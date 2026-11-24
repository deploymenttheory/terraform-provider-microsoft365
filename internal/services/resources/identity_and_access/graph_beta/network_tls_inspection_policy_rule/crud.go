// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0, networkaccess.TlsInspectionPoliciesItemPolicyRulesPolicyRuleItemRequestBuilder.Get
//     — the only generated GET for a TLS inspection rule is typed as the base policyRule and picks the concrete model from
//     @odata.type. There is no request builder that returns models/networkaccess.TlsInspectionRule.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/networkaccess-tlsinspectionrule-get?view=graph-rest-beta
//   - Observed: when the rule was added, Graph returned TLS inspection rules without @odata.type, so the generated GET fell
//     back to PolicyRule and dropped priority, settings and matchingConditions (recorded in the former requests.go; not
//     re-checked live on 2026-11-24 because the test tenant has no Global Secure Access licence). The GET therefore goes
//     through custom_requests.GetRequest with models.CreateTlsInspectionRuleFromDiscriminatorValue, with APIVersion left empty
//     so it uses the beta client's own base URL as the generated builder does. POST, PATCH and DELETE use the generated builders.
//
// Removal condition: Graph returns @odata.type on TLS inspection rule responses, or the SDK adds a TLS-typed GET for
// tlsInspectionPolicies/{id}/policyRules/{id}, so the rule can be read through the generated builder.
//
// SDK deviation: missing property
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0, models/networkaccess.TlsInspectionRule — the model has no
//     action property, so the generated serializer cannot send it and the parser leaves it in additional data.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/resources/networkaccess-tlsinspectionrule?view=graph-rest-beta
//     documents action (bypass or inspect).
//   - Observed: checked by reading the SDK source on 2026-11-24 (no live call: the test tenant has no Global Secure Access
//     licence). constructResource writes action into additional data and MapRemoteStateToTerraform reads it back from there.
//
// Removal condition: TlsInspectionRule gains GetAction/SetAction, so action can be set and read through the generated model.
package graphBetaNetworkTLSInspectionPolicyRule

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

func (r *NetworkTLSInspectionPolicyRuleResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var object NetworkTLSInspectionPolicyRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		object.Timeouts.Create,
		CreateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()
	body, err := constructResource(ctx, &object, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing TLS inspection resource", err.Error())
		return
	}
	created, err := r.client.NetworkAccess().TlsInspectionPolicies().ByTlsInspectionPolicyId(object.TLSInspectionPolicyID.ValueString()).PolicyRules().Post(ctx, body, nil)
	if err == nil && created == nil {
		err = errEmptyResponse
	}
	if err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationCreate,
			r.WritePermissions,
		)
		return
	}
	if created == nil || created.GetId() == nil || *created.GetId() == "" {
		resp.Diagnostics.AddError(
			"Invalid create response",
			"The API returned no resource id. Check Graph before retrying creation.",
		)
		return
	}
	// Persist the response body identity before readback. The API Location header may contain a placeholder.
	object.ID = types.StringValue(*created.GetId())
	object.Status = types.StringNull()
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(
			resp.Identity.Set(
				ctx,
				TLSInspectionPolicyRuleIdentity{
					ID:                    object.ID.ValueString(),
					TLSInspectionPolicyID: object.TLSInspectionPolicyID.ValueString(),
				},
			)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	// The generated GET would fall back to PolicyRule when Graph omits @odata.type (see the SDK deviation above).
	endpoint := "/networkAccess/tlsInspectionPolicies/" + object.TLSInspectionPolicyID.ValueString() + "/policyRules"
	getConfig := customrequests.GetRequestConfig{Endpoint: endpoint, ResourceID: object.ID.ValueString()}
	errorMappings := abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}
	result, err := customrequests.GetRequest(ctx, r.client.GetAdapter(), getConfig, models.CreateTlsInspectionRuleFromDiscriminatorValue, errorMappings)
	remote, ok := result.(models.TlsInspectionRuleable)
	switch {
	case err != nil:
	case result == nil:
		err = errEmptyResponse
	case !ok:
		err = fmt.Errorf("%w: received %T", errInvalidResponse, result)
	}
	// This is a create readback, not a refresh: even a 404 must retain the created id.
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		resp.Diagnostics.AddError("Invalid TLS inspection resource response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

func (r *NetworkTLSInspectionPolicyRuleResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var object NetworkTLSInspectionPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
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
	// The generated GET would fall back to PolicyRule when Graph omits @odata.type (see the SDK deviation above).
	endpoint := "/networkAccess/tlsInspectionPolicies/" + object.TLSInspectionPolicyID.ValueString() + "/policyRules"
	getConfig := customrequests.GetRequestConfig{Endpoint: endpoint, ResourceID: object.ID.ValueString()}
	errorMappings := abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}
	result, err := customrequests.GetRequest(ctx, r.client.GetAdapter(), getConfig, models.CreateTlsInspectionRuleFromDiscriminatorValue, errorMappings)
	remote, ok := result.(models.TlsInspectionRuleable)
	switch {
	case err != nil:
	case result == nil:
		err = errEmptyResponse
	case !ok:
		err = fmt.Errorf("%w: received %T", errInvalidResponse, result)
	}
	if err != nil {
		errors.HandleKiotaGraphErrorWithOptions(
			ctx,
			err,
			resp,
			constants.TfOperationRead,
			r.ReadPermissions,
			errors.GraphErrorOptions{PreserveStateOnReadBadRequest: true},
		)
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		resp.Diagnostics.AddError("Invalid TLS inspection resource response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(
			resp.Identity.Set(
				ctx,
				TLSInspectionPolicyRuleIdentity{
					ID:                    object.ID.ValueString(),
					TLSInspectionPolicyID: object.TLSInspectionPolicyID.ValueString(),
				},
			)...)
	}
}

func (r *NetworkTLSInspectionPolicyRuleResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan, state NetworkTLSInspectionPolicyRuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		plan.Timeouts.Update,
		UpdateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()
	body, err := constructResource(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing TLS inspection update", err.Error())
		return
	}
	changed := !plan.Name.Equal(state.Name) || !plan.Description.Equal(state.Description) || !plan.Action.Equal(state.Action) ||
		!plan.Priority.Equal(state.Priority) || !plan.Enabled.Equal(state.Enabled) || !plan.Destinations.Equal(state.Destinations)
	if changed {
		policy := r.client.NetworkAccess().TlsInspectionPolicies().ByTlsInspectionPolicyId(state.TLSInspectionPolicyID.ValueString())
		rule := policy.PolicyRules().ByPolicyRuleId(state.ID.ValueString())
		if _, err := rule.Patch(ctx, body, nil); err != nil {
			errors.HandleKiotaGraphError(
				ctx,
				err,
				resp,
				constants.TfOperationUpdate,
				r.WritePermissions,
			)
			return
		}
	}
	// Keep the previous known state if readback fails after a successful PATCH.
	resp.State = req.State
	object := plan
	object.ID = state.ID
	// The generated GET would fall back to PolicyRule when Graph omits @odata.type (see the SDK deviation above).
	endpoint := "/networkAccess/tlsInspectionPolicies/" + object.TLSInspectionPolicyID.ValueString() + "/policyRules"
	getConfig := customrequests.GetRequestConfig{Endpoint: endpoint, ResourceID: object.ID.ValueString()}
	errorMappings := abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}
	result, err := customrequests.GetRequest(ctx, r.client.GetAdapter(), getConfig, models.CreateTlsInspectionRuleFromDiscriminatorValue, errorMappings)
	remote, ok := result.(models.TlsInspectionRuleable)
	switch {
	case err != nil:
	case result == nil:
		err = errEmptyResponse
	case !ok:
		err = fmt.Errorf("%w: received %T", errInvalidResponse, result)
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.ReadPermissions)
		return
	}
	if err := MapRemoteStateToTerraform(ctx, &object, remote); err != nil {
		resp.Diagnostics.AddError("Invalid TLS inspection resource response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(
			resp.Identity.Set(
				ctx,
				TLSInspectionPolicyRuleIdentity{
					ID:                    object.ID.ValueString(),
					TLSInspectionPolicyID: object.TLSInspectionPolicyID.ValueString(),
				},
			)...)
	}
}

func (r *NetworkTLSInspectionPolicyRuleResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var object NetworkTLSInspectionPolicyRuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(
		ctx,
		object.Timeouts.Delete,
		DeleteTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()
	policy := r.client.NetworkAccess().TlsInspectionPolicies().ByTlsInspectionPolicyId(object.TLSInspectionPolicyID.ValueString())
	rule := policy.PolicyRules().ByPolicyRuleId(object.ID.ValueString())
	if err := rule.Delete(ctx, nil); err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationDelete,
			r.WritePermissions,
		)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	resp.State.RemoveResource(ctx)
}

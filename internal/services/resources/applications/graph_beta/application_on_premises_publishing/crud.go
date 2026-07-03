package graphBetaApplicationsOnPremisesPublishing

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	kiotaerrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphapplications "github.com/microsoftgraph/msgraph-beta-sdk-go/applications"
	grapherrors "github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

// Create handles the Create operation for On-Premises Publishing resources.
//
// Operation: Configure on-premises publishing for an application
// API Calls:
//   - PATCH /applications/{application-id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta
func (r *OnPremisesPublishingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object OnPremisesPublishingResourceModel

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

	// PATCH the application with on-premises publishing settings
	if err := r.patchOnPremisesPublishing(ctx, object); err != nil {
		kiotaerrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully configured on-premises publishing for application: %s", object.ApplicationID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = onPremisesPublishingStateKnown

	err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for On-Premises Publishing resources.
//
// Operation: Retrieve on-premises publishing configuration
// API Calls:
//   - GET /applications/{application-id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-get?view=graph-rest-beta
func (r *OnPremisesPublishingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object OnPremisesPublishingResourceModel
	var identity sharedmodels.ResourceIdentity

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	operation := constants.TfOperationRead
	if ctxOp := ctx.Value("retry_operation"); ctxOp != nil {
		if opStr, ok := ctxOp.(string); ok {
			operation = opStr
		}
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	identity.ID = object.ApplicationID.ValueString()

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	application, err := r.client.
		Applications().
		ByApplicationId(object.ApplicationID.ValueString()).
		Get(ctx, applicationReadRequestConfiguration())

	if err != nil {
		kiotaerrors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	object = MapRemoteStateToTerraform(ctx, object, application)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for On-Premises Publishing resources.
//
// Operation: Update on-premises publishing configuration
// API Calls:
//   - PATCH /applications/{application-id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta
func (r *OnPremisesPublishingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OnPremisesPublishingResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if err := r.patchOnPremisesPublishing(ctx, plan); err != nil {
		kiotaerrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = onPremisesPublishingStateKnown

	err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation for On-Premises Publishing resources.
//
// Operation: Remove on-premises publishing configuration
// API Calls:
//   - PATCH /applications/{application-id} (with null/empty onPremisesPublishing)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/application-update?view=graph-rest-beta
func (r *OnPremisesPublishingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OnPremisesPublishingResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if err := r.patchRawApplication(ctx, data.ApplicationID.ValueString(), map[string]any{"onPremisesPublishing": nil}); err != nil {
		kiotaerrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

func (r *OnPremisesPublishingResource) patchOnPremisesPublishing(ctx context.Context, data OnPremisesPublishingResourceModel) error {
	return r.patchRawApplication(ctx, data.ApplicationID.ValueString(), constructOnPremisesPublishingPatchPayload(data))
}

func onPremisesPublishingStateKnown(ctx context.Context, state tfsdk.State) bool {
	var data OnPremisesPublishingResourceModel
	if diags := state.Get(ctx, &data); diags.HasError() {
		return false
	}

	stringValues := []types.String{
		data.ApplicationID,
		data.AlternateUrl,
		data.ApplicationServerTimeout,
		data.ApplicationType,
		data.ExternalAuthenticationType,
		data.InternalUrl,
		data.ExternalUrl,
		data.TrafficRoutingMethod,
		data.WafProvider,
	}
	for _, value := range stringValues {
		if value.IsUnknown() {
			return false
		}
	}

	boolValues := []types.Bool{
		data.IsAccessibleViaZTNAClient,
		data.IsBackendCertificateValidationEnabled,
		data.IsContinuousAccessEvaluationEnabled,
		data.IsDnsResolutionEnabled,
		data.IsHttpOnlyCookieEnabled,
		data.IsOnPremPublishingEnabled,
		data.IsPersistentCookieEnabled,
		data.IsSecureCookieEnabled,
		data.IsStateSessionEnabled,
		data.IsTranslateHostHeaderEnabled,
		data.IsTranslateLinksInBodyEnabled,
		data.UseAlternateUrlForTranslationAndRedirect,
	}
	for _, value := range boolValues {
		if value.IsUnknown() {
			return false
		}
	}

	return true
}

func applicationReadRequestConfiguration() *graphapplications.ApplicationItemRequestBuilderGetRequestConfiguration {
	// onPremisesPublishing is the managed property for this resource. Select it
	// explicitly because Graph application GET responses do not reliably include
	// this nested property in the default projection.
	return &graphapplications.ApplicationItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &graphapplications.ApplicationItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "onPremisesPublishing"},
		},
	}
}

func (r *OnPremisesPublishingResource) patchRawApplication(ctx context.Context, applicationID string, payload map[string]any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request payload: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("PATCH %s raw application payload", ResourceName), map[string]any{
		"application_id": applicationID,
		"json":           string(body),
	})

	requestInfo := abstractions.NewRequestInformation()
	requestInfo.Method = abstractions.PATCH
	requestInfo.UrlTemplate = "{+baseurl}/applications/{application%2Did}"
	requestInfo.PathParameters = map[string]string{
		"baseurl":          r.client.GetAdapter().GetBaseUrl(),
		"application%2Did": applicationID,
	}
	requestInfo.Headers.Add("Content-Type", "application/json")
	requestInfo.Headers.Add("Accept", "application/json")
	requestInfo.Content = body

	// The generated Application.Patch method serializes a top-level @odata.type,
	// and Graph rejects that payload when only updating onPremisesPublishing.
	// Use Kiota RequestInformation directly so authentication, middleware, and
	// OData error mapping still go through the SDK adapter while the JSON body
	// stays limited to the properties Graph accepts.
	errorMapping := abstractions.ErrorMappings{
		"XXX": grapherrors.CreateODataErrorFromDiscriminatorValue,
	}

	return r.client.GetAdapter().SendNoContent(ctx, requestInfo, errorMapping)
}

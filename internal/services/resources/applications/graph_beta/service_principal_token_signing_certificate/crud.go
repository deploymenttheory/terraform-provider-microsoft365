package graphBetaServicePrincipalTokenSigningCertificate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	// staleReadMaxAttempts and staleReadRetryDelay bound the wait for a consistent read
	// of the credential lists before the destructive full-list PATCH on delete.
	staleReadMaxAttempts = 12
	staleReadRetryDelay  = 5 * time.Second
)

// credentialPatchMutex serializes the GET -> filter -> PATCH critical section of Delete.
// The PATCH replaces the service principal's full credential lists, so two Deletes for
// different certificates on the same (or any) service principal running in one apply
// would otherwise clobber each other's removals with their earlier snapshots.
var credentialPatchMutex sync.Mutex

// Create handles the Create operation for the token signing certificate resource.
//
// Operation: Generates a self-signed token signing certificate on the service principal
// API Calls:
//   - POST /servicePrincipals/{servicePrincipalId}/addTokenSigningCertificate
//
// Reference: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-addtokensigningcertificate?view=graph-rest-beta
func (r *ServicePrincipalTokenSigningCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object ServicePrincipalTokenSigningCertificateResourceModel

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
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s", err.Error()),
		)
		return
	}

	servicePrincipalID := object.ServicePrincipalID.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Generating token signing certificate on service principal: %s", servicePrincipalID))

	certificate, err := r.client.
		ServicePrincipals().
		ByServicePrincipalId(servicePrincipalID).
		AddTokenSigningCertificate().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	MapSelfSignedCertificateToTerraform(ctx, &object, certificate)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = tokenSigningCertificateConsistencyPredicate(&object)

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for the token signing certificate resource.
//
// Operation: Retrieves the service principal and locates the signing key credential by key_id
// API Calls:
//   - GET /servicePrincipals/{servicePrincipalId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-get?view=graph-rest-beta
func (r *ServicePrincipalTokenSigningCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object ServicePrincipalTokenSigningCertificateResourceModel
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

	identity.ID = object.Id.ValueString()

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	servicePrincipalID := object.ServicePrincipalID.ValueString()

	servicePrincipal, err := r.client.
		ServicePrincipals().
		ByServicePrincipalId(servicePrincipalID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	found := MapRemoteResourceStateToTerraform(ctx, &object, servicePrincipal)

	// A missing key credential on a normal refresh means the certificate was removed
	// remotely. During create retries the id is cleared instead so the consistency
	// predicate fails and the read is retried until propagation completes.
	if !found {
		if operation == constants.TfOperationRead {
			tflog.Debug(ctx, fmt.Sprintf("Key credential %s not found on service principal %s, removing %s from state", object.KeyId.ValueString(), servicePrincipalID, ResourceName))
			resp.State.RemoveResource(ctx)
			return
		}
		object.Id = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for the token signing certificate resource.
//
// Operation: No API call - every Graph-backed attribute forces recreation, so Update is
// only invoked for provider-side changes such as the timeouts block, which are persisted
// to state directly.
func (r *ServicePrincipalTokenSigningCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object ServicePrincipalTokenSigningCertificateResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s (state-only, no API call)", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation for the token signing certificate resource.
//
// Operation: Removes the certificate's key credentials (Sign and Verify) and the associated
// password credential from the service principal.
// API Calls:
//   - GET /servicePrincipals/{servicePrincipalId} (to retrieve existing credentials)
//   - PATCH /servicePrincipals/{servicePrincipalId} (to update without the certificate's credentials)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-update?view=graph-rest-beta
func (r *ServicePrincipalTokenSigningCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object ServicePrincipalTokenSigningCertificateResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	servicePrincipalID := object.ServicePrincipalID.ValueString()
	thumbprint := object.Thumbprint.ValueString()
	keyId := object.KeyId.ValueString()
	url := r.httpClient.GetBaseURL() + r.ResourcePath + "/" + servicePrincipalID

	credentialPatchMutex.Lock()
	defer credentialPatchMutex.Unlock()

	// The credential lists are fetched and patched as raw JSON so retained
	// credentials round-trip byte-for-byte (see constructDeleteBody).
	//
	// Because the PATCH replaces the full credential lists, a read served by a stale
	// Entra replica would silently drop credentials added by concurrent operations
	// (e.g. the replacement certificate in a create_before_destroy rotation). The
	// same GET's preferredTokenSigningKeyThumbprint is used as a freshness marker:
	// while it still points at the certificate being deleted, the response predates
	// the activation of its successor (or the certificate is still active), so the
	// read is retried before any destructive PATCH is issued.
	var rawServicePrincipal map[string]any
	notFound := false
	for attempt := 1; ; attempt++ {
		rawServicePrincipal = nil

		getReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating HTTP request",
				fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		getResp, err := client.DoWithRetry(ctx, r.httpClient, getReq, 10)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error retrieving service principal credentials",
				fmt.Sprintf("Could not retrieve service principal: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		getBody, err := io.ReadAll(getResp.Body)
		getResp.Body.Close()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error reading response body",
				fmt.Sprintf("Could not read response body: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		if getResp.StatusCode == http.StatusNotFound {
			notFound = true
			break
		}
		if getResp.StatusCode >= 400 {
			resp.Diagnostics.AddError(
				"Error retrieving service principal credentials",
				fmt.Sprintf("Graph API returned status %d for %s: %s", getResp.StatusCode, ResourceName, string(getBody)),
			)
			return
		}

		if err := json.Unmarshal(getBody, &rawServicePrincipal); err != nil {
			resp.Diagnostics.AddError(
				"Error parsing service principal response",
				fmt.Sprintf("Could not parse response: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		// The marker only applies when the certificate's thumbprint is known; imported
		// state may lack it, in which case there is no freshness signal to wait on.
		preferredThumbprint, _ := rawServicePrincipal["preferredTokenSigningKeyThumbprint"].(string)
		if thumbprint == "" || !strings.EqualFold(preferredThumbprint, thumbprint) {
			break
		}

		if attempt >= staleReadMaxAttempts {
			resp.Diagnostics.AddError(
				"Token signing certificate is still the preferred signing key",
				fmt.Sprintf("After %d attempts, preferredTokenSigningKeyThumbprint on service principal %s still references the certificate being deleted (%s). "+
					"Either the certificate is still active — clear or repoint preferredTokenSigningKeyThumbprint first (e.g. via the "+
					"microsoft365_graph_beta_applications_service_principal_preferred_token_signing_key_thumbprint resource) — or Microsoft Entra "+
					"replication is delayed; retry the destroy. Removing the credentials of the active signing certificate would break SAML sign-in.",
					attempt, servicePrincipalID, thumbprint),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("preferredTokenSigningKeyThumbprint still references the certificate being deleted (attempt %d), waiting for replication before removing credentials for %s", attempt, ResourceName))
		select {
		case <-ctx.Done():
			resp.Diagnostics.AddError(
				"Context cancelled while waiting for replication",
				fmt.Sprintf("Timed out waiting for a consistent read before deleting %s", ResourceName),
			)
			return
		case <-time.After(staleReadRetryDelay):
		}
	}

	if notFound {
		tflog.Debug(ctx, fmt.Sprintf("Service principal %s not found during delete, nothing to remove", servicePrincipalID))
		return
	}

	patchBody, err := constructDeleteBody(ctx, rawServicePrincipal, keyId, thumbprint)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing delete request",
			fmt.Sprintf("Could not construct request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	patchReq, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(patchBody))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating HTTP request",
			fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	patchReq.Header.Set("Content-Type", "application/json")

	patchResp, err := client.DoWithRetry(ctx, r.httpClient, patchReq, 10)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error removing token signing certificate credentials",
			fmt.Sprintf("Could not update service principal: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	defer patchResp.Body.Close()

	if patchResp.StatusCode >= 400 {
		patchRespBody, _ := io.ReadAll(patchResp.Body)
		resp.Diagnostics.AddError(
			"Error removing token signing certificate credentials",
			fmt.Sprintf("Graph API returned status %d for %s: %s", patchResp.StatusCode, ResourceName, string(patchRespBody)),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully removed token signing certificate credentials from service principal: %s", servicePrincipalID))
	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

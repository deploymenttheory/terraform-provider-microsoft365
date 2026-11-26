package graphBetaAuthorizationPolicy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
)

// mapPermissionGrantPolicyIDs preserves configured prefix casing when Graph returns
// ManagePermissionGrants instead of managePermissionGrants. Policy IDs remain case-sensitive,
// and remote additions or removals are retained so Terraform can detect drift.
func mapPermissionGrantPolicyIDs(
	ctx context.Context,
	remote []string,
	configured types.Set,
) types.Set {
	configuredValues := make(map[string]string, len(configured.Elements()))
	for _, element := range configured.Elements() {
		if value, ok := element.(types.String); ok && !value.IsNull() && !value.IsUnknown() {
			configuredValues[permissionGrantPolicyKey(value.ValueString())] = value.ValueString()
		}
	}

	values := make([]string, 0, len(remote))
	for _, value := range remote {
		if configuredValue, ok := configuredValues[permissionGrantPolicyKey(value)]; ok {
			value = configuredValue
		}
		values = append(values, value)
	}
	return convert.GraphToFrameworkStringSetPreserveEmpty(ctx, values)
}

func permissionGrantPolicyKey(value string) string {
	prefix, id, found := strings.Cut(value, ".")
	if !found {
		return value
	}
	return strings.ToLower(prefix) + "." + id
}

// patchSettings retries only transient singleton replication misses through the
// common write helper. All other Graph errors remain immediately visible.
func (r *AuthorizationPolicyResource) patchSettings(ctx context.Context, body graphmodels.AuthorizationPolicyable) error {
	err := crud.WriteWithPreconditionRetry(ctx, func(ctx context.Context) error {
		_, err := r.client.Policies().AuthorizationPolicy().ByAuthorizationPolicyId(singletonID).Patch(ctx, body, nil)
		if err != nil {
			return fmt.Errorf("patch authorization policy: %w", err)
		}
		return nil
	}, crud.PreconditionRetryOptions{
		MaxRetries: 15, RetryInterval: 2 * time.Second, StatusCodes: []int{404},
		ResourceTypeName: ResourceName, Operation: constants.TfOperationUpdate,
	})
	if err != nil {
		return fmt.Errorf("write authorization policy settings: %w", err)
	}
	return nil
}

// waitForPolicySettings uses the normal Read mapping without replacing the desired
// final state with intermediate invitation settings returned by Graph.
func (r *AuthorizationPolicyResource) waitForPolicySettings(ctx context.Context, desired tfsdk.State) error {
	names := make([]string, 0, len(readbackAttributes))
	for _, name := range readbackAttributes {
		if name != "allow_invites_from" {
			names = append(names, name)
		}
	}
	match := readbackMatch
	match.Paths = crud.RootPaths(names...)
	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName,
		crud.ReadWithConsistencyPredicate(crud.RequireStableReads(crud.AttributesMatchPredicate(desired, match))))
	if deadline, ok := ctx.Deadline(); ok {
		opts.MaxRetries = int(time.Until(deadline) / opts.RetryInterval)
	}
	pending := &resource.UpdateResponse{State: desired}
	err := crud.ReadWithRetry(ctx, r.Read, resource.ReadRequest{State: desired},
		&crud.UpdateResponseContainer{UpdateResponse: pending}, opts)
	if err != nil {
		return fmt.Errorf("confirm authorization settings: %w", err)
	}
	return nil
}

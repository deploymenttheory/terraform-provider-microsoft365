package graphBetaNetworkExplicitForwardProxy

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
)

var (
	errInvalidProxyResponse = errors.New(
		"graph omitted required explicit proxy settings or returned unsupported session affinity options; existing state is retained",
	)
	errUnknownProxySettings = errors.New(
		"internet_access attributes must be known and non-null, with affinity enabled exactly when options is not none",
	)
)

func (r *NetworkExplicitForwardProxyResource) get(
	ctx context.Context,
) (*proxyResponse, error) {
	var response proxyResponse
	if err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.GET,
		r.ResourcePath,
		nil,
		nil,
		&response,
	); err != nil {
		return nil, fmt.Errorf("get explicit proxy: %w", err)
	}
	return &response, nil
}

func proxyPatch(
	plan NetworkExplicitForwardProxyResourceModel,
	state *NetworkExplicitForwardProxyResourceModel,
) (map[string]any, error) {
	p := plan.InternetAccess
	if p == nil || p.SourceIPSessionAffinityOptions.IsNull() ||
		p.SourceIPSessionAffinityOptions.IsUnknown() {
		return nil, errUnknownProxySettings
	}
	for _, value := range []types.Bool{p.IsEnabled, p.IsSourceIPSessionAffinityEnabled} {
		if value.IsNull() || value.IsUnknown() {
			return nil, errUnknownProxySettings
		}
	}
	if !validAffinity(
		p.IsSourceIPSessionAffinityEnabled.ValueBool(),
		p.SourceIPSessionAffinityOptions.ValueString(),
	) {
		return nil, errUnknownProxySettings
	}
	var previous *InternetAccessResourceModel
	if state != nil {
		previous = state.InternetAccess
	}
	access := map[string]any{}
	if previous == nil || !p.IsEnabled.Equal(previous.IsEnabled) {
		access["isEnabled"] = p.IsEnabled.ValueBool()
	}
	if previous == nil ||
		!p.IsSourceIPSessionAffinityEnabled.Equal(previous.IsSourceIPSessionAffinityEnabled) ||
		!p.SourceIPSessionAffinityOptions.Equal(previous.SourceIPSessionAffinityOptions) {
		access["isSourceIpSessionAffinityEnabled"] = p.IsSourceIPSessionAffinityEnabled.ValueBool()
		if p.IsSourceIPSessionAffinityEnabled.ValueBool() {
			access["sourceIpSessionAffinityOptions"] = p.SourceIPSessionAffinityOptions.ValueString()
		}
		// Disabling affinity normalizes options to none in Graph.
	}
	if len(access) == 0 {
		return map[string]any{}, nil
	}
	return map[string]any{"internetAccess": access}, nil
}

func (r *NetworkExplicitForwardProxyResource) patch(
	ctx context.Context,
	plan NetworkExplicitForwardProxyResourceModel,
	state *NetworkExplicitForwardProxyResourceModel,
) error {
	body, err := proxyPatch(plan, state)
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return nil
	}
	if err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.PATCH,
		r.ResourcePath,
		nil,
		body,
		nil,
	); err != nil {
		return fmt.Errorf("patch explicit proxy: %w", err)
	}
	return nil
}

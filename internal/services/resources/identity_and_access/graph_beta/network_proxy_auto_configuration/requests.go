package graphBetaNetworkProxyAutoConfiguration

import (
	"context"
	"errors"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
)

var (
	errInvalidPACResponse = errors.New(
		"graph omitted required custom PAC fields; existing state is retained",
	)
	errInvalidPACConfig = errors.New(
		"name, content and is_enabled must be explicitly configured; name and content must be nonempty",
	)
)

func pacBody(
	plan NetworkProxyAutoConfigurationResourceModel,
	state *NetworkProxyAutoConfigurationResourceModel,
) (map[string]any, error) {
	if plan.Name.IsNull() || plan.Name.IsUnknown() || plan.Name.ValueString() == "" ||
		plan.Content.IsNull() ||
		plan.Content.IsUnknown() ||
		plan.Content.ValueString() == "" ||
		plan.IsEnabled.IsNull() ||
		plan.IsEnabled.IsUnknown() {
		return nil, errInvalidPACConfig
	}
	body := map[string]any{}
	if state == nil || !plan.Name.Equal(state.Name) {
		body["name"] = plan.Name.ValueString()
	}
	if state == nil || !plan.Content.Equal(state.Content) {
		body["content"] = plan.Content.ValueString()
	}
	if state == nil || !plan.IsEnabled.Equal(state.IsEnabled) {
		body["isEnabled"] = plan.IsEnabled.ValueBool()
	}
	return body, nil
}

func (r *NetworkProxyAutoConfigurationResource) create(
	ctx context.Context,
	plan NetworkProxyAutoConfigurationResourceModel,
) (*pacResponse, error) {
	body, err := pacBody(plan, nil)
	if err != nil {
		return nil, err
	}
	var response pacResponse
	if err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.POST,
		r.ResourcePath,
		nil,
		body,
		&response,
	); err != nil {
		return nil, fmt.Errorf("create custom PAC: %w", err)
	}
	if response.ID == nil || *response.ID == "" {
		return nil, errInvalidPACResponse
	}
	return &response, nil
}

func (r *NetworkProxyAutoConfigurationResource) get(
	ctx context.Context,
	id string,
) (*pacResponse, error) {
	var response pacResponse
	if err := customrequests.JSONRequest(
		ctx,
		r.client.GetAdapter(),
		abstractions.GET,
		r.ResourcePath+"/{pacId}",
		map[string]string{"pacId": id},
		nil,
		&response,
	); err != nil {
		return nil, fmt.Errorf("get custom PAC: %w", err)
	}
	return &response, nil
}

func (r *NetworkProxyAutoConfigurationResource) patch(
	ctx context.Context,
	plan NetworkProxyAutoConfigurationResourceModel,
	state *NetworkProxyAutoConfigurationResourceModel,
) error {
	body, err := pacBody(plan, state)
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
		r.ResourcePath+"/{pacId}",
		map[string]string{"pacId": plan.ID.ValueString()},
		body,
		nil,
	); err != nil {
		return fmt.Errorf("update custom PAC: %w", err)
	}
	return nil
}

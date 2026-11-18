package graphBetaConditionalAccessPolicy

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
	"github.com/stretchr/testify/require"

	providerclient "github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
)

const testPolicyURL = "https://graph.microsoft.com/beta/identity/conditionalAccess/policies/11111111-1111-1111-1111-111111111111"

func TestUnitGetConditionalAccessPolicyKeepsEveryTimeStrict(t *testing.T) {
	transport := httpmock.NewMockTransport()
	transport.RegisterResponder("GET", testPolicyURL, httpmock.NewStringResponder(200, `{
		"id": "11111111-1111-1111-1111-111111111111",
		"displayName": "strict",
		"grantControls": null,
		"sessionControls": {"signInFrequency": {"isEnabled": true, "authenticationType": "primaryAndSecondaryAuthentication", "frequencyInterval": "everyTimeStrict"}}
	}`))
	client := providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient

	policy, body, err := getConditionalAccessPolicy(context.Background(), client, "11111111-1111-1111-1111-111111111111")
	require.NoError(t, err)
	require.Equal(t, "strict", *policy.GetDisplayName())
	require.Nil(t, policy.GetSessionControls().GetSignInFrequency().GetFrequencyInterval(), "the SDK enum drops everyTimeStrict")

	interval, err := rawSignInFrequencyInterval(body)
	require.NoError(t, err)
	require.Equal(t, signInFrequencyIntervalEveryTimeStrict, interval.ValueString())
}

func TestUnitGetConditionalAccessPolicyReturnsODataErrorOnNotFound(t *testing.T) {
	transport := httpmock.NewMockTransport()
	transport.RegisterResponder("GET", testPolicyURL, httpmock.NewJsonResponderOrPanic(404, map[string]any{
		"error": map[string]any{"code": "ResourceNotFound", "message": "Resource not found"},
	}))
	client := providerclient.NewMockGraphClients(&http.Client{Transport: transport}).MockBetaClient

	_, _, err := getConditionalAccessPolicy(context.Background(), client, "11111111-1111-1111-1111-111111111111")

	var odataErr *odataerrors.ODataError
	require.True(t, errors.As(err, &odataErr), "a 404 must surface as an OData error so Read removes the resource from state")
	require.Equal(t, 404, odataErr.ResponseStatusCode)
}

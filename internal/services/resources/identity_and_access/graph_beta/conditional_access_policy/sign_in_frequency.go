package graphBetaConditionalAccessPolicy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

// signInFrequencyIntervalEveryTimeStrict is the Graph beta value behind the "Every time strict (preview)"
// sign-in frequency option. The SDK's signInFrequencyInterval enum doesn't include it, so it's
// written through additional data and read from the raw response.
const signInFrequencyIntervalEveryTimeStrict = "everyTimeStrict"

var errInvalidPolicyResponse = errors.New("invalid conditional access policy response")

// setSignInFrequencyInterval writes frequencyInterval, bypassing the SDK enum for values it doesn't know.
func setSignInFrequencyInterval(value types.String, signInFrequency graphmodels.SignInFrequencySessionControlable) error {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	if value.ValueString() == signInFrequencyIntervalEveryTimeStrict {
		additionalData := signInFrequency.GetAdditionalData()
		if additionalData == nil {
			additionalData = map[string]any{}
		}
		additionalData["frequencyInterval"] = signInFrequencyIntervalEveryTimeStrict
		signInFrequency.SetAdditionalData(additionalData)
		return nil
	}

	parsed, err := graphmodels.ParseSignInFrequencyInterval(value.ValueString())
	if err != nil {
		return fmt.Errorf("parse frequency_interval: %w", err)
	}
	if interval, ok := parsed.(*graphmodels.SignInFrequencyInterval); ok {
		signInFrequency.SetFrequencyInterval(interval)
	}
	return nil
}

// getConditionalAccessPolicy reads the policy as raw bytes through the SDK adapter and deserializes it
// with the SDK model. The raw body is also returned so values missing from SDK enums (e.g.
// everyTimeStrict) can be read from the same response.
//
//nolint:wrapcheck // Preserve SDK errors for shared Graph diagnostics.
func getConditionalAccessPolicy(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, policyID string) (graphmodels.ConditionalAccessPolicyable, []byte, error) {
	requestInfo, err := client.
		Identity().
		ConditionalAccess().
		Policies().
		ByConditionalAccessPolicyId(policyID).
		ToGetRequestInformation(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	result, err := client.GetAdapter().SendPrimitive(ctx, requestInfo, "[]byte",
		abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue})
	if err != nil {
		return nil, nil, err
	}
	body, ok := result.([]byte)
	if !ok || len(body) == 0 {
		return nil, nil, errInvalidPolicyResponse
	}

	parseNode, err := jsonserialization.NewJsonParseNode(body)
	if err != nil {
		return nil, nil, fmt.Errorf("parse conditional access policy response: %w", err)
	}
	parsed, err := parseNode.GetObjectValue(graphmodels.CreateConditionalAccessPolicyFromDiscriminatorValue)
	if err != nil {
		return nil, nil, fmt.Errorf("deserialize conditional access policy response: %w", err)
	}
	policy, ok := parsed.(graphmodels.ConditionalAccessPolicyable)
	if !ok || policy == nil {
		return nil, nil, errInvalidPolicyResponse
	}
	return policy, body, nil
}

// rawSignInFrequencyInterval reads sessionControls.signInFrequency.frequencyInterval from the raw policy JSON.
func rawSignInFrequencyInterval(body []byte) (types.String, error) {
	var policy struct {
		SessionControls *struct {
			SignInFrequency *struct {
				FrequencyInterval *string `json:"frequencyInterval"`
			} `json:"signInFrequency"`
		} `json:"sessionControls"`
	}
	if err := json.Unmarshal(body, &policy); err != nil {
		return types.StringNull(), fmt.Errorf("parse conditional access policy response: %w", err)
	}

	if policy.SessionControls == nil || policy.SessionControls.SignInFrequency == nil ||
		policy.SessionControls.SignInFrequency.FrequencyInterval == nil {
		return types.StringNull(), nil
	}
	return types.StringValue(*policy.SessionControls.SignInFrequency.FrequencyInterval), nil
}

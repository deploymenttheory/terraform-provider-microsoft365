package graphBetaConditionalAccessPolicy

import (
	"context"
	"encoding/json"
	"fmt"

	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// signInFrequencyIntervalEveryTimeStrict is the Graph beta value behind the "Every time strict (preview)"
// sign-in frequency option. The SDK's signInFrequencyInterval enum doesn't include it, so it's
// written through additional data and read back from the raw response.
const signInFrequencyIntervalEveryTimeStrict = "everyTimeStrict"

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
		return err
	}
	if interval, ok := parsed.(*graphmodels.SignInFrequencyInterval); ok {
		signInFrequency.SetFrequencyInterval(interval)
	}
	return nil
}

// signInFrequencyIntervalDropped reports whether the SDK dropped frequencyInterval while deserializing,
// which happens when Graph returns a value missing from the SDK enum (e.g. everyTimeStrict).
func signInFrequencyIntervalDropped(remoteResource graphmodels.ConditionalAccessPolicyable) bool {
	sessionControls := remoteResource.GetSessionControls()
	if sessionControls == nil || sessionControls.GetSignInFrequency() == nil {
		return false
	}
	return sessionControls.GetSignInFrequency().GetFrequencyInterval() == nil
}

// fetchRawSignInFrequencyInterval reads sessionControls.signInFrequency.frequencyInterval from the raw policy JSON.
func fetchRawSignInFrequencyInterval(ctx context.Context, adapter abstractions.RequestAdapter, policyID string) (types.String, error) {
	body, err := customrequests.GetRequestByResourceId(ctx, adapter, customrequests.GetRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          "/identity/conditionalAccess/policies",
		ResourceIDPattern: "/id",
		ResourceID:        policyID,
	})
	if err != nil {
		return types.StringNull(), err
	}

	var policy struct {
		SessionControls *struct {
			SignInFrequency *struct {
				FrequencyInterval *string `json:"frequencyInterval"`
			} `json:"signInFrequency"`
		} `json:"sessionControls"`
	}
	if err := json.Unmarshal(body, &policy); err != nil {
		return types.StringNull(), fmt.Errorf("failed to parse conditional access policy response: %w", err)
	}

	if policy.SessionControls == nil || policy.SessionControls.SignInFrequency == nil ||
		policy.SessionControls.SignInFrequency.FrequencyInterval == nil {
		return types.StringNull(), nil
	}
	return types.StringValue(*policy.SessionControls.SignInFrequency.FrequencyInterval), nil
}

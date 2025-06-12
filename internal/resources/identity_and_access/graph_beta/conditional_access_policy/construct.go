// Main entry point to construct the conditional access root resource for the Terraform provider.
package graphBetaConditionalAccessPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the conditional access root resource for the Terraform provider.
func constructResource(ctx context.Context, data *ConditionalAccessPolicyResourceModel) (graphmodels.ConditionalAccessRootable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewConditionalAccessRoot()

	constructors.SetStringProperty(data.Id, requestBody.SetId)

	// Construct Authentication Context Class References
	if data.AuthenticationContextClassReferences != nil {
		authContextRefs, err := constructAuthenticationContextClassReferences(ctx, data.AuthenticationContextClassReferences)
		if err != nil {
			return nil, fmt.Errorf("failed to construct authentication context class references: %s", err)
		}
		requestBody.SetAuthenticationContextClassReferences(authContextRefs)
	}

	// Construct Authentication Strength
	if data.AuthenticationStrength != nil {
		authStrength, err := constructAuthenticationStrengthRoot(ctx, data.AuthenticationStrength)
		if err != nil {
			return nil, fmt.Errorf("failed to construct authentication strength: %s", err)
		}
		requestBody.SetAuthenticationStrength(authStrength)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructAuthenticationContextClassReferences constructs a slice of AuthenticationContextClassReference from the model
func constructAuthenticationContextClassReferences(ctx context.Context, data []*AuthenticationContextClassReferenceResourceModel) ([]graphmodels.AuthenticationContextClassReferenceable, error) {
	if data == nil {
		return nil, nil
	}

	references := make([]graphmodels.AuthenticationContextClassReferenceable, 0, len(data))

	for _, refData := range data {
		ref := graphmodels.NewAuthenticationContextClassReference()

		// Set basic properties using constructors package helpers
		constructors.SetStringProperty(refData.Id, ref.SetId)
		constructors.SetStringProperty(refData.Description, ref.SetDescription)
		constructors.SetStringProperty(refData.DisplayName, ref.SetDisplayName)
		constructors.SetBoolProperty(refData.IsAvailable, ref.SetIsAvailable)

		references = append(references, ref)
	}

	return references, nil
}

// constructAuthenticationStrengthRoot constructs AuthenticationStrengthRoot from the model
func constructAuthenticationStrengthRoot(ctx context.Context, data *AuthenticationStrengthRootResourceModel) (graphmodels.AuthenticationStrengthRootable, error) {
	authStrengthRoot := graphmodels.NewAuthenticationStrengthRoot()

	// Set basic properties
	constructors.SetStringProperty(data.Id, authStrengthRoot.SetId)

	// Set Authentication Combinations using SetObjectsFromStringSet
	if err := constructors.SetObjectsFromStringSet(ctx, data.AuthenticationCombinations,
		func(ctx context.Context, values []string) []graphmodels.AuthenticationMethodModes {
			combinations := make([]graphmodels.AuthenticationMethodModes, 0, len(values))
			for _, value := range values {
				if authMethod, err := graphmodels.ParseAuthenticationMethodModes(value); err == nil {
					if authMethodPtr, ok := authMethod.(*graphmodels.AuthenticationMethodModes); ok {
						combinations = append(combinations, *authMethodPtr)
					}
				}
			}
			return combinations
		}, authStrengthRoot.SetAuthenticationCombinations); err != nil {
		return nil, fmt.Errorf("failed to set authentication combinations: %s", err)
	}

	// Set Combinations (alternative property) using SetObjectsFromStringSet
	if err := constructors.SetObjectsFromStringSet(ctx, data.Combinations,
		func(ctx context.Context, values []string) []graphmodels.AuthenticationMethodModes {
			combinations := make([]graphmodels.AuthenticationMethodModes, 0, len(values))
			for _, value := range values {
				if authMethod, err := graphmodels.ParseAuthenticationMethodModes(value); err == nil {
					if authMethodPtr, ok := authMethod.(*graphmodels.AuthenticationMethodModes); ok {
						combinations = append(combinations, *authMethodPtr)
					}
				}
			}
			return combinations
		}, authStrengthRoot.SetCombinations); err != nil {
		return nil, fmt.Errorf("failed to set combinations: %s", err)
	}

	// Set Authentication Method Modes
	if data.AuthenticationMethodModes != nil {
		methodModes, err := constructAuthenticationMethodModeDetails(ctx, data.AuthenticationMethodModes)
		if err != nil {
			return nil, fmt.Errorf("failed to construct authentication method modes: %s", err)
		}
		authStrengthRoot.SetAuthenticationMethodModes(methodModes)
	}

	// Set Policies
	if data.Policies != nil {
		policies, err := constructAuthenticationStrengthPolicies(ctx, data.Policies)
		if err != nil {
			return nil, fmt.Errorf("failed to construct authentication strength policies: %s", err)
		}
		authStrengthRoot.SetPolicies(policies)
	}

	return authStrengthRoot, nil
}

// constructAuthenticationMethodModeDetails constructs AuthenticationMethodModeDetail slice from the model
func constructAuthenticationMethodModeDetails(ctx context.Context, data []*AuthenticationMethodModeDetailResourceModel) ([]graphmodels.AuthenticationMethodModeDetailable, error) {
	if data == nil {
		return nil, nil
	}

	methodModes := make([]graphmodels.AuthenticationMethodModeDetailable, 0, len(data))

	for _, modeData := range data {
		mode := graphmodels.NewAuthenticationMethodModeDetail()

		// Set basic properties using constructors package helpers
		constructors.SetStringProperty(modeData.Id, mode.SetId)
		constructors.SetStringProperty(modeData.DisplayName, mode.SetDisplayName)

		// Set Authentication Method using SetBitmaskEnumProperty
		if err := constructors.SetBitmaskEnumProperty(modeData.AuthenticationMethod,
			graphmodels.ParseBaseAuthenticationMethod, mode.SetAuthenticationMethod); err != nil {
			return nil, fmt.Errorf("failed to set authentication method: %s", err)
		}

		methodModes = append(methodModes, mode)
	}

	return methodModes, nil
}

// constructAuthenticationStrengthPolicies constructs AuthenticationStrengthPolicy slice from the model
func constructAuthenticationStrengthPolicies(ctx context.Context, data []*AuthenticationStrengthPolicyResourceModel) ([]graphmodels.AuthenticationStrengthPolicyable, error) {
	if data == nil {
		return nil, nil
	}

	policies := make([]graphmodels.AuthenticationStrengthPolicyable, 0, len(data))

	for _, policyData := range data {
		// Note: AuthenticationStrengthPolicy structure not provided in documents
		// This is a placeholder implementation
		tflog.Warn(ctx, "AuthenticationStrengthPolicy construction not fully implemented", map[string]interface{}{
			"policy_id": policyData.Id.ValueString(),
		})
	}

	return policies, nil
}

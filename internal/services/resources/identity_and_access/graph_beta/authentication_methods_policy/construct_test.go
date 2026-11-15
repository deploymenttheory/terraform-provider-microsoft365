package graphBetaIdentityAndAccessAuthenticationMethodsPolicy

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// unknownModel returns a plan in which every optional and computed value is unknown, as on create with an empty configuration.
func unknownModel() AuthenticationMethodsPolicyResourceModel {
	return AuthenticationMethodsPolicyResourceModel{
		ID:                               types.StringUnknown(),
		DisplayName:                      types.StringUnknown(),
		Description:                      types.StringUnknown(),
		PolicyMigrationState:             types.StringUnknown(),
		RegistrationEnforcement:          types.ObjectUnknown(registrationEnforcementAttrTypes),
		ReportSuspiciousActivitySettings: types.ObjectUnknown(reportSuspiciousActivitySettingsAttrTypes),
		SystemCredentialPreferences:      types.ObjectUnknown(systemCredentialPreferencesAttrTypes),
		Email:                            types.ObjectUnknown(emailMethodAttrTypes),
		TemporaryAccessPass:              types.ObjectUnknown(temporaryAccessPassMethodAttrTypes),
		Sms:                              types.ObjectUnknown(smsMethodAttrTypes),
		SoftwareOath:                     types.ObjectUnknown(oathMethodAttrTypes),
		HardwareOath:                     types.ObjectUnknown(oathMethodAttrTypes),
	}
}

func groupMethodTargets(ids ...string) types.Set {
	elements := make([]attr.Value, 0, len(ids))
	for _, id := range ids {
		elements = append(elements, types.ObjectValueMust(methodTargetAttrTypes, map[string]attr.Value{
			"id":                       types.StringValue(id),
			"target_type":              types.StringValue("group"),
			"is_registration_required": types.BoolValue(false),
		}))
	}
	return types.SetValueMust(types.ObjectType{AttrTypes: methodTargetAttrTypes}, elements)
}

func emptyExcludeTargets() types.Set {
	return types.SetValueMust(types.ObjectType{AttrTypes: targetAttrTypes}, []attr.Value{})
}

func emailObject(include types.Set, exclude types.Set) types.Object {
	return types.ObjectValueMust(emailMethodAttrTypes, map[string]attr.Value{
		"state":                              types.StringValue("enabled"),
		"allow_external_id_to_use_email_otp": types.StringUnknown(),
		"include_targets":                    include,
		"exclude_targets":                    exclude,
	})
}

func temporaryAccessPassObject(length int32) types.Object {
	return types.ObjectValueMust(temporaryAccessPassMethodAttrTypes, map[string]attr.Value{
		"state":                       types.StringValue("enabled"),
		"default_lifetime_in_minutes": types.Int32Value(60),
		"default_length":              types.Int32Value(length),
		"minimum_lifetime_in_minutes": types.Int32Value(60),
		"maximum_lifetime_in_minutes": types.Int32Value(480),
		"is_usable_once":              types.BoolValue(false),
		"include_targets":             groupMethodTargets("all_users"),
		"exclude_targets":             emptyExcludeTargets(),
	})
}

func TestUnitConstructResource_CreateWithNothingConfiguredSendsNothing(t *testing.T) {
	plan := unknownModel()

	bodies, err := constructResource(context.Background(), &plan, nil)

	require.NoError(t, err)
	assert.Nil(t, bodies.Policy)
	assert.Empty(t, bodies.Methods)
}

func TestUnitConstructResource_CreateSendsOnlyConfiguredMethods(t *testing.T) {
	plan := unknownModel()
	plan.Email = emailObject(groupMethodTargets("33333333-3333-3333-3333-333333333333"), types.SetUnknown(types.ObjectType{AttrTypes: targetAttrTypes}))

	bodies, err := constructResource(context.Background(), &plan, nil)

	require.NoError(t, err)
	assert.Nil(t, bodies.Policy)
	require.Len(t, bodies.Methods, 1)
	assert.Equal(t, methodIDEmail, bodies.Methods[0].ID)

	email, ok := bodies.Methods[0].Body.(*graphmodels.EmailAuthenticationMethodConfiguration)
	require.True(t, ok)
	require.NotNil(t, email.GetOdataType())
	assert.Equal(t, "#microsoft.graph.emailAuthenticationMethodConfiguration", *email.GetOdataType())
	assert.Equal(t, graphmodels.ENABLED_AUTHENTICATIONMETHODSTATE, *email.GetState())
	assert.Nil(t, email.GetAllowExternalIdToUseEmailOtp(), "an unknown value must be omitted")
	assert.Nil(t, email.GetExcludeTargets(), "an unknown collection must be omitted, not sent as null")
	require.Len(t, email.GetIncludeTargets(), 1)
	assert.Equal(t, "33333333-3333-3333-3333-333333333333", *email.GetIncludeTargets()[0].GetId())
	assert.Equal(t, graphmodels.GROUP_AUTHENTICATIONMETHODTARGETTYPE, *email.GetIncludeTargets()[0].GetTargetType())
}

func TestUnitConstructResource_EmptyCollectionIsSentAsEmptyList(t *testing.T) {
	plan := unknownModel()
	plan.Email = emailObject(groupMethodTargets(), emptyExcludeTargets())

	bodies, err := constructResource(context.Background(), &plan, nil)

	require.NoError(t, err)
	require.Len(t, bodies.Methods, 1)
	email, ok := bodies.Methods[0].Body.(*graphmodels.EmailAuthenticationMethodConfiguration)
	require.True(t, ok)
	assert.NotNil(t, email.GetIncludeTargets())
	assert.Empty(t, email.GetIncludeTargets())
	assert.NotNil(t, email.GetExcludeTargets())
	assert.Empty(t, email.GetExcludeTargets())
}

func TestUnitConstructResource_UpdateSendsOnlyChangedMethods(t *testing.T) {
	state := unknownModel()
	state.ID = types.StringValue(singletonID)
	state.PolicyMigrationState = types.StringValue("preMigration")
	state.Email = emailObject(groupMethodTargets("33333333-3333-3333-3333-333333333333"), emptyExcludeTargets())
	state.TemporaryAccessPass = temporaryAccessPassObject(8)

	plan := state
	plan.TemporaryAccessPass = temporaryAccessPassObject(12)

	bodies, err := constructResource(context.Background(), &plan, &state)

	require.NoError(t, err)
	assert.Nil(t, bodies.Policy, "unchanged policy-level settings must not be sent")
	require.Len(t, bodies.Methods, 1)
	assert.Equal(t, methodIDTemporaryAccessPass, bodies.Methods[0].ID)
	temporaryAccessPass, ok := bodies.Methods[0].Body.(*graphmodels.TemporaryAccessPassAuthenticationMethodConfiguration)
	require.True(t, ok)
	assert.Equal(t, int32(12), *temporaryAccessPass.GetDefaultLength())
}

func TestUnitConstructResource_UpdateWithNoChangesSendsNothing(t *testing.T) {
	state := unknownModel()
	state.ID = types.StringValue(singletonID)
	state.Email = emailObject(groupMethodTargets("33333333-3333-3333-3333-333333333333"), emptyExcludeTargets())
	plan := state

	bodies, err := constructResource(context.Background(), &plan, &state)

	require.NoError(t, err)
	assert.Nil(t, bodies.Policy)
	assert.Empty(t, bodies.Methods)
}

func TestUnitConstructResource_PolicyMigrationStateChange(t *testing.T) {
	state := unknownModel()
	state.ID = types.StringValue(singletonID)
	state.PolicyMigrationState = types.StringValue("preMigration")
	plan := state
	plan.PolicyMigrationState = types.StringValue("migrationComplete")

	bodies, err := constructResource(context.Background(), &plan, &state)

	require.NoError(t, err)
	require.NotNil(t, bodies.Policy)
	assert.Equal(t, graphmodels.MIGRATIONCOMPLETE_AUTHENTICATIONMETHODSPOLICYMIGRATIONSTATE, *bodies.Policy.GetPolicyMigrationState())
	assert.Nil(t, bodies.Policy.GetRegistrationEnforcement())
	assert.Empty(t, bodies.Methods)
}

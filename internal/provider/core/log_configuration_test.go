package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-log/tflogtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_RedactedProviderConfiguration(t *testing.T) {
	ctx := context.Background()
	var response frameworkprovider.SchemaResponse
	(&M365Provider{}).Schema(ctx, frameworkprovider.SchemaRequest{}, &response)
	secrets := make(map[string]string)
	var populate func(map[string]schema.Attribute, string) types.Object
	populate = func(attributes map[string]schema.Attribute, prefix string) types.Object {
		values := make(map[string]attr.Value)
		for name, attribute := range attributes {
			if attribute.IsSensitive() {
				secret := prefix + name + "-secret-\"\\\n"
				secrets[prefix+name] = secret
				values[name] = types.StringValue(secret)
				continue
			}
			switch attribute := attribute.(type) {
			case schema.SingleNestedAttribute:
				values[name] = populate(attribute.Attributes, name+".")
			case schema.StringAttribute:
				values[name] = types.StringValue("visible-" + name)
			case schema.BoolAttribute:
				values[name] = types.BoolValue(true)
			case schema.Int64Attribute:
				values[name] = types.Int64Value(3)
			case schema.ListAttribute:
				values[name] = types.ListValueMust(
					types.StringType,
					[]attr.Value{types.StringValue("visible-tenant")},
				)
			default:
				t.Fatalf("unsupported provider attribute %s", name)
			}
		}
		return types.ObjectValueMust(schemaToAttrTypes(attributes), values)
	}
	object := populate(response.Schema.Attributes, "")
	var model M365ProviderModel
	diags := object.As(ctx, &model, basetypes.ObjectAsOptions{})
	require.False(t, diags.HasError(), "%v", diags)
	beforeEntra := model.EntraIDOptions.String()
	beforeClient := model.ClientOptions.String()
	beforeTenant := model.TenantID.ValueString()
	redacted := redactedProviderConfiguration(ctx, model)
	var output bytes.Buffer
	logContext := tflogtest.RootLogger(ctx, &output)
	tflog.Debug(logContext, "Provider configuration completed", map[string]any{"config": redacted})
	encoded := output.Bytes()
	var logged map[string]any
	require.NoError(t, json.Unmarshal(encoded, &logged))
	assert.Contains(t, string(encoded), "visible-cloud")
	for name, secret := range secrets {
		assert.NotContains(t, string(encoded), secret, name)
		assert.NotContains(t, string(encoded), name+"-secret", name)
	}
	assert.Equal(t, "[REDACTED]", redacted["tenant_id"])
	assert.Equal(t, "visible-cloud", redacted["cloud"])
	assert.Equal(t, true, redacted["debug_mode"])
	entra := redacted["entra_id_options"].(map[string]any)
	assert.Equal(t, "[REDACTED]", entra["client_certificate_base64"])
	assert.Equal(t, true, entra["send_certificate_chain"])
	options := redacted["client_options"].(map[string]any)
	assert.Equal(t, "[REDACTED]", options["proxy_password"])
	assert.Equal(t, int64(3), options["max_retries"])
	assert.Equal(t, beforeEntra, model.EntraIDOptions.String(), "redaction must not change credentials")
	assert.Equal(t, beforeClient, model.ClientOptions.String(), "redaction must not change client options")
	assert.Equal(t, beforeTenant, model.TenantID.ValueString())
}

func TestUnit_RedactedProviderConfiguration_NullAndUnknown(t *testing.T) {
	for _, value := range []types.Object{
		types.ObjectNull(schemaToAttrTypes(EntraIDOptionsSchema())),
		types.ObjectUnknown(schemaToAttrTypes(EntraIDOptionsSchema())),
	} {
		result := redactedProviderConfiguration(context.Background(), M365ProviderModel{
			EntraIDOptions: value,
			ClientOptions:  types.ObjectNull(schemaToAttrTypes(ClientOptionsSchema())),
		})
		if value.IsNull() {
			assert.Nil(t, result["entra_id_options"])
		} else {
			assert.Equal(t, "[unknown]", result["entra_id_options"])
		}
	}
}

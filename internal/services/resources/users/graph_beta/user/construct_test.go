package graphBetaUsersUser

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
)

func serializeUser(t *testing.T, data *UserResourceModel) string {
	t.Helper()
	user, err := constructResource(context.Background(), data)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}
	writer := jsonserialization.NewJsonSerializationWriter()
	if err := writer.WriteObjectValue("", user); err != nil {
		t.Fatalf("serialize user: %v", err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatalf("read serialized user: %v", err)
	}
	return string(content)
}

// Graph rejects null for these non-nullable collections, so unset values must be left out.
func TestConstructResourceOmitsUnsetStringCollections(t *testing.T) {
	data := &UserResourceModel{
		DisplayName:    types.StringValue("Minimal User"),
		BusinessPhones: types.SetUnknown(types.StringType),
		OtherMails:     types.SetNull(types.StringType),
		ProxyAddresses: types.SetUnknown(types.StringType),
	}

	body := serializeUser(t, data)
	for _, field := range []string{"businessPhones", "otherMails", "proxyAddresses"} {
		if strings.Contains(body, `"`+field+`"`) {
			t.Errorf("expected %s to be omitted, got %s", field, body)
		}
	}
}

func TestConstructResourceSendsEmptyStringCollections(t *testing.T) {
	empty := types.SetValueMust(types.StringType, []attr.Value{})
	data := &UserResourceModel{
		DisplayName:    types.StringValue("Minimal User"),
		BusinessPhones: empty,
		OtherMails:     empty,
		ProxyAddresses: empty,
	}

	body := serializeUser(t, data)
	for _, field := range []string{"businessPhones", "otherMails", "proxyAddresses"} {
		if !strings.Contains(body, `"`+field+`":[]`) {
			t.Errorf("expected %s to be sent as [], got %s", field, body)
		}
	}
}

func TestConstructResourceRejectsNonStringCollections(t *testing.T) {
	numbers := types.SetValueMust(types.Int64Type, []attr.Value{types.Int64Value(1)})
	cases := map[string]func(*UserResourceModel){
		"business phones": func(data *UserResourceModel) { data.BusinessPhones = numbers },
		"other mails":     func(data *UserResourceModel) { data.OtherMails = numbers },
		"proxy addresses": func(data *UserResourceModel) { data.ProxyAddresses = numbers },
	}

	for name, apply := range cases {
		t.Run(name, func(t *testing.T) {
			data := &UserResourceModel{DisplayName: types.StringValue("Minimal User")}
			apply(data)
			_, err := constructResource(context.Background(), data)
			if err == nil || !strings.Contains(err.Error(), "error converting "+name) {
				t.Fatalf("expected error converting %s, got %v", name, err)
			}
		})
	}
}

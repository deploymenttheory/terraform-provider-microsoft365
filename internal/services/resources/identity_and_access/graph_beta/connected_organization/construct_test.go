package graphBetaIdentityAndAccessConnectedOrganization

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func serializeConnectedOrganization(t *testing.T, body graphmodels.ConnectedOrganizationable) string {
	t.Helper()
	writer := jsonserialization.NewJsonSerializationWriter()
	if err := writer.WriteObjectValue("", body); err != nil {
		t.Fatalf("serialize connected organization: %v", err)
	}
	content, err := writer.GetSerializedContent()
	if err != nil {
		t.Fatalf("read serialized connected organization: %v", err)
	}
	return string(content)
}

func domainSource(domain string) IdentitySourceModel {
	return IdentitySourceModel{
		ODataType:   types.StringValue(ODataTypeDomainIdentitySource),
		DisplayName: types.StringValue(domain),
		DomainName:  types.StringValue(domain),
		TenantID:    types.StringNull(),
		IssuerURI:   types.StringNull(),
	}
}

func minimalModel() *ConnectedOrganizationResourceModel {
	return &ConnectedOrganizationResourceModel{
		DisplayName:     types.StringValue("connected organization"),
		Description:     types.StringNull(),
		State:           types.StringValue("configured"),
		IdentitySources: []IdentitySourceModel{domainSource("example.com")},
	}
}

func assertOmitted(t *testing.T, body string, fields ...string) {
	t.Helper()
	for _, field := range fields {
		if strings.Contains(body, `"`+field+`"`) {
			t.Errorf("expected %s to be omitted, got %s", field, body)
		}
	}
}

func assertPresent(t *testing.T, body string, fragments ...string) {
	t.Helper()
	for _, fragment := range fragments {
		if !strings.Contains(body, fragment) {
			t.Errorf("expected %s in body, got %s", fragment, body)
		}
	}
}

// Unset optional fields are omitted on create, never sent as null.
func TestConstructCreateOmitsUnsetFields(t *testing.T) {
	body, err := constructResource(context.Background(), minimalModel(), false)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	json := serializeConnectedOrganization(t, body)
	assertOmitted(t, json, "description", "tenantId", "issuerUri")
	assertPresent(t, json, `"displayName":"connected organization"`, `"state":"configured"`, `"identitySources":[`, `"@odata.type":"#microsoft.graph.domainIdentitySource"`, `"domainName":"example.com"`)
	if strings.Contains(json, "null") {
		t.Errorf("expected no null values, got %s", json)
	}
}

// Unknown optional values are also omitted on create.
func TestConstructCreateOmitsUnknownFields(t *testing.T) {
	data := minimalModel()
	data.Description = types.StringUnknown()
	data.State = types.StringUnknown()

	body, err := constructResource(context.Background(), data, false)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	assertOmitted(t, serializeConnectedOrganization(t, body), "description", "state")
}

// Update never sends identity sources; they are create-only and force replacement.
func TestConstructUpdateOmitsIdentitySources(t *testing.T) {
	data := minimalModel()
	data.Description = types.StringValue("description")

	body, err := constructResource(context.Background(), data, true)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	json := serializeConnectedOrganization(t, body)
	assertOmitted(t, json, "identitySources")
	assertPresent(t, json, `"description":"description"`)
}

// Removing the description on update clears it with an empty string rather than null.
func TestConstructUpdateClearsDescription(t *testing.T) {
	body, err := constructResource(context.Background(), minimalModel(), true)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	json := serializeConnectedOrganization(t, body)
	assertPresent(t, json, `"description":""`)
	if strings.Contains(json, "null") {
		t.Errorf("expected no null values, got %s", json)
	}
}

// Each identity source type only carries its own attributes.
func TestConstructIdentitySourceTypes(t *testing.T) {
	sources := []IdentitySourceModel{
		{
			ODataType:   types.StringValue(ODataTypeAzureActiveDirectoryTenant),
			DisplayName: types.StringValue("Contoso"),
			DomainName:  types.StringNull(),
			TenantID:    types.StringValue("11111111-1111-1111-1111-111111111111"),
			IssuerURI:   types.StringNull(),
		},
		{
			ODataType:   types.StringValue(ODataTypeExternalDomainFederation),
			DisplayName: types.StringValue("Fabrikam"),
			DomainName:  types.StringValue("fabrikam.example"),
			TenantID:    types.StringNull(),
			IssuerURI:   types.StringValue("https://idp.fabrikam.example/saml"),
		},
	}

	data := minimalModel()
	data.IdentitySources = sources
	body, err := constructResource(context.Background(), data, false)
	if err != nil {
		t.Fatalf("constructResource returned error: %v", err)
	}

	json := serializeConnectedOrganization(t, body)
	assertPresent(t, json,
		`"@odata.type":"#microsoft.graph.azureActiveDirectoryTenant"`,
		`"tenantId":"11111111-1111-1111-1111-111111111111"`,
		`"@odata.type":"#microsoft.graph.externalDomainFederation"`,
		`"issuerUri":"https://idp.fabrikam.example/saml"`,
	)
	if strings.Contains(json, "null") {
		t.Errorf("expected no null values, got %s", json)
	}
}

// Identity sources whose attributes do not match their type are rejected before any request.
func TestConstructRejectsMismatchedIdentitySource(t *testing.T) {
	cases := map[string]struct {
		source IdentitySourceModel
		want   string
	}{
		"domain without domain_name": {
			source: IdentitySourceModel{ODataType: types.StringValue(ODataTypeDomainIdentitySource), DisplayName: types.StringValue("x"), DomainName: types.StringNull(), TenantID: types.StringNull(), IssuerURI: types.StringNull()},
			want:   "domain_name is required",
		},
		"domain with tenant_id": {
			source: IdentitySourceModel{ODataType: types.StringValue(ODataTypeDomainIdentitySource), DisplayName: types.StringValue("x"), DomainName: types.StringValue("example.com"), TenantID: types.StringValue("11111111-1111-1111-1111-111111111111"), IssuerURI: types.StringNull()},
			want:   "tenant_id is not supported",
		},
		"tenant without tenant_id": {
			source: IdentitySourceModel{ODataType: types.StringValue(ODataTypeAzureActiveDirectoryTenant), DisplayName: types.StringValue("x"), DomainName: types.StringNull(), TenantID: types.StringNull(), IssuerURI: types.StringNull()},
			want:   "tenant_id is required",
		},
		"federation without issuer_uri": {
			source: IdentitySourceModel{ODataType: types.StringValue(ODataTypeExternalDomainFederation), DisplayName: types.StringValue("x"), DomainName: types.StringValue("fabrikam.example"), TenantID: types.StringNull(), IssuerURI: types.StringNull()},
			want:   "issuer_uri is required",
		},
		"unsupported type": {
			source: IdentitySourceModel{ODataType: types.StringValue("#microsoft.graph.socialIdentitySource"), DisplayName: types.StringValue("x"), DomainName: types.StringNull(), TenantID: types.StringNull(), IssuerURI: types.StringNull()},
			want:   "unsupported odata_type",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			data := minimalModel()
			data.IdentitySources = []IdentitySourceModel{tc.source}
			_, err := constructResource(context.Background(), data, false)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected error containing %q, got %v", tc.want, err)
			}
		})
	}
}

// An empty description from Graph maps to null so a cleared description does not drift.
func TestMapRemoteStateEmptyDescriptionIsNull(t *testing.T) {
	remote := graphmodels.NewConnectedOrganization()
	id := "00000000-0000-0000-0000-000000000001"
	empty := ""
	remote.SetId(&id)
	remote.SetDescription(&empty)

	state := MapRemoteResourceStateToTerraform(context.Background(), ConnectedOrganizationResourceModel{}, remote)
	if !state.Description.IsNull() {
		t.Fatalf("expected null description, got %q", state.Description.ValueString())
	}
	if state.IdentitySources != nil {
		t.Fatalf("expected no identity sources, got %v", state.IdentitySources)
	}
}

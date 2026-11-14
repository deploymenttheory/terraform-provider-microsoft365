package graphBetaIdentityAndAccessConnectedOrganization

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	errUnsupportedIdentitySourceType = errors.New("unsupported odata_type")
	errIdentitySourceFieldRequired   = errors.New("is required")
	errIdentitySourceFieldNotAllowed = errors.New("is not supported")
)

// identitySourceFields lists, per identity source type, the type-specific attributes that must be set.
// Any type-specific attribute not listed for a type must be left unset.
var identitySourceFields = map[string][]string{
	ODataTypeDomainIdentitySource:       {"domain_name"},
	ODataTypeAzureActiveDirectoryTenant: {"tenant_id"},
	ODataTypeExternalDomainFederation:   {"domain_name", "issuer_uri"},
}

// ValidateConfig reports identity sources whose type-specific attributes do not match their odata_type at plan time.
func (r *ConnectedOrganizationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data ConnectedOrganizationResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for i, source := range data.IdentitySources {
		if err := validateIdentitySource(source); err != nil {
			resp.Diagnostics.AddAttributeError(path.Root("identity_sources"), "Invalid identity source", fmt.Sprintf("identity_sources[%d]: %s", i, err.Error()))
		}
	}
}

// validateRequest validates the connected organization request before it is sent to Microsoft Graph
func validateRequest(ctx context.Context, data *ConnectedOrganizationResourceModel) error {
	tflog.Debug(ctx, "Starting connected organization request validation")

	for i, source := range data.IdentitySources {
		if err := validateIdentitySource(source); err != nil {
			return fmt.Errorf("identity_sources[%d]: %w", i, err)
		}
	}

	tflog.Debug(ctx, "Connected organization request validation completed successfully")
	return nil
}

// validateIdentitySource checks that an identity source sets exactly the attributes its odata_type uses.
// Unknown values are treated as set so validation does not fail before dependencies are resolved.
func validateIdentitySource(source IdentitySourceModel) error {
	if source.ODataType.IsNull() || source.ODataType.IsUnknown() {
		return nil
	}

	odataType := source.ODataType.ValueString()
	required, ok := identitySourceFields[odataType]
	if !ok {
		return fmt.Errorf("%w %q", errUnsupportedIdentitySourceType, odataType)
	}

	values := map[string]types.String{
		"domain_name": source.DomainName,
		"tenant_id":   source.TenantID,
		"issuer_uri":  source.IssuerURI,
	}

	for _, name := range []string{"domain_name", "tenant_id", "issuer_uri"} {
		isRequired := false
		for _, field := range required {
			if field == name {
				isRequired = true
				break
			}
		}

		if isRequired && values[name].IsNull() {
			return fmt.Errorf("%s %w when odata_type is %q", name, errIdentitySourceFieldRequired, odataType)
		}
		if !isRequired && !values[name].IsNull() {
			return fmt.Errorf("%s %w when odata_type is %q", name, errIdentitySourceFieldNotAllowed, odataType)
		}
	}

	return nil
}

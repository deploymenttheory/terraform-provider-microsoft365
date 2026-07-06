package graphBetaNetworkWebContentFilteringPolicyRule

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
)

// constructResource builds the portal-observed webFilteringRule payload.
//
// Microsoft Learn documents the generic filteringRule hierarchy, but not
// /networkaccess/webFilteringPolicies/{id}/policyRules or the webFiltering*
// OData discriminator values used by the Entra Global Secure Access Web content
// filtering blade:
// https://learn.microsoft.com/graph/api/resources/networkaccess-filteringrule
//
// The shape below intentionally follows captured portal traffic:
//   - URL/FQDN and category destinations are sibling entries in
//     matchingConditions.destinations.targets.
//   - URL/FQDN destinations are sent to Graph as a values array. The portal
//     offers a comma-delimited text box for usability, but Terraform keeps the
//     API shape as set(string).
//   - HTTP methods and session types are serialized back to comma-separated
//     strings because that is what the endpoint returns and accepts.
//   - Custom request header insertions are nested under
//     action.headerSettings.modifications for allow rules.
func constructResource(ctx context.Context, data *NetworkWebContentFilteringPolicyRuleResourceModel) (s.Parsable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	urlsOrFqdns, err := stringSetValues(ctx, data.UrlsOrFqdns)
	if err != nil {
		return nil, fmt.Errorf("failed to read urls_or_fqdns: %w", err)
	}
	webCategories, err := stringSetValues(ctx, data.WebCategories)
	if err != nil {
		return nil, fmt.Errorf("failed to read web_categories: %w", err)
	}
	httpMethods, err := stringSetValues(ctx, data.HTTPMethods)
	if err != nil {
		return nil, fmt.Errorf("failed to read http_methods: %w", err)
	}
	sessionTypes, err := stringSetValues(ctx, data.SessionTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to read session_types: %w", err)
	}
	customHeaders, err := customHeaderValues(ctx, data.CustomHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to read custom_headers: %w", err)
	}
	if data.Action.ValueString() != "allow" && len(customHeaders) > 0 {
		return nil, fmt.Errorf("custom_headers can only be used when action is allow")
	}
	if len(urlsOrFqdns) == 0 && len(webCategories) == 0 {
		return nil, fmt.Errorf("at least one destination must be specified using urls_or_fqdns or web_categories")
	}

	// The Entra portal and Graph endpoint reject rules without a destination.
	// Portal UI text asks for at least one URL/FQDN or web category. Terraform
	// follows the Graph payload shape and models URL/FQDN values as set(string)
	// even though the portal combines them into one comma-delimited text box.
	targets := make([]s.Parsable, 0, 2)
	if len(urlsOrFqdns) > 0 {
		targets = append(targets, &destinationTargetRequestBody{
			odataType: webFilteringURLDestinationODataType,
			values:    urlsOrFqdns,
		})
	}
	if len(webCategories) > 0 {
		targets = append(targets, &destinationTargetRequestBody{
			odataType: webFilteringWebCategoryDestinationODataType,
			values:    webCategories,
		})
	}

	action := &actionRequestBody{
		odataType: graphActionODataType(data.Action.ValueString()),
	}
	if len(customHeaders) > 0 {
		action.headerSettings = &headerSettingsRequestBody{
			modifications: customHeaders,
		}
	}

	requestBody := &webContentFilteringPolicyRuleRequestBody{
		odataType:   webFilteringRuleODataType,
		name:        data.Name.ValueStringPointer(),
		description: data.Description.ValueStringPointer(),
		action:      action,
		priority:    data.Priority.ValueInt64Pointer(),
		settings: &ruleSettingsRequestBody{
			status: data.Status.ValueStringPointer(),
		},
		matchingConditions: &matchingConditionsRequestBody{
			destinations: &destinationsRequestBody{
				targets:           targets,
				httpRequestMethod: commaStringPointer(httpMethods),
			},
		},
	}

	if len(sessionTypes) > 0 {
		requestBody.matchingConditions.sources = &sourcesRequestBody{
			sessionType: commaStringPointer(sessionTypes),
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

func stringSetValues(ctx context.Context, value types.Set) ([]string, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}

	var result []string
	diags := value.ElementsAs(ctx, &result, false)
	if diags.HasError() {
		return nil, fmt.Errorf("%s", diags.Errors()[0].Detail())
	}

	return result, nil
}

// customHeaderValues preserves the Entra portal's header insertion model. These
// are request headers added to traffic that matches an allow rule, not response
// headers returned by the destination. The portal rejects CR/LF injection both
// as literal control characters and as common escaped forms; schema validation
// handles literal control characters and hasCustomHeaderValueLineBreakEscape
// handles escaped forms observed in the portal validator.
func customHeaderValues(ctx context.Context, value types.List) ([]s.Parsable, error) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}

	var headers []customHeaderModel
	diags := value.ElementsAs(ctx, &headers, false)
	if diags.HasError() {
		return nil, fmt.Errorf("%s", diags.Errors()[0].Detail())
	}

	result := make([]s.Parsable, 0, len(headers))
	for _, header := range headers {
		headerName := header.HeaderName.ValueString()
		headerValue := header.HeaderValue.ValueString()
		if hasCustomHeaderValueLineBreakEscape(headerValue) {
			return nil, fmt.Errorf("custom_headers header_value for %q must not contain escaped CR or LF sequences", headerName)
		}
		result = append(result, &customHeaderRequestBody{
			odataType:   headerModificationAddODataType,
			headerName:  header.HeaderName.ValueStringPointer(),
			headerValue: header.HeaderValue.ValueStringPointer(),
		})
	}

	return result, nil
}

func hasCustomHeaderValueLineBreakEscape(value string) bool {
	return customHeaderLineBreakEscapePattern.MatchString(value)
}

// The Entra portal custom header validator rejects literal CR/LF characters and
// common escaped forms such as %0d, %0a, \x0d, and \u000a. Literal CR/LF are
// already rejected by the schema's printable-ASCII validator; this pattern keeps
// Terraform from accepting escaped values that the portal blocks.
var customHeaderLineBreakEscapePattern = regexp.MustCompile(`(?i)(%0[da]|\\x0[da]|\\u000[da])`)

func commaStringPointer(values []string) *string {
	if len(values) == 0 {
		return nil
	}

	value := strings.Join(values, ", ")
	return &value
}

type webContentFilteringPolicyRuleRequestBody struct {
	odataType          string
	name               *string
	description        *string
	action             *actionRequestBody
	priority           *int64
	settings           *ruleSettingsRequestBody
	matchingConditions *matchingConditionsRequestBody
}

func (b *webContentFilteringPolicyRuleRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("@odata.type", &b.odataType); err != nil {
		return err
	}
	if err := writer.WriteStringValue("name", b.name); err != nil {
		return err
	}
	if err := writer.WriteStringValue("description", b.description); err != nil {
		return err
	}
	if err := writer.WriteObjectValue("action", b.action); err != nil {
		return err
	}
	if err := writer.WriteInt64Value("priority", b.priority); err != nil {
		return err
	}
	if err := writer.WriteObjectValue("settings", b.settings); err != nil {
		return err
	}
	if err := writer.WriteObjectValue("matchingConditions", b.matchingConditions); err != nil {
		return err
	}

	return nil
}

func (b *webContentFilteringPolicyRuleRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type actionRequestBody struct {
	odataType      string
	headerSettings *headerSettingsRequestBody
}

func (b *actionRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("@odata.type", &b.odataType); err != nil {
		return err
	}
	if b.headerSettings != nil {
		if err := writer.WriteObjectValue("headerSettings", b.headerSettings); err != nil {
			return err
		}
	}

	return nil
}

func (b *actionRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type headerSettingsRequestBody struct {
	modifications []s.Parsable
}

func (b *headerSettingsRequestBody) Serialize(writer s.SerializationWriter) error {
	return writer.WriteCollectionOfObjectValues("modifications", b.modifications)
}

func (b *headerSettingsRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type ruleSettingsRequestBody struct {
	status *string
}

func (b *ruleSettingsRequestBody) Serialize(writer s.SerializationWriter) error {
	return writer.WriteStringValue("status", b.status)
}

func (b *ruleSettingsRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type matchingConditionsRequestBody struct {
	destinations *destinationsRequestBody
	sources      *sourcesRequestBody
}

func (b *matchingConditionsRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteObjectValue("destinations", b.destinations); err != nil {
		return err
	}
	if b.sources != nil {
		if err := writer.WriteObjectValue("sources", b.sources); err != nil {
			return err
		}
	}

	return nil
}

func (b *matchingConditionsRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type destinationsRequestBody struct {
	targets           []s.Parsable
	httpRequestMethod *string
}

func (b *destinationsRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteCollectionOfObjectValues("targets", b.targets); err != nil {
		return err
	}
	if err := writer.WriteStringValue("httpRequestMethod", b.httpRequestMethod); err != nil {
		return err
	}

	return nil
}

func (b *destinationsRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type destinationTargetRequestBody struct {
	odataType string
	values    []string
}

func (b *destinationTargetRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("@odata.type", &b.odataType); err != nil {
		return err
	}
	return writer.WriteCollectionOfStringValues("values", b.values)
}

func (b *destinationTargetRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type sourcesRequestBody struct {
	sessionType *string
}

func (b *sourcesRequestBody) Serialize(writer s.SerializationWriter) error {
	return writer.WriteStringValue("sessionType", b.sessionType)
}

func (b *sourcesRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

type customHeaderRequestBody struct {
	odataType   string
	headerName  *string
	headerValue *string
}

func (b *customHeaderRequestBody) Serialize(writer s.SerializationWriter) error {
	if err := writer.WriteStringValue("@odata.type", &b.odataType); err != nil {
		return err
	}
	if err := writer.WriteStringValue("headerName", b.headerName); err != nil {
		return err
	}
	return writer.WriteStringValue("headerValue", b.headerValue)
}

func (b *customHeaderRequestBody) GetFieldDeserializers() map[string]func(s.ParseNode) error {
	return map[string]func(s.ParseNode) error{}
}

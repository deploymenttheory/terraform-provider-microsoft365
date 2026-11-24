package customrequests

import (
	"context"
	"net/http"
	"testing"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryExpansionTemplate(t *testing.T) {
	tests := []struct {
		name          string
		urlTemplate   string
		query         map[string]string
		wantTemplate  string
		wantVariables map[string]string
	}{
		{name: "none", urlTemplate: "https://graph.microsoft.com/beta/x", wantTemplate: ""},
		{name: "OData options are sorted and percent-encoded", urlTemplate: "https://graph.microsoft.com/beta/x", query: map[string]string{"$select": "id", "$expand": "children"}, wantTemplate: "{?%24expand,%24select}", wantVariables: map[string]string{"%24expand": "children", "%24select": "id"}},
		{name: "plain and hyphenated names", urlTemplate: "https://graph.microsoft.com/beta/x", query: map[string]string{"api-version": "1", "top_n": "5"}, wantTemplate: "{?api%2Dversion,top_n}", wantVariables: map[string]string{"api%2Dversion": "1", "top_n": "5"}},
		{name: "literal query string continues with &", urlTemplate: "https://graph.microsoft.com/beta/x?a=b", query: map[string]string{"$top": "1"}, wantTemplate: "{&%24top}", wantVariables: map[string]string{"%24top": "1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTemplate, gotVariables := queryExpansionTemplate(tt.urlTemplate, tt.query)
			assert.Equal(t, tt.wantTemplate, gotTemplate)
			assert.Equal(t, tt.wantVariables, gotVariables)
		})
	}
}

// TestQueryParametersReachTheURL proves QueryParameters are put on the wire by every verb. Kiota only expands query
// variables the URL template names, so before #4941 these were silently dropped. The path, including each resource
// ID format, is unchanged; only the query string is added.
func TestQueryParametersReachTheURL(t *testing.T) {
	ctx := context.Background()
	query := map[string]string{"$expand": "children", "$select": "id,displayName", "$filter": "displayName eq 'a b'"}
	wantQuery := "?%24expand=children&%24filter=displayName%20eq%20%27a%20b%27&%24select=id%2CdisplayName"

	tests := []struct {
		name, wantURL string
		send          func(adapter abstractions.RequestAdapter) error
	}{
		{name: "GetRequest", wantURL: "GET https://graph.microsoft.com/beta/deviceManagement/deviceCategories/" + testResourceID + wantQuery, send: func(adapter abstractions.RequestAdapter) error {
			_, err := GetRequest(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", ResourceID: testResourceID, QueryParameters: query}, models.CreateDeviceCategoryFromDiscriminatorValue, nil)
			return err
		}},
		{name: "GetRequestByResourceId parentheses", wantURL: "GET https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('" + testResourceID + "')/settings" + wantQuery, send: func(adapter abstractions.RequestAdapter) error {
			_, err := GetRequestByResourceId(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/configurationPolicies", ResourceIDPattern: ResourceIDPatternParentheses, ResourceID: testResourceID, EndpointSuffix: "/settings", QueryParameters: query})
			return err
		}},
		{name: "GetRequestByResourceId segment", wantURL: "GET https://graph.microsoft.com/beta/applications/" + testResourceID + "/owners" + wantQuery, send: func(adapter abstractions.RequestAdapter) error {
			_, err := GetRequestByResourceId(ctx, adapter, GetRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "applications", ResourceIDPattern: ResourceIDPatternSegment, ResourceID: testResourceID, EndpointSuffix: "/owners", QueryParameters: query})
			return err
		}},
		{name: "PostRequest", wantURL: "POST https://graph.microsoft.com/beta/deviceManagement/deviceCategories" + wantQuery, send: func(adapter abstractions.RequestAdapter) error {
			_, err := PostRequest(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", RequestBody: newTestBody(), QueryParameters: query}, models.CreateDeviceCategoryFromDiscriminatorValue, nil)
			return err
		}},
		{name: "PostRequestNoContent", wantURL: "POST https://graph.microsoft.com/beta/deviceManagement/deviceCategories" + wantQuery, send: func(adapter abstractions.RequestAdapter) error {
			return PostRequestNoContent(ctx, adapter, PostRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", RequestBody: newTestBody(), QueryParameters: query})
		}},
		{name: "PatchRequestByResourceId", wantURL: "PATCH https://graph.microsoft.com/beta/deviceManagement/deviceCategories/" + testResourceID + wantQuery, send: func(adapter abstractions.RequestAdapter) error {
			return PatchRequestByResourceId(ctx, adapter, PatchRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/deviceCategories", ResourceIDPattern: ResourceIDPatternSegment, ResourceID: testResourceID, RequestBody: newTestBody(), QueryParameters: query})
		}},
		{name: "PutRequestByResourceId", wantURL: "PUT https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('" + testResourceID + "')" + wantQuery, send: func(adapter abstractions.RequestAdapter) error {
			return PutRequestByResourceId(ctx, adapter, PutRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/configurationPolicies", ResourceID: testResourceID, RequestBody: newTestBody(), QueryParameters: query})
		}},
		{name: "DeleteRequestByResourceId", wantURL: "DELETE https://graph.microsoft.com/beta/deviceManagement/configurationPolicies('" + testResourceID + "')" + wantQuery, send: func(adapter abstractions.RequestAdapter) error {
			return DeleteRequestByResourceId(ctx, adapter, DeleteRequestConfig{APIVersion: GraphAPIBeta, Endpoint: "deviceManagement/configurationPolicies", ResourceIDPattern: ResourceIDPatternParentheses, ResourceID: testResourceID, QueryParameters: query})
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter, transport := newTestAdapter(t, "https://graph.microsoft.com/beta")
			captured := captureRequests(transport, jsonResponse(http.StatusOK, `{}`))

			require.NoError(t, tt.send(adapter))
			require.Len(t, *captured, 1)
			assert.Equal(t, tt.wantURL, (*captured)[0].method+" "+(*captured)[0].url)
		})
	}
}

package mocks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jarcoal/httpmock"
)

// AuthenticationMocks provides mock responses for Azure authentication flows
type AuthenticationMocks struct {
	TenantID                string
	ClientID                string
	TokenExpiryMinutes      int
	AccessToken             string
	RefreshToken            string
	IDToken                 string
	InstanceDiscoveryResult *InstanceDiscoveryResponse
}

// NewAuthenticationMocks creates a new instance of AuthenticationMocks with default values
func NewAuthenticationMocks() *AuthenticationMocks {
	return &AuthenticationMocks{
		TenantID:           "00000000-0000-0000-0000-000000000000",
		ClientID:           "11111111-1111-1111-1111-111111111111",
		TokenExpiryMinutes: 60,
		AccessToken:        "mock_access_token",
		RefreshToken:       "mock_refresh_token",
		IDToken:            "mock_id_token",
		InstanceDiscoveryResult: &InstanceDiscoveryResponse{
			TenantDiscoveryEndpoint: "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000/v2.0/.well-known/openid-configuration",
			ApiVersion:              "1.1",
			Metadata: []MetadataEntry{
				{
					Preferred: true,
					Aliases: []string{
						"login.microsoftonline.com",
						"login.windows.net",
						"login.microsoft.com",
						"sts.windows.net",
					},
				},
			},
		},
	}
}

// RegisterMocks registers all mock responses with httpmock
func (a *AuthenticationMocks) RegisterMocks() {
	a.registerInstanceDiscovery()
	a.registerTokenEndpoint()
	a.registerOpenIDConfiguration()
	a.registerDeviceCodeFlow()
	a.registerManagedIdentityFlow()
}

// registerInstanceDiscovery registers the instance discovery endpoint mock
func (a *AuthenticationMocks) registerInstanceDiscovery() {
	httpmock.RegisterResponder(
		"GET",
		"https://login.microsoftonline.com/common/discovery/instance",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, a.InstanceDiscoveryResult)
			return resp, err
		},
	)
}

// registerTokenEndpoint registers the token endpoint mock for various auth flows
func (a *AuthenticationMocks) registerTokenEndpoint() {
	// Client credentials flow
	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", a.TenantID),
		func(req *http.Request) (*http.Response, error) {
			err := req.ParseForm()
			if err != nil {
				return httpmock.NewStringResponse(400, "Bad request"), nil
			}

			grantType := req.FormValue("grant_type")
			clientID := req.FormValue("client_id")

			// Validate the request
			if clientID != a.ClientID {
				return httpmock.NewStringResponse(401, "Invalid client"), nil
			}

			var tokenResp TokenResponse
			switch grantType {
			case "client_credentials":
				tokenResp = a.createTokenResponse([]string{"https://graph.microsoft.com/.default"})
			case "refresh_token":
				tokenResp = a.createTokenResponse([]string{"https://graph.microsoft.com/.default"})
				tokenResp.RefreshToken = a.RefreshToken
			case "password":
				tokenResp = a.createTokenResponse([]string{"https://graph.microsoft.com/.default"})
				tokenResp.RefreshToken = a.RefreshToken
				tokenResp.IDToken = a.IDToken
			case "urn:ietf:params:oauth:grant-type:jwt-bearer":
				// On-behalf-of flow
				tokenResp = a.createTokenResponse([]string{"https://graph.microsoft.com/.default"})
				tokenResp.RefreshToken = a.RefreshToken
			case "authorization_code":
				tokenResp = a.createTokenResponse([]string{"https://graph.microsoft.com/.default"})
				tokenResp.RefreshToken = a.RefreshToken
				tokenResp.IDToken = a.IDToken
			case "urn:ietf:params:oauth:grant-type:device_code":
				tokenResp = a.createTokenResponse([]string{"https://graph.microsoft.com/.default"})
				tokenResp.RefreshToken = a.RefreshToken
				tokenResp.IDToken = a.IDToken
			default:
				return httpmock.NewStringResponse(400, "Unsupported grant type"), nil
			}

			return httpmock.NewJsonResponse(200, tokenResp)
		},
	)
}

// registerOpenIDConfiguration registers the OpenID configuration endpoint mock
func (a *AuthenticationMocks) registerOpenIDConfiguration() {
	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0/.well-known/openid-configuration", a.TenantID),
		func(req *http.Request) (*http.Response, error) {
			openIDConfig := OpenIDConfigurationResponse{
				TokenEndpoint:                    fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", a.TenantID),
				TokenEndpointAuthMethods:         []string{"client_secret_post", "private_key_jwt", "client_secret_basic"},
				JwksURI:                          fmt.Sprintf("https://login.microsoftonline.com/%s/discovery/v2.0/keys", a.TenantID),
				ResponseModesSupported:           []string{"query", "fragment", "form_post"},
				SubjectTypesSupported:            []string{"pairwise"},
				IDTokenSigningAlgValuesSupported: []string{"RS256"},
				ResponseTypesSupported:           []string{"code", "id_token", "code id_token", "id_token token"},
				ScopesSupported:                  []string{"openid", "profile", "email", "offline_access"},
				Issuer:                           fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", a.TenantID),
				RequestURIParameterSupported:     false,
				UserInfoEndpoint:                 fmt.Sprintf("https://graph.microsoft.com/oidc/userinfo"),
				AuthorizationEndpoint:            fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", a.TenantID),
				DeviceAuthorizationEndpoint:      fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/devicecode", a.TenantID),
				HTTPLogoutSupported:              true,
				FrontchannelLogoutSupported:      true,
				EndSessionEndpoint:               fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/logout", a.TenantID),
				ClaimsSupported:                  []string{"sub", "iss", "cloud_instance_name", "cloud_instance_host_name", "cloud_graph_host_name", "msgraph_host", "aud", "exp", "iat", "auth_time", "acr", "nonce", "preferred_username", "name", "tid", "ver", "at_hash", "c_hash", "email"},
				TenantRegionScope:                "NA",
				CloudGraphHostName:               "graph.windows.net",
				MsgraphHost:                      "graph.microsoft.com",
				RbacURL:                          "https://pas.windows.net",
			}

			return httpmock.NewJsonResponse(200, openIDConfig)
		},
	)
}

// registerDeviceCodeFlow registers the device code flow endpoints
func (a *AuthenticationMocks) registerDeviceCodeFlow() {
	httpmock.RegisterResponder(
		"POST",
		fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/devicecode", a.TenantID),
		func(req *http.Request) (*http.Response, error) {
			err := req.ParseForm()
			if err != nil {
				return httpmock.NewStringResponse(400, "Bad request"), nil
			}

			clientID := req.FormValue("client_id")
			if clientID != a.ClientID {
				return httpmock.NewStringResponse(401, "Invalid client"), nil
			}

			deviceCodeResp := DeviceCodeResponse{
				UserCode:        "ABCDEFGH",
				DeviceCode:      "mock_device_code_value",
				VerificationURI: "https://microsoft.com/devicelogin",
				ExpiresIn:       900,
				Interval:        5,
				Message:         "To sign in, use a web browser to open the page https://microsoft.com/devicelogin and enter the code ABCDEFGH to authenticate.",
			}

			return httpmock.NewJsonResponse(200, deviceCodeResp)
		},
	)
}

// registerManagedIdentityFlow registers the managed identity endpoints
func (a *AuthenticationMocks) registerManagedIdentityFlow() {
	// Azure VM IMDS endpoint
	httpmock.RegisterResponder(
		"GET",
		"http://169.254.169.254/metadata/identity/oauth2/token",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Metadata") != "true" {
				return httpmock.NewStringResponse(400, "Metadata header required"), nil
			}

			resource := req.URL.Query().Get("resource")
			if resource == "" {
				return httpmock.NewStringResponse(400, "Resource required"), nil
			}

			tokenResp := a.createTokenResponse([]string{resource})
			// Managed identity response format is slightly different
			managedIdentityResp := map[string]interface{}{
				"access_token": tokenResp.AccessToken,
				"expires_on":   fmt.Sprintf("%d", time.Now().Add(time.Duration(a.TokenExpiryMinutes)*time.Minute).Unix()),
				"resource":     resource,
				"token_type":   "Bearer",
				"client_id":    a.ClientID,
			}

			return httpmock.NewJsonResponse(200, managedIdentityResp)
		},
	)

	// App Service MSI endpoint
	httpmock.RegisterResponder(
		"GET",
		"http://localhost:8081/msi/token",
		func(req *http.Request) (*http.Response, error) {
			resource := req.URL.Query().Get("resource")
			if resource == "" {
				return httpmock.NewStringResponse(400, "Resource required"), nil
			}

			tokenResp := a.createTokenResponse([]string{resource})
			// App Service MSI response format
			msiResp := map[string]interface{}{
				"access_token": tokenResp.AccessToken,
				"expires_on":   fmt.Sprintf("%d", time.Now().Add(time.Duration(a.TokenExpiryMinutes)*time.Minute).Unix()),
				"resource":     resource,
				"token_type":   "Bearer",
			}

			return httpmock.NewJsonResponse(200, msiResp)
		},
	)
}

// createTokenResponse creates a token response with the specified scopes
func (a *AuthenticationMocks) createTokenResponse(scopes []string) TokenResponse {
	now := time.Now()
	expiryTime := now.Add(time.Duration(a.TokenExpiryMinutes) * time.Minute)

	return TokenResponse{
		TokenType:    "Bearer",
		ExpiresIn:    a.TokenExpiryMinutes * 60,
		ExtExpiresIn: a.TokenExpiryMinutes * 60,
		AccessToken:  a.AccessToken,
		Scope:        strings.Join(scopes, " "),
		ExpiresOn:    fmt.Sprintf("%d", expiryTime.Unix()),
	}
}

// Response Structures

// InstanceDiscoveryResponse represents the response from the instance discovery endpoint
type InstanceDiscoveryResponse struct {
	TenantDiscoveryEndpoint string          `json:"tenant_discovery_endpoint"`
	ApiVersion              string          `json:"api-version"`
	Metadata                []MetadataEntry `json:"metadata"`
}

// MetadataEntry represents a metadata entry in the instance discovery response
type MetadataEntry struct {
	Preferred bool     `json:"preferred"`
	Aliases   []string `json:"aliases"`
}

// TokenResponse represents the response from the token endpoint
type TokenResponse struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	ExpiresOn    string `json:"expires_on,omitempty"`
}

// OpenIDConfigurationResponse represents the response from the OpenID configuration endpoint
type OpenIDConfigurationResponse struct {
	TokenEndpoint                    string   `json:"token_endpoint"`
	TokenEndpointAuthMethods         []string `json:"token_endpoint_auth_methods_supported"`
	JwksURI                          string   `json:"jwks_uri"`
	ResponseModesSupported           []string `json:"response_modes_supported"`
	SubjectTypesSupported            []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
	ResponseTypesSupported           []string `json:"response_types_supported"`
	ScopesSupported                  []string `json:"scopes_supported"`
	Issuer                           string   `json:"issuer"`
	RequestURIParameterSupported     bool     `json:"request_uri_parameter_supported"`
	UserInfoEndpoint                 string   `json:"userinfo_endpoint"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	DeviceAuthorizationEndpoint      string   `json:"device_authorization_endpoint"`
	HTTPLogoutSupported              bool     `json:"http_logout_supported"`
	FrontchannelLogoutSupported      bool     `json:"frontchannel_logout_supported"`
	EndSessionEndpoint               string   `json:"end_session_endpoint"`
	ClaimsSupported                  []string `json:"claims_supported"`
	TenantRegionScope                string   `json:"tenant_region_scope"`
	CloudGraphHostName               string   `json:"cloud_graph_host_name"`
	MsgraphHost                      string   `json:"msgraph_host"`
	RbacURL                          string   `json:"rbac_url"`
}

// DeviceCodeResponse represents the response from the device code endpoint
type DeviceCodeResponse struct {
	UserCode        string `json:"user_code"`
	DeviceCode      string `json:"device_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	Message         string `json:"message"`
}

func (a *AuthenticationMocks) GraphBetaTokenResponse() []byte {
	tokenResp := a.createTokenResponse([]string{"https://graph.microsoft.com/.default"})
	respBytes, _ := json.Marshal(tokenResp)
	return respBytes
}

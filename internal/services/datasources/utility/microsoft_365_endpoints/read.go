package utilityMicrosoft365Endpoints

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	endpointsAPIBaseURL = "https://endpoints.office.com/endpoints"
	versionAPIURL       = "https://endpoints.office.com/version"
	apiTimeout          = 30 * time.Second
)

// instanceToAPIEndpoint maps user-friendly instance names to API endpoint names
var instanceToAPIEndpoint = map[string]string{
	"worldwide":     "worldwide",
	"usgov-dod":     "USGOVDoD",
	"usgov-gcchigh": "USGOVGCCHigh",
	"china":         "China",
}

// Microsoft365EndpointResponse represents the JSON response from the Microsoft 365 endpoints API
type Microsoft365EndpointResponse struct {
	ID                     int64    `json:"id"`
	ServiceArea            string   `json:"serviceArea"`
	ServiceAreaDisplayName string   `json:"serviceAreaDisplayName"`
	URLs                   []string `json:"urls,omitempty"`
	IPs                    []string `json:"ips,omitempty"`
	TCPPorts               string   `json:"tcpPorts,omitempty"`
	UDPPorts               string   `json:"udpPorts,omitempty"`
	ExpressRoute           bool     `json:"expressRoute"`
	Category               string   `json:"category"`
	Required               bool     `json:"required"`
	Notes                  string   `json:"notes,omitempty"`
}

func (d *microsoft365EndpointsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data Microsoft365EndpointsDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", DataSourceName))

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	instance := data.Instance.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Reading Microsoft 365 endpoints for instance: %s", instance))

	apiInstance, ok := instanceToAPIEndpoint[instance]
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid instance",
			fmt.Sprintf("Unknown Microsoft 365 cloud instance: %s", instance),
		)
		return
	}

	endpoints, err := getAllowListEndpoints(ctx, apiInstance)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to fetch Microsoft 365 endpoints",
			fmt.Sprintf("Error calling Microsoft 365 endpoints API: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetched %d endpoints from API", len(endpoints)))

	filteredEndpoints := filterEndpoints(ctx, endpoints, data)
	tflog.Debug(ctx, fmt.Sprintf("After filtering: %d endpoints", len(filteredEndpoints)))

	if err := mapResponseToState(ctx, &data, filteredEndpoints); err != nil {
		resp.Diagnostics.AddError(
			"Failed to populate datasource state",
			fmt.Sprintf("Error converting endpoints to Terraform state: %s", err.Error()),
		)
		return
	}

	// Generate deterministic ID based on instance and filters
	data.Id = generateID(data)

	tflog.Info(ctx, fmt.Sprintf("Successfully read %d Microsoft 365 endpoints", len(filteredEndpoints)))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// getAllowListEndpoints calls the Microsoft 365 endpoints API
func getAllowListEndpoints(ctx context.Context, instance string) ([]Microsoft365EndpointResponse, error) {

	clientRequestID := uuid.New().String()

	url := fmt.Sprintf("%s/%s?clientrequestid=%s", endpointsAPIBaseURL, instance, clientRequestID)

	tflog.Debug(ctx, fmt.Sprintf("Calling Microsoft 365 endpoints API: %s", url))

	client := &http.Client{
		Timeout: apiTimeout,
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "terraform-provider-microsoft365")

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", httpResp.StatusCode, string(body))
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	tflog.Trace(ctx, fmt.Sprintf("API response body length: %d bytes", len(body)))

	var endpoints []Microsoft365EndpointResponse
	if err := json.Unmarshal(body, &endpoints); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return endpoints, nil
}

// filterEndpoints applies user-specified filters to the endpoint list
func filterEndpoints(ctx context.Context, endpoints []Microsoft365EndpointResponse, data Microsoft365EndpointsDataSourceModel) []Microsoft365EndpointResponse {
	if data.ServiceAreas.IsNull() && data.Categories.IsNull() && data.RequiredOnly.IsNull() && data.ExpressRoute.IsNull() {
		return endpoints
	}

	serviceAreasFilter := make(map[string]bool)
	if !data.ServiceAreas.IsNull() {
		var serviceAreas []string
		data.ServiceAreas.ElementsAs(ctx, &serviceAreas, false)
		for _, sa := range serviceAreas {
			serviceAreasFilter[sa] = true
		}
	}

	categoriesFilter := make(map[string]bool)
	if !data.Categories.IsNull() {
		var categories []string
		data.Categories.ElementsAs(ctx, &categories, false)
		for _, cat := range categories {
			categoriesFilter[cat] = true
		}
	}

	requiredOnly := false
	if !data.RequiredOnly.IsNull() {
		requiredOnly = data.RequiredOnly.ValueBool()
	}

	expressRouteOnly := false
	if !data.ExpressRoute.IsNull() {
		expressRouteOnly = data.ExpressRoute.ValueBool()
	}

	filtered := make([]Microsoft365EndpointResponse, 0, len(endpoints))
	for _, endpoint := range endpoints {
		if len(serviceAreasFilter) > 0 && !serviceAreasFilter[endpoint.ServiceArea] {
			continue
		}

		if len(categoriesFilter) > 0 && !categoriesFilter[endpoint.Category] {
			continue
		}

		if requiredOnly && !endpoint.Required {
			continue
		}

		if expressRouteOnly && !endpoint.ExpressRoute {
			continue
		}

		filtered = append(filtered, endpoint)
	}

	return filtered
}

// generateID creates a deterministic ID based on the datasource configuration
func generateID(data Microsoft365EndpointsDataSourceModel) types.String {

	h := sha256.New()
	h.Write([]byte(data.Instance.ValueString()))

	if !data.ServiceAreas.IsNull() {
		var serviceAreas []string
		data.ServiceAreas.ElementsAs(context.Background(), &serviceAreas, false)
		for _, sa := range serviceAreas {
			h.Write([]byte(sa))
		}
	}

	if !data.Categories.IsNull() {
		var categories []string
		data.Categories.ElementsAs(context.Background(), &categories, false)
		for _, cat := range categories {
			h.Write([]byte(cat))
		}
	}

	if !data.RequiredOnly.IsNull() && data.RequiredOnly.ValueBool() {
		h.Write([]byte("required"))
	}

	if !data.ExpressRoute.IsNull() && data.ExpressRoute.ValueBool() {
		h.Write([]byte("expressroute"))
	}

	hashStr := fmt.Sprintf("%x", h.Sum(nil))[:16]
	return types.StringValue(fmt.Sprintf("%s_%s", data.Instance.ValueString(), hashStr))
}

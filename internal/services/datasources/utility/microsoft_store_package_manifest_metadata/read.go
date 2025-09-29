package utilityMicrosoftStorePackageManifest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Microsoft Store API URLs
const (
	MSStoreManifestBaseURL = "https://storeedgefd.dsx.mp.microsoft.com/v9.0"
	MSStoreSearchURL       = MSStoreManifestBaseURL + "/manifestSearch"
	MSStoreManifestURL     = MSStoreManifestBaseURL + "/packageManifests"
)

// API Response structures
type MSStoreAPIResponse struct {
	Data interface{} `json:"Data"`
}

type MSStoreSearchRequest struct {
	Query MSStoreQuery `json:"Query"`
}

type MSStoreQuery struct {
	KeyWord   string `json:"KeyWord"`
	MatchType string `json:"MatchType"`
}

type MSStoreSearchResult struct {
	PackageIdentifier string `json:"PackageIdentifier"`
	PackageName       string `json:"PackageName"`
}

// Read fetches Microsoft Store package manifests and sets them in the data source state
func (d *MicrosoftStorePackageManifestDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config MicrosoftStorePackageManifestDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s_%s", d.ProviderTypeName, d.TypeName))

	// Get the configuration
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	packageIdProvided := !config.PackageIdentifier.IsNull() && config.PackageIdentifier.ValueString() != ""
	searchTermProvided := !config.SearchTerm.IsNull() && config.SearchTerm.ValueString() != ""

	tflog.Debug(ctx, fmt.Sprintf("Reading %s_%s with package ID provided: %t, search term provided: %t",
		d.ProviderTypeName, d.TypeName, packageIdProvided, searchTermProvided))

	// Validate inputs - must have either a package identifier or search term, but not both
	if !packageIdProvided && !searchTermProvided {
		resp.Diagnostics.AddError(
			"Missing Input Parameter",
			"Either package_identifier or search_term must be provided",
		)
		return
	}

	if packageIdProvided && searchTermProvided {
		resp.Diagnostics.AddError(
			"Multiple Input Parameters",
			"Only one of package_identifier or search_term can be provided",
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, config.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var manifests []interface{}
	var err error

	// Get package manifests
	if packageIdProvided {
		packageId := config.PackageIdentifier.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Getting specific package manifest: %s", packageId))

		manifest, err := d.getPackageManifestById(ctx, packageId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Getting Package Manifest",
				fmt.Sprintf("Unable to get package manifest for ID %s: %s", packageId, err),
			)
			return
		}

		if manifest != nil {
			manifests = []interface{}{manifest}
		}
	} else {
		searchTerm := config.SearchTerm.ValueString()
		tflog.Debug(ctx, fmt.Sprintf("Searching for packages with term: %s", searchTerm))

		manifests, err = d.searchAndGetManifests(ctx, searchTerm)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Searching Package Manifests",
				fmt.Sprintf("Unable to search for packages with term %s: %s", searchTerm, err),
			)
			return
		}
	}

	// Create a new state model with the results
	var state MicrosoftStorePackageManifestDataSourceModel

	// Copy the configuration values
	state.PackageIdentifier = config.PackageIdentifier
	state.SearchTerm = config.SearchTerm
	state.Timeouts = config.Timeouts

	// Convert API response to Terraform model
	terraformManifests, diags := d.mapRemoteStateToTerraformState(ctx, manifests)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Manifests = terraformManifests

	tflog.Debug(ctx, fmt.Sprintf("Successfully retrieved %d package manifest(s)", len(terraformManifests)))

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// getPackageManifestById retrieves a specific package manifest by ID
func (d *MicrosoftStorePackageManifestDataSource) getPackageManifestById(ctx context.Context, packageId string) (interface{}, error) {
	url := fmt.Sprintf("%s/%s", MSStoreManifestURL, packageId)

	tflog.Debug(ctx, fmt.Sprintf("Making GET request to: %s", url))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	httpResp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", httpResp.StatusCode)
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var apiResponse MSStoreAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	return apiResponse.Data, nil
}

// searchAndGetManifests searches for packages and retrieves their manifests
func (d *MicrosoftStorePackageManifestDataSource) searchAndGetManifests(ctx context.Context, searchTerm string) ([]interface{}, error) {
	// First, search for packages
	searchResults, err := d.searchPackages(ctx, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("error searching packages: %w", err)
	}

	if len(searchResults) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("No packages found for search term: %s", searchTerm))
		return []interface{}{}, nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Found %d package(s), retrieving manifests", len(searchResults)))

	var manifests []interface{}
	for _, result := range searchResults {
		manifest, err := d.getPackageManifestById(ctx, result.PackageIdentifier)
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Failed to get manifest for package %s: %s", result.PackageIdentifier, err))
			continue
		}

		if manifest != nil {
			manifests = append(manifests, manifest)
		}
	}

	return manifests, nil
}

// searchPackages searches for packages using the search term
func (d *MicrosoftStorePackageManifestDataSource) searchPackages(ctx context.Context, searchTerm string) ([]MSStoreSearchResult, error) {
	searchRequest := MSStoreSearchRequest{
		Query: MSStoreQuery{
			KeyWord:   searchTerm,
			MatchType: "Substring",
		},
	}

	requestBody, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshaling search request: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Making POST request to: %s", MSStoreSearchURL))

	req, err := http.NewRequestWithContext(ctx, "POST", MSStoreSearchURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating search request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	httpResp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making search request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search API returned status %d", httpResp.StatusCode)
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading search response body: %w", err)
	}

	var apiResponse MSStoreAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("error unmarshaling search response: %w", err)
	}

	// Convert the interface{} to []MSStoreSearchResult
	var searchResults []MSStoreSearchResult
	if data, ok := apiResponse.Data.([]interface{}); ok {
		for _, item := range data {
			if itemMap, ok := item.(map[string]any); ok {
				result := MSStoreSearchResult{}
				if packageId, ok := itemMap["PackageIdentifier"].(string); ok {
					result.PackageIdentifier = packageId
				}
				if packageName, ok := itemMap["PackageName"].(string); ok {
					result.PackageName = packageName
				}
				searchResults = append(searchResults, result)
			}
		}
	}

	return searchResults, nil
}

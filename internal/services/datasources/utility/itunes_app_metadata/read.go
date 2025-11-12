package itunes_app_metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Read refreshes the Terraform state with the latest data
func (d *itunesAppMetadataDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ItunesAppMetadataDataSourceModel

	tflog.Debug(ctx, "Reading iTunes App Metadata data source")

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Validate required attributes
	if state.SearchTerm.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("search_term"),
			"Missing Search Term",
			"The search_term attribute must be set.",
		)
		return
	}

	if state.CountryCode.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("country_code"),
			"Missing Country Code",
			"The country_code attribute must be set.",
		)
		return
	}

	state.Id = types.StringValue(fmt.Sprintf("%s_%s", state.CountryCode.ValueString(), state.SearchTerm.ValueString()))

	request, err := constructRequest(ctx, state.CountryCode.ValueString(), state.SearchTerm.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Constructing iTunes Search Request",
			fmt.Sprintf("Could not construct iTunes search request: %s", err),
		)
		return
	}

	searchResponse, diags := executeItunesSearchRequest(ctx, request)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Map response to state
	diags = mapResponseToState(ctx, searchResponse, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// executeItunesSearchRequest performs the HTTP request and returns the parsed response
func executeItunesSearchRequest(ctx context.Context, req *http.Request) (ItunesSearchResponse, diag.Diagnostics) {
	var diags diag.Diagnostics
	var searchResponse ItunesSearchResponse

	tflog.Debug(ctx, fmt.Sprintf("Executing iTunes app metadata request to: %s", req.URL.String()))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		diags.AddError(
			"Error Making HTTP Request",
			fmt.Sprintf("Could not make HTTP request: %s", err),
		)
		return searchResponse, diags
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		diags.AddError(
			"Error from iTunes API",
			fmt.Sprintf("Received non-OK status code: %d", resp.StatusCode),
		)
		return searchResponse, diags
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		diags.AddError(
			"Error Reading Response Body",
			fmt.Sprintf("Could not read response body: %s", err),
		)
		return searchResponse, diags
	}

	if err := json.Unmarshal(body, &searchResponse); err != nil {
		diags.AddError(
			"Error Parsing JSON Response",
			fmt.Sprintf("Could not parse JSON response: %s", err),
		)
		return searchResponse, diags
	}

	tflog.Debug(ctx, fmt.Sprintf("Found %d apps matching search criteria", searchResponse.ResultCount))
	return searchResponse, diags
}

// REF: https://learn.microsoft.com/en-us/graph/api/application-get?view=graph-rest-beta
package graphBetaApplication

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/applications"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// lookupMethod represents the different ways to look up an application
type lookupMethod int

const (
	lookupByODataQuery lookupMethod = iota
	lookupByObjectId
	lookupByAppId
	lookupByDisplayName
)

// Read handles the Read operation.
func (d *ApplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object ApplicationDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for datasource: %s", DataSourceName))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	var application graphmodels.Applicationable
	var err error

	method := determineLookupMethod(object)
	switch method {
	case lookupByODataQuery:
		application, err = d.getApplicationByODataQuery(ctx, object)
	case lookupByObjectId:
		application, err = d.getApplicationByObjectId(ctx, object)
	case lookupByAppId:
		application, err = d.getApplicationByAppId(ctx, object)
	case lookupByDisplayName:
		application, err = d.getApplicationByDisplayName(ctx, object)
	}

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	if application == nil || application.GetId() == nil {
		resp.Diagnostics.AddError(
			"Application Not Found",
			"The application lookup did not return a valid application with an ID. The application may not exist or may not have fully propagated in the directory.",
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully found application with ID: %s", *application.GetId()))

	ownerUserIds, err := d.listAllApplicationOwners(ctx, *application.GetId())
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to retrieve application owners: %v", err))
	}

	mappedState := MapRemoteStateToDataSource(ctx, application, ownerUserIds, object)

	resp.Diagnostics.Append(resp.State.Set(ctx, &mappedState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s", DataSourceName))
}

// determineLookupMethod determines which lookup method to use based on provided attributes
func determineLookupMethod(object ApplicationDataSourceModel) lookupMethod {
	switch {
	case !object.ODataQuery.IsNull() && object.ODataQuery.ValueString() != "":
		return lookupByODataQuery
	case !object.ObjectId.IsNull() && object.ObjectId.ValueString() != "":
		return lookupByObjectId
	case !object.AppId.IsNull() && object.AppId.ValueString() != "":
		return lookupByAppId
	case !object.DisplayName.IsNull() && object.DisplayName.ValueString() != "":
		return lookupByDisplayName
	default:
		return lookupByObjectId // This should never happen due to schema validators
	}
}

// getApplicationByObjectId retrieves an application by its object ID
// Includes retry logic because even direct GET can return 404 immediately after creation
func (d *ApplicationDataSource) getApplicationByObjectId(ctx context.Context, object ApplicationDataSourceModel) (graphmodels.Applicationable, error) {
	objectId := object.ObjectId.ValueString()

	maxRetries := 6
	retryDelay := 10 * time.Second
	startTime := time.Now()

	var application graphmodels.Applicationable
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		application, err = d.client.
			Applications().
			ByApplicationId(objectId).
			Get(ctx, nil)

		if err == nil && application != nil && application.GetId() != nil {
			return application, nil
		}

		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errors.IsNonRetryableReadError(&errorInfo) {
				return nil, err
			}
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("application not found after %d attempts (%v total wait): %w", maxRetries, time.Since(startTime), err)
	}
	return nil, fmt.Errorf("application not found or invalid after %d attempts (%v total wait)", maxRetries, time.Since(startTime))
}

// getApplicationByAppId retrieves an application by its app ID (client ID)
func (d *ApplicationDataSource) getApplicationByAppId(ctx context.Context, object ApplicationDataSourceModel) (graphmodels.Applicationable, error) {
	filter := fmt.Sprintf("appId eq '%s'", object.AppId.ValueString())
	return d.executeOdataQueryWithRetry(ctx, filter, fmt.Sprintf("app_id: %s", object.AppId.ValueString()))
}

// getApplicationByDisplayName retrieves an application by display name
func (d *ApplicationDataSource) getApplicationByDisplayName(ctx context.Context, object ApplicationDataSourceModel) (graphmodels.Applicationable, error) {
	filter := fmt.Sprintf("displayName eq '%s'", object.DisplayName.ValueString())
	return d.executeOdataQueryWithRetry(ctx, filter, fmt.Sprintf("display_name: %s", object.DisplayName.ValueString()))
}

// getApplicationByODataQuery retrieves an application using a custom OData query
func (d *ApplicationDataSource) getApplicationByODataQuery(ctx context.Context, object ApplicationDataSourceModel) (graphmodels.Applicationable, error) {
	filter := object.ODataQuery.ValueString()
	return d.executeOdataQueryWithRetry(ctx, filter, fmt.Sprintf("OData query: %s", filter))
}

// executeOdataQueryWithRetry executes a filtered query with retry logic for eventual consistency
func (d *ApplicationDataSource) executeOdataQueryWithRetry(ctx context.Context, filter string, description string) (graphmodels.Applicationable, error) {
	maxRetries := 6
	retryDelay := 10 * time.Second

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &applications.ApplicationsRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &applications.ApplicationsRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		applicationsResponse, err := d.client.Applications().Get(ctx, requestConfig)
		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errors.IsNonRetryableReadError(&errorInfo) {
				return nil, err
			}
		} else if len(applicationsResponse.GetValue()) > 0 {
			return validateSingleApplication(applicationsResponse.GetValue(), description)
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	return validateSingleApplication([]graphmodels.Applicationable{}, description)
}

// validateSingleApplication ensures exactly one application was returned
func validateSingleApplication(applicationList []graphmodels.Applicationable, criteria string) (graphmodels.Applicationable, error) {
	switch len(applicationList) {
	case 0:
		return nil, fmt.Errorf("no application found with %s", criteria)
	case 1:
		return applicationList[0], nil
	default:
		return nil, fmt.Errorf("found %d applications with %s. The query must return exactly one application", len(applicationList), criteria)
	}
}

// listAllApplicationOwners retrieves all owners of an application using pagination
func (d *ApplicationDataSource) listAllApplicationOwners(ctx context.Context, applicationId string) ([]string, error) {
	var ownerIds []string

	ownersResponse, err := d.client.
		Applications().
		ByApplicationId(applicationId).
		Owners().
		Get(ctx, nil)

	if err != nil {
		return nil, err
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.DirectoryObjectable](
		ownersResponse,
		d.client.GetAdapter(),
		graphmodels.CreateDirectoryObjectCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator for owners: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(item graphmodels.DirectoryObjectable) bool {
		if item != nil && item.GetId() != nil {
			ownerIds = append(ownerIds, *item.GetId())
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate owner pages: %w", err)
	}

	return ownerIds, nil
}

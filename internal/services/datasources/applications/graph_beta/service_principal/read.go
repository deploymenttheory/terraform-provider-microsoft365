// REF: https://learn.microsoft.com/en-us/graph/api/serviceprincipal-get?view=graph-rest-beta
package graphBetaServicePrincipal

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
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/serviceprincipals"
)

// lookupMethod represents the different ways to look up a service principal
type lookupMethod int

const (
	lookupByODataQuery lookupMethod = iota
	lookupByObjectId
	lookupByAppId
	lookupByDisplayName
)

// Read handles the Read operation.
func (d *ServicePrincipalDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object ServicePrincipalDataSourceModel

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

	var servicePrincipal graphmodels.ServicePrincipalable
	var err error

	method := determineLookupMethod(object)
	switch method {
	case lookupByODataQuery:
		servicePrincipal, err = d.getServicePrincipalByODataQuery(ctx, object)
	case lookupByObjectId:
		servicePrincipal, err = d.getServicePrincipalByObjectId(ctx, object)
	case lookupByAppId:
		servicePrincipal, err = d.getServicePrincipalByAppId(ctx, object)
	case lookupByDisplayName:
		servicePrincipal, err = d.getServicePrincipalByDisplayName(ctx, object)
	}

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	if servicePrincipal == nil || servicePrincipal.GetId() == nil {
		resp.Diagnostics.AddError(
			"Service Principal Not Found",
			"The service principal lookup did not return a valid service principal with an ID. The service principal may not exist or may not have fully propagated in the directory.",
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully found service principal with ID: %s", *servicePrincipal.GetId()))

	mappedState := MapRemoteStateToDataSource(ctx, servicePrincipal, object)

	resp.Diagnostics.Append(resp.State.Set(ctx, &mappedState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s", DataSourceName))
}

// determineLookupMethod determines which lookup method to use based on provided attributes
func determineLookupMethod(object ServicePrincipalDataSourceModel) lookupMethod {
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

// getServicePrincipalByObjectId retrieves a service principal by its object ID
// Includes retry logic because even direct GET can return 404 immediately after creation
func (d *ServicePrincipalDataSource) getServicePrincipalByObjectId(ctx context.Context, object ServicePrincipalDataSourceModel) (graphmodels.ServicePrincipalable, error) {
	objectId := object.ObjectId.ValueString()

	maxRetries := 6
	retryDelay := 10 * time.Second
	startTime := time.Now()

	var servicePrincipal graphmodels.ServicePrincipalable
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		servicePrincipal, err = d.client.
			ServicePrincipals().
			ByServicePrincipalId(objectId).
			Get(ctx, nil)

		if err == nil && servicePrincipal != nil && servicePrincipal.GetId() != nil {
			return servicePrincipal, nil
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
		return nil, fmt.Errorf("service principal not found after %d attempts (%v total wait): %w", maxRetries, time.Since(startTime), err)
	}
	return nil, fmt.Errorf("service principal not found or invalid after %d attempts (%v total wait)", maxRetries, time.Since(startTime))
}

// getServicePrincipalByAppId retrieves a service principal by its app ID (client ID)
func (d *ServicePrincipalDataSource) getServicePrincipalByAppId(ctx context.Context, object ServicePrincipalDataSourceModel) (graphmodels.ServicePrincipalable, error) {
	filter := fmt.Sprintf("appId eq '%s'", object.AppId.ValueString())
	return d.executeOdataQueryWithRetry(ctx, filter, fmt.Sprintf("app_id: %s", object.AppId.ValueString()))
}

// getServicePrincipalByDisplayName retrieves a service principal by display name
func (d *ServicePrincipalDataSource) getServicePrincipalByDisplayName(ctx context.Context, object ServicePrincipalDataSourceModel) (graphmodels.ServicePrincipalable, error) {
	filter := fmt.Sprintf("displayName eq '%s'", object.DisplayName.ValueString())
	return d.executeOdataQueryWithRetry(ctx, filter, fmt.Sprintf("display_name: %s", object.DisplayName.ValueString()))
}

// getServicePrincipalByODataQuery retrieves a service principal using a custom OData query
func (d *ServicePrincipalDataSource) getServicePrincipalByODataQuery(ctx context.Context, object ServicePrincipalDataSourceModel) (graphmodels.ServicePrincipalable, error) {
	filter := object.ODataQuery.ValueString()
	return d.executeOdataQueryWithRetry(ctx, filter, fmt.Sprintf("OData query: %s", filter))
}

// executeOdataQueryWithRetry executes a filtered query with retry logic for eventual consistency
func (d *ServicePrincipalDataSource) executeOdataQueryWithRetry(ctx context.Context, filter string, description string) (graphmodels.ServicePrincipalable, error) {
	maxRetries := 6
	retryDelay := 10 * time.Second

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &serviceprincipals.ServicePrincipalsRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &serviceprincipals.ServicePrincipalsRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		servicePrincipalsResponse, err := d.client.ServicePrincipals().Get(ctx, requestConfig)
		if err != nil {
			errorInfo := errors.GraphError(ctx, err)
			if errors.IsNonRetryableReadError(&errorInfo) {
				return nil, err
			}
		} else if len(servicePrincipalsResponse.GetValue()) > 0 {
			return validateSingleServicePrincipal(servicePrincipalsResponse.GetValue(), description)
		}

		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	return validateSingleServicePrincipal([]graphmodels.ServicePrincipalable{}, description)
}

// validateSingleServicePrincipal ensures exactly one service principal was returned
func validateSingleServicePrincipal(servicePrincipalList []graphmodels.ServicePrincipalable, criteria string) (graphmodels.ServicePrincipalable, error) {
	switch len(servicePrincipalList) {
	case 0:
		return nil, fmt.Errorf("no service principal found with %s", criteria)
	case 1:
		return servicePrincipalList[0], nil
	default:
		return nil, fmt.Errorf("found %d service principals with %s. The query must return exactly one service principal", len(servicePrincipalList), criteria)
	}
}

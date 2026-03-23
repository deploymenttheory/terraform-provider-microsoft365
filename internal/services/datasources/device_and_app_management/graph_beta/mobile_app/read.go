// REF: https://learn.microsoft.com/en-us/graph/api/intune-apps-mobileapp-list?view=graph-rest-beta
package graphBetaMobileApp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// lookupMethod represents the different ways to look up mobile apps
type lookupMethod int

const (
	lookupByODataQuery lookupMethod = iota
	lookupByAppId
	lookupByDisplayName
	lookupByPublisher
	lookupByDeveloper
	lookupByCategory
	lookupListAll
)

// Read handles the Read operation.
func (d *MobileAppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object MobileAppDataSourceModel

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

	var apps []graphmodels.MobileAppable
	var err error

	method := determineLookupMethod(object)
	switch method {
	case lookupByODataQuery:
		apps, err = d.getAppsByODataQuery(ctx, object)
	case lookupByAppId:
		apps, err = d.getAppById(ctx, object)
	case lookupByDisplayName:
		apps, err = d.getAppsByDisplayName(ctx, object)
	case lookupByPublisher:
		apps, err = d.getAppsByPublisher(ctx, object)
	case lookupByDeveloper:
		apps, err = d.getAppsByDeveloper(ctx, object)
	case lookupByCategory:
		apps, err = d.getAppsByCategory(ctx, object)
	case lookupListAll:
		apps, err = d.listAllApps(ctx, object)
	}

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	if len(apps) == 0 {
		resp.Diagnostics.AddWarning(
			"No Mobile Apps Found",
			"The lookup did not return any mobile apps matching the specified criteria.",
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully found %d mobile app(s)", len(apps)))

	// Apply app type filter if specified
	filteredItems := make([]MobileAppModel, 0, len(apps))
	appTypeFilter := object.AppTypeFilter.ValueString()

	for _, app := range apps {
		if appTypeFilter != "" {
			currentAppType := getAppTypeFromMobileApp(app)
			if currentAppType != appTypeFilter {
				continue
			}
		}

		appItem := MapRemoteStateToDataSource(ctx, app)
		filteredItems = append(filteredItems, appItem)
	}

	object.Items = filteredItems

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s", DataSourceName))
}

// determineLookupMethod determines which lookup method to use based on provided attributes
func determineLookupMethod(object MobileAppDataSourceModel) lookupMethod {
	switch {
	case !object.ODataQuery.IsNull() && object.ODataQuery.ValueString() != "":
		return lookupByODataQuery
	case !object.AppId.IsNull() && object.AppId.ValueString() != "":
		return lookupByAppId
	case !object.DisplayName.IsNull() && object.DisplayName.ValueString() != "":
		return lookupByDisplayName
	case !object.Publisher.IsNull() && object.Publisher.ValueString() != "":
		return lookupByPublisher
	case !object.Developer.IsNull() && object.Developer.ValueString() != "":
		return lookupByDeveloper
	case !object.Category.IsNull() && object.Category.ValueString() != "":
		return lookupByCategory
	case !object.ListAll.IsNull() && object.ListAll.ValueBool():
		return lookupListAll
	default:
		return lookupByAppId
	}
}

// getAppById retrieves a mobile app by its app ID
func (d *MobileAppDataSource) getAppById(ctx context.Context, object MobileAppDataSourceModel) ([]graphmodels.MobileAppable, error) {
	appId := object.AppId.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up mobile app by ID: %s", appId))

	app, err := d.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(appId).
		Get(ctx, nil)

	if err != nil {
		return nil, err
	}

	if app == nil {
		return []graphmodels.MobileAppable{}, nil
	}

	return []graphmodels.MobileAppable{app}, nil
}

// getAppsByDisplayName retrieves mobile apps by display name using OData server-side filtering
func (d *MobileAppDataSource) getAppsByDisplayName(ctx context.Context, object MobileAppDataSourceModel) ([]graphmodels.MobileAppable, error) {
	displayName := object.DisplayName.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up mobile apps by display name: %s", displayName))

	filter := fmt.Sprintf("contains(tolower(displayName), '%s')", strings.ToLower(displayName))

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &deviceappmanagement.MobileAppsRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.getAllMobileAppsWithPageIterator(ctx, requestConfig)
}

// getAppsByPublisher retrieves mobile apps by publisher name using OData server-side filtering
func (d *MobileAppDataSource) getAppsByPublisher(ctx context.Context, object MobileAppDataSourceModel) ([]graphmodels.MobileAppable, error) {
	publisher := object.Publisher.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up mobile apps by publisher: %s", publisher))

	filter := fmt.Sprintf("contains(tolower(publisher), '%s')", strings.ToLower(publisher))

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &deviceappmanagement.MobileAppsRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.getAllMobileAppsWithPageIterator(ctx, requestConfig)
}

// getAppsByDeveloper retrieves mobile apps by developer name using OData server-side filtering
func (d *MobileAppDataSource) getAppsByDeveloper(ctx context.Context, object MobileAppDataSourceModel) ([]graphmodels.MobileAppable, error) {
	developer := object.Developer.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up mobile apps by developer: %s", developer))

	filter := fmt.Sprintf("contains(tolower(developer), '%s')", strings.ToLower(developer))

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &deviceappmanagement.MobileAppsRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.getAllMobileAppsWithPageIterator(ctx, requestConfig)
}

// getAppsByCategory retrieves mobile apps by category name using local filtering
func (d *MobileAppDataSource) getAppsByCategory(ctx context.Context, object MobileAppDataSourceModel) ([]graphmodels.MobileAppable, error) {
	category := object.Category.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up mobile apps by category: %s", category))

	requestParameters := &deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration{
		QueryParameters: &deviceappmanagement.MobileAppsRequestBuilderGetQueryParameters{},
	}

	allApps, err := d.getAllMobileAppsWithPageIterator(ctx, requestParameters)
	if err != nil {
		return nil, err
	}

	var filteredApps []graphmodels.MobileAppable
	for _, app := range allApps {
		appId := app.GetId()
		if appId == nil {
			continue
		}

		categories, err := d.client.
			DeviceAppManagement().
			MobileApps().
			ByMobileAppId(*appId).
			Categories().
			Get(ctx, nil)

		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Failed to fetch categories for app %s: %v", *appId, err))
			continue
		}

		if categories == nil || categories.GetValue() == nil {
			continue
		}

		for _, cat := range categories.GetValue() {
			if cat.GetDisplayName() != nil && strings.Contains(
				strings.ToLower(*cat.GetDisplayName()),
				strings.ToLower(category)) {
				app.SetCategories(categories.GetValue())
				filteredApps = append(filteredApps, app)
				break
			}
		}
	}

	return filteredApps, nil
}

// getAppsByODataQuery retrieves mobile apps using a custom OData query
func (d *MobileAppDataSource) getAppsByODataQuery(ctx context.Context, object MobileAppDataSourceModel) ([]graphmodels.MobileAppable, error) {
	filter := object.ODataQuery.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up mobile apps with OData query: %s", filter))

	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestParameters := &deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &deviceappmanagement.MobileAppsRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.getAllMobileAppsWithPageIterator(ctx, requestParameters)
}

// listAllApps retrieves all mobile apps
func (d *MobileAppDataSource) listAllApps(ctx context.Context, object MobileAppDataSourceModel) ([]graphmodels.MobileAppable, error) {
	tflog.Debug(ctx, "Listing all mobile apps")

	requestConfig := &deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration{
		QueryParameters: &deviceappmanagement.MobileAppsRequestBuilderGetQueryParameters{},
	}

	return d.getAllMobileAppsWithPageIterator(ctx, requestConfig)
}

// getAllMobileAppsWithPageIterator uses Microsoft Graph SDK's PageIterator to handle pagination
func (d *MobileAppDataSource) getAllMobileAppsWithPageIterator(
	ctx context.Context,
	requestConfig *deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration,
) ([]graphmodels.MobileAppable, error) {
	var allApps []graphmodels.MobileAppable

	resp, err := d.client.
		DeviceAppManagement().
		MobileApps().
		Get(ctx, requestConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to get initial page of mobile apps: %w", err)
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.MobileAppable](
		resp,
		d.client.GetAdapter(),
		graphmodels.CreateMobileAppCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(app graphmodels.MobileAppable) bool {
		allApps = append(allApps, app)
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("error during pagination: %w", err)
	}

	return allApps, nil
}

// REF: https://learn.microsoft.com/en-us/graph/api/user-list?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/user-get?view=graph-rest-beta
package graphBetaUser

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphusers "github.com/microsoftgraph/msgraph-beta-sdk-go/users"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

type lookupMethod int

const (
	lookupByObjectId lookupMethod = iota
	lookupByDisplayName
	lookupByEmployeeId
	lookupByGivenName
	lookupByUserPrincipalName
	lookupByOnPremisesImmutableId
	lookupByOnPremisesDistinguishedName
	lookupByODataQuery
	lookupListAll
)

// Read handles the Read operation.
func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object UserDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	method := d.determineLookupMethod(object)
	var users []graphmodels.Userable
	var err error

	switch method {
	case lookupByObjectId:
		users, err = d.getUserByObjectId(ctx, object)
	case lookupByDisplayName:
		users, err = d.getUsersByFilter(ctx, fmt.Sprintf("displayName eq '%s'", object.DisplayName.ValueString()))
	case lookupByEmployeeId:
		users, err = d.getUsersByFilter(ctx, fmt.Sprintf("employeeId eq '%s'", object.EmployeeId.ValueString()))
	case lookupByGivenName:
		users, err = d.getUsersByFilter(ctx, fmt.Sprintf("givenName eq '%s'", object.GivenName.ValueString()))
	case lookupByUserPrincipalName:
		users, err = d.getUsersByFilter(ctx, fmt.Sprintf("userPrincipalName eq '%s'", object.UserPrincipalName.ValueString()))
	case lookupByOnPremisesImmutableId:
		users, err = d.getUsersByFilter(ctx, fmt.Sprintf("onPremisesImmutableId eq '%s'", object.OnPremisesImmutableId.ValueString()))
	case lookupByOnPremisesDistinguishedName:
		users, err = d.getUsersByFilter(ctx, fmt.Sprintf("onPremisesDistinguishedName eq '%s'", object.OnPremisesDistinguishedName.ValueString()))
	case lookupByODataQuery:
		users, err = d.getUsersByFilter(ctx, object.ODataQuery.ValueString())
	case lookupListAll:
		users, err = d.listAllUsers(ctx)
	}

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	object.Items = ConstructUserItems(users)
	object.ID = types.StringValue(fmt.Sprintf("user-datasource-%d", time.Now().Unix()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

// determineLookupMethod determines which lookup method to use based on the provided attributes
func (d *UserDataSource) determineLookupMethod(object UserDataSourceModel) lookupMethod {
	switch {
	case !object.ObjectId.IsNull() && object.ObjectId.ValueString() != "":
		return lookupByObjectId
	case !object.DisplayName.IsNull() && object.DisplayName.ValueString() != "":
		return lookupByDisplayName
	case !object.EmployeeId.IsNull() && object.EmployeeId.ValueString() != "":
		return lookupByEmployeeId
	case !object.GivenName.IsNull() && object.GivenName.ValueString() != "":
		return lookupByGivenName
	case !object.UserPrincipalName.IsNull() && object.UserPrincipalName.ValueString() != "":
		return lookupByUserPrincipalName
	case !object.OnPremisesImmutableId.IsNull() && object.OnPremisesImmutableId.ValueString() != "":
		return lookupByOnPremisesImmutableId
	case !object.OnPremisesDistinguishedName.IsNull() && object.OnPremisesDistinguishedName.ValueString() != "":
		return lookupByOnPremisesDistinguishedName
	case !object.ODataQuery.IsNull() && object.ODataQuery.ValueString() != "":
		return lookupByODataQuery
	case !object.ListAll.IsNull() && object.ListAll.ValueBool():
		return lookupListAll
	default:
		return lookupListAll
	}
}

// getUserByObjectId retrieves a single user by its object ID
func (d *UserDataSource) getUserByObjectId(ctx context.Context, object UserDataSourceModel) ([]graphmodels.Userable, error) {
	objectId := object.ObjectId.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Looking up user by object ID: %s", objectId))

	user, err := d.client.Users().ByUserId(objectId).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return []graphmodels.Userable{}, nil
	}

	return []graphmodels.Userable{user}, nil
}

// getUsersByFilter retrieves users using an OData $filter expression. Advanced query
// capabilities (ConsistencyLevel: eventual and $count=true) are enabled so that filters
// on properties such as employeeId and onPremisesDistinguishedName are supported.
func (d *UserDataSource) getUsersByFilter(ctx context.Context, filter string) ([]graphmodels.Userable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Looking up users with filter: %s", filter))

	count := true
	headers := abstractions.NewRequestHeaders()
	headers.Add("ConsistencyLevel", "eventual")

	requestConfig := &graphusers.UsersRequestBuilderGetRequestConfiguration{
		Headers: headers,
		QueryParameters: &graphusers.UsersRequestBuilderGetQueryParameters{
			Filter: &filter,
			Count:  &count,
		},
	}

	return d.listAllUsersWithPageIterator(ctx, requestConfig)
}

// listAllUsers retrieves all users in the tenant
func (d *UserDataSource) listAllUsers(ctx context.Context) ([]graphmodels.Userable, error) {
	tflog.Debug(ctx, "Listing all users")

	return d.listAllUsersWithPageIterator(ctx, nil)
}

// listAllUsersWithPageIterator handles pagination for user list requests
func (d *UserDataSource) listAllUsersWithPageIterator(ctx context.Context, requestConfig *graphusers.UsersRequestBuilderGetRequestConfiguration) ([]graphmodels.Userable, error) {
	var allUsers []graphmodels.Userable

	result, err := d.client.Users().Get(ctx, requestConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get initial page of users: %w", err)
	}

	if result == nil {
		return allUsers, nil
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.Userable](
		result,
		d.client.GetAdapter(),
		graphmodels.CreateUserCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(user graphmodels.Userable) bool {
		if user != nil {
			allUsers = append(allUsers, user)
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("error during pagination: %w", err)
	}

	return allUsers, nil
}

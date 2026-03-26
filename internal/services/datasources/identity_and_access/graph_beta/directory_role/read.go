// REF: https://learn.microsoft.com/en-us/graph/api/directoryrole-list?view=graph-rest-beta
// REF: https://learn.microsoft.com/en-us/graph/api/directoryrole-get?view=graph-rest-beta

package graphBetaDirectoryRole

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphdirectoryroles "github.com/microsoftgraph/msgraph-beta-sdk-go/directoryroles"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

type lookupMethod int

const (
	lookupByRoleID lookupMethod = iota
	lookupByDisplayName
	lookupListAll
)

func (d *DirectoryRoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object DirectoryRoleDataSourceModel

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

	var roles []graphmodels.DirectoryRoleable
	var err error

	method := determineLookupMethod(object)
	switch method {
	case lookupByRoleID:
		roles, err = d.getRoleByID(ctx, object)
	case lookupByDisplayName:
		roles, err = d.getRolesByDisplayName(ctx, object)
	case lookupListAll:
		roles, err = d.listAllRoles(ctx)
	}

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	if len(roles) == 0 {
		resp.Diagnostics.AddWarning(
			"No directory roles found",
			fmt.Sprintf("No activated directory roles matched the specified filter for datasource: %s", DataSourceName),
		)
	}

	object.Items = make([]DirectoryRoleModel, 0, len(roles))
	for _, role := range roles {
		object.Items = append(object.Items, MapRemoteStateToDataSource(ctx, role))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Datasource Read Method: %s, found %d items", DataSourceName, len(object.Items)))
}

func determineLookupMethod(object DirectoryRoleDataSourceModel) lookupMethod {
	if !object.RoleID.IsNull() && !object.RoleID.IsUnknown() && object.RoleID.ValueString() != "" {
		return lookupByRoleID
	}
	if !object.DisplayName.IsNull() && !object.DisplayName.IsUnknown() && object.DisplayName.ValueString() != "" {
		return lookupByDisplayName
	}
	return lookupListAll
}

func (d *DirectoryRoleDataSource) getRoleByID(ctx context.Context, object DirectoryRoleDataSourceModel) ([]graphmodels.DirectoryRoleable, error) {
	roleID := object.RoleID.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("Fetching directory role by ID: %s", roleID))

	role, err := d.client.DirectoryRoles().ByDirectoryRoleId(roleID).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	return []graphmodels.DirectoryRoleable{role}, nil
}

func (d *DirectoryRoleDataSource) getRolesByDisplayName(ctx context.Context, object DirectoryRoleDataSourceModel) ([]graphmodels.DirectoryRoleable, error) {
	displayName := object.DisplayName.ValueString()
	filter := fmt.Sprintf("displayName eq '%s'", displayName)
	tflog.Debug(ctx, fmt.Sprintf("Fetching directory roles with OData filter: %s", filter))

	requestConfig := &graphdirectoryroles.DirectoryRolesRequestBuilderGetRequestConfiguration{
		QueryParameters: &graphdirectoryroles.DirectoryRolesRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	return d.listAllRolesRaw(ctx, requestConfig)
}

func (d *DirectoryRoleDataSource) listAllRoles(ctx context.Context) ([]graphmodels.DirectoryRoleable, error) {
	tflog.Debug(ctx, "Listing all activated directory roles")
	return d.listAllRolesRaw(ctx, nil)
}

func (d *DirectoryRoleDataSource) listAllRolesRaw(ctx context.Context, requestConfig *graphdirectoryroles.DirectoryRolesRequestBuilderGetRequestConfiguration) ([]graphmodels.DirectoryRoleable, error) {
	rolesResp, err := d.client.DirectoryRoles().Get(ctx, requestConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to list directory roles: %w", err)
	}

	pageIterator, err := graphcore.NewPageIterator[graphmodels.DirectoryRoleable](
		rolesResp,
		d.client.GetAdapter(),
		graphmodels.CreateDirectoryRoleCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	var allRoles []graphmodels.DirectoryRoleable
	err = pageIterator.Iterate(ctx, func(role graphmodels.DirectoryRoleable) bool {
		if role != nil {
			allRoles = append(allRoles, role)
		}
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate directory roles: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Retrieved %d activated directory roles", len(allRoles)))
	return allRoles, nil
}

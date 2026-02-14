package graphBetaUsersUser

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/users"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// listAllResourcesWithPageIterator retrieves all users using the PageIterator
// This handles pagination automatically and returns ALL users across all pages
func (r *UserListResource) listAllResourcesWithPageIterator(
	ctx context.Context,
	requestConfig *users.UsersRequestBuilderGetRequestConfiguration,
) ([]models.Userable, error) {
	var allUsers []models.Userable

	tflog.Debug(ctx, "Fetching first page of users")

	usersResponse, err := r.client.
		Users().
		Get(ctx, requestConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	pageIterator, err := graphcore.NewPageIterator[models.Userable](
		usersResponse,
		r.client.GetAdapter(),
		models.CreateUserCollectionResponseFromDiscriminatorValue,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	pageCount := 0
	err = pageIterator.Iterate(ctx, func(item models.Userable) bool {
		if item != nil {
			allUsers = append(allUsers, item)

			if len(allUsers)%100 == 0 {
				pageCount++
				tflog.Debug(ctx, fmt.Sprintf("PageIterator: collected %d users (estimated page %d)", len(allUsers), pageCount))
			}
		}
		return true
	})

	if err != nil {
		return nil, fmt.Errorf("failed to iterate pages: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("PageIterator complete: collected %d total users", len(allUsers)))

	return allUsers, nil
}

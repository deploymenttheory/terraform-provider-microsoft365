package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// DoesNotExistInGraph validates that the specified resource does not exist within Microsoft Graph
func DoesNotExistInGraph(testResource types.TestResource, resourceName string) resource.TestCheckFunc {
	return existsFunc(false)(testResource, resourceName)
}

// ExistsInGraph validates that the specified resource exists within Microsoft Graph
func ExistsInGraph(testResource types.TestResource, resourceName string) resource.TestCheckFunc {
	return existsFunc(true)(testResource, resourceName)
}

func existsFunc(shouldExist bool) func(types.TestResource, string) resource.TestCheckFunc {
	return func(testResource types.TestResource, resourceName string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			// Exists check should never take more than 5 minutes
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()

			rs, ok := s.RootModule().Resources[resourceName]
			if !ok {
				return fmt.Errorf("%q was not found in the state", resourceName)
			}

			result, err := testResource.Exists(ctx, nil, rs.Primary)
			if err != nil {
				return fmt.Errorf("running exists func for %q: %+v", resourceName, err)
			}

			if result == nil {
				return fmt.Errorf("received nil for exists for %q", resourceName)
			}

			if *result != shouldExist {
				if !shouldExist {
					return fmt.Errorf("%q still exists", resourceName)
				}
				return fmt.Errorf("%q did not exist", resourceName)
			}

			return nil
		}
	}
}

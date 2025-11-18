package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestResource is an interface which should be implemented by all Terraform Resources
// for testing purposes - this allows a consistent approach to testing Resources.
type TestResource interface {
	// Exists should check whether the resource exists within the remote API
	// returning true if the resource exists, false if it doesn't exist and an
	// error if something went wrong when checking for the existence of the resource.
	Exists(ctx context.Context, client any, state *terraform.InstanceState) (*bool, error)
}

// TestResourceVerifyingRemoved is an interface which extends TestResource for resources
// which also need to verify they can be manually removed (for testing destroy operations)
type TestResourceVerifyingRemoved interface {
	TestResource

	// Destroy manually destroys this resource
	Destroy(ctx context.Context, client any, state *terraform.InstanceState) (*bool, error)
}

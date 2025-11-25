package graphBetaGroupLifecycleExpirationPolicyAssignment

// No acceptance tests are configured for this resource due to the following constraints:
//
// 1. This resource depends on the group_lifecycle_expiration_policy resource, which is a
//    tenant-wide singleton (only one lifecycle policy can exist per Microsoft 365 tenant).
//
// 2. This assignment resource specifically requires the lifecycle policy to have
//    managed_group_types = "Selected". If the policy is configured with managed_group_types = "All",
//    this resource cannot be used (all groups are automatically managed by the policy).
//
// 3. The global nature of the lifecycle policy setting creates a conflict with lab/testing
//    environments. If acceptance tests configure the policy with managed_group_types = "Selected"
//    for testing, but the tenant is actively using the policy with different settings
//    (e.g., managed_group_types = "All" for production use), there would be conflicts as both
//    configurations cannot co-exist.
//
// 4. Modifying the tenant-wide lifecycle policy during acceptance tests would disrupt any
//    existing group lifecycle management in the tenant, potentially affecting production workloads.
//
// For these reasons, only unit tests with mocked API responses are implemented for this resource.

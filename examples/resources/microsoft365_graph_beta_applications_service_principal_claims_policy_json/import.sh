#!/bin/sh
# Import by the owning service principal UUID, not the returned custom policy UUID.
terraform import microsoft365_graph_beta_applications_service_principal_claims_policy_json.example '11111111-1111-4111-8111-111111111111'

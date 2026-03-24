#!/bin/bash
# Import with soft delete (default)
terraform import microsoft365_graph_beta_identity_and_access_administrative_unit.example 00000000-0000-0000-0000-000000000000

# Import with hard delete enabled
terraform import microsoft365_graph_beta_identity_and_access_administrative_unit.example 00000000-0000-0000-0000-000000000000:hard_delete=true

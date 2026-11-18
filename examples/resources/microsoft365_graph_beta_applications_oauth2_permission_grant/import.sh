#!/bin/sh
# Grant IDs are opaque strings, not UUIDs.
terraform import microsoft365_graph_beta_applications_oauth2_permission_grant.tenant_consent 'OPAQUE_GRAPH_GRANT_ID'

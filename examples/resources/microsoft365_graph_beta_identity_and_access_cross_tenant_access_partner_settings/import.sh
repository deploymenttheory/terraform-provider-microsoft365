#!/bin/bash
# Import using the partner tenant ID
terraform import microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings.example 12345678-1234-1234-1234-123456789012

# Import with hard_delete enabled
terraform import microsoft365_graph_beta_identity_and_access_cross_tenant_access_partner_settings.example "12345678-1234-1234-1234-123456789012:hard_delete=true"

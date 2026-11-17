#!/bin/bash
# Import using the directory property name. The key is case-sensitive: use the exact casing Microsoft Graph returns
# from GET /admin/people/profileCardProperties (for example CustomAttribute1, not customattribute1).
terraform import microsoft365_graph_beta_m365_admin_profile_card_property.example CustomAttribute1

#!/bin/bash

# Import a group-specific directory setting
# Format: {group_id}/{setting_id}
# Replace {group_id} with the actual group ID and {setting_id} with the setting ID

terraform import microsoft365_graph_beta_groups_group_settings.guest_settings "{group_id}/{setting_id}" 
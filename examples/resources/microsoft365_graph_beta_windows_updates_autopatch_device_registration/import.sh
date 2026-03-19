#!/bin/bash

# Import existing Windows Autopatch device registration by update category
# The ID format is the update category: "quality", "feature", or "driver"

# Import quality updates registration
terraform import microsoft365_graph_beta_windows_updates_autopatch_device_registration.quality_updates "quality"

# Import feature updates registration
terraform import microsoft365_graph_beta_windows_updates_autopatch_device_registration.feature_updates "feature"

# Import driver updates registration
terraform import microsoft365_graph_beta_windows_updates_autopatch_device_registration.driver_updates "driver"

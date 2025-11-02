#!/bin/bash

# Import an existing IP application segment by its ID
terraform import microsoft365_graph_beta_applications_ip_application_segment.example_ip_address "00000000-0000-0000-0000-000000000000"

# The ID format is the segment's unique identifier (GUID)
# You can find the segment ID in the Azure Portal or via Microsoft Graph API:
# GET https://graph.microsoft.com/beta/applications/{application-id}/onPremisesPublishing/segmentsConfiguration/microsoft.graph.ipSegmentConfiguration/applicationSegments

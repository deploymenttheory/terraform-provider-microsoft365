---
page_title: "Data Source: microsoft365_graph_cloud_pc_device_image"
description: |-
  
---

# Data Source: microsoft365_graph_cloud_pc_device_image





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The unique identifier (ID) of the image resource on the Cloud PC.

### Read-Only

- `display_name` (String) The display name of the associated device image.
- `error_code` (String) The error code of the status of the image that indicates why the upload failed, if applicable.
- `expiration_date` (String) The date when the image became unavailable.
- `last_modified_date_time` (String) The date and time when the image was last modified.
- `operating_system` (String) The operating system (OS) of the image.
- `os_build_number` (String) The OS build version of the image.
- `os_status` (String) The OS status of this image.
- `source_image_resource_id` (String) The unique identifier (ID) of the source image resource on Azure.
- `status` (String) The status of the image on the Cloud PC.
- `version` (String) The image version.

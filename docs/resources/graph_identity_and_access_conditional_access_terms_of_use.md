---
page_title: "microsoft365_graph_identity_and_access_conditional_access_terms_of_use Resource - terraform-provider-microsoft365"
subcategory: "Identity and Access"
description: |-
  Manages Microsoft 365 Terms of Use Agreements using the /agreements endpoint. This resource is used to terms of use agreements allow organizations to present information that users must accept before accessing data or applications. These agreements can be used to ensure compliance with legal or regulatory requirements..
---

# microsoft365_graph_identity_and_access_conditional_access_terms_of_use (Resource)

Manages Microsoft 365 Terms of Use Agreements using the `/agreements` endpoint. This resource is used to terms of use agreements allow organizations to present information that users must accept before accessing data or applications. These agreements can be used to ensure compliance with legal or regulatory requirements..

## Microsoft Documentation

- [conditionalAccessTermsOfUse resource type](https://learn.microsoft.com/en-us/graph/api/resources/agreement?view=graph-rest-1.0)
- [Create conditionalAccessTermsOfUse](https://learn.microsoft.com/en-us/graph/api/termsofusecontainer-post-agreements?view=graph-rest-1.0&tabs=http)
- [Update conditionalAccessTermsOfUse](https://learn.microsoft.com/en-us/graph/api/agreement-update?view=graph-rest-1.0&tabs=http)
- [Delete conditionalAccessTermsOfUse](https://learn.microsoft.com/en-us/graph/api/agreement-delete?view=graph-rest-1.0&tabs=http)

## Microsoft Graph API Permissions

The following client `application` permissions are needed in order to use this resource:

**Required:**
- `DeviceManagementConfiguration.Read.All`
- `DeviceManagementConfiguration.ReadWrite.All`

**Optional:**
- `None` `[N/A]`

## Version History

| Version | Status | Notes |
|---------|--------|-------|
| v0.34.0-alpha | Experimental | Initial release |

## Example Usage

```terraform
resource "microsoft365_graph_identity_and_access_conditional_access_terms_of_use" "example" {
  display_name                          = "Example Conditional Access Terms of Use"
  is_viewing_before_acceptance_required = true
  is_per_device_acceptance_required     = false
  user_reaccept_required_frequency      = "P10D"

  terms_expiration = {
    start_date_time = "2025-11-06"
    frequency       = "P180D"
  }

  file = {
    localizations = [
      {
        file_name    = "test.pdf"
        display_name = "Terms of Use - English (US)"
        language     = "en-US"
        is_default   = true
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - English"
        language         = "en"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - English (UK)"
        language         = "en-GB"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Afrikaans"
        language         = "af"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Amharic"
        language         = "am"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Arabic (Saudi Arabia)"
        language         = "ar-SA"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Armenian"
        language         = "hy"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Assamese"
        language         = "as"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Azerbaijani"
        language         = "az"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Belarusian"
        language         = "be"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Bangla"
        language         = "bn"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Bangla (India)"
        language         = "bn-IN"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Basque"
        language         = "eu"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Central Kurdish"
        language         = "ku-Arab"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Japanese"
        language         = "ja"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Zulu"
        language         = "zu"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Telugu"
        language         = "te"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Thai"
        language         = "th"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Tigrinya"
        language         = "ti"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Turkish"
        language         = "tr"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Turkmen"
        language         = "tk"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Ukrainian"
        language         = "uk"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Urdu"
        language         = "ur"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Uyghur"
        language         = "ug"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Uzbek"
        language         = "uz"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Valencian (Spain)"
        language         = "ca-ES-valencia"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Vietnamese"
        language         = "vi"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Welsh"
        language         = "cy"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Wolof"
        language         = "wo"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      },
      {
        file_name        = "test.pdf"
        display_name     = "Terms of Use - Yoruba"
        language         = "yo"
        is_default       = false
        is_major_version = false
        file_data = {
          data = filebase64("${path.module}/test.pdf")
        }
      }
    ]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `display_name` (String) Display name of the agreement. The display name is used for internal tracking of the agreement but isn't shown to end users who view the agreement.
- `file` (Attributes) Default PDF linked to this agreement. This is required when creating a new agreement. (see [below for nested schema](#nestedatt--file))
- `is_per_device_acceptance_required` (Boolean) This setting enables you to require end users to accept this agreement on every device that they're accessing it from. The end user is required to register their device in Microsoft Entra ID, if they haven't already done so.
- `is_viewing_before_acceptance_required` (Boolean) Indicates whether the user has to expand the agreement before accepting.

### Optional

- `terms_expiration` (Attributes) Expiration schedule and frequency of agreement for all users. (see [below for nested schema](#nestedatt--terms_expiration))
- `timeouts` (Attributes) (see [below for nested schema](#nestedatt--timeouts))
- `user_reaccept_required_frequency` (String) The duration after which the user must reaccept the terms of use. Must be in ISO 8601 duration format (e.g., `P90D`).

**Common values:**
- `P365D` - Annually (365 days)
- `P180D` - Bi-annually (180 days)
- `P90D` - Quarterly (90 days)
- `P30D` - Monthly (30 days)
- `P270D` - 270 days

### Read-Only

- `id` (String) The unique identifier for the agreement. This is automatically generated when the agreement is created.

<a id="nestedatt--file"></a>
### Nested Schema for `file`

Required:

- `localizations` (Attributes Set) The localized version of the terms of use agreement files attached to the agreement. (see [below for nested schema](#nestedatt--file--localizations))

<a id="nestedatt--file--localizations"></a>
### Nested Schema for `file.localizations`

Required:

- `display_name` (String) Localized display name of the policy file of an agreement. The localized display name is shown to end users who view the agreement.
- `file_data` (Attributes) Data that represents the terms of use PDF document. Must be provided during creation but is not returned by the API and will always be null in state. This field is intentionally not persisted for security reasons. (see [below for nested schema](#nestedatt--file--localizations--file_data))
- `file_name` (String) Name of the agreement file (for example, TOU.pdf).
- `is_default` (Boolean) If none of the languages matches the client preference, indicates whether this is the default agreement file. If none of the files are marked as default, the first one is treated as the default. Must be true if the language is 'en-US'.
- `language` (String) The language of the agreement file. When `is_default` is `true`, must be `en-US` (full format with country code). When `is_default` is `false`, must use only the two-letter language code (e.g., `en`, `fr`, `de`). The language code is a lowercase two-letter code derived from ISO 639-1.

Optional:

- `is_major_version` (Boolean) Indicates whether the agreement file is a major version update. Major version updates invalidate the agreement's acceptances on the corresponding language.

<a id="nestedatt--file--localizations--file_data"></a>
### Nested Schema for `file.localizations.file_data`

Required:

- `data` (String) Data that represents the terms of use PDF document as raw bytes (base64 encoded).




<a id="nestedatt--terms_expiration"></a>
### Nested Schema for `terms_expiration`

Optional:

- `frequency` (String) Represents the frequency at which the terms will expire, after its first expiration as set in startDateTime. Must be in ISO 8601 duration format.

**Accepted values:**
- `P365D` - Annually (365 days)
- `P180D` - Bi-annually (180 days)
- `P90D` - Quarterly (90 days)
- `P30D` - Monthly (30 days)
- `start_date_time` (String) The date when the agreement is set to expire for all users. Must be in YYYY-MM-DD format (e.g., `2025-12-31`) and is always in UTC time. The time portion (T00:00:00Z) will be automatically appended.


<a id="nestedatt--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).
- `delete` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Setting a timeout for a Delete operation is only applicable if changes are saved into state before the destroy operation occurs.
- `read` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours). Read operations occur during any refresh or planning operation when refresh is enabled.
- `update` (String) A string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as "30s" or "2h45m". Valid time units are "s" (seconds), "m" (minutes), "h" (hours).

## Import

Import is supported using the following syntax:

```shell
#!/bin/bash

terraform import microsoft365_graph_identity_and_access_conditional_access_terms_of_use.example 00000000-0000-0000-0000-000000000000 
# where 00000000-0000-0000-0000-000000000000 is the Conditional Access Terms of Use ID
``` 
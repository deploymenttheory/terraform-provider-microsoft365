# Intune template coverage

This checklist follows the user-supplied Intune template list (2026-09-29). A UI label is not proof of an API mapping; the notes distinguish live chooser observations, curl captures, and Terraform lifecycle results. The resource manages `/beta/deviceManagement/deviceConfigurations`; separate API families use their existing resources rather than a discriminator-based endpoint switch.

The 72 supplied entries comprise **67 entries in the deviceConfigurations family** and **five separate-API entries**. The latter are BIOS, current Delivery Optimization, Imported Administrative templates, Properties catalog, and OEMConfig. The legacy Delivery Optimization profile is an additional supported API scenario.

## Evidence levels

- Cases 01–63: prior sequential acceptance results (54 live cases), including assignment transitions and JSON encoding/relationship regressions.
- Cases 65–105 extend the lifecycle coverage with 41 scenarios. Captures and tests cover create, GET, full settings PATCH, import, empty plans, and deletion. Final validation results are recorded below.
- Cases 98–103 cover standalone AOSP, Android owner, and Android work trusted-root/SCEP lifecycles. Cases 74 and 82 create their certificate dependencies in the same configuration.
- Live UI inspection distinguishes legacy profiles from templates migrated to other API families. Metadata inventory alone does not establish current UI availability.
- Tests create unassigned profiles unless an explicit assignment scenario deploys empty test groups. API lifecycle success does not verify device-side application, certificate issuance, VPN connectivity, BIOS execution, or Android migration.

## Windows

| UI template | API type / collection | Evidence / remaining work |
| --- | --- | --- |
| BIOS configurations and other settings | `hardwareConfigurations` | Separate resource: windows_bios_configurations_and_other_settings_template |
| Custom | `windows10CustomConfiguration` | 07,10,22–28,63 |
| Delivery Optimization | `configurationPolicies` (current UI); `windowsDeliveryOptimizationConfiguration` (legacy) | Current UI uses template `132f1027-0325-45e0-854a-6955cd3c68c0_1`; separate `settings_catalog_configuration_policy` resource. Case 65 covers the legacy API. |
| Device firmware configuration interface | `windows10DeviceFirmwareConfigurationInterface` | 66 |
| Device restrictions | `windows10GeneralConfiguration` | 67 |
| Domain join | `windowsDomainJoinConfiguration` | 68 |
| Edition upgrade and mode switch | `editionUpgradeConfiguration` | 69 |
| Email | `windows10EasEmailProfileConfiguration` | 70 |
| Endpoint protection | `windows10EndpointProtectionConfiguration` | 71 |
| Imported Administrative templates | `groupPolicyConfigurations` | Separate resource: group_policy_configuration |
| Kiosk | `windowsKioskConfiguration` | 72 |
| PKCS certificate | `windows10PkcsCertificateProfile` | 49 |
| PKCS imported certificate | `windows10ImportedPFXCertificateProfile` | 73 |
| Properties catalog | `inventoryPolicies` | Separate resource: settings_catalog_inventory_policy |
| SCEP certificate | `windows81SCEPCertificateProfile` | 74 |
| Secure assessment (Education) | `windows10SecureAssessmentConfiguration` | 75 |
| Shared multi-user device | `sharedPCConfiguration` | 76 |
| Trusted certificate | `windows81TrustedRootCertificate` | 08 |
| VPN | `windows10VpnConfiguration` | 44 |
| Wi-Fi | `windowsWifiConfiguration / windowsWifiEnterpriseEAPConfiguration` | 32,38 |
| Windows health monitoring | `windowsHealthMonitoringConfiguration` | 77 |
| Wired network | `windowsWiredNetworkConfiguration` | 78 |

## iOS / iPadOS

| UI template | API type / collection | Evidence / remaining work |
| --- | --- | --- |
| Custom | `iosCustomConfiguration` | 03,21,61,62 |
| Derived credential | `iosDerivedCredentialAuthenticationConfiguration` | 79 |
| Device features | `iosDeviceFeaturesConfiguration` | 02,11 |
| Device restrictions | `iosGeneralDeviceConfiguration` | 01,15 |
| Edition upgrade and mode switch | `iosUpdateConfiguration` | 105; live chooser opens `policyType/IosUpdate` |
| Email | `iosEasEmailProfileConfiguration` | 46 |
| PKCS certificate | `iosPkcsCertificateProfile` | 47 |
| PKCS imported certificate | `iosImportedPFXCertificateProfile` | 80 |
| SCEP certificate | `iosScepCertificateProfile` | 05,54 |
| Secure assessment (Education) | `iosEduDeviceConfiguration` | 81; live chooser opens `policyType/IosEducation`, titled Education |
| Trusted certificate | `iosTrustedRootCertificate` | 04 |
| VPN | `iosVpnConfiguration` | 42 |
| Wi-Fi | `iosWiFiConfiguration / iosEnterpriseWiFiConfiguration` | 30,36,55 |
| Wired network | `iosWiredNetworkConfiguration` | 82 |

## macOS

| UI template | API type / collection | Evidence / remaining work |
| --- | --- | --- |
| Custom | `macOSCustomConfiguration` | 06 |
| Device features | `macOSDeviceFeaturesConfiguration` | 83 |
| Device restrictions | `macOSGeneralDeviceConfiguration` | 84 |
| Endpoint protection (Deprecated) | `macOSEndpointProtectionConfiguration` | 85; live chooser marks Deprecated; legacy API lifecycle |
| Extensions (Deprecated) | `macOSExtensionsConfiguration` | 86; live chooser marks Deprecated; legacy API lifecycle |
| PKCS certificate | `macOSPkcsCertificateProfile` | 48 |
| PKCS imported certificate | `macOSImportedPFXCertificateProfile` | 87 |
| Preference file | `macOSCustomAppConfiguration` | 51 |
| SCEP certificate | `macOSScepCertificateProfile` | 13 |
| Trusted certificate | `macOSTrustedRootCertificate` | 12 |
| VPN | `macOSVpnConfiguration` | 43 |
| Wi-Fi | `macOSWiFiConfiguration / macOSEnterpriseWiFiConfiguration` | 31,37 |
| Wired network | `macOSWiredNetworkConfiguration` | 88 |

## Android AOSP

| UI template | API type / collection | Evidence / remaining work |
| --- | --- | --- |
| Device restrictions | `aospDeviceOwnerDeviceConfiguration` | 53 |
| PKCS certificate | `aospDeviceOwnerPkcsCertificateProfile` | 89 |
| SCEP certificate | `aospDeviceOwnerScepCertificateProfile` | 41,99 |
| Trusted certificate | `aospDeviceOwnerTrustedRootCertificate` | 41,98 |
| Wi-Fi | `aospDeviceOwnerWiFiConfiguration / aospDeviceOwnerEnterpriseWiFiConfiguration` | 35,41 |

## Android Enterprise

| UI template | API type / collection | Evidence / remaining work |
| --- | --- | --- |
| OEMConfig | `mobileAppConfigurations` | Separate endpoint; existing Android managed-device app configuration resource explicitly excludes `supportsOemConfig=true`; separate resource gap |

## Android Enterprise — fully managed / dedicated / corporate work profile

| UI template | API type / collection | Evidence / remaining work |
| --- | --- | --- |
| Device restrictions | `androidDeviceOwnerGeneralDeviceConfiguration` | 09 |
| Wi-Fi | `androidDeviceOwnerWiFiConfiguration / androidDeviceOwnerEnterpriseWiFiConfiguration` | 33,39,57 |
| VPN | `androidDeviceOwnerVpnConfiguration` | 45 |
| Trusted certificate | `androidDeviceOwnerTrustedRootCertificate` | 39,100 |
| PKCS certificate | `androidDeviceOwnerPkcsCertificateProfile` | 90 |
| PKCS imported certificate | `androidDeviceOwnerImportedPFXCertificateProfile` | 91 |
| Derived credential | `androidDeviceOwnerDerivedCredentialAuthenticationConfiguration` | 92 |
| SCEP certificate | `androidDeviceOwnerScepCertificateProfile` | 39,101 |

## Android Enterprise — personally owned work profile

| UI template | API type / collection | Evidence / remaining work |
| --- | --- | --- |
| Device restrictions | `androidWorkProfileGeneralDeviceConfiguration` | 93 |
| Email | `androidWorkProfileGmailEasConfiguration / androidWorkProfileNineWorkEasConfiguration` | 94,95 |
| Wi-Fi | `androidWorkProfileWiFiConfiguration / androidWorkProfileEnterpriseWiFiConfiguration` | 34,40 |
| VPN | `androidWorkProfileVpnConfiguration` | 96 |
| SCEP certificate | `androidWorkProfileScepCertificateProfile` | 40,103 |
| PKCS imported certificate | `androidForWorkImportedPFXCertificateProfile` | 97; current work-profile chooser opens `policyType/AndroidForWorkImportedPFX` |
| PKCS certificate | `androidWorkProfilePkcsCertificateProfile` | 50 |
| Trusted certificate | `androidWorkProfileTrustedRootCertificate` | 40,102 |
| Move to Android Management API | `androidWorkProfileMigrationConfiguration` | 104; created unassigned in the UI, read through Graph, full typed curl updates in both directions verified, DELETE followed by GET 404 |


## Captured API behavior discovered during the expanded tests

- Shared-PC account-manager fields are coupled: `disableAccountManager: true` reads back with `setAccountManager: "disabled"`; the enabled combination is `false` with `"enabled"`. Curl verified both updates and deletion. Complete JSON must provide a consistent pair; the provider does not suppress this difference in state.
- Imported PKCS profiles require a valid explicit `intendedPurpose`, such as `smimeEncryption` or `smimeSigning`. Omitting the purpose can create a profile whose returned enum value is rejected on a subsequent full-payload PATCH.
- The authenticated application received HTTP 403 reading `/deviceManagement/inventoryPolicies`. BIOS, device configurations, configuration policies, and group-policy reads succeeded with the same credentials. Inventory lifecycle proof remains blocked; this is not an authentication failure.
- The tenant's managed Google Play inventory currently contains no OEMConfig profile and no identified OEMConfig app prerequisite. No app was added and no existing profile modified.

## UI-discovered migration type

The live Android Enterprise chooser exposes **Move to Android Management API**, and its creation wizard uses `policyType/AndroidForWorkMigrationPolicy`. On 2026-09-29, a uniquely named, unassigned profile created through that wizard was located with a filtered Graph GET under `/deviceManagement/deviceConfigurations`. Its root type was `#microsoft.graph.androidWorkProfileMigrationConfiguration`, with the writable field `disableMigration: false`. This type was absent from the public metadata snapshot. A subsequent disposable curl profile verified complete typed PATCH requests setting the flag to both `true` and `false`. A partial untyped PATCH was rejected. Both temporary profiles were deleted and subsequent GETs returned 404. The resource already sends complete typed settings, so no type allowlist or endpoint switch is needed. Case 104 and its example retain empty assignments; they do not trigger device migration.

## Verified relationship and state behavior

- iOS wired network profiles reject certificate `@odata.bind` fields in the create body with HTTP 400. Creating the profile first and PUTing each typed navigation `$ref` succeeds; GET verifies both certificate IDs. Case 82 creates its own trusted root and SCEP dependencies.
- Android owner/work imported-PFX profiles expose `rootCertificate` in inherited metadata, but their navigation GET routes return HTTP 400 (no matching route). These types do not use that read path. Cases 91 and 97 reproduce the rejection in mocks.
- Those imported-PFX endpoints append an identical Any Purpose key usage on create and full PATCH. Curl verified the behavior across repeated PATCHes. State removes identical duplicates for these types and compares the complete settings payload; case 106 checks that a different key usage still produces drift.

## UI observations (2026-09-29)

- The Windows chooser contains all 22 supplied labels. Delivery Optimization opens `Microsoft_Intune_Workflows/TemplateWizard.ReactView/templateId/132f1027-0325-45e0-854a-6955cd3c68c0`. Graph GET of configuration-policy template `132f1027-0325-45e0-854a-6955cd3c68c0_1` returns HTTP 200, display name Delivery Optimization, family `deviceConfigurationPolicies`, and 29 setting templates. The existing `settings_catalog_configuration_policy` resource accepts this template family and template reference. Its lifecycle is outside this flat resource. Windows Endpoint protection opens `policyType/Windows10EndpointProtection`.
- The iOS chooser contains all 14 supplied labels. Edition upgrade and mode switch opens `policyType/IosUpdate`. Secure assessment (Education) opens `policyType/IosEducation`, titled Education, with Teacher certificates and Student certificates settings. Those labels do not mean Windows edition upgrades or the Windows secure-assessment payload.

- The macOS chooser contains all 13 supplied labels and explicitly marks Endpoint protection and Extensions as Deprecated. Their APIs remain creatable and are covered as legacy profiles; the examples recommend Settings Catalog for new deployments.

## Validation results

- All **106 resource unit cases** passed on the final implementation, with **83.8% resource statement coverage**. This includes the imported-certificate drift regression.

- All **95 acceptance cases** passed individually at Go and Terraform parallelism 1. Cases 65–105 add 41 live scenarios to the prior 54. Imports, full updates, empty plans, and destruction checks are included. Cases 74 and 82 were rerun after adding their certificate dependencies.
- The Android migration profile remains unassigned throughout; its test changes only the policy flag and verifies import, empty plans, and deletion.
- The AOSP chooser contains all five supplied labels: Device restrictions, PKCS certificate, SCEP certificate, Trusted certificate, and Wi-Fi.

- The Android Enterprise chooser matches OEMConfig plus all eight owner-profile and nine personally owned work-profile entries in the supplied list. The migration wizard/API mapping was captured separately as described above.

The cases exercise representative writable configurations for each template family, not every possible field combination. Derived-credential profiles require an external issuer for device deployment; their lifecycle updates metadata. The kiosk example configures Calculator with automatic logon. Classroom education includes its nested certificate configuration.

- All **96 standalone Registry examples** pass `terraform validate` individually using the rebuilt local provider. The template references each file once in Android, iOS/iPadOS, macOS, Windows order; all 145 resource addresses across the collected examples are unique. Focused correctness linting reports 0 issues at concurrency 1.

- `make userdocs` regenerated the Registry page. It contains every example verbatim, is **313,344 bytes**, and passes scoped `tfplugindocs validate` against the freshly exported provider schema. Full-repository docs validation still fails on three pre-existing oversized pages (targeted managed app configuration, Settings Catalog configuration policy, and Settings Catalog inventory policy); unrelated generated changes were restored.

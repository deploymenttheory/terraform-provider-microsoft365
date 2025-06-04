# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.15.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.14.1-alpha...v0.15.0-alpha) (2025-06-04)


### Features

* Add macOS LOB app publisher validation and update gitignore ([c9cd922](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c9cd922080ec768367bc79f1fc73632c9248d2f4))
* added group license assignment, user license assignment resources and datasource license subscribed SKUs ([bede181](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/bede1812ddd65df921206baf67602e02693e1f2a))
* added group ownership assignment with examples ([18810ea](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/18810ea49955060d9bfb87947a2b25a4e560c3ba))
* added group ownership assignment with examples / renamed user license assignment to follow resource naming pattern ([#497](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/497)) ([6e2fdf3](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6e2fdf3e6411537f679f5201abf9db426f82fd5b))
* added GroupMemberAssignment resource with examples ([21dde58](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/21dde582a0048b75e462dd881831d08706ef8069))
* added GroupOwnerAssignmentResource ([4f9bfd6](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/4f9bfd60e1013b7030ac31187b42fdd5224d78a7))
* added GroupOwnerAssignmentResource ([#499](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/499)) ([9bc5dee](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9bc5dee1a016aa39f78acbee5b91236d4c95d267))
* added groups with examples ([ffc38e3](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ffc38e3db9b488cbfb86b062d40bc946c9ebe22a))
* added groups, and group membership assignment, updates to doc descriptions ([#495](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/495)) ([0601507](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/0601507bf998791c5f74ae433005a492d0ecf1b0))
* added macos_dmg_app fixes and updated resource examples ([1da2c8d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/1da2c8d26bd0963fdd48d72ff384cb0245a62be2))
* added multiple new resources for license assignment and fixes for macos dmg app ([#491](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/491)) ([8d09b67](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8d09b67ec036a30817aecb7745d050adc96c880c))
* added multiple onboarding resources ([602135e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/602135e85aeab9ad2f8923c437004be647357d58))
* added windows autopilot and assignments thereof ([#488](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/488)) ([5ffe081](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5ffe08175062fa490c4b7babc5e8ecf658cc3ba3))
* half implementation of cloudPC as a datasource ([68e7a43](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/68e7a43dc062a7808d2265f7b76c7a354fff5d85))
* initial implementation of mac_dmg_app ([35b6fe6](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/35b6fe60dd97db4c75b442387d26c26680f61afd))
* refactored retry mechanism for create update func calls to read and fixed resource name issue in tflog ([6f3842b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6f3842b684eab3bcde146a90a9afd4b7b1c1cb34))
* refactored retry mechanism for create update func calls to read and fixed resource name issue in tflog ([#498](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/498)) ([5257e5e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5257e5ef5653b0ed0ca2ba5fe4ff1c70c3a5cf29))


### Bug Fixes

* for examples ([244a04b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/244a04b7ef7521d33ffcf42e0af0015d4d2791c8))
* improved markdown descriptions for documentation generation ([16bc1d2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/16bc1d2d7324c2a5c0a224a274e73100648b4caa))
* rearraging provider docs ([a0aab65](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a0aab656a60e99ef2554e053d01fdad5c4539016))
* renamed user license assignment to follow resource naming pattern ([29f4c8b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/29f4c8b163dd7e1d95dfe045831fc64e125db76f))

## [0.14.1-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.14.0-alpha...v0.14.1-alpha) (2025-05-26)


### Bug Fixes

* added close stale issues maintainence pipeline ([1d90efd](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/1d90efd55aab386abff032eda94948987f69d50f))
* added close stale issues maintainence pipeline ([#475](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/475)) ([40affe2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/40affe22d27c45128e8b4f7887419f10d41efc58))
* docs ([0a676f8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/0a676f87825c2d75b2c3ed5175d48474e23ec28d))
* for doc gen pipeline ([f959ac8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f959ac87f1d6ea3e0d66b9979502a5c1f914466f))
* for doc gen pipeline ([d3fecdf](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d3fecdf388ed6ad8e82fffffce38cbdd0126553a))
* for doc gen pipeline ([#471](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/471)) ([c93301e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c93301e7f95b32b56ed812ba5f373d2c2d3b0d80))
* for doc gen pipeline ([#473](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/473)) ([3c58522](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3c58522f6f1522eca9da2b4ee613b32b96319368))
* unit test trigger ([2bf41c4](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2bf41c4a70e4a47e268bb9d68950dda48d24bfe0))
* unit test trigger ([#474](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/474)) ([e426630](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e42663025764a290342fde1a4f489a19931e94a1))

## [0.14.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.13.0-alpha...v0.14.0-alpha) (2025-05-25)


### Features

* added assignments for windows updates as a seperate set of resource types ([837ee80](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/837ee809bc587496f2e8be44996db786cbd0d10f))
* added datasource microsoft365_utility_macos_pkg_app_metadata ([2c2cbca](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2c2cbca7feb8d1336c9327e487767ecec1e4bc35))
* added managed device clean up rule with examples ([fc33266](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/fc332665c817df572c0d349ec4d68ccf1b037570))
* added managed device clean up rule with examples ([#463](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/463)) ([f1a8777](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f1a87770fa8e61a0edbc8370876cf86b4cb07850))
* added microsoft365_graph_beta_device_management_terms_and_conditions ([335e4a1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/335e4a1ec8441e7f49b980af60de4d07e2afea8f))
* added microsoft365_graph_beta_device_management_terms_and_conditions ([#466](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/466)) ([16849a7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/16849a77e9262aeddbae84bc43b55552da62c146))
* added microsoft365_graph_beta_device_management_terms_and_conditions_assignment with examples ([c3db9ac](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c3db9acb7c2fb11122658bdf817c2ff57af01fec))
* added microsoft365_graph_beta_device_management_terms_and_conditions_assignment with examples ([#468](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/468)) ([6e407c0](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6e407c067f71cbc2f01cef185adaff9f214774db))
* added microsoft365_graph_beta_device_management_windows_feature_update_profile_assignment with examples ([74f286e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/74f286eb668a342f700b1ce3c5da9671695ffdcb))
* added microsoft365_graph_beta/microsoft365_graph_beta_device_management_windows_driver_update_profile_assignment with examples ([863774c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/863774cbb6a0b9944c1d79b2e87ac2f4ad66a6b5))
* added mobile app data source with odata queries ([3a1c0f4](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3a1c0f414e70ea9c55f36b31c1c48e85feb2aa02))
* added operation approval policy ([40b307c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/40b307cfd8c1ec89fd342c7a29f20254d2cacaef))
* added resource operation for intune rbac with examples ([4f36c5f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/4f36c5f5b50c9d19f3cdcdcc2507d8d8929a5931))
* added resource operations with examples for intune rbac ([#465](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/465)) ([1c31571](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/1c315714f220e136300bf4da2f0c029dd39e2c11))
* added setting catalog in native hcl ([596d6c9](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/596d6c9f0454810004ff30a3177a5c4ec2bc4ebc))
* added settings catalog in native hcl ([#461](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/461)) ([6e5e693](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6e5e6937324d83df7fe5f96445dc45f851bf771f))
* added validation for mobile app assignments for 8 invalid scenarios ([cc88fd1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/cc88fd104a35cf071e4a5693a84420f57d1151d8))
* added windows quality update profile assignment resource with examples ([a49a109](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a49a109323d23f0c8e0ed2be4afea53a77b1b05d))
* mobile app assignments ([80c892e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/80c892ec95186025afc2c6f7edbd61f7ed1863e4))
* mobile app assignments resource with examples ([#453](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/453)) ([498ff50](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/498ff5087f942dffa701be9da736c8d3e1a5dde2))


### Bug Fixes

* final fixes for macos assignments ([8d998cf](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8d998cfc27ccb4aa1b4bab3be60823f479b7cf85))
* final fixes for macos assignments ([#451](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/451)) ([7e3b853](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/7e3b853731d22cc832840ceab3ca9c7a256418c4))
* for docs path ([2fce669](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2fce66923162574af62abe0c7236cf9da8bd76ab))
* for mutex handling for SetCommittedContentVersion ([058bf32](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/058bf32f301bc9b5082f668427aaf4ca3cff3878))
* for successful stating logging with resource name and resource id ([8c81932](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8c819320ff467a553640c8ff5ce9c78c5f7cf76a))
* for winget assignments ([8a2bdd0](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8a2bdd03111f263f47db1614ba1e0fb8deb5a650))
* migrating winget app assignments to sets ([2272cc6](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2272cc64166db8c02accd2a243d79a646fc3dfe4))
* minor fix for sorting order of assignments ([fb50e0f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/fb50e0fbac88e72d4d784f400131449be1c4ff07))
* moved mobile app assignment validation to create and update funcs to fail early if any hcl mistakes before any rresources get deployed ([8c79958](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8c79958155e1445d959487a27cd66ddca38947b3))
* removed mutex, going with parallism of 1 instead ([dd0e6d8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/dd0e6d822dcf3eb6fd43bbdea105d62020dc488f))
* typo fix for datasource ([5c60f13](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5c60f13a9ce6d8701d9310ea41f07e1104524967))
* win_get_app now skips property 'InstallExperience' and 'PackageIdentifier' as it cannot be patched during update operations. and added support for app categories ([8b5a4a5](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8b5a4a59a43e6bddb6e82a3763713a85bcfe4047))

## [0.13.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.12.0-alpha...v0.13.0-alpha) (2025-05-12)


### Features

* mobile app assignments can now be set in any order with modify plan (diff supression) ([c8f97ef](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c8f97ef7201c1a79efc3850add74eac0e78468be))


### Bug Fixes

* added mutex to all api calls to defender against Kiota's middleware (e.g., HeadersInspectionHandler) when modifing shared header maps during HTTP request processing. appears to be a bug in the latest version ([976ad22](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/976ad22dea6b1ddf6e3027714a04c449b73ac11c))
* added mutex to all api calls to defender against Kiota's middleware (e.g., HeadersInspectionHandler) when modifing shared header maps during HTTP request processing. appears to be a bug in the latest version ([#450](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/450)) ([10327c8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/10327c856db27c79ae6883677a48049c6caf2a1e))
* base resource updates for macos_pkg_app doesnt use redundant file size checks ([6162a1d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6162a1d4dc9860de2aa812ff386289bd6fc387f7))
* defined default values and plan modifier for role_scope_tag_ids ([e8909e0](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e8909e0fc876dbe03422bfb26943ac872cfab432))
* for parrallism issues with concurrent header map settings in client mobile apps ([49561f6](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/49561f64cea87bd23eab9dda180824ddfd8934de))
* for role_scope_tag_ids ([#444](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/444)) ([b1e0fd8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b1e0fd819d4956ee065e96d4300fb24785719dc2))
* removed conditional access. complete refactor required with latest sdk version ([4cb550d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/4cb550d56f376cb30bde71c71b4699e63653e07c))

## [0.12.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.11.0-alpha...v0.12.0-alpha) (2025-05-11)


### Features

* added refactored unit tests for provider build and credential factory ([bd7691c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/bd7691c76d6403dadb4eb68d76ba37093f506b9b))
* added refactored unit tests for provider build and credential factory ([#439](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/439)) ([c191a93](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c191a9345c962f5f19a1bfd76dcafb3fb1deec01))
* mocking scaffolding ([f1be6d1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f1be6d1f6fddfca091d34df19750da364545c1b1))
* mocks scaffolding ([#440](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/440)) ([9ee57df](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9ee57dfc4eecebb95c35c0e60686925d3228c2c6))
* refactored data source for maospkgapp ([#442](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/442)) ([fec56b7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/fec56b72934dc5c9fbd1dad05e8d690fab0793ab))
* refactored macospkg app datasource with example ([cbc1681](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/cbc1681feca36a096ea63b69dc61540178b21447))


### Bug Fixes

* examples so that device_management examples are reflective ([d6856b4](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d6856b480ef5a0c9f7277592bc391eb47487b443))
* examples so that device_management examples are reflective ([#441](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/441)) ([6062154](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6062154d0eedb88177fec252a4c5fb5f925ca24a))
* for enum handling within GetDeviceEnrollmentConfigurationType ([#437](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/437)) ([896c291](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/896c2918aa799ec20127a79538a06cfe1cff892e))
* for handling enum ([ce13f69](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ce13f693c89723f9766c02519f17596ea328d7e2))
* for handling enum ([6e2f1ff](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6e2f1ff38ef16ddbcd405eff88c380371e315856))
* refactored docs for repo restructure ([546ebb4](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/546ebb4ff5dfe06b9b6c9febad4223a4acf71abf))
* refactored docs for repo restructure ([#435](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/435)) ([3b14e78](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3b14e78f1db1659b176983abced03e04b32195d6))

## [0.11.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.10.0-alpha...v0.11.0-alpha) (2025-05-07)


### Features

* added 1st itteration of oidc support for github actions and azdo ([3a2c14b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3a2c14bdcc51b039fc9269e9febf8edae89fa866))
* added application category resource and datasource with examples ([55a099a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/55a099a087f713f4e4c01558e4aec5fdfaf303f8))
* added datasource graphBetaDeviceAndAppManagementWindowsQualityUpdateProfile ([9f6f79d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9f6f79d7472eef5d7f76973c8195046b65c4462b))
* added device_and_app_management_windows_update_catalog_item datasource with examples ([655e804](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/655e804d7ea7989c549fa86a8925ea82b0b33d46))
* added graph_beta_device_and_app_management_windows_driver_update_profile_assignment ([73dc481](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/73dc481446cb3c88284ef6b5fcac719da08bfb81))
* added graph_beta_device_and_app_management_windows_feature_update_profile with example ([050cec1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/050cec1724781fa5bdf81b14c9c5890f5a0386e4))
* added helper SetObjectsFromStringSet for object construction ([e2b110e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e2b110e4f7b4aa73aa2e823a5ae9b58f5b597e8c))
* added macos pkg app datasource ([a65ca03](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a65ca038103e064225b7d136205179b7fbb04fad))
* added macos pkg datasource ([f0f6144](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f0f61447b307f221cdf6f0b0d65f74d8e4f6beb6))
* added microsoft365_graph_beta_device_and_app_management_windows_update_profile_assignment with example ([44e3539](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/44e3539f943013c5a0414e2053f56a1c02757f04))
* added timeouts to datasources ([#397](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/397)) ([100840f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/100840f1bff363a598fd9a76e26183c16c086150))
* added windows driver update profile and inventory resources and datasource ([#416](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/416)) ([8e3272d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8e3272dbc77296cdd874a11e7c7693103173419b))
* added windows driver update profile and inventory resources and datasources ([85ff58f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/85ff58f605f939be963c3aba58ebd44c0ac41f0f))
* device_enrollment_configuration with examples ([#420](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/420)) ([88516db](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/88516db1ff74964069d9965f29c21972bb147144))


### Bug Fixes

* added examples and more fixes for windows_remediation_script_assignments ([9827467](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9827467ac9972380540ccb29490670d5db960c87))
* added read timeouts for datasources ([243e6b2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/243e6b2fdc92441d54949d0d1bd9544b4b479c1e))
* categories fix for macos app categories and stating issues for infured included app values ([defa936](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/defa93679fa6bed0bc6065678987c76429383133))
* content version ([74c417b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/74c417b8312e6353bbbdbf78fd012d9fa96b39e1))
* for deps bump of go sdk version ([54a119e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/54a119edcca117fa697b21c59f2ba35cdda12ae2))
* for example path ([771bb1c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/771bb1c2792a16955e05f124cbe7a4e6d88cd3bb))
* for macos pkg apps ([#395](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/395)) ([deb9309](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/deb9309edc8a8aff4d33a022102f30f9c90eb679))
* for mobile app assignment id's and content version stating ([a57a269](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a57a26945a184c294e4c37c41b6b3457dbdad180))
* Go Lint configuration yml for provider project types ([#424](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/424)) ([b92b8ac](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b92b8ac28100a40611b0cca6e191750db9b8d277))
* go-lint ([#423](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/423)) ([57ba99d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/57ba99dd57bf8b26e9785767cf05449df1b2fca2))
* m365 app installl options ([e2db1a7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e2db1a7982f5cced9b3f346b7f2df8a30c92695a))
* migrated microsoft365_graph_device_and_app_management_m365_apps_installation_options to graph v1.0 ([a3c9682](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a3c968250d67787711d46b2589facd62cf3ab8b3))
* numerous edits and tweaks ([dccd487](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/dccd4873fc4dec676b2a0e3859c7a2a42d595b4f))
* numerous schema fixes ([ddb6e87](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ddb6e8780360a63c31354f0899314fd4ed51c111))
* pipeline step naming ([3fade5e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3fade5e96ab87097c3b8b18dc3e2f0ea5850a7a3))
* pipeline step naming ([#425](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/425)) ([8067017](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/80670170abf8aa817432805c0d4a786859655951))
* refining logic for handling app metadata ([9ecf223](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9ecf2236482cffbe5ef3cc67cf9e10e94068f345))
* restructured repo to follow api endpoint paths more precisely ([1b2c110](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/1b2c1108add87c7395eaed7cfce79887c7899f12))
* restructured repo to follow api endpoint paths more precisely ([#433](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/433)) ([6aed990](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6aed9906d8f76430c646909e6111c6d91b29f1b3))
* tidy up pipeline step naming ([102e356](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/102e35695ac0629381b740e08aec981cfdaf84a1))
* tidy up repo restructure ([b5b35cf](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b5b35cf4b001d4ba9ebc0cfed8913bc8a7054e12))
* windows remediation script assignments ([04243b6](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/04243b62ecb12ed869edec6caf334e6671c8fd90))

## [0.10.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.9.0-alpha...v0.10.0-alpha) (2025-04-08)


### Features

* 1st itteration of intune applications. ([1a83d27](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/1a83d277be16d0e08fd61d32919d07e39d776313))
* add debug logging for Graph API request bodies in constructResource function ([3fdbba4](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3fdbba463403254f7dbd8e3473d48f7a3d71e3a3))
* add macOS PKG metadata extraction functionality and update dependencies ([bedd812](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/bedd8121ef7b3c1bca5e5022ec2db27c88bb2f37))
* add package installer file source and extract metadata for MacOS PKG applications ([7cfd15a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/7cfd15a8b37a61d5a5cc18abb4f6c9c86739798c))
* add support for extracting files from CPIO archives and enhance payload extraction logic ([609fc22](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/609fc22d130a2480dbd533d737a22c08192138f1))
* add validation for MacOSPkgApp publisher field and implement requiredWith validator ([6000864](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6000864632b107c490e6dbe0d899766af23aafce))
* added 3 secuirty baseline settings catalog templates ([ff7b7de](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ff7b7de4f4d984bce801a182cd6394b5bcbcbd05))
* added debug logging to uploadToAzureStorage func ([ec57f66](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ec57f66dce52a0153ed2f8ead56154913a9710a5))
* added Encrypted file analysis for troubleshooting for macos packages ([10886cb](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/10886cb97abb6b533e22f7ef38eca7681e86a06c))
* added first itteration of schema builder for settings catalog ([7590e6f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/7590e6fabb8fd49a3b7f5111cd476ddcf5a6fff2))
* added intune device categories resource and data source with examples ([#389](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/389)) ([b8d33fb](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b8d33fba7b11121bf61514164b77d01858a9a74c))
* added macospkg resource to provider ([#347](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/347)) ([10001df](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/10001df8134fe806c661b94850531f98b245643f))
* added new type conversion helpers that handle null values ([8aaa14e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8aaa14ed15e4805bd63309d8de32db5b1266b284))
* added ps scripts for testing dmb, pkg and macoslob apps. also added some getters to help with data model definitions ([37877de](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/37877defb6ae270df8059273e31322bace8708e6))
* added schema Mutually Exclusive Attribute validator ([e2c7deb](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e2c7deb3b6b008c0556c296bf2f31eb8b7c3c2a9))
* added secuirty baseline templates support to settings catalog template resource ([8e4649b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8e4649b21793674f196d27c11608c73ffb87be38))
* added support for url based and local file based icons and pkg installers ([c75ca2c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c75ca2c77727f7ce35e3f5990eb6b90e488ab71f))
* added web sources for icon uploads for macos app ([dbc36d6](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/dbc36d6f5642ad011adaedf0729ab0fb87b88407))
* completed macos pkg app creation flow ([99a96fb](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/99a96fb43efc1b09483c5f98c29f21138027a338))
* enhance schema for package installer with validation and improved descriptions ([5488a54](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5488a54f7c3fbacf79ac359ab14c5109f3c38ab8))
* implement XAR file extraction and field parsing functionality ([c6dab71](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c6dab71a99f18b5a08708c60a3771cbe4914c8af))
* implement XAR file parsing and metadata extraction for macOS installers ([b91e24c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b91e24c2030639be1052b4b97292e821018b0b8c))
* implement XAR package reader and metadata extraction for MacOS installers ([5d70348](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5d703488093e25721b0ffae8b797ec6ad70b5aa7))
* implemented device categories resource and datasource with examples ([acf2626](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/acf262635ca9aeb5175f0744fc108dfd0694b027))
* updated macos pkg example ([6ecc553](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6ecc553ce50049b5afd46e72a4671dd27211c342))
* updated read and stating logic ([82a2cfd](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/82a2cfd8f2d39eb1ce758e39428170e77a94e91e))


### Bug Fixes

* added content version assignment ([0418ecb](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/0418ecb57f547d216f7121e1aa187a12cf3758be))
* added initial stating logic to mascos create func ([e99d915](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e99d915beef8bf56424eab4197663a712c0fe92f))
* added planmodifiers.UseStateForUnknownString(), for id ([dccde7e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/dccde7e78d92bf90a23bbdf7ad5673391baafc65))
* changed assignment optional settings to a pointer ([b3bb2d7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b3bb2d74af6c3732af8b3eba83cd95c42b255d36))
* extended file upload chunk time ([6f7ba33](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6f7ba3328934b9d7256c887034dc3f102b16d59e))
* for package imports ([95c9a15](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/95c9a156c5fbd9c53e6dc467c0150c011b84eeeb))
* for stating icons to use stringunknwon logic ([0006aec](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/0006aec1068ea0ceb599998b9183096525c585e8))
* improve comments and formatting in constants and assignment filter read functions ([e48f5c7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e48f5c7a4ca32e3e5db93abe0a243a47e81e159d))
* nil pointer dereference in func validateInstallTimeSettings ([b52d393](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b52d39323ffb12a8b76d2fee02a9994fb571be43))
* numerous fixes for constructors and stating logic ([4d09816](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/4d09816644f027bdbac3de734ff2984f6c0624a3))
* removed redundant struct field ([a7a8094](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a7a80945b91f07d4b5fb87871a388ef39e90c6df))
* testing new encryption strategy for pkg's with intune ([203c69e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/203c69e877f275ba56dce3e8673db86bdac26994))
* timeout usage of consts ([d94aec2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d94aec2d5e2da25d9de54f256477671590ea6582))

## [0.9.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.8.0-alpha...v0.9.0-alpha) (2025-01-18)


### Features

* add device management template type attribute with validation options ([d07cbb9](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d07cbb92532ebead04ed8f956b53017ef0b7e798))
* add new Windows Defender Antivirus policy templates to device management configuration ([45834c3](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/45834c33cb439eb5166976e58abba68f3bfb5f76))
* add settings catalog template and update related configurations ([a3b9a33](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a3b9a33347942f8545b4878d7de88dc195ab831e))
* add template for graph beta device and app management reusable policy setting documentation ([01334e2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/01334e26b1fa9092bc9bd141ca161f7bdc0f0bef))
* add validation for settings templates and extend policy configuration map with new templates ([b240513](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b24051395001714144e98e3c9d149164bca61db9))
* add Windows Firewall rules template and remove deprecated script ([4b264f7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/4b264f781e65bf776da40f17037a2a683b3860b1))
* added additional endpoint security templates to settings catalog templates resource + numerous doc edits ([#318](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/318)) ([d9f3ff2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d9f3ff22346ccffdc635070d316cd7c6323b03b4))
* added all remaining settings catalog templates for Intune Endpoint security ([#319](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/319)) ([e9a9320](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e9a9320ae6e7d7ba7eb68ebe86a60e7452948b0a))
* added settings catalog template to provider with examples ([#314](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/314)) ([3f062f1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3f062f113536da431e3abd606174cb33768e3698))
* enhanced setting catalog construction logic ([eec7669](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/eec7669ad0cb63f0561f4617ed0509ef2eb21f35))
* extend policy configuration map with new settings_catalog_templates for endpoint security ([1d76a96](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/1d76a96e6bcbed86574d59c4378406d9bbca8fd6))


### Bug Fixes

* update example file path for reusable policy setting documentation ([a944199](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a94419961b97490b0b187dc3b29c41b16ce48405))
* update settings key in resource configuration for Endpoint Privilege Management ([763c8ae](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/763c8aed4b7b9dedbc32d7a83397b3d08f1cc547))
* update Terraform version requirements in documentation and configuration ([97a200c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/97a200cbe4e9805c166e3690322e885d507eabd3))
* updated schema in settings catalog templates to reflect all secuirty templates supported with descriptions ([d9c1e7a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d9c1e7a255a0adb8483981fc149ea3c9d5dd61a5))

## [0.8.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.7.0-alpha...v0.8.0-alpha) (2025-01-15)


### Features

* add endpoint privilege management resource and example usage ([b7b4b3b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b7b4b3b0ff5d3fed63cd4aa52dce7f173690baa7))
* add Linux platform script resource and example ([bcd6530](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/bcd65300c18bed23339a10c1ac8eaf90ec148f3d))
* add plan modifiers to use state values for unknown attributes in reusable policy settings ([f8f2c14](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f8f2c1423a437ccc300857436af72cd6c6b16c37))
* add reusable policy settings data source and example usage for Endpoint Privilege Management ([d5a5d10](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d5a5d107b5fe54d3f9a35a8f9551577ba95e6898))
* add reusable policy settings models and modify plan handling ([5029495](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5029495db1a3cbd5465235d103aa40b9488f7acd))
* added endpoint privilege management resource and example usage ([#303](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/303)) ([48a2514](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/48a2514f4cccca52df60431c9699d2eda3d80b76))
* added linux platform script with examples ([#302](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/302)) ([6324ebe](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6324ebeafdb411527b0801664df04b20c85efe11))
* enhance reusable policy settings API calls with additional select parameters and improved debug logging ([fc03249](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/fc0324948032fe79fd83742eed43a7a530dafcf3))
* enhance setting instance handling in reusable policy settings resource ([40f98e7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/40f98e7e99025e64998ae645acf0cd23108001ec))
* enhance state handling for reusable policy settings and normalize JSON responses ([6c90118](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6c901185f05926f5e2a99f8bd9921ce6d7417ec0))
* implement custom DELETE request handling and refactor URL template configuration ([ec7c0a2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ec7c0a24f2a4fc16d6d996006c172a0e90fcbf51))
* implemented data source for reuseable policies for epm ([#312](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/312)) ([6b0bbbe](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6b0bbbed729938659e8193a82021bf1a65a20697))
* improve error handling and logging in StateReusablePolicySettings function ([0361f31](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/0361f314aedf74f34261ad94788995d39bacd3ec))
* refined logic for reuseable policies and updates to schema ([#305](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/305)) ([00bd476](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/00bd4769d1d44d85bcd7b1344492e11675ef2cb7))
* rename settingsDetails to settings for consistency across models and resources and to align with intune gui exports ([937cdbd](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/937cdbdbbef05e8f64ae25ea122916b485d1be8b))
* update reusable policy settings documentation and add example resource for Endpoint Privilege Management ([a0aa693](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a0aa6936c2098b84ab0700ce90e67441b0cef877))
* update reusable policy settings model and integrate into provider resources ([f1e528d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f1e528df99b148ac7d3887ff8d71ae59cd242715))


### Bug Fixes

* add comment to clarify ConfigurationPolicyTemplates requirement in state_base_resource.go ([f38ebe7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f38ebe7c314a50359b1fc3830b68b31df1b23976))
* correct resource naming for reusable policy settings to singular form ([5ccd617](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5ccd617bb4ce09e98c0079656afd4192ec908b15))
* correct resource naming for reusable policy settings to singular form ([#313](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/313)) ([db7ab1e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/db7ab1e53e04586f0a89b161243cec18abbc62c3))
* remove unused plan modifiers for created and last modified date attributes in reusable policy settings ([878f800](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/878f8000a328967da6ffe20f0ba9195729ee3dd0))
* Update Graph Metadata - 2025-01-05_00-07-29 ([#298](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/298)) ([43a831d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/43a831d447bda7122e00c8a7b31ec18ad8cba93d))
* update PowerShell script links in Markdown descriptions for reusable policy settings and settings catalog ([a6c0292](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a6c02923de3fdfa33bd412746c52109e30544a5a))

## [0.7.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.6.0-alpha...v0.7.0-alpha) (2025-01-04)


### Features

* add mapping functions for various remote assignment settings to Terraform ([e24fb0f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e24fb0f0851a662a562a210a1902d74343641ed9))
* add resource documentation and modify plan handling for various device management scripts ([644dbc9](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/644dbc996981837466b0d30d9bf00e6f2721a1e9))
* add StringListToTypeList function for converting string slices ([#293](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/293)) ([cd62438](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/cd62438c70bdee68a2b1ec30184d098c1c72d17a))
* add StringListToTypeList function for converting string slices to types.List ([dd28488](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/dd28488c7fd29d58c8331a9e904155acaa94f5aa))
* add validation for mobile app assignment ordering ([c0fdeb7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c0fdeb744e9ac2218a7fc54d07c40310cc2ead84))
* added mobile app assignment schema for all app types ([ec8edc1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ec8edc1c22f91bab217489985541be0fb14410c7))
* added the option to manually define winget app metadata along side auto generation + plan modifers ([#281](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/281)) ([6049e8e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6049e8ec1b32d42c1449f263b9723dc4f8d6c372))
* enhance mobile app assignment configuration with new settings and sorting logic ([81dd048](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/81dd0486edf1bda8e3b59490faf464059b4834ec))
* implement mobile app assignment validation and update related constructors ([5bbc084](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5bbc084ac6f42644c94b8bd073c4866c04fdacbb))
* implement validation for mobile app assignment settings and restart timing relationships ([28aa0b3](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/28aa0b34a27577a132e6c4067efe88da3fbd2ee8))
* implement validation for mobile app assignment settings and restart timing relationships ([#292](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/292)) ([ed375cc](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ed375cce3138b1b1bd35844ff59d0842904843c2))
* refactored stating and constructor func patterns to be more concise and leverage correct lib and project helpers ([#277](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/277)) ([694ade9](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/694ade95bf2d9d761e0c187f868836db8054c96d))
* standardized stating structure and constructor approach ([1db44c8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/1db44c86b3fd9afc5c8fa72a916c461c3af0bf4f))
* standardized stating structure and constructor approach ([#286](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/286)) ([e798fd5](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e798fd567519686f46aec1bbd738edc303af456b))


### Bug Fixes

* add target type field to AssignmentTargetResourceModel for improved clarity ([41cbea2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/41cbea24ae15fcb023281076179bdb14529c601d))
* added function comments for mobile app assignments ([#294](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/294)) ([d742bd4](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d742bd4efd35b70be286acaadf94e0e35df36a2b))
* bug fixes for conditional access policies and fixes for crud permissions ([#267](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/267)) ([4d7f8d8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/4d7f8d86e44e9d1ee88723f25d8279fc4de514ee))
* centralised settings catalog assignments ([#287](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/287)) ([751e2f3](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/751e2f3277a1adfc4f19a9b428a3eb9b1938400a))
* for constructAssignmentTarget within mobile app assignments ([#289](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/289)) ([5ee4531](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5ee45311814e9a2bc56810be63c4f13f78dba663))
* numerous fixes in docs and pipelines ([#268](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/268)) ([b243425](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b243425bb0bbc2ac00a69d0bbce03c27f2f51f3f))
* refined gorelease pipeline validation ([#269](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/269)) ([18ff5aa](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/18ff5aadafd3449983fa655deb5fbab6cd767427))
* remove obsolete mobile app assignment resource from provider ([#297](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/297)) ([af9b4cf](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/af9b4cf84eebb03e25c84868501bacd1c2be5d73))
* removed icons plan from plans ([47bc42a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/47bc42a81d109f0a411456b22a4c92645219d1b5))
* rename struct for clarity and add configuration policy assignment constructor ([32b5b65](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/32b5b65b48aa167e4ae24ddd2778e928fe6f1d6e))
* reorganised repo to use graph_beta and and graph_api consistently for all package naming ([#290](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/290)) ([54fc3aa](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/54fc3aa0e9091bb56af3142ab2bd1a13a92e442f))
* replace constructAssignment function calls with specific constructors and remove obsolete construct_assignment.go file ([#288](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/288)) ([2475c9f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2475c9f954dc9ab06b29aa18bb04a29d00c4ec3c))
* sorting mobile app assignment stating logic ([#295](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/295)) ([9de00c2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9de00c2792a64a55b8c0d7cfd4541a1e8242048b))
* standardised use of object throughout crud functions ([#259](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/259)) ([bbc3dc3](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/bbc3dc3c0a8b3a3b26d16b86b65d0300912995db))
* streamline property handling in WinGetApp resource mapping ([#278](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/278)) ([025a5b1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/025a5b144d8b0255e38e7d430fa751e911c5c7eb))
* tidied up repo and add validation for mobile app assignment order ([#296](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/296)) ([036e95a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/036e95a1b0b4364c5b27fcbe02ba26b291e0d954))
* update import paths for device and app management resource to use graph_beta and graph_v1.0 ([#291](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/291)) ([62aa968](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/62aa9686904eca813e6e921b34265a8354a34074))
* update import paths from 'construct' to 'constructors' for consistency ([d878700](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d8787001898571f07e5a1d65fabcef6c679e00ef))
* update import paths from 'construct' to 'constructors' for consistency ([#284](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/284)) ([d878700](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d8787001898571f07e5a1d65fabcef6c679e00ef))
* updated docs ([#283](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/283)) ([e76a368](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e76a368632c7a8e077c016f89e35978c5d9299aa))
* updated win_get examples ([#282](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/282)) ([e54f1fc](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e54f1fc540368183230afed9dfe0e214ba99b353))
* various small fixes to docs and pipelines ([#271](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/271)) ([d244e16](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d244e1633fcbe824f19a8d3e82ffe5ff8f3d8a0c))

## [0.6.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.5.0-alpha...v0.6.0-alpha) (2024-12-16)


### Features

* added tf-docs auto generation pipeline for merge into main ([#257](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/257)) ([2b4e9ad](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2b4e9add53640a4c950248b23deec223c60e181b))
* intune role scope tags + added release please ([#246](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/246)) ([f09c60c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f09c60cbfac38b0bf78e83aaf1b5d18903714d34))
* refactored datasource examples to support search by name or by resource id ([#255](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/255)) ([a07bd03](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a07bd035ac0e6b0a2324673adc59ee31eff7136e))


### Bug Fixes

* added retry logic from sdkv2 to settings catalog resource types ([#250](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/250)) ([deb384d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/deb384d335dba7cf219f099676213d957c7663c5))

## [Unreleased]

### Added

- Added xyz [@your_username](https://github.com/your_username)

### Fixed

- Fixed zyx [@your_username](https://github.com/your_username)

## [1.1.0] - 2021-06-23

### Added

- Added x [@your_username](https://github.com/your_username)

### Changed

- Changed y [@your_username](https://github.com/your_username)

## [1.0.0] - 2021-06-20

### Added

- Inititated y [@your_username](https://github.com/your_username)
- Inititated z [@your_username](https://github.com/your_username)

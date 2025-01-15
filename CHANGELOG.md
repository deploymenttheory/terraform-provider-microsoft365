# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

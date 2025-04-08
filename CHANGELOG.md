# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

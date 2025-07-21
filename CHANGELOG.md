# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.22.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.21.3-alpha...v0.22.0-alpha) (2025-07-21)


### Features

* added settings catalog to hcl generator script ([#615](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/615)) ([762ab6f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/762ab6f349d9c9ce111deecd4ba953c42c164a84))


### Bug Fixes

* final set of fixes to handle a list assignment response as a set ([#623](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/623)) ([47fb051](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/47fb051ea448c5d6da025d9aaf5f80a02594d643))
* for settings catalog simple settings with secret values and added schema validator to only accept "notEncrypted" for requests ([#617](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/617)) ([f40a8dd](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f40a8ddd8753cdeba4d755d6c1972a78c4456959))
* for state handling of assignments for macos_platform_script , macos_custom_attribute_script , windows_remediation_script ([#620](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/620)) ([decfee9](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/decfee95177344aa29c5aa58e863a22a0c166f1d))
* specify provider name in tfplugindocs go:generate directive ([#621](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/621)) ([5cffa73](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5cffa737074d85cc0f7f03716b58ffd98912d6c7))
* updated powershell hcl exporter for settings catalog to support complex types ([#619](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/619)) ([9d20c88](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9d20c88e5aba1a2adc7e238affb1e2f2e865af40))

## [0.21.3-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.21.2-alpha...v0.21.3-alpha) (2025-07-16)


### Bug Fixes

* added assignment removal to deletion of CloudPcProvisioningPolicy resource ([#610](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/610)) ([aeb93c7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/aeb93c78439f676010ef1d2cdb3cf1ca71ff830c))
* multiple assignment fixes for macos scripts and cloudpc provisioning policies ([#612](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/612)) ([769bb85](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/769bb85cc8d573418652b3e21646f1488886445b))

## [0.21.2-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.21.1-alpha...v0.21.2-alpha) (2025-07-16)


### Bug Fixes

* doc gen for device_management_windows_driver_update_profile and device_management_windows_driver_update_inventory datasources ([#607](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/607)) ([bacccc7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/bacccc7cc7d6305753c15085e0724d215f3c889f))

## [0.21.1-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.21.0-alpha...v0.21.1-alpha) (2025-07-16)


### Bug Fixes

* added page itteration to settings_catalog to handle &gt;25 config items within hcl. ([e9bfa49](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e9bfa4956acc56e650b1d4d845f5b902e69ed2fe))
* added page itteration to settings_catalog to handle &gt;25 config items within hcl. ([#603](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/603)) ([50947df](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/50947df300c70b48cc20ff59424cc382ab98fb06))
* added template_reference block to device_management_settings_catalog_configuration_policy and updated examples ([0460890](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/046089081783545b364328eed01d2d535075f9f8))
* for macos_custom_attribute_script resource updates ([c175a29](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c175a29db1dd7d5461a525d834fed444cfbe096c))
* for macos_custom_attribute_script resource updates and added template_reference block to device_management_settings_catalog_configuration_policy and updated examples ([#601](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/601)) ([38e1a4c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/38e1a4c1e87455be5bd33a1a1d50dafc56c87eba))
* setValueTemplateReference unreachable case matching ([45709ca](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/45709ca36b7a74f4f10f007e6f0612ca0403a731))

## [0.21.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.20.0-alpha...v0.21.0-alpha) (2025-07-14)


### Features

* added 1st itteration of enterprise app catalog ([65bcd68](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/65bcd68cb69394888fe00c71a470e30c388f634b))
* added datasource utility_windows_msi_app_metadata with examples ([d547d5c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d547d5c2121ec4e5cf04d72ea023df1cd9c6f29a))
* added datasource utility_windows_msi_app_metadata with examples ([#593](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/593)) ([370093b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/370093bd1359d99983d5b74d2e64bedd91ee85f9))
* added ios_store_app with examples ([e6aa49f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e6aa49f6297a50afbbb0829f087c5ec4708b4cd6))
* added ios_store_app with examples ([#586](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/586)) ([0f9342d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/0f9342dd220c2e09582edaf78d7c0b11434cb49d))
* added mobile_app_relationship datasource and mobile_app_supersedence resource with example ([b5bdc2d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b5bdc2d80edca31c8d55b8a9b9ae9309cf72ac4c))
* added mobile_app_relationship datasource and mobile_app_supersedence resource with examples ([#591](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/591)) ([4531c1b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/4531c1b5b992fe74b379a6dffc07a4f592c83662))
* added utility_microsoft_store_package_manifest_metadata with examples ([0fb7804](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/0fb78041580686f9a4935d79e0da705e31d920b1))
* added utility_microsoft_store_package_manifest_metadata with examples ([#592](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/592)) ([28187fa](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/28187fa4b8c8b8e0d96e1f2775a5907312e57203))
* added win32_lob_app_with_examples ([ba42395](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ba423951160590485e02af4f9bd1582733aa0410))
* added win32_lob_app_with_examples ([#594](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/594)) ([ef9ee54](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ef9ee54e11b586aab828df896668625962258659))
* added windows_autopilot_device_identity with example ([6eb785a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6eb785a71fd73c0dbf87f4e2aed7e0aea7499432))
* added windows_autopilot_device_identity with example ([#596](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/596)) ([c3f2425](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c3f2425566e7300693b47cff3fb418ca6bf83f00))
* added windows_web_app resource and ios_ipados_web_clip with examples ([#590](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/590)) ([2b34fe1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2b34fe1d1ce469226c2e7b55aff4b7b76b2058f7))
* added windows_web_app resource with examples ([d59e19c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d59e19c7e6b685c43b5a83cbd900fbc1da4a662d))
* initial commit of vpp apps ([58b00f3](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/58b00f3169719f3d003fe523577a16f404ba0d09))
* initial implementation of microsoft_teams resources ([4ee45cf](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/4ee45cfc924e1f6de39d426b779e22090e6b049a))
* initial implementation of microsoft_teams resources ([#584](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/584)) ([912c14a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/912c14a61fb0f346d1ac55ad0cbc1a6c1eef140b))


### Bug Fixes

* for macos platform script assignments ([32a67de](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/32a67ded912795e02aebd4c107ba821fa6c44088))
* for windows 32 log app ([69d45c1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/69d45c1c757b7d406933798c91fb87739f99bd24))
* macos script assignments ([#581](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/581)) ([2635444](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2635444fd9ac8beb68aeac0c0f4eacd1152a9528))
* updated all mobile apps to support all image types and added http[s] validators with constants consistently ([89d6277](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/89d6277b081e1973e392a6b4cd08d5fb8aaa7443))

## [0.20.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.19.0-alpha...v0.20.0-alpha) (2025-07-03)


### Features

* added cloudpc alert rule ([3f04efd](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3f04efdb66497d71c94f2a60bbe2d83fc35b7d13))


### Bug Fixes

* updated docs ([#576](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/576)) ([4340d62](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/4340d62e2a1e26ee7faf478b1abd73773a1df0eb))

## [0.19.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.18.2-alpha...v0.19.0-alpha) (2025-07-02)


### Features

* added graph_beta_windows_365_azure_network_connection, graph_beta_windows_365_cloud_pc_organization_settings, graph_beta_windows_365_cloud_pc_provisioning_policy and graph_beta_windows_365_user_setting with examples ([#567](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/567)) ([005dbf8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/005dbf8585f2bbe06f9b263751373f563693f8d2))
* added windows 365 examples and bugies for move to graph beta ([5da3eaa](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5da3eaa34501da056131aa27982b95d5c67e79af))


### Bug Fixes

* added managed devices 1st itteration ([2355602](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/235560283c5b4344b676d23cc35bfaec6a4231c8))
* added numerous cloud pc datasources ([d1eef9d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d1eef9d883ea610b5c213bef5f65d6f118d00407))
* added numerous cloud pc datasources ([#557](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/557)) ([3e576e1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3e576e1e0c39052d70702b31fb42e7b580192a0f))
* added python venv to gitignore ([be5d039](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/be5d039574d85cf84432e8ea7a6e4a4d186040f4))
* added windows_365 datasources and added odata support ([08b3b8f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/08b3b8f76c46f1bb82ac7bca9554d6a0942aaf43))
* added windows_365 datasources and added odata support ([#559](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/559)) ([74c69e0](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/74c69e0203b0aa2793b50562d112b59775a4c321))
* example ([fd0bae5](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/fd0bae56cb1c529e158cebe7b981b608fe192aa3))
* for assignments with cloud provisoning polciies and tester ps script ([bde2160](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/bde2160635e9e6d3d0adcfa19a9ced1b394bfb2e))
* for branch merge ([07f4a5e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/07f4a5e0b0a5d6aeda02c64b6518914887a75eb7))
* for branch merge ([#560](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/560)) ([2481027](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2481027cda4a473df85e2102ad765fe829c66d37))
* for cloud pc provisioning profiles graph beta ([da6a74a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/da6a74af856acf168f4bca4a583a21ef467ad595))
* for conditional access policies ([8d7f5f0](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8d7f5f01304079faf2c6b567378432caf10f3e8c))
* for GroupSettings unit tests ([ffb8f71](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ffb8f719dedf3155f2d13242dda63dcb1dfaa580))
* for plan and state definitions within update funcs ([1cfce7a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/1cfce7a1a2d04f5eee804b14cbae506d3d8acb67))
* for pr merge ([bb7b5a2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/bb7b5a2e9de1026da924fc22c2d31125dca764e1))
* for pr merge ([#568](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/568)) ([a4cf05a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a4cf05a3d8aaf1a1098e5c68895b6a7d0a2c81da))
* for read with retry loop. correct ctx handling ([d2bf94b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d2bf94bc304006bc7b01333f1c59ca4c7f58e089))
* graph beta additions for windows 365 ([de472bb](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/de472bbaeb29dbbaf63fd18732ba99076540d152))
* implemented win365 provisoning policy fixes with frontline support and assignments. ([#561](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/561)) ([a0d66fe](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a0d66fe8c314d9b27d13598ba6c2bed493613059))
* implemented win365 provisoning policy fixes with frontline support. ([fe1ff4f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/fe1ff4f7b2d8b5c5041b1aabd8c205d7b08661c9))
* more fixes for conditional access unit tests ([a916947](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a9169476781401a45e6a540e7348b0ec21e671d5))
* numerous fixes for conditional access policies. uses a new non kiota based sdk http client ([f62c5cc](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f62c5cc95be1e96192929098e70d910ad55e96f8))
* numerous fixes for conditional access policies. uses a new non kiota based sdk http client ([#553](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/553)) ([ab62461](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ab624617b6f329c42cb17df17866970b9412f372))
* numerous fixes for read with retry logic to now correctly handle operation type ([#562](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/562)) ([16275b8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/16275b84a72254328a0813e19630fd03c87baf87))
* partial fix for conditional access introduced with sdk v0.137.0 ([#550](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/550)) ([82452cc](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/82452cc393627f84402e5e8181bc91584d1eb421))
* re-up ([9fa33bf](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9fa33bf3a68a935a28208942de6b319742a24430))
* re-up ([#569](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/569)) ([e30b87c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e30b87cb931cddc4d4848515a49bd706452c913e))
* set unit test timeout to 120 minutes ([e39a0af](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e39a0afa4d436c9b528a8fe31bd0039ea2e7f6d3))
* state removal logic and logging consistency across all resources ([b77d94a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b77d94ac10c980e5daf860c3c40bdd14fdc4e9f9))
* tflog logs for update func ([109e2cc](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/109e2ccadaaf6617ecdf0c5d93f82ccb24f3c446))
* tidy up ([a378df2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a378df2d4f1e4579ebe3ef2ae3cb47a42d5cea73))
* unit tests for conditional access policies ([2fcabab](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2fcababe1f4cfa2924b0530f5ade95df5a2222ac))
* unit tests only run now on changes to go files ([47ff1f5](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/47ff1f54d7ddd1b799da8ae44679cb8e2b52896a))
* updated docs for groups resources ([f733e4c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f733e4c3fe4ae71a0e99237d4a6e983e689a2c6d))

## [0.18.2-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.18.1-alpha...v0.18.2-alpha) (2025-06-25)


### Bug Fixes

* tf-registry-goreleaser.yml permission fix ([#548](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/548)) ([ef2e54e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ef2e54e157242ef55afe5e9aac1e141932bead17))

## [0.18.1-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.18.0-alpha...v0.18.1-alpha) (2025-06-25)


### Bug Fixes

* added group acceptance and unit tests ([6b0bc50](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6b0bc5029011509f60c167911672f34712ff28d9))
* added unit and acc tests for groups and groups member assignment ([aef90f6](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/aef90f6cc2397a7f0ee066083ba7fbee638e14ab))
* added unit and acc tests for groups and groups member assignment ([#545](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/545)) ([99759a4](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/99759a48b342313f9d6e26a7d2eda40f769282e5))

## [0.18.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.17.1-alpha...v0.18.0-alpha) (2025-06-24)


### Features

* added numerous acc and unit tests ([#541](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/541)) ([5be07af](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5be07af61c225192c817139f90c36a0bfb71d986))
* added unit and acc testing for users, m365_admin_install_options , license assignment ([a558b64](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a558b6458aa0c1f2a2192ddc2d60f064c777dbcc))
* added unit and acc tests for resource user and resource license assignment ([76058c0](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/76058c034b568a2cdc081aafc06e432c3cc9a647))


### Bug Fixes

* for true to false ([369f494](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/369f49448fb9ed595c7c4736c8728ad19c93174f))
* for unit tests for microsoft365_graph_beta_device_management_macos_software_update_configuration ([ced892f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ced892f18331710ee05be492c1788c9d47788675))

## [0.17.1-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.17.0-alpha...v0.17.1-alpha) (2025-06-23)


### Bug Fixes

* added docs for device config assignment graph v1.0 ([ef66aab](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ef66aab7ae4c7720e18016d2380cf87bd6bad5d0))
* added docs for device config assignment graph v1.0 ([#536](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/536)) ([788f7f1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/788f7f1d38fa728ab864987d2bb6c06460decb54))
* increased timeout for goreleaser ([709cbb2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/709cbb29a33a7062d699c1e82427a3c5f9ea269a))
* increased timeout for goreleaser to 90 minutes ([#538](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/538)) ([a528317](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a5283177fe997a903a4f072d159d118d32f8044a))

## [0.17.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.16.1-alpha...v0.17.0-alpha) (2025-06-23)


### Features

* Add macOS LOB app publisher validation and update gitignore ([c9cd922](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c9cd922080ec768367bc79f1fc73632c9248d2f4))
* added 1st itteration of oidc support for github actions and azdo ([3a2c14b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3a2c14bdcc51b039fc9269e9febf8edae89fa866))
* added additional extraction points for graph metadata ([3893926](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3893926ae66dca2ed555653745049b530d0fbc32))
* added ai instructions set ([af64371](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/af64371837294cfc9e11d78ed2bf2395c661fedd))
* added application category resource and datasource with examples ([55a099a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/55a099a087f713f4e4c01558e4aec5fdfaf303f8))
* added assignments for windows updates as a seperate set of resource types ([837ee80](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/837ee809bc587496f2e8be44996db786cbd0d10f))
* added configuration policy unit tests ([cc5a938](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/cc5a93863fa5112832619bdf1f897e0ba6966b66))
* added datasource graphBetaDeviceAndAppManagementWindowsQualityUpdateProfile ([9f6f79d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9f6f79d7472eef5d7f76973c8195046b65c4462b))
* added datasource microsoft365_utility_macos_pkg_app_metadata ([2c2cbca](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2c2cbca7feb8d1336c9327e487767ecec1e4bc35))
* added dev docs ([d31361f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d31361fa979172bed9709402f76d8fbc3616b32d))
* added developer docs ([6b1717a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6b1717a81351b40460835e8f9da819d32a7cf7c2))
* added developer docs for resource development ([#521](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/521)) ([2b5178c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2b5178c7cfd28ef3df52ae9fe9672a4d0623d999))
* added device_and_app_management_windows_update_catalog_item datasource with examples ([655e804](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/655e804d7ea7989c549fa86a8925ea82b0b33d46))
* added graph metadata extraction script ([cef8591](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/cef8591dc8904959fb4844c06aa72d1770f1fe6b))
* added graph_beta_device_and_app_management_windows_driver_update_profile_assignment ([73dc481](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/73dc481446cb3c88284ef6b5fcac719da08bfb81))
* added graph_beta_device_and_app_management_windows_feature_update_profile with example ([050cec1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/050cec1724781fa5bdf81b14c9c5890f5a0386e4))
* added graph_beta_device_management_apple_user_initiated_enrollment_profile_assignment with examples ([5ae2eca](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5ae2eca5809e0083d525cd9695002aab18ed886e))
* added graph_beta_device_management_apple_user_initiated_enrollment_profile_assignment with examples ([#503](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/503)) ([e4cade6](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e4cade6e0009400e6ec09ffe9b8935fe02b5cadf))
* added graph_beta_device_management_macos_software_update_configuration with examples ([33247a9](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/33247a91e23b12381eaa41c2364448c26d5e502a))
* added graph_beta_device_management_macos_software_update_configuration with examples ([#506](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/506)) ([d82a919](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d82a9197b67a1b1625a32f800d153dd738daf860))
* added group license assignment, user license assignment resources and datasource license subscribed SKUs ([bede181](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/bede1812ddd65df921206baf67602e02693e1f2a))
* added group ownership assignment with examples ([18810ea](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/18810ea49955060d9bfb87947a2b25a4e560c3ba))
* added group ownership assignment with examples / renamed user license assignment to follow resource naming pattern ([#497](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/497)) ([6e2fdf3](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6e2fdf3e6411537f679f5201abf9db426f82fd5b))
* added GroupMemberAssignment resource with examples ([21dde58](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/21dde582a0048b75e462dd881831d08706ef8069))
* added GroupOwnerAssignmentResource ([4f9bfd6](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/4f9bfd60e1013b7030ac31187b42fdd5224d78a7))
* added GroupOwnerAssignmentResource ([#499](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/499)) ([9bc5dee](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9bc5dee1a016aa39f78acbee5b91236d4c95d267))
* added groups with examples ([ffc38e3](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ffc38e3db9b488cbfb86b062d40bc946c9ebe22a))
* added groups, and group membership assignment, updates to doc descriptions ([#495](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/495)) ([0601507](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/0601507bf998791c5f74ae433005a492d0ecf1b0))
* added helper SetObjectsFromStringSet for object construction ([e2b110e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e2b110e4f7b4aa73aa2e823a5ae9b58f5b597e8c))
* added intune device categories resource and data source with examples ([#389](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/389)) ([b8d33fb](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b8d33fba7b11121bf61514164b77d01858a9a74c))
* added macos pkg app datasource ([a65ca03](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a65ca038103e064225b7d136205179b7fbb04fad))
* added macos pkg datasource ([f0f6144](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f0f61447b307f221cdf6f0b0d65f74d8e4f6beb6))
* added macos_dmg_app fixes and updated resource examples ([1da2c8d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/1da2c8d26bd0963fdd48d72ff384cb0245a62be2))
* added makefile for tf development ([36b7a0a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/36b7a0a57f687dbd96bd1c521c745078cbed2977))
* added makefile, ai instruction set, refactored data type conversion helpers, scaffolding for unit and acc tests with mocking ([#528](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/528)) ([6788408](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6788408661ed58258ee10589b6415f7ddab18800))
* added managed device clean up rule with examples ([fc33266](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/fc332665c817df572c0d349ec4d68ccf1b037570))
* added managed device clean up rule with examples ([#463](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/463)) ([f1a8777](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f1a87770fa8e61a0edbc8370876cf86b4cb07850))
* added microsoft365_graph_beta_device_and_app_management_windows_update_profile_assignment with example ([44e3539](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/44e3539f943013c5a0414e2053f56a1c02757f04))
* added microsoft365_graph_beta_device_management_terms_and_conditions ([335e4a1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/335e4a1ec8441e7f49b980af60de4d07e2afea8f))
* added microsoft365_graph_beta_device_management_terms_and_conditions ([#466](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/466)) ([16849a7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/16849a77e9262aeddbae84bc43b55552da62c146))
* added microsoft365_graph_beta_device_management_terms_and_conditions_assignment with examples ([c3db9ac](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c3db9acb7c2fb11122658bdf817c2ff57af01fec))
* added microsoft365_graph_beta_device_management_terms_and_conditions_assignment with examples ([#468](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/468)) ([6e407c0](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6e407c067f71cbc2f01cef185adaff9f214774db))
* added microsoft365_graph_beta_device_management_windows_feature_update_profile_assignment with examples ([74f286e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/74f286eb668a342f700b1ce3c5da9671695ffdcb))
* added microsoft365_graph_beta/microsoft365_graph_beta_device_management_windows_driver_update_profile_assignment with examples ([863774c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/863774cbb6a0b9944c1d79b2e87ac2f4ad66a6b5))
* added mobile app data source with odata queries ([3a1c0f4](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3a1c0f414e70ea9c55f36b31c1c48e85feb2aa02))
* added multiple new resources for license assignment and fixes for macos dmg app ([#491](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/491)) ([8d09b67](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8d09b67ec036a30817aecb7745d050adc96c880c))
* added multiple onboarding resources ([602135e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/602135e85aeab9ad2f8923c437004be647357d58))
* added operation approval policy ([40b307c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/40b307cfd8c1ec89fd342c7a29f20254d2cacaef))
* added ps graph metadata extraction script ([#522](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/522)) ([a3ca194](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a3ca1944e6bbd89c00dca90a327a1188e4074aa1))
* added refactored unit tests for provider build and credential factory ([bd7691c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/bd7691c76d6403dadb4eb68d76ba37093f506b9b))
* added refactored unit tests for provider build and credential factory ([#439](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/439)) ([c191a93](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c191a9345c962f5f19a1bfd76dcafb3fb1deec01))
* added resource operation for intune rbac with examples ([4f36c5f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/4f36c5f5b50c9d19f3cdcdcc2507d8d8929a5931))
* added resource operations with examples for intune rbac ([#465](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/465)) ([1c31571](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/1c315714f220e136300bf4da2f0c029dd39e2c11))
* added setting catalog in native hcl ([596d6c9](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/596d6c9f0454810004ff30a3177a5c4ec2bc4ebc))
* added settings catalog in native hcl ([#461](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/461)) ([6e5e693](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6e5e6937324d83df7fe5f96445dc45f851bf771f))
* added timeouts to datasources ([#397](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/397)) ([100840f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/100840f1bff363a598fd9a76e26183c16c086150))
* added unit tests for configuration policy settings go and json builders ([#511](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/511)) ([3293255](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/32932553d3ce83264e3c0d618f5ca0db5681bc09))
* added validation for mobile app assignments for 8 invalid scenarios ([cc88fd1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/cc88fd104a35cf071e4a5693a84420f57d1151d8))
* added windows autopilot and assignments thereof ([#488](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/488)) ([5ffe081](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5ffe08175062fa490c4b7babc5e8ecf658cc3ba3))
* added windows driver update profile and inventory resources and datasource ([#416](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/416)) ([8e3272d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8e3272dbc77296cdd874a11e7c7693103173419b))
* added windows driver update profile and inventory resources and datasources ([85ff58f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/85ff58f605f939be963c3aba58ebd44c0ac41f0f))
* added windows quality update profile assignment resource with examples ([a49a109](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a49a109323d23f0c8e0ed2be4afea53a77b1b05d))
* device_enrollment_configuration with examples ([#420](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/420)) ([88516db](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/88516db1ff74964069d9965f29c21972bb147144))
* finalised device preparation policies until msft fix their api ([b9949b1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b9949b141afcd3822012fbaf850c417ed4b51b93))
* graph_beta_device_management_device_enrollment_notification_configuration with examples ([ff7fecb](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ff7fecbdb1f77970dabf1994fc47269cd2cc58c1))
* graph_beta_device_management_device_enrollment_notification_configuration with examples ([#501](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/501)) ([145625d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/145625d2510a28544419c69beab4eb777da6e3cc))
* half implementation of cloudPC as a datasource ([68e7a43](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/68e7a43dc062a7808d2265f7b76c7a354fff5d85))
* implemented device categories resource and datasource with examples ([acf2626](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/acf262635ca9aeb5175f0744fc108dfd0694b027))
* implemented refactor state and constructor funcs ([307fb91](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/307fb91f0966100ac6517915ac263a6db1dc4b96))
* initial implementation of mac_dmg_app ([35b6fe6](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/35b6fe60dd97db4c75b442387d26c26680f61afd))
* mobile app assignments ([80c892e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/80c892ec95186025afc2c6f7edbd61f7ed1863e4))
* mobile app assignments can now be set in any order with modify plan (diff supression) ([c8f97ef](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c8f97ef7201c1a79efc3850add74eac0e78468be))
* mobile app assignments resource with examples ([#453](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/453)) ([498ff50](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/498ff5087f942dffa701be9da736c8d3e1a5dde2))
* mocking scaffolding ([f1be6d1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f1be6d1f6fddfca091d34df19750da364545c1b1))
* mocks scaffolding ([#440](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/440)) ([9ee57df](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9ee57dfc4eecebb95c35c0e60686925d3228c2c6))
* moved all client initization to package client and refined client init logic with provider configure func ([45de96e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/45de96e0ef0052544f38055b483c138db8698d84))
* moved all client related logic for building to the client package from the provider package ([9e2561c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9e2561cbcaa78b5458eae138ca025609a3416f87))
* moved custom requests to it's own package within common out of the client package ([356f090](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/356f09017278074550b35a4b4f5e4a3aac175356))
* refactored data source for maospkgapp ([#442](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/442)) ([fec56b7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/fec56b72934dc5c9fbd1dad05e8d690fab0793ab))
* refactored directory hirachy so that all client related logic is within the client package. moved all custom request logic to it's own package within custom from client package ([#529](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/529)) ([8754a07](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8754a071caed0fd816f7080b840e4f5003c07981))
* refactored macospkg app datasource with example ([cbc1681](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/cbc1681feca36a096ea63b69dc61540178b21447))
* refactored retry mechanism for create update func calls to read and fixed resource name issue in tflog ([6f3842b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6f3842b684eab3bcde146a90a9afd4b7b1c1cb34))
* refactored retry mechanism for create update func calls to read and fixed resource name issue in tflog ([#498](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/498)) ([5257e5e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5257e5ef5653b0ed0ca2ba5fe4ff1c70c3a5cf29))
* refactored unit and acc test strategy to use local responders and centralised factories for common opperations ([c878bda](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c878bdaaf8fa7c65ce69323323a27dec6ff340bd))
* refactored unit and acc test strategy to use local responders and centralised factories for common opperations ([#534](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/534)) ([8179a7d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8179a7d00b43e34f21e7b941fea16b90ca0d4d76))
* repo structure to move datasoources and resources within a common services folder. moved common to services folder in folder hierarchy ([0799a5b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/0799a5b52ca34d78f186c7c6068b7d0abf47313c))
* repo structure to move datasources and resources within a common ./internal/services folder. Moved the common folder with shared assets to ./internal/services from /resources ([#524](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/524)) ([22520ec](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/22520ec6ed4492addcc0a0ac1d205a2fc1720da1))


### Bug Fixes

* added close stale issues maintainence pipeline ([1d90efd](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/1d90efd55aab386abff032eda94948987f69d50f))
* added close stale issues maintainence pipeline ([#475](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/475)) ([40affe2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/40affe22d27c45128e8b4f7887419f10d41efc58))
* added discord info to readme ([88c0a53](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/88c0a53746578ceb5479c71a5c9e61bb83485f7a))
* added discord info to readme ([#519](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/519)) ([aded08c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/aded08ca82a4a35e3ba2e111b96f9ed5e140571b))
* added examples and more fixes for windows_remediation_script_assignments ([9827467](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9827467ac9972380540ccb29490670d5db960c87))
* added mutex to all api calls to defender against Kiota's middleware (e.g., HeadersInspectionHandler) when modifing shared header maps during HTTP request processing. appears to be a bug in the latest version ([976ad22](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/976ad22dea6b1ddf6e3027714a04c449b73ac11c))
* added mutex to all api calls to defender against Kiota's middleware (e.g., HeadersInspectionHandler) when modifing shared header maps during HTTP request processing. appears to be a bug in the latest version ([#450](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/450)) ([10327c8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/10327c856db27c79ae6883677a48049c6caf2a1e))
* added new regex for ISO8601Duration and comment fix ups ([1e11d71](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/1e11d71d4890d1adaf6f3901c6ed2268264b597a))
* added planmodifiers.UseStateForUnknownString(), for id ([dccde7e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/dccde7e78d92bf90a23bbdf7ad5673391baafc65))
* added read timeouts for datasources ([243e6b2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/243e6b2fdc92441d54949d0d1bd9544b4b479c1e))
* base resource updates for macos_pkg_app doesnt use redundant file size checks ([6162a1d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6162a1d4dc9860de2aa812ff386289bd6fc387f7))
* categories fix for macos app categories and stating issues for infured included app values ([defa936](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/defa93679fa6bed0bc6065678987c76429383133))
* changed all actions to use commit sha's ([cd48ede](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/cd48ede424a0363694107d44fb6fe64823bce198))
* content version ([74c417b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/74c417b8312e6353bbbdbf78fd012d9fa96b39e1))
* defined default values and plan modifier for role_scope_tag_ids ([e8909e0](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e8909e0fc876dbe03422bfb26943ac872cfab432))
* determined that the api for AutopilotDevicePreparationPolicy is broken.sigh ([481dcdd](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/481dcddaab9b8510b0740960919ff1c7725d98e4))
* docs ([0a676f8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/0a676f87825c2d75b2c3ed5175d48474e23ec28d))
* examples so that device_management examples are reflective ([d6856b4](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d6856b480ef5a0c9f7277592bc391eb47487b443))
* examples so that device_management examples are reflective ([#441](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/441)) ([6062154](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6062154d0eedb88177fec252a4c5fb5f925ca24a))
* extended go lint timeout to 20minutes ([e5f9b3b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e5f9b3bbf42b420ba1568a9042930c281d82523b))
* final fixes for macos assignments ([8d998cf](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8d998cfc27ccb4aa1b4bab3be60823f479b7cf85))
* final fixes for macos assignments ([#451](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/451)) ([7e3b853](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/7e3b853731d22cc832840ceab3ca9c7a256418c4))
* finalised device preparation policies until msft fix their api ([#505](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/505)) ([089d398](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/089d398d18efd2a065d5dcdd4ed191bda458fab5))
* for deps bump of go sdk version ([54a119e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/54a119edcca117fa697b21c59f2ba35cdda12ae2))
* for doc gen pipeline ([f959ac8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f959ac87f1d6ea3e0d66b9979502a5c1f914466f))
* for doc gen pipeline ([d3fecdf](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d3fecdf388ed6ad8e82fffffce38cbdd0126553a))
* for doc gen pipeline ([#471](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/471)) ([c93301e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/c93301e7f95b32b56ed812ba5f373d2c2d3b0d80))
* for doc gen pipeline ([#473](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/473)) ([3c58522](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3c58522f6f1522eca9da2b4ee613b32b96319368))
* for docs path ([2fce669](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2fce66923162574af62abe0c7236cf9da8bd76ab))
* for enum handling within GetDeviceEnrollmentConfigurationType ([#437](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/437)) ([896c291](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/896c2918aa799ec20127a79538a06cfe1cff892e))
* for example path ([771bb1c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/771bb1c2792a16955e05f124cbe7a4e6d88cd3bb))
* for examples ([244a04b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/244a04b7ef7521d33ffcf42e0af0015d4d2791c8))
* for go lint failure handling ([e3e52be](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e3e52be832b1ca02cb214c3c0bbb7d6d5b8c57e9))
* for GraphToFrameworkISODuration helpers to resolve unit tests ([e8338a7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e8338a764a864067ad78aef8ed9e14be86f4c592))
* for handling enum ([ce13f69](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ce13f693c89723f9766c02519f17596ea328d7e2))
* for handling enum ([6e2f1ff](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6e2f1ff38ef16ddbcd405eff88c380371e315856))
* for macos pkg apps ([#395](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/395)) ([deb9309](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/deb9309edc8a8aff4d33a022102f30f9c90eb679))
* for mobile app assignment id's and content version stating ([a57a269](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a57a26945a184c294e4c37c41b6b3457dbdad180))
* for mutex handling for SetCommittedContentVersion ([058bf32](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/058bf32f301bc9b5082f668427aaf4ca3cff3878))
* for parrallism issues with concurrent header map settings in client mobile apps ([49561f6](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/49561f64cea87bd23eab9dda180824ddfd8934de))
* for pipeline permissions ([a8b6990](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a8b69906db84a5a90db4b9f9373a64e4dabf9a08))
* for role_scope_tag_ids ([#444](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/444)) ([b1e0fd8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b1e0fd819d4956ee065e96d4300fb24785719dc2))
* for successful stating logging with resource name and resource id ([8c81932](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8c819320ff467a553640c8ff5ce9c78c5f7cf76a))
* for winget assignments ([8a2bdd0](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8a2bdd03111f263f47db1614ba1e0fb8deb5a650))
* Go Lint configuration yml for provider project types ([#424](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/424)) ([b92b8ac](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b92b8ac28100a40611b0cca6e191750db9b8d277))
* go unit test pipeline trigger ([e972c36](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e972c3658f3f1a4c74f29fa55657ccef63117d5b))
* go-lint ([#423](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/423)) ([57ba99d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/57ba99dd57bf8b26e9785767cf05449df1b2fca2))
* improved markdown descriptions for documentation generation ([16bc1d2](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/16bc1d2d7324c2a5c0a224a274e73100648b4caa))
* m365 app installl options ([e2db1a7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e2db1a7982f5cced9b3f346b7f2df8a30c92695a))
* migrated microsoft365_graph_device_and_app_management_m365_apps_installation_options to graph v1.0 ([a3c9682](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a3c968250d67787711d46b2589facd62cf3ab8b3))
* migrating winget app assignments to sets ([2272cc6](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2272cc64166db8c02accd2a243d79a646fc3dfe4))
* minor fix for sorting order of assignments ([fb50e0f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/fb50e0fbac88e72d4d784f400131449be1c4ff07))
* misc naming fix ups ([502ba24](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/502ba24503a3473bd3351a9ae31cb667ee5b409c))
* moved mobile app assignment validation to create and update funcs to fail early if any hcl mistakes before any rresources get deployed ([8c79958](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8c79958155e1445d959487a27cd66ddca38947b3))
* numerous bug fixes and unit test additions ([44bb11a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/44bb11a80d63c5d3f2b1a7a726d32f39387591a7))
* numerous edits and tweaks ([dccd487](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/dccd4873fc4dec676b2a0e3859c7a2a42d595b4f))
* numerous fixes and scaffolding for mock tests ([f951cda](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f951cdaeec2b4e5401f447731cd383011963f9b7))
* numerous schema fixes ([ddb6e87](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ddb6e8780360a63c31354f0899314fd4ed51c111))
* pipeline permissions and tidy up ([#532](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/532)) ([475c686](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/475c686ffce21927d8b72d3419d5fada6e79766c))
* pipeline step naming ([3fade5e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3fade5e96ab87097c3b8b18dc3e2f0ea5850a7a3))
* pipeline step naming ([#425](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/425)) ([8067017](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/80670170abf8aa817432805c0d4a786859655951))
* python test scripts ([591c402](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/591c40260095757c86431a7888715f3b3f4b50c3))
* read with retry test ([3053bd5](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3053bd530f85440aa1ac41911e6b94d4a5347ed7))
* readme's and schema docs for tf registry ([afbab65](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/afbab654a7ad3e83868e0b7f52fb7d0849444a28))
* readme's and schema docs for tf registry ([#520](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/520)) ([bc16e69](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/bc16e692736e681bf81d45f348156e3e18763c1f))
* rearraging provider docs ([a0aab65](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a0aab656a60e99ef2554e053d01fdad5c4539016))
* refactored docs for repo restructure ([546ebb4](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/546ebb4ff5dfe06b9b6c9febad4223a4acf71abf))
* refactored docs for repo restructure ([#435](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/435)) ([3b14e78](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3b14e78f1db1659b176983abced03e04b32195d6))
* refining logic for handling app metadata ([9ecf223](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9ecf2236482cffbe5ef3cc67cf9e10e94068f345))
* regex const naming patterns ([54b5af1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/54b5af124f7a2b97901be6c7b3144707f13de508))
* removed conditional access. complete refactor required with latest sdk version ([4cb550d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/4cb550d56f376cb30bde71c71b4699e63653e07c))
* removed mutex, going with parallism of 1 instead ([dd0e6d8](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/dd0e6d822dcf3eb6fd43bbdea105d62020dc488f))
* renamed user license assignment to follow resource naming pattern ([29f4c8b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/29f4c8b163dd7e1d95dfe045831fc64e125db76f))
* restructured repo to follow api endpoint paths more precisely ([1b2c110](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/1b2c1108add87c7395eaed7cfce79887c7899f12))
* restructured repo to follow api endpoint paths more precisely ([#433](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/433)) ([6aed990](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6aed9906d8f76430c646909e6111c6d91b29f1b3))
* tidy up pipeline step naming ([102e356](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/102e35695ac0629381b740e08aec981cfdaf84a1))
* tidy up repo restructure ([b5b35cf](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b5b35cf4b001d4ba9ebc0cfed8913bc8a7054e12))
* typo fix for datasource ([5c60f13](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5c60f13a9ce6d8701d9310ea41f07e1104524967))
* unit test trigger ([2bf41c4](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2bf41c4a70e4a47e268bb9d68950dda48d24bfe0))
* unit test trigger ([#474](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/474)) ([e426630](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e42663025764a290342fde1a4f489a19931e94a1))
* unit tests ([e0dd719](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e0dd7194b4dd58cf18a241f927119e8dd185c630))
* unit tests ([5ffcd30](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5ffcd30efd1beb34be608df2095a0eda5e632c9f))
* updated create and delete to use the crud.ReadWithRetry pattern ([83003fd](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/83003fd8275e537495cf3a94865b425e9712f7d8))
* win_get_app now skips property 'InstallExperience' and 'PackageIdentifier' as it cannot be patched during update operations. and added support for app categories ([8b5a4a5](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8b5a4a59a43e6bddb6e82a3763713a85bcfe4047))
* windows remediation script assignments ([04243b6](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/04243b62ecb12ed869edec6caf334e6671c8fd90))

## [0.16.1-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.16.0-alpha...v0.16.1-alpha) (2025-06-21)


### Bug Fixes

* for pipeline permissions ([a8b6990](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a8b69906db84a5a90db4b9f9373a64e4dabf9a08))
* pipeline permissions and tidy up ([#532](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/532)) ([475c686](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/475c686ffce21927d8b72d3419d5fada6e79766c))

## [0.16.0-alpha](https://github.com/deploymenttheory/terraform-provider-microsoft365/compare/v0.15.0-alpha...v0.16.0-alpha) (2025-06-21)


### Features

* added additional extraction points for graph metadata ([3893926](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3893926ae66dca2ed555653745049b530d0fbc32))
* added ai instructions set ([af64371](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/af64371837294cfc9e11d78ed2bf2395c661fedd))
* added configuration policy unit tests ([cc5a938](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/cc5a93863fa5112832619bdf1f897e0ba6966b66))
* added dev docs ([d31361f](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d31361fa979172bed9709402f76d8fbc3616b32d))
* added developer docs ([6b1717a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6b1717a81351b40460835e8f9da819d32a7cf7c2))
* added developer docs for resource development ([#521](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/521)) ([2b5178c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/2b5178c7cfd28ef3df52ae9fe9672a4d0623d999))
* added graph metadata extraction script ([cef8591](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/cef8591dc8904959fb4844c06aa72d1770f1fe6b))
* added graph_beta_device_management_apple_user_initiated_enrollment_profile_assignment with examples ([5ae2eca](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5ae2eca5809e0083d525cd9695002aab18ed886e))
* added graph_beta_device_management_apple_user_initiated_enrollment_profile_assignment with examples ([#503](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/503)) ([e4cade6](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e4cade6e0009400e6ec09ffe9b8935fe02b5cadf))
* added graph_beta_device_management_macos_software_update_configuration with examples ([33247a9](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/33247a91e23b12381eaa41c2364448c26d5e502a))
* added graph_beta_device_management_macos_software_update_configuration with examples ([#506](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/506)) ([d82a919](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/d82a9197b67a1b1625a32f800d153dd738daf860))
* added makefile for tf development ([36b7a0a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/36b7a0a57f687dbd96bd1c521c745078cbed2977))
* added makefile, ai instruction set, refactored data type conversion helpers, scaffolding for unit and acc tests with mocking ([#528](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/528)) ([6788408](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/6788408661ed58258ee10589b6415f7ddab18800))
* added ps graph metadata extraction script ([#522](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/522)) ([a3ca194](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/a3ca1944e6bbd89c00dca90a327a1188e4074aa1))
* added unit tests for configuration policy settings go and json builders ([#511](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/511)) ([3293255](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/32932553d3ce83264e3c0d618f5ca0db5681bc09))
* finalised device preparation policies until msft fix their api ([b9949b1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/b9949b141afcd3822012fbaf850c417ed4b51b93))
* graph_beta_device_management_device_enrollment_notification_configuration with examples ([ff7fecb](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/ff7fecbdb1f77970dabf1994fc47269cd2cc58c1))
* graph_beta_device_management_device_enrollment_notification_configuration with examples ([#501](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/501)) ([145625d](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/145625d2510a28544419c69beab4eb777da6e3cc))
* implemented refactor state and constructor funcs ([307fb91](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/307fb91f0966100ac6517915ac263a6db1dc4b96))
* moved all client initization to package client and refined client init logic with provider configure func ([45de96e](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/45de96e0ef0052544f38055b483c138db8698d84))
* moved all client related logic for building to the client package from the provider package ([9e2561c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/9e2561cbcaa78b5458eae138ca025609a3416f87))
* moved custom requests to it's own package within common out of the client package ([356f090](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/356f09017278074550b35a4b4f5e4a3aac175356))
* refactored directory hirachy so that all client related logic is within the client package. moved all custom request logic to it's own package within custom from client package ([#529](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/529)) ([8754a07](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/8754a071caed0fd816f7080b840e4f5003c07981))
* repo structure to move datasoources and resources within a common services folder. moved common to services folder in folder hierarchy ([0799a5b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/0799a5b52ca34d78f186c7c6068b7d0abf47313c))
* repo structure to move datasources and resources within a common ./internal/services folder. Moved the common folder with shared assets to ./internal/services from /resources ([#524](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/524)) ([22520ec](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/22520ec6ed4492addcc0a0ac1d205a2fc1720da1))


### Bug Fixes

* added discord info to readme ([88c0a53](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/88c0a53746578ceb5479c71a5c9e61bb83485f7a))
* added discord info to readme ([#519](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/519)) ([aded08c](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/aded08ca82a4a35e3ba2e111b96f9ed5e140571b))
* added new regex for ISO8601Duration and comment fix ups ([1e11d71](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/1e11d71d4890d1adaf6f3901c6ed2268264b597a))
* changed all actions to use commit sha's ([cd48ede](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/cd48ede424a0363694107d44fb6fe64823bce198))
* determined that the api for AutopilotDevicePreparationPolicy is broken.sigh ([481dcdd](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/481dcddaab9b8510b0740960919ff1c7725d98e4))
* extended go lint timeout to 20minutes ([e5f9b3b](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e5f9b3bbf42b420ba1568a9042930c281d82523b))
* finalised device preparation policies until msft fix their api ([#505](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/505)) ([089d398](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/089d398d18efd2a065d5dcdd4ed191bda458fab5))
* for go lint failure handling ([e3e52be](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e3e52be832b1ca02cb214c3c0bbb7d6d5b8c57e9))
* for GraphToFrameworkISODuration helpers to resolve unit tests ([e8338a7](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e8338a764a864067ad78aef8ed9e14be86f4c592))
* go unit test pipeline trigger ([e972c36](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/e972c3658f3f1a4c74f29fa55657ccef63117d5b))
* misc naming fix ups ([502ba24](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/502ba24503a3473bd3351a9ae31cb667ee5b409c))
* numerous bug fixes and unit test additions ([44bb11a](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/44bb11a80d63c5d3f2b1a7a726d32f39387591a7))
* numerous fixes and scaffolding for mock tests ([f951cda](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/f951cdaeec2b4e5401f447731cd383011963f9b7))
* python test scripts ([591c402](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/591c40260095757c86431a7888715f3b3f4b50c3))
* read with retry test ([3053bd5](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/3053bd530f85440aa1ac41911e6b94d4a5347ed7))
* readme's and schema docs for tf registry ([afbab65](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/afbab654a7ad3e83868e0b7f52fb7d0849444a28))
* readme's and schema docs for tf registry ([#520](https://github.com/deploymenttheory/terraform-provider-microsoft365/issues/520)) ([bc16e69](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/bc16e692736e681bf81d45f348156e3e18763c1f))
* regex const naming patterns ([54b5af1](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/54b5af124f7a2b97901be6c7b3144707f13de508))
* unit tests ([5ffcd30](https://github.com/deploymenttheory/terraform-provider-microsoft365/commit/5ffcd30efd1beb34be608df2095a0eda5e632c9f))

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
